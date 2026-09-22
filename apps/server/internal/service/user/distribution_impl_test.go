// distribution_impl_test.go 分销与资金（017 批次 11）——红线场景 SC-3 全测试化。
// 重点: 两级封顶/自邀拒绝/计提幂等/提现状态机不可逆/打款渠道单号幂等/并发提现余额条件更新。
package user

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// ---- fixture（精确名自清; 子表先删） ----

func distCleanupUser(ctx context.Context, t *gtest.T, handle string) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE nickname LIKE ?", handle+"%")
	res, err := g.DB().Exec(ctx,
		"INSERT INTO `user`(nickname,phone,phone_hash,growth_value,status) VALUES(?,?,?,0,1)",
		handle+"U", handle+"-p", handle+"-ph")
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

func distCleanupAll(ctx context.Context, ids ...int64) {
	for _, id := range ids {
		_, _ = g.DB().Exec(ctx, "DELETE FROM account_log WHERE user_id=?", id)
		_, _ = g.DB().Exec(ctx, "DELETE FROM user_account WHERE user_id=?", id)
		_, _ = g.DB().Exec(ctx, "DELETE FROM withdraw_order WHERE user_id=?", id)
		_, _ = g.DB().Exec(ctx, "DELETE FROM commission_record WHERE beneficiary_user_id=? OR order_no LIKE 'TF-DIST%'", id)
		_, _ = g.DB().Exec(ctx, "DELETE FROM invite_record WHERE inviter_id=? OR new_user_id=?", id, id)
		_, _ = g.DB().Exec(ctx, "DELETE FROM user_relation WHERE user_id=? OR inviter_id=?", id, id)
		_, _ = g.DB().Exec(ctx, "DELETE FROM distribution_user WHERE user_id=?", id)
		_, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE id=?", id)
	}
}

// seedRelation 直连 A→B（A 为 B 的直接上级）。
func seedRelation(ctx context.Context, t *gtest.T, inviterId, userId int64) {
	_, err := g.DB().Exec(ctx,
		"INSERT INTO user_relation(user_id,inviter_id,bind_channel) VALUES(?,?,1)", userId, inviterId)
	t.AssertNil(err)
}

// seedDistributor 直插通过态推广员。
func seedDistributor(ctx context.Context, t *gtest.T, userId int64, status int) int64 {
	res, err := g.DB().Exec(ctx,
		"INSERT INTO distribution_user(user_id,status,audit_time) VALUES(?,?,NOW())", userId, status)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

// seedAccount 直插账户（余额/冻结可控）。
func seedAccount(ctx context.Context, t *gtest.T, userId int64, balance, frozen string) {
	_, err := g.DB().Exec(ctx,
		"INSERT INTO user_account(user_id,balance,frozen) VALUES(?,?,?) ON DUPLICATE KEY UPDATE balance=VALUES(balance), frozen=VALUES(frozen)",
		userId, balance, frozen)
	t.AssertNil(err)
}

// ---- US1 关系链与资质 ----

// TestDistApplyAndAudit 申请防重 60001 / 审核状态机（重放拒绝）/ 冻结解冻（FR-1）。
func TestDistApplyAndAudit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		a := distCleanupUser(ctx, t, "TF-DA-ADM")
		b := distCleanupUser(ctx, t, "TF-DA-B")
		defer distCleanupAll(ctx, a, b)
		logic := NewDistributionLogic()
		admin := NewDistributionAdminLogic()

		// 申请 → 待审核; 重复申请 → 60001
		t.AssertNil(logic.Apply(ctx, b))
		t.Assert(errCode(logic.Apply(ctx, b)), errcode.CodeDistAlreadyApplied)

		// 后台通过 → 状态 2; 终态重放审核 → 拒绝
		dist, err := g.DB().GetValue(ctx, "SELECT id FROM distribution_user WHERE user_id=?", b)
		t.AssertNil(err)
		t.AssertNil(admin.AdminDistributorAudit(ctx, dist.Int64(), true))
		st, err := logic.Status(ctx, b)
		t.AssertNil(err)
		t.Assert(st.Status, 2)
		t.Assert(errCode(admin.AdminDistributorAudit(ctx, dist.Int64(), false)), errcode.CodeStatusNotAllowed)

		// 冻结/解冻
		t.AssertNil(admin.AdminDistributorFreeze(ctx, dist.Int64(), true))
		st, _ = logic.Status(ctx, b)
		t.Assert(st.Status, 3)
		t.AssertNil(admin.AdminDistributorFreeze(ctx, dist.Int64(), false))
		st, _ = logic.Status(ctx, b)
		t.Assert(st.Status, 2)

		_ = a
	})
}

// TestDistBindRelation 绑定自邀拒绝 / 一人一链唯一键兜底 / 关系页脱敏 / 两级解析到顶（FR-2, 红线）。
func TestDistBindRelation(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		a := distCleanupUser(ctx, t, "TF-DR-A")
		b := distCleanupUser(ctx, t, "TF-DR-B")
		c := distCleanupUser(ctx, t, "TF-DR-C")
		defer distCleanupAll(ctx, a, b, c)
		logic := NewDistributionLogic()

		// 自邀 → 拒绝（红线）
		t.Assert(errCode(logic.BindRelation(ctx, b, b, 1)), errcode.CodeInvalidParam)
		// 正常绑定 B→A; 二次绑定他人 → 一人一链拒绝
		t.AssertNil(logic.BindRelation(ctx, b, a, 1))
		t.Assert(errCode(logic.BindRelation(ctx, b, c, 1)), errcode.CodeInvalidParam)
		// C→B 形成 a←b←c 两级链; c 无上级（到顶即止, 无祖父可查）
		t.AssertNil(logic.BindRelation(ctx, c, b, 1))

		// 关系页: B 见上级 A 与下级 C（脱敏——昵称不含完整 handle 逻辑由 mask 保证, 此处断言非空+分页正确）
		rel, err := logic.Relations(ctx, b, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(rel.Total, int64(1))
		t.Assert(rel.Inviter != nil && rel.Inviter["nickname"] != "", true)

		// 两级封顶: C 的链 = B(一级) → A(二级), 不存在第三级可解析
		row, err := g.DB().GetValue(ctx,
			"SELECT r2.inviter_id FROM user_relation r1 JOIN user_relation r2 ON r2.user_id=r1.inviter_id WHERE r1.user_id=?", c)
		t.AssertNil(err)
		t.Assert(row.Int64(), a) // C 的一级=B、二级=A; 表结构无第三级存取路径
	})
}

// ---- US3 账户与提现状态机 ----

// TestDistWithdrawStateMachine 提现全状态机（FR-5, 红线）:
// 余额校验/并发双申请/审核通过拒绝回退/打款成功幂等/失败回退/不可逆迁移。
func TestDistWithdrawStateMachine(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		a := distCleanupUser(ctx, t, "TF-DW-A")
		defer distCleanupAll(ctx, a)
		logic := NewDistributionLogic()
		admin := NewDistributionAdminLogic()
		seedAccount(ctx, t, a, "100.00", "0.00")

		// 金额校验: 超余额/非正/>2位小数 → 拒绝
		_, err := logic.WithdrawApply(ctx, a, "200.00")
		t.Assert(errCode(err), errcode.CodeInvalidParam)
		_, err = logic.WithdrawApply(ctx, a, "0.00")
		t.Assert(errCode(err), errcode.CodeInvalidParam)
		_, err = logic.WithdrawApply(ctx, a, "10.005")
		t.Assert(errCode(err), errcode.CodeInvalidParam)

		// 申请 60 → 余额 40 / 冻结 60 / 单据 10 / 流水 biz_type=2
		no, err := logic.WithdrawApply(ctx, a, "60.00")
		t.AssertNil(err)
		acc, _ := g.DB().GetOne(ctx, "SELECT balance, frozen FROM user_account WHERE user_id=?", a)
		t.Assert(acc["balance"].String(), "40.00")
		t.Assert(acc["frozen"].String(), "60.00")
		st, _ := g.DB().GetValue(ctx, "SELECT status FROM withdraw_order WHERE withdraw_no=?", no)
		t.Assert(st.Int(), 10)

		// 并发双申请（每笔 60, 余额只够一笔冻结）→ 恰一笔成功
		seedAccount(ctx, t, a, "100.00", "0.00")
		_, _ = g.DB().Exec(ctx, "DELETE FROM withdraw_order WHERE user_id=?", a)
		var wg sync.WaitGroup
		errs := make([]error, 2)
		nos := make([]string, 2)
		start := make(chan struct{})
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				<-start
				noI, e := logic.WithdrawApply(ctx, a, "60.00")
				nos[idx], errs[idx] = noI, e
			}(i)
		}
		close(start)
		wg.Wait()
		okCnt := 0
		win := ""
		for i, e := range errs {
			if e == nil {
				okCnt++
				win = nos[i]
			} else {
				t.Assert(errCode(e), errcode.CodeInvalidParam) // 余额不足业务码
			}
		}
		t.Assert(okCnt, 1)
		t.Assert(win != "", true)

		// 中段流程以并发赢家单为主线: 审核 20 → 打款成功 40 → 不可逆 → 渠道单号幂等
		t.AssertNil(admin.AdminWithdrawAudit(ctx, win, true, ""))
		st, _ = g.DB().GetValue(ctx, "SELECT status FROM withdraw_order WHERE withdraw_no=?", win)
		t.Assert(st.Int(), 20)

		// 拒绝回退: 另造一单（余额剩 40）→ 50 + 冻结回退 + 流水 biz_type=4
		no2, err := logic.WithdrawApply(ctx, a, "10.00")
		t.AssertNil(err)
		t.AssertNil(admin.AdminWithdrawAudit(ctx, no2, false, "资料不全"))
		st, _ = g.DB().GetValue(ctx, "SELECT status FROM withdraw_order WHERE withdraw_no=?", no2)
		t.Assert(st.Int(), 50)
		acc, _ = g.DB().GetOne(ctx, "SELECT balance, frozen FROM user_account WHERE user_id=?", a)
		t.Assert(acc["frozen"].String(), "60.00") // 只剩赢家单的冻结
		lg, _ := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM account_log WHERE user_id=? AND biz_type=4", a)
		t.Assert(lg.Int() >= 1, true)

		// 打款: 成功 → 40 + 渠道单号 + 冻结清零
		t.AssertNil(admin.AdminWithdrawPay(ctx, win, true, "CH-ORDER-X", ""))
		st, _ = g.DB().GetValue(ctx, "SELECT status FROM withdraw_order WHERE withdraw_no=?", win)
		t.Assert(st.Int(), 40)
		acc, _ = g.DB().GetOne(ctx, "SELECT balance, frozen FROM user_account WHERE user_id=?", a)
		t.Assert(acc["frozen"].String(), "0.00")

		// 渠道单号幂等: 同渠道单号给另一笔过审单登记 → 拒绝（防重复打款红线）
		no3, err := logic.WithdrawApply(ctx, a, "10.00")
		t.AssertNil(err)
		t.AssertNil(admin.AdminWithdrawAudit(ctx, no3, true, ""))
		t.Assert(errCode(admin.AdminWithdrawPay(ctx, no3, true, "CH-ORDER-X", "")), errcode.CodeStatusNotAllowed)

		// 状态机不可逆: 40 的单再审核 → 拒绝
		t.Assert(errCode(admin.AdminWithdrawAudit(ctx, win, false, "")), errcode.CodeStatusNotAllowed)

		// 失败回退: 20 的单打款失败 → 60 + 冻结回退
		t.AssertNil(admin.AdminWithdrawPay(ctx, no3, false, "", "渠道维护"))
		st, _ = g.DB().GetValue(ctx, "SELECT status FROM withdraw_order WHERE withdraw_no=?", no3)
		t.Assert(st.Int(), 60)
		acc, _ = g.DB().GetOne(ctx, "SELECT balance, frozen FROM user_account WHERE user_id=?", a)
		t.Assert(acc["frozen"].String(), "0.00")

		// 列表分页
		wl, err := logic.WithdrawList(ctx, a, 0, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(wl.Total >= 3, true)
	})
}

// ---- US2 佣金计提/结算/冲销 ----

// seedDistOrder 造"已完成订单"（订单头+单项, pay_amount 可控; fixture 名 TF-DIST- 前缀自清）。
func seedDistOrder(ctx context.Context, t *gtest.T, orderNo string, buyer int64, payYuan string) int64 {
	res, err := g.DB().Exec(ctx,
		"INSERT INTO trade_order(order_no,user_id,status,total_amount,promotion_amount,pay_amount,currency,"+
			"receiver_name,receiver_phone,receiver_province,receiver_city,receiver_detail) "+
			"VALUES(?,?,40,?,0,?,'CNY','分销测试','13800000000','浙江省','杭州市','T路')",
		orderNo, buyer, payYuan, payYuan)
	t.AssertNil(err)
	oid, _ := res.LastInsertId()
	ires, err := g.DB().Exec(ctx,
		"INSERT INTO trade_order_item(order_no,order_id,spu_id,sku_id,sku_no,spu_name,sku_name,sku_specs,"+
			"quantity,original_price,price,pay_amount) VALUES(?,?,1,?,'TF-D-SKU','分销商品','分销商品','{}',1,?,?,?)",
		orderNo, oid, buyer, payYuan, payYuan, payYuan)
	t.AssertNil(err)
	iid, _ := ires.LastInsertId()
	return iid
}

func cleanupDistOrder(ctx context.Context, orderNo string) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order_item WHERE order_no=?", orderNo)
	_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order WHERE order_no=?", orderNo)
}

// TestDistCommissionCycle 佣金全链（FR-3/FR-4, SC-3）:
// 规则命中计提 → 计提幂等（重放不重复）→ 保护期结算入账（余额+流水）→ 退款冲销（负额+扣回可负）。
func TestDistCommissionCycle(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		a := distCleanupUser(ctx, t, "TF-DC-A") // 一级受益人
		b := distCleanupUser(ctx, t, "TF-DC-B") // 买家
		defer distCleanupAll(ctx, a, b)
		const no = "TF-DIST-ORD-1"
		defer cleanupDistOrder(ctx, no)

		// 轻量 fixture: 分类 + SPU + SKU + 库存（user 包无 shop 的 trade fixture, 自建最小集）
		catId, err := g.DB().Model("product_category").Ctx(ctx).Data(g.Map{
			"parent_id": 0, "name": "TF-DIST-分类", "level": 3, "status": 1,
		}).InsertAndGetId()
		t.AssertNil(err)
		spuId, err := g.DB().Model("product_spu").Ctx(ctx).Data(g.Map{
			"spu_no": "TF-D-SPU-1", "name": "TF-DIST-商品", "category_id": catId,
			"images": `[]`, "spec_definitions": `[]`, "status": 1,
		}).InsertAndGetId()
		t.AssertNil(err)
		skuId, err := g.DB().Model("product_sku").Ctx(ctx).Data(g.Map{
			"sku_no": "TF-D-SKU", "spu_id": spuId, "name": "TF-DIST-商品 默认",
			"specs": "{}", "price": "100.00", "status": 1,
		}).InsertAndGetId()
		t.AssertNil(err)
		_, _ = g.DB().Exec(ctx,
			"INSERT INTO inventory(sku_id,total,locked,warn_count) VALUES(?,100,0,0) ON DUPLICATE KEY UPDATE total=100", skuId)
		defer func() {
			_, _ = g.DB().Exec(ctx, "DELETE FROM commission_rule WHERE scope_id=?", catId)
			_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id=?", skuId)
			_, _ = g.DB().Exec(ctx, "DELETE FROM product_sku WHERE id=?", skuId)
			_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE id=?", spuId)
			_, _ = g.DB().Exec(ctx, "DELETE FROM product_category WHERE id=?", catId)
		}()
		// trade_order_item.sku_id 用真实 sku（ sku_no 列非空已由 seedDistOrder 写死——此处回填对齐 spu/sku）
		logic := NewDistributionLogic()
		seedRelation(ctx, t, a, b) // B 的上级 = A
		seedDistributor(ctx, t, a, 2)
		// 规则: 该分类 L1=10%
		_, _ = g.DB().Exec(ctx, "DELETE FROM commission_rule WHERE scope_type=1 AND scope_id=?", catId)
		_, err = g.DB().Exec(ctx,
			"INSERT INTO commission_rule(scope_type,scope_id,level1_rate,level2_rate,status) VALUES(1,?,10.00,0.00,1)",
			catId)
		t.AssertNil(err)

		itemId := seedDistOrder(ctx, t, no, b, "100.00")
		_, _ = g.DB().Exec(ctx, "UPDATE trade_order_item SET spu_id=?, sku_id=? WHERE id=?", spuId, skuId, itemId)
		t.AssertNil(logic.SettleOrder(ctx, no))

		// 计提: A 待结算 10.00（基数 100 × 10%）
		rec, err := g.DB().GetOne(ctx,
			"SELECT id, amount, status FROM commission_record WHERE order_item_id=? AND beneficiary_user_id=?", itemId, a)
		t.AssertNil(err)
		t.Assert(rec["amount"].String(), "10.00")
		t.Assert(rec["status"].Int(), 1)

		// 计提幂等: 事件重放不重复落账
		t.AssertNil(logic.SettleOrder(ctx, no))
		cnt, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM commission_record WHERE order_item_id=? AND beneficiary_user_id=?", itemId, a)
		t.AssertNil(err)
		t.Assert(cnt.Int(), 1)

		// 保护期未满: 结算任务不入账
		n, err := logic.ConfirmSettle(ctx)
		t.AssertNil(err)
		t.Assert(n, int64(0))
		st, _ := g.DB().GetValue(ctx, "SELECT status FROM commission_record WHERE id=?", rec["id"].Int64())
		t.Assert(st.Int(), 1)

		// 保护期满（直接改 created_at）→ 结算入账: 余额 10 + 流水 biz_type=1 + 状态 2
		_, _ = g.DB().Exec(ctx,
			"UPDATE commission_record SET created_at=DATE_SUB(NOW(), INTERVAL 8 DAY) WHERE id=?", rec["id"].Int64())
		n, err = logic.ConfirmSettle(ctx)
		t.AssertNil(err)
		t.Assert(n, int64(1))
		acc, _ := g.DB().GetOne(ctx, "SELECT balance FROM user_account WHERE user_id=?", a)
		t.Assert(acc["balance"].String(), "10.00")
		lg, _ := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM account_log WHERE user_id=? AND biz_type=1", a)
		t.Assert(lg.Int(), 1)

		// 退款冲销: 已结算 → 负额冲销记录 + 余额扣回（10-10=0; 若多退则可负）
		t.AssertNil(logic.ReverseOnRefund(ctx, itemId))
		rev, err := g.DB().GetOne(ctx,
			"SELECT amount, status, reversal_of_id FROM commission_record WHERE reversal_of_id=?", rec["id"].Int64())
		t.AssertNil(err)
		t.Assert(rev["amount"].String(), "-10.00")
		t.Assert(rev["status"].Int(), 4)
		acc, _ = g.DB().GetOne(ctx, "SELECT balance FROM user_account WHERE user_id=?", a)
		t.Assert(acc["balance"].String(), "0.00")
		lg, _ = g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM account_log WHERE user_id=? AND biz_type=5", a)
		t.Assert(lg.Int(), 1)

		// 冲销幂等: 重复触发不再扣回
		t.AssertNil(logic.ReverseOnRefund(ctx, itemId))
		acc, _ = g.DB().GetOne(ctx, "SELECT balance FROM user_account WHERE user_id=?", a)
		t.Assert(acc["balance"].String(), "0.00")
	})
}

// ---- US4 推广码 ----

// TestDistShareCode 推广码稳定（首访生成, 再访一致; FR-6）。
func TestDistShareCode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		a := distCleanupUser(ctx, t, "TF-DSC-A")
		defer distCleanupAll(ctx, a)
		code1, link1, err := NewDistributionLogic().ShareCode(ctx, a)
		t.AssertNil(err)
		t.Assert(code1 != "", true)
		t.Assert(strings.Contains(link1, code1), true)
		code2, _, err := NewDistributionLogic().ShareCode(ctx, a)
		t.AssertNil(err)
		t.Assert(code2, code1) // 稳定
	})
}

// ---- 修复轮 helper: user 包内轻量商品 fixture（分类+SPU+SKU+库存, 自清） ----

type distFixture struct {
	catId int64
	spuId int64
	skuId int64
}

func newDistFixture(t *gtest.T) *distFixture {
	ctx := context.Background()
	f := &distFixture{}
	// 前置清理（幂等; spu_no/sku_no 唯一键——重复跑测试不留死行）
	_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id IN (SELECT id FROM product_sku WHERE sku_no='TF-D2-SKU')")
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_sku WHERE sku_no='TF-D2-SKU'")
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE spu_no='TF-D2-SPU'")
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_category WHERE name='TF-DIST2-分类'")
	var err error
	f.catId, err = g.DB().Model("product_category").Ctx(ctx).Data(g.Map{
		"parent_id": 0, "name": "TF-DIST2-分类", "level": 3, "status": 1,
	}).InsertAndGetId()
	t.AssertNil(err)
	f.spuId, err = g.DB().Model("product_spu").Ctx(ctx).Data(g.Map{
		"spu_no": "TF-D2-SPU", "name": "TF-DIST2-商品", "category_id": f.catId,
		"images": `[]`, "spec_definitions": `[]`, "status": 1,
	}).InsertAndGetId()
	t.AssertNil(err)
	f.skuId, err = g.DB().Model("product_sku").Ctx(ctx).Data(g.Map{
		"sku_no": "TF-D2-SKU", "spu_id": f.spuId, "name": "TF-DIST2-商品 默认",
		"specs": "{}", "price": "100.00", "status": 1,
	}).InsertAndGetId()
	t.AssertNil(err)
	_, _ = g.DB().Exec(ctx,
		"INSERT INTO inventory(sku_id,total,locked,warn_count) VALUES(?,100,0,0) ON DUPLICATE KEY UPDATE total=100", f.skuId)
	return f
}

func (f *distFixture) teardown(ctx context.Context) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id=?", f.skuId)
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_sku WHERE id=?", f.skuId)
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE id=?", f.spuId)
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_category WHERE id=?", f.catId)
}

// seedDistOrderSku 造订单项（sku/spu 可指定, 配合归因/规则命中测试）。
func seedDistOrderSku(ctx context.Context, t *gtest.T, orderNo string, buyer int64, payYuan string, skuId, spuId int64) int64 {
	res, err := g.DB().Exec(ctx,
		"INSERT INTO trade_order(order_no,user_id,status,total_amount,promotion_amount,pay_amount,currency,"+
			"receiver_name,receiver_phone,receiver_province,receiver_city,receiver_detail) "+
			"VALUES(?,?,40,?,0,?,'CNY','分销测试','13800000000','浙江省','杭州市','T路')",
		orderNo, buyer, payYuan, payYuan)
	t.AssertNil(err)
	oid, _ := res.LastInsertId()
	ires, err := g.DB().Exec(ctx,
		"INSERT INTO trade_order_item(order_no,order_id,spu_id,sku_id,sku_no,spu_name,sku_name,sku_specs,"+
			"quantity,original_price,price,pay_amount) "+
			"VALUES(?,?,?,?,?, '分销商品2','分销商品2 默认','{}',1,?,?,?)",
		orderNo, oid, spuId, skuId, "TF-D-SKU-"+orderNo, payYuan, payYuan, payYuan)
	t.AssertNil(err)
	iid, _ := ires.LastInsertId()
	return iid
}
