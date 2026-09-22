// aftersale_impl_test.go 售后域（013-after-sale / 批次 07）——状态机每一跳 + 幂等 + 边界。
// fixture: 一笔"已完成(40)"订单 + 单个订单项（数量/单价可控，便于验证部分退的按比例折算与尾差）
// + 一行库存（用于验证退货退款的回补）。全部自清，判据为精确值。
package shop

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/library/money"
	"ecboot/internal/library/paychannel"
	"ecboot/internal/model"
)

// afterSaleFixture 售后测试数据（订单/订单项/库存/用户）。
type afterSaleFixture struct {
	UserId   int64
	OrderId  int64
	OrderNo  string
	ItemId   int64
	SkuId    int64
	Quantity int    // 行购买数量
	PayYuan  string // 行实付（元，= 单价 × 数量）
	phone    string
}

// seedAfterSaleFixture 建 fixture（qty=行数量, skuId=库存定位用哨兵, unitYuan=成交单价）。
func seedAfterSaleFixture(
	ctx context.Context, t *gtest.T, phone, orderNo string, skuId int64, qty int, unitYuan string,
) *afterSaleFixture {
	cleanupAfterSaleByOrder(ctx, orderNo)
	uid := seedPointUser(ctx, t, phone)
	f := &afterSaleFixture{UserId: uid, OrderNo: orderNo, SkuId: skuId, Quantity: qty, phone: phone}

	unitFen, err := money.FromYuanString(unitYuan)
	t.AssertNil(err)
	payFen := unitFen * int64(qty)
	f.PayYuan = money.ToYuanString(payFen)

	res, err := g.DB().Exec(ctx,
		"INSERT INTO trade_order(order_no,user_id,status,total_amount,promotion_amount,pay_amount,currency,"+
			"receiver_name,receiver_phone,receiver_province,receiver_city,receiver_detail,refund_status) "+
			"VALUES(?,?,40,?,0,?,'CNY','售后测试','13800000000','浙江省','杭州市','T路售后',0)",
		orderNo, uid, f.PayYuan, f.PayYuan)
	t.AssertNil(err)
	oid, _ := res.LastInsertId()
	f.OrderId = oid

	res, err = g.DB().Exec(ctx,
		"INSERT INTO trade_order_item(order_no,order_id,spu_id,sku_id,sku_no,spu_name,sku_name,sku_specs,"+
			"quantity,original_price,price,coupon_amount,full_reduction_amount,point_amount,promotion_amount,pay_amount) "+
			"VALUES(?,?,1,?,'TF-AS-SKU','售后商品','售后商品 默认','{}',?,?,?,0,0,0,0,?)",
		orderNo, oid, skuId, qty, unitYuan, unitYuan, f.PayYuan)
	t.AssertNil(err)
	iid, _ := res.LastInsertId()
	f.ItemId = iid

	_, err = g.DB().Exec(ctx,
		"INSERT INTO inventory(sku_id,total,locked,warn_count) VALUES(?,100,0,0) "+
			"ON DUPLICATE KEY UPDATE total=100, locked=0", skuId)
	t.AssertNil(err)
	return f
}

// cleanupAfterSaleByOrder 按订单号清场（子表先删；after_sale_order 优先）。
func cleanupAfterSaleByOrder(ctx context.Context, orderNo string) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM after_sale_order WHERE order_no=?", orderNo)
	_, _ = g.DB().Exec(ctx, "DELETE i FROM trade_order_item i JOIN trade_order o ON i.order_id=o.id WHERE o.order_no=?", orderNo)
	_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order WHERE order_no=?", orderNo)
}

func cleanupAfterSaleFixture(ctx context.Context, t *gtest.T, f *afterSaleFixture) {
	cleanupAfterSaleByOrder(ctx, f.OrderNo)
	_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id=?", f.SkuId)
	cleanupPointUser(ctx, t, f.phone)
}

// seedAfterSaleRow 直接造一条售后单（用于绕开尚未实现的方法去验证"累计额度释放"等口径）。
func seedAfterSaleRow(
	ctx context.Context, t *gtest.T, f *afterSaleFixture, no string, status, qty int, refundYuan string,
) {
	_, err := g.DB().Exec(ctx,
		"INSERT INTO after_sale_order(after_sale_no,order_id,order_no,order_item_id,user_id,type,status,currency,"+
			"quantity,reason,description,refund_amount) VALUES(?,?,?,?,?,1,?,'CNY',?,'测试','',?)",
		no, f.OrderId, f.OrderNo, f.ItemId, f.UserId, status, qty, refundYuan)
	t.AssertNil(err)
}

// afterSaleRecord 取该订单下唯一的售后单（未找到返回 nil）。
func afterSaleRecord(ctx context.Context, orderNo string) gdb.Record {
	rec, err := g.DB().GetOne(ctx,
		"SELECT after_sale_no, status, type, quantity, refund_amount, return_logistics_no, refund_no, "+
			"reject_reason, operator_id, fail_reason, audit_time, refund_time "+
			"FROM after_sale_order WHERE order_no=?", orderNo)
	if err != nil || rec.IsEmpty() {
		return nil
	}
	return rec
}

// ---------- US1: 会员申请售后 ----------

// TestAfterSaleApply 申请校验矩阵（FR-001~004）。
func TestAfterSaleApply(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-APPLY-1", "T-AS-A1", 980000001, 2, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)

		// ① 部分申请（行实付 20.00, 数量 2 → 退 1 件 = 10.00）
		no, err := NewAfterSaleLogic().Apply(ctx, f.UserId, model.AfterSaleApplyInput{
			OrderItemId: f.ItemId, Type: afterSaleTypeRefundOnly, Quantity: 1, Reason: "不想要了",
		})
		t.AssertNil(err)
		t.Assert(no != "", true)
		row := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(row["after_sale_no"].String(), no)
		t.Assert(row["status"].Int(), afterSalePendingAudit)
		t.Assert(row["quantity"].Int(), 1)
		t.Assert(row["refund_amount"].String(), "10.00")

		// ② 剩余可退 1: 再申请 2 → 累计超额 40008
		_, err = NewAfterSaleLogic().Apply(ctx, f.UserId, model.AfterSaleApplyInput{
			OrderItemId: f.ItemId, Type: afterSaleTypeRefundOnly, Quantity: 2, Reason: "不想要了",
		})
		t.Assert(errCode(err), errcode.CodeAfterSaleDenied)

		// ③ 剩余 1 可以申请（累计口径正确）→ 走退货退款, 落在 10 待审核
		no2, err := NewAfterSaleLogic().Apply(ctx, f.UserId, model.AfterSaleApplyInput{
			OrderItemId: f.ItemId, Type: afterSaleTypeReturnGoods, Quantity: 1, Reason: "质量问题",
		})
		t.AssertNil(err)
		t.Assert(no2 != no, true)

		// ④ 数量超过行数量 → 40008
		_, err = NewAfterSaleLogic().Apply(ctx, f.UserId, model.AfterSaleApplyInput{
			OrderItemId: f.ItemId, Type: afterSaleTypeRefundOnly, Quantity: 3, Reason: "不想要了",
		})
		t.Assert(errCode(err), errcode.CodeAfterSaleDenied)

		// ⑤ 他人订单项 → 40009（按不存在处理, 不泄露他人资源）
		_, err = NewAfterSaleLogic().Apply(ctx, f.UserId+999, model.AfterSaleApplyInput{
			OrderItemId: f.ItemId, Type: afterSaleTypeRefundOnly, Quantity: 1, Reason: "不想要了",
		})
		t.Assert(errCode(err), errcode.CodeAfterSaleNotFound)
	})
}

// TestAfterSaleApplyReleasedQuota 已拒绝/已撤销**不占用**可退额度（research D5）。
func TestAfterSaleApplyReleasedQuota(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-APPLY-2", "T-AS-A2", 980000002, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)

		// 先造一条"已拒绝"的满额申请 → 额度应已释放
		seedAfterSaleRow(ctx, t, f, "AS-REJ-1", afterSaleRejected, 1, "10.00")
		no, err := NewAfterSaleLogic().Apply(ctx, f.UserId, model.AfterSaleApplyInput{
			OrderItemId: f.ItemId, Type: afterSaleTypeRefundOnly, Quantity: 1, Reason: "重试申请",
		})
		t.AssertNil(err)
		t.Assert(no != "", true)

		// 换成"已撤销"同样释放
		_, _ = g.DB().Exec(ctx, "DELETE FROM after_sale_order WHERE order_no=?", f.OrderNo)
		seedAfterSaleRow(ctx, t, f, "AS-CAN-1", afterSaleCanceled, 1, "10.00")
		_, err = NewAfterSaleLogic().Apply(ctx, f.UserId, model.AfterSaleApplyInput{
			OrderItemId: f.ItemId, Type: afterSaleTypeRefundOnly, Quantity: 1, Reason: "重试申请",
		})
		t.AssertNil(err)

		// 对照: 处于"待审核"的申请**占用**额度 → 再申请 40008
		_, err = NewAfterSaleLogic().Apply(ctx, f.UserId, model.AfterSaleApplyInput{
			OrderItemId: f.ItemId, Type: afterSaleTypeRefundOnly, Quantity: 1, Reason: "再来一单",
		})
		t.Assert(errCode(err), errcode.CodeAfterSaleDenied)
	})
}

// TestAfterSaleApplyFullRefundAmount 全额退的金额口径: 退完剩余数量时取剩余全额（消尾差）。
func TestAfterSaleApplyFullRefundAmount(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		// 行实付 10.00 / 数量 3 → 单件 3.33（分, 向下取整）; 前两笔 3.33, 末笔取剩余 3.34
		f := seedAfterSaleFixture(ctx, t, "AS-APPLY-3", "T-AS-A3", 980000003, 3, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		// 把行实付改成 10.00 而数量仍为 3（构造"行实付不被数量整除"的尾差场景; 单价×数量本身必须整除）
		_, _ = g.DB().Exec(ctx, "UPDATE trade_order_item SET pay_amount='10.00' WHERE id=?", f.ItemId)
		_, _ = g.DB().Exec(ctx, "UPDATE trade_order SET pay_amount='10.00', total_amount='10.00' WHERE id=?", f.OrderId)
		f.PayYuan = "10.00"

		for i := 0; i < 3; i++ {
			_, err := NewAfterSaleLogic().Apply(ctx, f.UserId, model.AfterSaleApplyInput{
				OrderItemId: f.ItemId, Type: afterSaleTypeRefundOnly, Quantity: 1, Reason: "分次退",
			})
			t.AssertNil(err)
		}
		rows, err := g.DB().Model("after_sale_order").Ctx(ctx).
			Where("order_no", f.OrderNo).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(rows), 3)
		// 前两笔按比例取整 3.33, 末笔取剩余 10.00-3.33-3.33 = 3.34 → 合计恒等于行实付
		t.Assert(rows[0]["refund_amount"].String(), "3.33")
		t.Assert(rows[1]["refund_amount"].String(), "3.33")
		t.Assert(rows[2]["refund_amount"].String(), "3.34")

		// 全部可退数量已用完 → 再申请 40008
		_, err = NewAfterSaleLogic().Apply(ctx, f.UserId, model.AfterSaleApplyInput{
			OrderItemId: f.ItemId, Type: afterSaleTypeRefundOnly, Quantity: 1, Reason: "第四笔",
		})
		t.Assert(errCode(err), errcode.CodeAfterSaleDenied)
	})
}

// TestAfterSaleApplyOrderState 订单状态守卫: 未完成（待发货 20）的订单项不可申请。
func TestAfterSaleApplyOrderState(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-APPLY-4", "T-AS-A4", 980000004, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)

		_, _ = g.DB().Exec(ctx, "UPDATE trade_order SET status=20 WHERE id=?", f.OrderId)
		_, err := NewAfterSaleLogic().Apply(ctx, f.UserId, model.AfterSaleApplyInput{
			OrderItemId: f.ItemId, Type: afterSaleTypeRefundOnly, Quantity: 1, Reason: "不想要了",
		})
		t.Assert(errCode(err), errcode.CodeAfterSaleDenied)
	})
}

// ---------- US2: 后台审核 ----------

// refundSpy 退款渠道桩（记录调用; 可切换为失败）。默认渠道为 mock, 测试内替换并还原。
type refundSpy struct {
	calls    int
	lastNo   string
	lastFen  int64
	failWith error
}

func (r *refundSpy) Name() string { return "spy" }
func (r *refundSpy) CreateParams(string, int64, string) (map[string]any, error) {
	return nil, nil
}
func (r *refundSpy) ParseNotify([]byte) (*paychannel.NotifyPayload, error) { return nil, nil }
func (r *refundSpy) Refund(outRefundNo string, amountFen int64) error {
	r.calls++
	r.lastNo, r.lastFen = outRefundNo, amountFen
	return r.failWith
}

// useRefundSpy 替换渠道并返回还原函数。
func useRefundSpy(spy *refundSpy) func() {
	orig := afterSaleRefundChannel
	afterSaleRefundChannel = spy
	return func() { afterSaleRefundChannel = orig }
}

// TestAfterSaleApproveRefundOnly 仅退款: 同意 → 30→发起退款→40; 操作人/审核时间落库; 渠道以单号为幂等键。
func TestAfterSaleApproveRefundOnly(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-APR-1", "T-AS-APR1", 980000011, 2, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow(ctx, t, f, "AS-APR-1", afterSalePendingAudit, 1, "10.00")
		seedAfterSaleRow2(ctx, t, f, "AS-APR-2", afterSalePendingAudit, afterSaleTypeReturnGoods, 1, "10.00")
		_, _ = g.DB().Exec(ctx, "DELETE FROM after_sale_order WHERE after_sale_no='AS-APR-2'")

		spy := &refundSpy{}
		defer useRefundSpy(spy)()

		err := NewAfterSaleLogic().Approve(ctx, "AS-APR-1", "admin:7")
		t.AssertNil(err)
		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSaleRefunding) // 40 退款中
		t.Assert(rec["operator_id"].String(), "admin:7")
		t.Assert(rec["audit_time"].String() != "", true)
		t.Assert(spy.calls, 1)
		t.Assert(spy.lastNo, "AS-APR-1") // 幂等键 = 售后单号
		t.Assert(spy.lastFen, int64(1000))
	})
}

// TestAfterSaleApproveReturnGoods 退货退款: 同意 → 20 待买家寄回, 且**不发起**渠道退款。
func TestAfterSaleApproveReturnGoods(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-APR-2", "T-AS-APR2", 980000012, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow2(ctx, t, f, "AS-APR-3", afterSalePendingAudit, afterSaleTypeReturnGoods, 1, "10.00")

		spy := &refundSpy{}
		defer useRefundSpy(spy)()

		t.AssertNil(NewAfterSaleLogic().Approve(ctx, "AS-APR-3", "admin:8"))
		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSalePendingReturn) // 20
		t.Assert(rec["operator_id"].String(), "admin:8")
		t.Assert(spy.calls, 0)
	})
}

// TestAfterSaleApproveStateGuard 状态机守卫: 非 10 → 40006; 不存在的单 → 40009。
func TestAfterSaleApproveStateGuard(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-APR-3", "T-AS-APR3", 980000013, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow2(ctx, t, f, "AS-APR-4", afterSalePendingReturn, afterSaleTypeReturnGoods, 1, "10.00")

		t.Assert(errCode(NewAfterSaleLogic().Approve(ctx, "AS-APR-4", "admin:9")), errcode.CodeStatusNotAllowed)
		t.Assert(errCode(NewAfterSaleLogic().Approve(ctx, "AS-NOTHING", "admin:9")), errcode.CodeAfterSaleNotFound)
	})
}

// TestAfterSaleReject 拒绝: 10→90, 原因与操作人落库; 原因必填（服务侧亦兜底）; 非 10 → 40006。
func TestAfterSaleReject(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-REJ-9", "T-AS-REJ9", 980000014, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow2(ctx, t, f, "AS-REJ-9", afterSalePendingAudit, afterSaleTypeRefundOnly, 1, "10.00")

		// 原因为空 → 参数错误
		t.Assert(errCode(NewAfterSaleLogic().Reject(ctx, "AS-REJ-9", "  ", "admin:3")), errcode.CodeInvalidParam)

		t.AssertNil(NewAfterSaleLogic().Reject(ctx, "AS-REJ-9", "凭证不足", "admin:3"))
		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSaleRejected)
		t.Assert(rec["reject_reason"].String(), "凭证不足")
		t.Assert(rec["operator_id"].String(), "admin:3")

		// 已拒绝再拒 → 40006
		t.Assert(errCode(NewAfterSaleLogic().Reject(ctx, "AS-REJ-9", "再拒一次", "admin:3")), errcode.CodeStatusNotAllowed)
	})
}

// TestAfterSaleAdminListDetail 后台列表（状态/类型筛选）与详情。
func TestAfterSaleAdminListDetail(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-ADM-1", "T-AS-ADM1", 980000015, 2, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow2(ctx, t, f, "AS-ADM-L1", afterSalePendingAudit, afterSaleTypeRefundOnly, 1, "5.00")
		seedAfterSaleRow2(ctx, t, f, "AS-ADM-L2", afterSaleFinished, afterSaleTypeReturnGoods, 1, "5.00")

		res, err := NewAfterSaleLogic().AdminList(ctx, afterSalePendingAudit, 0, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		// 精确断言: 列表内必含 L1, 且不含 L2（状态筛选生效）
		hit := map[string]bool{}
		for _, it := range res.List {
			hit[it.AfterSaleNo] = true
		}
		t.Assert(hit["AS-ADM-L1"], true)
		t.Assert(hit["AS-ADM-L2"], false)

		// 类型筛选
		res2, err := NewAfterSaleLogic().AdminList(ctx, 0, afterSaleTypeReturnGoods, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		hit2 := map[string]bool{}
		for _, it := range res2.List {
			hit2[it.AfterSaleNo] = true
		}
		t.Assert(hit2["AS-ADM-L2"], true)
		t.Assert(hit2["AS-ADM-L1"], false)

		d, err := NewAfterSaleLogic().AdminDetail(ctx, "AS-ADM-L1")
		t.AssertNil(err)
		t.Assert(d.AfterSaleNo, "AS-ADM-L1")
		t.Assert(d.Quantity, 1)
		t.Assert(d.RefundAmount, "5.00")
	})
}

// seedAfterSaleRow2 造售后单（可指定类型; seedAfterSaleRow 固定 type=1）。
func seedAfterSaleRow2(
	ctx context.Context, t *gtest.T, f *afterSaleFixture, no string, status, saleType, qty int, refundYuan string,
) {
	_, err := g.DB().Exec(ctx,
		"INSERT INTO after_sale_order(after_sale_no,order_id,order_no,order_item_id,user_id,type,status,currency,"+
			"quantity,reason,description,refund_amount) VALUES(?,?,?,?,?,?,?,'CNY',?,'测试','',?)",
		no, f.OrderId, f.OrderNo, f.ItemId, f.UserId, saleType, status, qty, refundYuan)
	t.AssertNil(err)
}

// ---------- US3: 退货寄回与确认收货 ----------

// TestAfterSaleSubmitReturn 填写寄回单号: 仅"退货退款 + 待寄回"可填; 幂等覆盖最后一次。
func TestAfterSaleSubmitReturn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-RET-1", "T-AS-RET1", 980000021, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow2(ctx, t, f, "AS-RET-1", afterSalePendingReturn, afterSaleTypeReturnGoods, 1, "10.00")

		t.AssertNil(NewAfterSaleLogic().SubmitReturn(ctx, f.UserId, "AS-RET-1", "SF123456"))
		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["return_logistics_no"].String(), "SF123456")
		t.Assert(rec["status"].Int(), afterSalePendingReturn)

		// 再填一次 → 以最后一次为准
		t.AssertNil(NewAfterSaleLogic().SubmitReturn(ctx, f.UserId, "AS-RET-1", "SF999999"))
		rec = afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["return_logistics_no"].String(), "SF999999")

		// 他人 → 40009
		t.Assert(errCode(NewAfterSaleLogic().SubmitReturn(ctx, f.UserId+999, "AS-RET-1", "SF1")), errcode.CodeAfterSaleNotFound)

		// 状态不符（改成待审核）→ 40006
		_, _ = g.DB().Exec(ctx, "UPDATE after_sale_order SET status=? WHERE after_sale_no='AS-RET-1'", afterSalePendingAudit)
		t.Assert(errCode(NewAfterSaleLogic().SubmitReturn(ctx, f.UserId, "AS-RET-1", "SF2")), errcode.CodeStatusNotAllowed)

		// 仅退款类型不可填寄回单号 → 40006
		_, _ = g.DB().Exec(ctx, "UPDATE after_sale_order SET status=?, type=? WHERE after_sale_no='AS-RET-1'",
			afterSalePendingReturn, afterSaleTypeRefundOnly)
		t.Assert(errCode(NewAfterSaleLogic().SubmitReturn(ctx, f.UserId, "AS-RET-1", "SF3")), errcode.CodeStatusNotAllowed)
	})
}

// TestAfterSaleConfirmReceipt 确认收货: 须已填寄回单号 → 30→发起退款→40; 缺单号 → 40008。
func TestAfterSaleConfirmReceipt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-CFR-1", "T-AS-CFR1", 980000022, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow2(ctx, t, f, "AS-CFR-1", afterSalePendingReturn, afterSaleTypeReturnGoods, 1, "10.00")

		spy := &refundSpy{}
		defer useRefundSpy(spy)()

		// 未填寄回单号 → 拒绝
		t.Assert(errCode(NewAfterSaleLogic().ConfirmReceipt(ctx, "AS-CFR-1", "admin:5")), errcode.CodeAfterSaleDenied)

		// 填单号后确认收货 → 40 退款中 + 渠道已发起
		_, _ = g.DB().Exec(ctx, "UPDATE after_sale_order SET return_logistics_no='SF777' WHERE after_sale_no='AS-CFR-1'")
		t.AssertNil(NewAfterSaleLogic().ConfirmReceipt(ctx, "AS-CFR-1", "admin:5"))
		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSaleRefunding)
		t.Assert(rec["operator_id"].String(), "admin:5")
		t.Assert(spy.calls, 1)
	})
}

// ---------- US4: 退款闭环与重试 ----------

// TestAfterSaleRetryRefund 重试退款: 30→40（清 fail_reason）; 40/50 不可重试。
func TestAfterSaleRetryRefund(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-RTY-1", "T-AS-RTY1", 980000031, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow2(ctx, t, f, "AS-RTY-1", afterSalePendingRefund, afterSaleTypeRefundOnly, 1, "10.00")
		_, _ = g.DB().Exec(ctx, "UPDATE after_sale_order SET fail_reason='渠道超时' WHERE after_sale_no='AS-RTY-1'")

		spy := &refundSpy{}
		defer useRefundSpy(spy)()

		t.AssertNil(NewAfterSaleLogic().RetryRefund(ctx, "AS-RTY-1", "admin:11"))
		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSaleRefunding)
		t.Assert(rec["fail_reason"].String(), "")
		t.Assert(rec["operator_id"].String(), "admin:11")
		t.Assert(spy.calls, 1)

		// 退款中（40）不可重试（避免重复出款）
		t.Assert(errCode(NewAfterSaleLogic().RetryRefund(ctx, "AS-RTY-1", "admin:11")), errcode.CodeStatusNotAllowed)
		t.Assert(spy.calls, 1)

		// 已完成后（50）不可重试
		_, _ = g.DB().Exec(ctx, "UPDATE after_sale_order SET status=? WHERE after_sale_no='AS-RTY-1'", afterSaleFinished)
		t.Assert(errCode(NewAfterSaleLogic().RetryRefund(ctx, "AS-RTY-1", "admin:11")), errcode.CodeStatusNotAllowed)
		t.Assert(spy.calls, 1)
	})
}

// TestAfterSaleRefundChannelFail 渠道失败: 停在 30 且留 fail_reason（不静默、不回滚状态机）; 可重试成功。
func TestAfterSaleRefundChannelFail(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-FAIL-1", "T-AS-FAIL1", 980000032, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow2(ctx, t, f, "AS-FAIL-1", afterSalePendingAudit, afterSaleTypeRefundOnly, 1, "10.00")

		spy := &refundSpy{failWith: errors.New("渠道返回: 余额不足")}
		defer useRefundSpy(spy)()

		t.AssertNil(NewAfterSaleLogic().Approve(ctx, "AS-FAIL-1", "admin:12"))
		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSalePendingRefund)
		t.Assert(rec["fail_reason"].String(), "渠道返回: 余额不足")

		// 渠道恢复 → 重试成功 → 40 且原因清空
		spy.failWith = nil
		t.AssertNil(NewAfterSaleLogic().RetryRefund(ctx, "AS-FAIL-1", "admin:12"))
		rec = afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSaleRefunding)
		t.Assert(rec["fail_reason"].String(), "")
		t.Assert(spy.calls, 2)
	})
}

// TestAfterSaleRefundCallback 渠道退款回调（复用既有 HandleRefundNotify）: 40→50 + refund_time; 重复回调幂等。
func TestAfterSaleRefundCallback(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-CB-1", "T-AS-CB1", 980000033, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow2(ctx, t, f, "AS-CB-1", afterSaleRefunding, afterSaleTypeRefundOnly, 1, "10.00")
		defer cleanupCallbackLogsAfter(ctx, maxCallbackLogId(ctx))

		body, _ := json.Marshal(map[string]any{
			"payNo": "PAY-AS-1", "outRefundNo": "AS-CB-1", "success": true,
		})
		t.AssertNil(NewPayLogic().HandleRefundNotify(ctx, "mock", body))
		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSaleFinished)
		t.Assert(rec["refund_time"].String() != "", true)

		// 重复回调: 状态不为 40 → 不匹配 → 报错且状态不变（幂等语义: 不二次推进）
		t.AssertNE(NewPayLogic().HandleRefundNotify(ctx, "mock", body), nil)
		rec = afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSaleFinished)

		// 未匹配单号 → 报错
		bad, _ := json.Marshal(map[string]any{"payNo": "PAY-AS-X", "outRefundNo": "AS-NONE", "success": true})
		t.Assert(errCode(NewPayLogic().HandleRefundNotify(ctx, "mock", bad)), errcode.CodeNotFound)
	})
}

// TestAfterSaleZeroRefund 0 元边界: 不调用渠道, 直接已完成且 refund_no 为空。
func TestAfterSaleZeroRefund(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-ZERO-1", "T-AS-ZERO1", 980000034, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow2(ctx, t, f, "AS-ZERO-1", afterSalePendingRefund, afterSaleTypeRefundOnly, 1, "0.00")

		spy := &refundSpy{}
		defer useRefundSpy(spy)()

		t.AssertNil(NewAfterSaleLogic().RetryRefund(ctx, "AS-ZERO-1", "admin:13"))
		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSaleFinished)
		t.Assert(rec["refund_no"].String(), "")
		t.Assert(spy.calls, 0)
	})
}

// ---------- US5: 会员查询与撤销 ----------

// TestAfterSaleListDetailOwn 列表（状态筛选+仅本人）与详情（他人 → 40009）。
func TestAfterSaleListDetailOwn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-OWN-1", "T-AS-OWN1", 980000041, 2, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		g := seedAfterSaleFixture(ctx, t, "AS-OWN-2", "T-AS-OWN2", 980000042, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, g)

		seedAfterSaleRow2(ctx, t, f, "AS-OWN-L1", afterSalePendingAudit, afterSaleTypeRefundOnly, 1, "5.00")
		seedAfterSaleRow2(ctx, t, f, "AS-OWN-L2", afterSaleFinished, afterSaleTypeRefundOnly, 1, "5.00")
		seedAfterSaleRow2(ctx, t, g, "AS-OWN-L3", afterSalePendingAudit, afterSaleTypeRefundOnly, 1, "5.00")

		res, err := NewAfterSaleLogic().List(ctx, f.UserId, afterSalePendingAudit, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		hit := map[string]bool{}
		for _, it := range res.List {
			hit[it.AfterSaleNo] = true
		}
		t.Assert(hit["AS-OWN-L1"], true)
		t.Assert(hit["AS-OWN-L2"], false)
		t.Assert(hit["AS-OWN-L3"], false) // 他人的不可见

		// 详情: 本人可见
		d, err := NewAfterSaleLogic().Detail(ctx, f.UserId, "AS-OWN-L1")
		t.AssertNil(err)
		t.Assert(d.AfterSaleNo, "AS-OWN-L1")
		t.Assert(d.RefundAmount, "5.00")

		// 他人的 → 40009
		_, err = NewAfterSaleLogic().Detail(ctx, f.UserId, "AS-OWN-L3")
		t.Assert(errCode(err), errcode.CodeAfterSaleNotFound)
	})
}

// TestAfterSaleCancel 撤销边界（D4）: {10,20,30} 可撤→91 且释放额度; 40/50 → 40006; 他人 → 40009。
func TestAfterSaleCancel(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-CAN-9", "T-AS-CAN9", 980000043, 2, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)

		// 10 可撤
		seedAfterSaleRow2(ctx, t, f, "AS-CAN-A", afterSalePendingAudit, afterSaleTypeRefundOnly, 1, "5.00")
		t.AssertNil(NewAfterSaleLogic().Cancel(ctx, f.UserId, "AS-CAN-A"))
		rec, _ := g.DB().GetOne(ctx, "SELECT status FROM after_sale_order WHERE after_sale_no='AS-CAN-A'")
		t.Assert(rec["status"].Int(), afterSaleCanceled)

		// 撤销后额度释放: 行数量 2, 已撤销 1 → 还能申请 2
		_, err := NewAfterSaleLogic().Apply(ctx, f.UserId, model.AfterSaleApplyInput{
			OrderItemId: f.ItemId, Type: afterSaleTypeRefundOnly, Quantity: 2, Reason: "重新申请",
		})
		t.AssertNil(err)

		// 20 / 30 可撤
		for _, tc := range []struct {
			no string
			st int
		}{{"AS-CAN-B", afterSalePendingReturn}, {"AS-CAN-C", afterSalePendingRefund}} {
			_, _ = g.DB().Exec(ctx, "DELETE FROM after_sale_order WHERE after_sale_no=?", tc.no)
			seedAfterSaleRow2(ctx, t, f, tc.no, tc.st, afterSaleTypeReturnGoods, 1, "5.00")
			t.AssertNil(NewAfterSaleLogic().Cancel(ctx, f.UserId, tc.no))
		}

		// 40 不可撤（在途资金）→ 40006
		seedAfterSaleRow2(ctx, t, f, "AS-CAN-D", afterSaleRefunding, afterSaleTypeRefundOnly, 1, "5.00")
		t.Assert(errCode(NewAfterSaleLogic().Cancel(ctx, f.UserId, "AS-CAN-D")), errcode.CodeStatusNotAllowed)
		// 50 不可撤
		seedAfterSaleRow2(ctx, t, f, "AS-CAN-E", afterSaleFinished, afterSaleTypeRefundOnly, 1, "5.00")
		t.Assert(errCode(NewAfterSaleLogic().Cancel(ctx, f.UserId, "AS-CAN-E")), errcode.CodeStatusNotAllowed)

		// 他人 → 40009
		t.Assert(errCode(NewAfterSaleLogic().Cancel(ctx, f.UserId+999, "AS-CAN-D")), errcode.CodeAfterSaleNotFound)
	})
}

// ---------- US6: 完成后的账实一致 ----------

// commissionSpy 佣金冲销事件出口桩。
type commissionSpy struct {
	calls      int
	lastOrder  string
	lastSaleNo string
	lastRefund int64
}

func (c *commissionSpy) ReverseForAfterSale(_ context.Context, orderNo, afterSaleNo string, refundFen int64) {
	c.calls++
	c.lastOrder, c.lastSaleNo, c.lastRefund = orderNo, afterSaleNo, refundFen
}

// refundDone 走完整链路到"已完成": 30→(发起退款)→40→(渠道回调)→50。
func refundDone(ctx context.Context, t *gtest.T, afterSaleNo string) {
	t.AssertNil(NewAfterSaleLogic().RetryRefund(ctx, afterSaleNo, "admin:21"))
	body, _ := json.Marshal(map[string]any{
		"payNo": "PAY-" + afterSaleNo, "outRefundNo": afterSaleNo, "success": true,
	})
	t.AssertNil(NewPayLogic().HandleRefundNotify(ctx, "mock", body))
}

// TestAfterSaleFinishedSideEffects 退货退款完成: 回补库存 + 全额退款状态 + 佣金冲销事件。
func TestAfterSaleFinishedSideEffects(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-SFX-1", "T-AS-SFX1", 980000051, 2, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		// 单行订单 + 退货退款 2 件（整行） → 完成后应"全额退款"且库存 +2
		seedAfterSaleRow2(ctx, t, f, "AS-SFX-1", afterSalePendingRefund, afterSaleTypeReturnGoods, 2, "20.00")
		defer cleanupCallbackLogsAfter(ctx, maxCallbackLogId(ctx))

		spy := &commissionSpy{}
		old := CommissionReverse
		CommissionReverse = spy
		defer func() { CommissionReverse = old }()

		refundDone(ctx, t, "AS-SFX-1")

		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSaleFinished)
		inv, err := g.DB().GetOne(ctx, "SELECT total, locked FROM inventory WHERE sku_id=?", f.SkuId)
		t.AssertNil(err)
		t.Assert(inv["total"].Int(), 102) // 回补 2
		t.Assert(inv["locked"].Int(), 0)  // 锁定不动
		st, err := g.DB().GetValue(ctx, "SELECT refund_status FROM trade_order WHERE id=?", f.OrderId)
		t.AssertNil(err)
		t.Assert(st.Int(), 2) // 全额退款
		t.Assert(spy.calls, 1)
		t.Assert(spy.lastOrder, f.OrderNo)
		t.Assert(spy.lastSaleNo, "AS-SFX-1")
		t.Assert(spy.lastRefund, int64(2000))
	})
}

// TestAfterSaleFinishedSideEffectsRefundOnly 仅退款完成: **不回补**库存, 但订单退款状态仍更新为部分退款。
func TestAfterSaleFinishedSideEffectsRefundOnly(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-SFX-2", "T-AS-SFX2", 980000052, 2, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow2(ctx, t, f, "AS-SFX-2", afterSalePendingRefund, afterSaleTypeRefundOnly, 1, "10.00")
		defer cleanupCallbackLogsAfter(ctx, maxCallbackLogId(ctx))

		spy := &commissionSpy{}
		old := CommissionReverse
		CommissionReverse = spy
		defer func() { CommissionReverse = old }()

		refundDone(ctx, t, "AS-SFX-2")
		inv, err := g.DB().GetOne(ctx, "SELECT total FROM inventory WHERE sku_id=?", f.SkuId)
		t.AssertNil(err)
		t.Assert(inv["total"].Int(), 100) // 仅退款不回补
		st, err := g.DB().GetValue(ctx, "SELECT refund_status FROM trade_order WHERE id=?", f.OrderId)
		t.AssertNil(err)
		t.Assert(st.Int(), 1) // 2 件只退 1 件 → 部分退款
		t.Assert(spy.calls, 1)
	})
}

// TestAfterSaleCallbackIdempotentSideEffect 重复回调不产生二次副作用（回补/事件各一次）。
func TestAfterSaleCallbackIdempotentSideEffect(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-SFX-3", "T-AS-SFX3", 980000053, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		seedAfterSaleRow2(ctx, t, f, "AS-SFX-3", afterSaleRefunding, afterSaleTypeReturnGoods, 1, "10.00")
		defer cleanupCallbackLogsAfter(ctx, maxCallbackLogId(ctx))

		spy := &commissionSpy{}
		old := CommissionReverse
		CommissionReverse = spy
		defer func() { CommissionReverse = old }()

		body, _ := json.Marshal(map[string]any{
			"payNo": "PAY-AS-3", "outRefundNo": "AS-SFX-3", "success": true,
		})
		t.AssertNil(NewPayLogic().HandleRefundNotify(ctx, "mock", body))
		inv, _ := g.DB().GetOne(ctx, "SELECT total FROM inventory WHERE sku_id=?", f.SkuId)
		t.Assert(inv["total"].Int(), 101)

		// 重复回调: 状态已 50 → 不匹配 → 报错; 副作用不重复
		t.AssertNE(NewPayLogic().HandleRefundNotify(ctx, "mock", body), nil)
		inv, _ = g.DB().GetOne(ctx, "SELECT total FROM inventory WHERE sku_id=?", f.SkuId)
		t.Assert(inv["total"].Int(), 101)
		t.Assert(spy.calls, 1)
	})
}

// TestAfterSaleZeroRefundSideEffects 0 元边界: 不调渠道、直接完成, 且**仍执行**完成副作用。
func TestAfterSaleZeroRefundSideEffects(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedAfterSaleFixture(ctx, t, "AS-SFX-4", "T-AS-SFX4", 980000054, 1, "10.00")
		defer cleanupAfterSaleFixture(ctx, t, f)
		// 全优惠行: 行实付 0.00 → 退款额 0
		seedAfterSaleRow2(ctx, t, f, "AS-SFX-4", afterSalePendingRefund, afterSaleTypeReturnGoods, 1, "0.00")

		spy := &refundSpy{}
		defer useRefundSpy(spy)()
		comm := &commissionSpy{}
		old := CommissionReverse
		CommissionReverse = comm
		defer func() { CommissionReverse = old }()

		t.AssertNil(NewAfterSaleLogic().RetryRefund(ctx, "AS-SFX-4", "admin:23"))
		rec := afterSaleRecord(ctx, f.OrderNo)
		t.Assert(rec["status"].Int(), afterSaleFinished)
		t.Assert(spy.calls, 0) // 0 元不调渠道
		inv, _ := g.DB().GetOne(ctx, "SELECT total FROM inventory WHERE sku_id=?", f.SkuId)
		t.Assert(inv["total"].Int(), 101) // 退货退款仍回补
		t.Assert(comm.calls, 1)
	})
}
