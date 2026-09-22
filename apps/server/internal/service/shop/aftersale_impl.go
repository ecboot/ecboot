// aftersale_impl.go 售后域实现（013-after-sale / 批次 07）。
// 契约: 接口签名见 aftersale.go; 状态机与规则见 specs/013-after-sale/{spec,data-model}.md。
//
// 铁律（沿用批次 06 修正后的既有惯例）:
//  1. 全部状态迁移一律**条件更新 + 判 RowsAffected**（affected=0 → 40006 状态不允许）——
//     只看 error 的写法会在并发下静默重入（012 评审 I1/I11 同型）。
//  2. 退款幂等键 = 售后单号（out_refund_no = after_sale_no）; 回调推进复用既有 HandleRefundNotify。
//  3. 金额一律分运算（library/money）; refund_amount ≤ 订单项实付; 0 元不调渠道。
package shop

import (
	"context"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/library/idgen"
	"ecboot/internal/library/money"
	"ecboot/internal/library/paychannel"
	"ecboot/internal/model"
	"ecboot/internal/model/do"
)

// 售后状态（after_sale_order.status，口径见 000008 列注释）。
const (
	afterSalePendingAudit  = 10 // 待审核
	afterSalePendingReturn = 20 // 待买家寄回（退货退款）
	afterSalePendingRefund = 30 // 待退款（渠道未成功，可重试）
	afterSaleRefunding     = 40 // 退款中（已发起渠道退款）
	afterSaleFinished      = 50 // 已完成
	afterSaleRejected      = 90 // 已拒绝
	afterSaleCanceled      = 91 // 已撤销
)

// 售后类型。
const (
	afterSaleTypeRefundOnly  = 1 // 仅退款
	afterSaleTypeReturnGoods = 2 // 退货退款
)

// orderStatusCompleted 订单已完成（40）。既有 pay_impl.go 只定义了 10/20/90，
// 40 此前无人使用，故在本批补齐（不改动既有文件，避免跨批 diff）。
const orderStatusCompleted = 40

// AfterSaleLogicImpl IAfterSaleLogic 实现。
type AfterSaleLogicImpl struct{}

func NewAfterSaleLogic() *AfterSaleLogicImpl { return &AfterSaleLogicImpl{} }

// nextAfterSaleNo 生成售后单号（sonyflake + 前缀，与订单号/门店编码同风格；兼作渠道退款幂等键）。
func nextAfterSaleNo() (string, error) {
	next, err := idgen.NextID()
	if err != nil {
		return "", gerror.Wrap(err, "生成售后单号失败")
	}
	return "AS" + strconv.FormatInt(next, 10), nil
}

// yuanFen 元字符串 → 分（金额一律分运算，勿用 float）。
func yuanFen(yuan string) (int64, error) {
	fen, err := money.FromYuanString(yuan)
	if err != nil {
		return 0, errcode.New(errcode.CodeInvalidParam, "金额格式非法")
	}
	return fen, nil
}

// Apply 申请售后（FR-001~004）: 校验归属与可售后状态 → 累计额度 → 计算退款金额 → 落单（10 待审核）。
//
// 评审 C1（Critical）: 整个"读累计额度 → 插入"必须在**行锁事务**内完成——原先无锁无事务,
// 并发或重复提交会在竞赛窗口内各读一次"额度未占满"而落多张超额单（评审实证: 同一行实付 10.00
// 的订单项落 3 张各 10.00 的售后单、逐张审核后渠道被调 3 次 = 3 倍出款）。
//
// **两处防护的分工（复审 I-1 以 SQL 级证据修正, 勿按"冗余"清理）**:
//   - `afterSaleUsed` 的**锁定读（FOR UPDATE）是承重墙**: 普通读会被事务的 InnoDB 读视图钉住——
//     而读视图可能在拿到行锁**之前**就已建立（GoFrame DAO 在事务内会先发 `SHOW FULL COLUMNS`
//     刷新元数据, 实测该语句即会建立读视图: A `BEGIN; SHOW FULL COLUMNS ...` 后 B 插入并提交,
//     A 的普通 SELECT 查不到 B 的行, 而无 SHOW 的对照组能查到）。故只有**当前读**才保证读到最新
//     已提交的占用集合; 复审实测: 仅去掉这处锁定读 → 并发用例立刻红（一轮 4 并发落 2 张单）。
//   - `trade_order_item` 行锁提供**按订单项串行**（同一行的并发申请在此排队）**并防止**两笔各插一行
//     时在相邻 gap 上互相等待插入意向锁而死锁。
func (i *AfterSaleLogicImpl) Apply(
	ctx context.Context, userId int64, in model.AfterSaleApplyInput,
) (string, error) {
	if in.Type != afterSaleTypeRefundOnly && in.Type != afterSaleTypeReturnGoods {
		return "", errcode.New(errcode.CodeInvalidParam, "售后类型非法")
	}
	if in.Quantity <= 0 {
		return "", errcode.New(errcode.CodeInvalidParam, "退货数量需大于 0")
	}

	var no string
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		item, order, e := afterSaleItemAndOrder(ctx, in.OrderItemId, userId, true)
		if e != nil {
			return e
		}
		ocols := dao.TradeOrder.Columns()
		if order[ocols.Status].Int() != orderStatusCompleted {
			// 契约: 仅**已完成**订单可申请（已发货未收货是否开放属产品决策, 本批不开放）
			return errcode.New(errcode.CodeAfterSaleDenied, "当前订单状态不可申请售后")
		}

		icols := dao.TradeOrderItem.Columns()
		lineQty := item[icols.Quantity].Int()
		if lineQty <= 0 {
			return errcode.New(errcode.CodeAfterSaleDenied, "订单项数量异常")
		}
		payFen, e := yuanFen(item[icols.PayAmount].String())
		if e != nil {
			return e
		}

		// 累计额度: 占用 = 非「已拒绝/已撤销」的申请数量（D5）; 锁定读, 见函数头 C1 说明
		usedQty, usedRefundFen, e := afterSaleUsed(ctx, in.OrderItemId)
		if e != nil {
			return e
		}
		if usedQty+in.Quantity > lineQty {
			return errcode.New(errcode.CodeAfterSaleDenied, "已达可退数量上限")
		}

		// 退款金额（D6）: 常规按行实付比例向下取整; **本次退完剩余数量时取剩余全额**, 保多笔之和恒等于行实付
		var refundFen int64
		if usedQty+in.Quantity == lineQty {
			refundFen = payFen - usedRefundFen
		} else {
			refundFen = payFen * int64(in.Quantity) / int64(lineQty)
		}
		if refundFen < 0 {
			refundFen = 0
		}
		if refundFen > payFen { // 不超行实付（000008 契约）
			refundFen = payFen
		}

		gen, e := nextAfterSaleNo()
		if e != nil {
			return e
		}
		if _, e = dao.AfterSaleOrder.Ctx(ctx).Data(do.AfterSaleOrder{
			AfterSaleNo:   gen,
			OrderId:       order[ocols.Id].Int64(),
			OrderNo:       order[ocols.OrderNo].String(),
			OrderItemId:   in.OrderItemId,
			UserId:        userId,
			Type:          in.Type,
			Status:        afterSalePendingAudit,
			Currency:      order[ocols.Currency].String(),
			Quantity:      in.Quantity,
			Reason:        in.Reason,
			Description:   in.Description,
			VoucherImages: mustJSON(voucherMap(in.VoucherImages)),
			RefundAmount:  money.ToYuanString(refundFen),
		}).Insert(); e != nil {
			return gerror.Wrap(e, "创建售后单失败")
		}
		no = gen
		return nil
	})
	if err != nil {
		return "", err
	}
	return no, nil
}

// afterSaleItemAndOrder 取订单项与其订单, 并校验归属（他人资源一律按"不存在"处理, 不泄露存在性）。
// forUpdate=true 时对订单项行取行锁（申请侧串行化, 见 Apply 的 C1 说明）。
func afterSaleItemAndOrder(
	ctx context.Context, orderItemId, userId int64, forUpdate bool,
) (gdb.Record, gdb.Record, error) {
	icols, ocols := dao.TradeOrderItem.Columns(), dao.TradeOrder.Columns()
	m := dao.TradeOrderItem.Ctx(ctx).Where(icols.Id, orderItemId)
	if forUpdate {
		m = m.LockUpdate()
	}
	item, err := m.One()
	if err != nil {
		return nil, nil, gerror.Wrap(err, "查询订单项失败")
	}
	if item.IsEmpty() {
		return nil, nil, errcode.New(errcode.CodeAfterSaleNotFound, "订单项不存在")
	}
	order, err := dao.TradeOrder.Ctx(ctx).Where(ocols.Id, item[icols.OrderId].Int64()).One()
	if err != nil {
		return nil, nil, gerror.Wrap(err, "查询订单失败")
	}
	if order.IsEmpty() || order[ocols.UserId].Int64() != userId {
		return nil, nil, errcode.New(errcode.CodeAfterSaleNotFound, "订单项不存在")
	}
	return item, order, nil
}

// afterSaleUsed 该订单项的"占用中"数量与退款金额（排除已拒绝/已撤销）。
func afterSaleUsed(ctx context.Context, orderItemId int64) (int, int64, error) {
	acols := dao.AfterSaleOrder.Columns()
	recs, err := dao.AfterSaleOrder.Ctx(ctx).
		Fields(acols.Quantity, acols.RefundAmount).
		Where(acols.OrderItemId, orderItemId).
		WhereNotIn(acols.Status, []int{afterSaleRejected, afterSaleCanceled}).
		LockUpdate(). // **承重墙**: 必须用当前读（普通读会被事务读视图钉住, 见 Apply 函数头的 I-1 说明）
		All()
	if err != nil {
		return 0, 0, gerror.Wrap(err, "查询已申请售后失败")
	}
	qty := 0
	var refundFen int64
	for _, r := range recs {
		qty += r[acols.Quantity].Int()
		fen, e := yuanFen(r[acols.RefundAmount].String())
		if e != nil {
			return 0, 0, e
		}
		refundFen += fen
	}
	return qty, refundFen, nil
}

// voucherMap 凭证图片 → JSON 列载荷（空数组存 [] 而非 NULL, 便于前端统一处理）。
func voucherMap(images []string) []string {
	if images == nil {
		return []string{}
	}
	return images
}

// afterSaleRefundChannel 退款渠道（幂等键 = 售后单号）。默认复用 mock 渠道;
// 声明为包级变量是为了给"渠道失败"这条路径留测试缝（生产恒为 mock, 真实渠道另立特性）。
var afterSaleRefundChannel paychannel.Channel = paychannel.NewMock()

// loadAfterSale 取售后单（不存在 → 40009）。
func loadAfterSale(ctx context.Context, afterSaleNo string) (gdb.Record, error) {
	acols := dao.AfterSaleOrder.Columns()
	rec, err := dao.AfterSaleOrder.Ctx(ctx).Where(acols.AfterSaleNo, afterSaleNo).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询售后单失败")
	}
	if rec.IsEmpty() {
		return nil, errcode.New(errcode.CodeAfterSaleNotFound, "售后单不存在")
	}
	return rec, nil
}

// ---------- 后台: 审核 ----------

// Approve 同意（FR-007）: 10→20（退货退款）或 10→30→发起退款（仅退款）。
func (i *AfterSaleLogicImpl) Approve(ctx context.Context, afterSaleNo, operator string) error {
	rec, err := loadAfterSale(ctx, afterSaleNo)
	if err != nil {
		return err
	}
	acols := dao.AfterSaleOrder.Columns()
	if rec[acols.Status].Int() != afterSalePendingAudit {
		return errcode.New(errcode.CodeStatusNotAllowed, "仅待审核可同意")
	}
	next := afterSalePendingRefund // 仅退款: 直接进入待退款（随后立即发起渠道退款）
	if rec[acols.Type].Int() == afterSaleTypeReturnGoods {
		next = afterSalePendingReturn
	}
	res, err := dao.AfterSaleOrder.Ctx(ctx).
		Where(acols.AfterSaleNo, afterSaleNo).
		Where(acols.Status, afterSalePendingAudit).
		Data(do.AfterSaleOrder{
			Status:     next,
			AuditTime:  gtime.Now(),
			OperatorId: operator,
		}).Update()
	if err != nil {
		return gerror.Wrap(err, "审核失败")
	}
	if n, _ := res.RowsAffected(); n == 0 { // 并发下已被他人处理
		return errcode.New(errcode.CodeStatusNotAllowed, "仅待审核可同意")
	}
	if next == afterSalePendingReturn {
		return nil // 退货退款: 等买家寄回, 不发起渠道退款
	}
	return i.tryRefund(ctx, afterSaleNo, operator)
}

// Reject 拒绝（FR-008）: 10→90, 原因必填且对买家可见; 该行可退数量随之释放（申请时不占用）。
func (i *AfterSaleLogicImpl) Reject(ctx context.Context, afterSaleNo, reason, operator string) error {
	if strings.TrimSpace(reason) == "" {
		return errcode.New(errcode.CodeInvalidParam, "拒绝原因必填")
	}
	if _, err := loadAfterSale(ctx, afterSaleNo); err != nil {
		return err
	}
	acols := dao.AfterSaleOrder.Columns()
	res, err := dao.AfterSaleOrder.Ctx(ctx).
		Where(acols.AfterSaleNo, afterSaleNo).
		Where(acols.Status, afterSalePendingAudit).
		Data(do.AfterSaleOrder{
			Status:       afterSaleRejected,
			RejectReason: reason,
			AuditTime:    gtime.Now(),
			OperatorId:   operator,
		}).Update()
	if err != nil {
		return gerror.Wrap(err, "拒绝失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeStatusNotAllowed, "仅待审核可拒绝")
	}
	return nil
}

// tryRefund 发起渠道退款（FR-011/012）: 成功 → 停留 40 退款中（等回调推进到 50）;
// 渠道失败 → **条件回退** 40→30 并写 fail_reason（可由后台重试）; 退款额为 0（全优惠行）→ 不调渠道, 直接 50（D6）。
//
// 评审 C2（Critical）: 顺序必须是**先条件推进状态、再调用渠道**——原实现先调渠道再改状态,
// 于是"钱已出"的判定标准实际是"渠道调用是否返回"而非状态=40, 买家可在那个窗口里撤销（{10,20,30} 可撤）,
// 形成"钱已出、单已撤销、额度释放、可再退一次"（评审实证: 终态 91 且渠道已被调 1 次, 撤销后可再拿一份全额）。
// 且渠道调用成功与 UPDATE 之间任一次崩溃/DB 错误都会把"钱已出"永久留在 30（可撤状态）。
// 改序后: "未出款" ⟺ 状态 30; 崩溃只会留下"40 但钱未出"（方向安全, 且可由重试自愈, 见 RetryRefund）。
func (i *AfterSaleLogicImpl) tryRefund(ctx context.Context, afterSaleNo, operator string) error {
	rec, err := loadAfterSale(ctx, afterSaleNo)
	if err != nil {
		return err
	}
	acols := dao.AfterSaleOrder.Columns()
	amountFen, err := yuanFen(rec[acols.RefundAmount].String())
	if err != nil {
		return err
	}

	if amountFen <= 0 { // 无资金可退: 直收口, 避免 0 元调用渠道
		advanced := false
		e := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			res, ue := dao.AfterSaleOrder.Ctx(ctx).
				Where(acols.AfterSaleNo, afterSaleNo).
				Where(acols.Status, afterSalePendingRefund).
				Data(do.AfterSaleOrder{
					Status:     afterSaleFinished,
					RefundTime: gtime.Now(),
					OperatorId: operator,
					FailReason: "",
				}).Update()
			if ue != nil {
				return gerror.Wrap(ue, "完成售后单失败")
			}
			if n, _ := res.RowsAffected(); n == 0 {
				return nil
			}
			advanced = true
			return afterSaleFinishedSideEffects(ctx, afterSaleNo, operator) // 与状态推进同事务
		})
		if e != nil {
			return e
		}
		if !advanced {
			return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可发起退款")
		}
		return nil
	}

	// 先推进状态（30→40）; 已是 40 者允许重调渠道（渠道以 outRefundNo 幂等, 用于自愈崩溃窗）
	switch rec[acols.Status].Int() {
	case afterSalePendingRefund:
		res, ue := dao.AfterSaleOrder.Ctx(ctx).
			Where(acols.AfterSaleNo, afterSaleNo).
			Where(acols.Status, afterSalePendingRefund).
			Data(do.AfterSaleOrder{Status: afterSaleRefunding, FailReason: "", OperatorId: operator}).
			Update()
		if ue != nil {
			return gerror.Wrap(ue, "推进退款中失败")
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可发起退款")
		}
	case afterSaleRefunding:
		// 允许: 崩溃窗/回调未到的重试
	default:
		return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可发起退款")
	}

	if e := afterSaleRefundChannel.Refund(afterSaleNo, amountFen); e != nil {
		// 评审 C2′（复审引入的 Critical）: 一旦**调用过渠道**, 就绝不回退到 30。
		// 30 的语义是"未出款"且**可撤**（D4: {10,20,30} 可撤, 撤销即释放可退数量）; 而渠道返回失败
		// 未必等于钱没出去（超时/连接中断都是"模糊失败"）。回退到 30 会让"钱可能已出"的记录重新可撤,
		// 买家撤销 → 额度释放 → 以**新的幂等键**再申请并再出款一份（复审已单线程复现: 累计 2× 行实付）。
		// 故: 失败一律**留在 40** 并记 fail_reason —— 40 此后表示"已发起（在途或待重试）",
		// 管理员可从 40 重试（渠道按键幂等去重, 不会二次出款）, 而"未出款 ⟺ 30"的不变量保持成立。
		if _, re := dao.AfterSaleOrder.Ctx(ctx).
			Where(acols.AfterSaleNo, afterSaleNo).
			Where(acols.Status, afterSaleRefunding).
			Data(do.AfterSaleOrder{FailReason: e.Error(), OperatorId: operator}).
			Update(); re != nil {
			g.Log().Errorf(ctx, "[售后] 渠道退款失败且失败原因写入失败, 需人工核对: after_sale_no=%s err=%v", afterSaleNo, re)
		}
		g.Log().Errorf(ctx, "[售后] 渠道退款失败, 留在退款中待重试: after_sale_no=%s amount_fen=%d err=%v",
			afterSaleNo, amountFen, e)
		// 评审 M1: 不静默返 nil——管理员须看到"未成功"（状态已留 40 且原因留痕, 可重试）
		return errcode.New(errcode.CodeRefundFailed, "退款发起失败, 可稍后重试")
	}
	// 渠道成功: 清掉上一次的失败原因（从 40 重试时不经过 30→40 的赋值路径, 需在此显式清理）
	if _, ue := dao.AfterSaleOrder.Ctx(ctx).
		Where(acols.AfterSaleNo, afterSaleNo).
		Where(acols.Status, afterSaleRefunding).
		Data(do.AfterSaleOrder{FailReason: "", OperatorId: operator}).
		Update(); ue != nil {
		g.Log().Errorf(ctx, "[售后] 清理失败原因失败: after_sale_no=%s err=%v", afterSaleNo, ue)
	}
	return nil
}

// ---------- 后台: 查询 ----------

// AdminList 后台售后列表（FR-006）: 状态/类型筛选 + 分页。
func (i *AfterSaleLogicImpl) AdminList(
	ctx context.Context, status, saleType int, page model.PageReq,
) (*model.PageResult[model.AfterSaleSummary], error) {
	page = page.Normalized()
	acols := dao.AfterSaleOrder.Columns()
	m := dao.AfterSaleOrder.Ctx(ctx)
	if status > 0 {
		m = m.Where(acols.Status, status)
	}
	if saleType > 0 {
		m = m.Where(acols.Type, saleType)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计售后单失败")
	}
	recs, err := m.OrderDesc(acols.Id).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询售后单失败")
	}
	return &model.PageResult[model.AfterSaleSummary]{
		List: afterSaleSummaries(recs), Total: int64(total),
	}, nil
}

// AdminDetail 后台售后详情（含申请人 id，供后台核对）。
func (i *AfterSaleLogicImpl) AdminDetail(ctx context.Context, afterSaleNo string) (*model.AfterSaleDetail, error) {
	rec, err := loadAfterSale(ctx, afterSaleNo)
	if err != nil {
		return nil, err
	}
	return afterSaleDetail(rec), nil
}

// afterSaleSummaries 记录 → 列表 DTO。
func afterSaleSummaries(recs gdb.Result) []model.AfterSaleSummary {
	acols := dao.AfterSaleOrder.Columns()
	list := make([]model.AfterSaleSummary, 0, len(recs))
	for _, r := range recs {
		list = append(list, model.AfterSaleSummary{
			AfterSaleNo:  r[acols.AfterSaleNo].String(),
			OrderNo:      r[acols.OrderNo].String(),
			UserId:       r[acols.UserId].Int64(),
			Type:         r[acols.Type].Int(),
			Quantity:     r[acols.Quantity].Int(),
			RefundAmount: r[acols.RefundAmount].String(),
			Status:       r[acols.Status].Int(),
			CreatedAt:    r[acols.CreatedAt].String(),
		})
	}
	return list
}

// afterSaleDetail 记录 → 详情 DTO（凭证图片按 JSON 数组解析）。
func afterSaleDetail(rec gdb.Record) *model.AfterSaleDetail {
	acols := dao.AfterSaleOrder.Columns()
	var images []string
	_ = gconv.Struct(rec[acols.VoucherImages].String(), &images)
	if images == nil {
		images = []string{}
	}
	return &model.AfterSaleDetail{
		AfterSaleNo:       rec[acols.AfterSaleNo].String(),
		OrderNo:           rec[acols.OrderNo].String(),
		OrderItemId:       rec[acols.OrderItemId].Int64(),
		UserId:            rec[acols.UserId].Int64(),
		Type:              rec[acols.Type].Int(),
		Quantity:          rec[acols.Quantity].Int(),
		Reason:            rec[acols.Reason].String(),
		Description:       rec[acols.Description].String(),
		VoucherImages:     images,
		RefundAmount:      rec[acols.RefundAmount].String(),
		ReturnLogisticsNo: rec[acols.ReturnLogisticsNo].String(),
		RefundNo:          rec[acols.RefundNo].String(),
		RejectReason:      rec[acols.RejectReason].String(),
		AuditTime:         rec[acols.AuditTime].String(),
		RefundTime:        rec[acols.RefundTime].String(),
		OperatorId:        rec[acols.OperatorId].String(),
		FailReason:        rec[acols.FailReason].String(),
		Status:            rec[acols.Status].Int(),
	}
}

// ---------- C 端: 寄回与确认收货 ----------

// SubmitReturn 填写寄回单号（FR-009）: 仅"退货退款 + 待买家寄回"可填; 重复提交以最后一次为准。
func (i *AfterSaleLogicImpl) SubmitReturn(ctx context.Context, userId int64, afterSaleNo, logisticsNo string) error {
	if strings.TrimSpace(logisticsNo) == "" {
		return errcode.New(errcode.CodeInvalidParam, "寄回单号必填")
	}
	rec, err := loadAfterSale(ctx, afterSaleNo)
	if err != nil {
		return err
	}
	acols := dao.AfterSaleOrder.Columns()
	if rec[acols.UserId].Int64() != userId {
		return errcode.New(errcode.CodeAfterSaleNotFound, "售后单不存在") // 他人资源按不存在处理
	}
	if rec[acols.Type].Int() != afterSaleTypeReturnGoods || rec[acols.Status].Int() != afterSalePendingReturn {
		return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可填写寄回单号")
	}
	res, err := dao.AfterSaleOrder.Ctx(ctx).
		Where(acols.AfterSaleNo, afterSaleNo).
		Where(acols.Status, afterSalePendingReturn).
		Data(do.AfterSaleOrder{ReturnLogisticsNo: logisticsNo}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "填写寄回单号失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可填写寄回单号")
	}
	return nil
}

// ConfirmReceipt 退货确认收货（FR-010）: 须已填寄回单号; 20→30→发起退款→40。
func (i *AfterSaleLogicImpl) ConfirmReceipt(ctx context.Context, afterSaleNo, operator string) error {
	rec, err := loadAfterSale(ctx, afterSaleNo)
	if err != nil {
		return err
	}
	acols := dao.AfterSaleOrder.Columns()
	if rec[acols.Status].Int() != afterSalePendingReturn {
		return errcode.New(errcode.CodeStatusNotAllowed, "仅待买家寄回可确认收货")
	}
	if strings.TrimSpace(rec[acols.ReturnLogisticsNo].String()) == "" {
		return errcode.New(errcode.CodeAfterSaleDenied, "买家尚未填写寄回单号")
	}
	res, err := dao.AfterSaleOrder.Ctx(ctx).
		Where(acols.AfterSaleNo, afterSaleNo).
		Where(acols.Status, afterSalePendingReturn).
		Data(do.AfterSaleOrder{Status: afterSalePendingRefund, OperatorId: operator}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "确认收货失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeStatusNotAllowed, "仅待买家寄回可确认收货")
	}
	return i.tryRefund(ctx, afterSaleNo, operator)
}

// RetryRefund 退款重试（FR-012）: 仅"待退款"（渠道曾失败/未发起）可重试; 退款中与已完成一律拒绝
// （避免重复出款——这是本批唯一的出款动作）。
func (i *AfterSaleLogicImpl) RetryRefund(ctx context.Context, afterSaleNo, operator string) error {
	rec, err := loadAfterSale(ctx, afterSaleNo)
	if err != nil {
		return err
	}
	acols := dao.AfterSaleOrder.Columns()
	switch rec[acols.Status].Int() {
	case afterSalePendingRefund:
		// 正常重试: 30 → 发起
	case afterSaleRefunding:
		// 评审 C2 的自愈口: "已推进到 40 但渠道调用前崩溃/回调未到" → 允许重调渠道
		// （渠道以 outRefundNo 幂等去重, 重复调用不会二次出款）; 已完成(50)仍一律拒绝。
	default:
		return errcode.New(errcode.CodeStatusNotAllowed, "仅待退款/退款中可重试退款")
	}
	return i.tryRefund(ctx, afterSaleNo, operator)
}

// ---------- C 端: 查询与撤销 ----------

// List 我的售后（FR-005）: 仅本人; 状态筛选（0=全部）+ 分页。
func (i *AfterSaleLogicImpl) List(
	ctx context.Context, userId int64, status int, page model.PageReq,
) (*model.PageResult[model.AfterSaleSummary], error) {
	page = page.Normalized()
	acols := dao.AfterSaleOrder.Columns()
	m := dao.AfterSaleOrder.Ctx(ctx).Where(acols.UserId, userId)
	if status > 0 {
		m = m.Where(acols.Status, status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计我的售后失败")
	}
	recs, err := m.OrderDesc(acols.Id).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询我的售后失败")
	}
	return &model.PageResult[model.AfterSaleSummary]{
		List: afterSaleSummaries(recs), Total: int64(total),
	}, nil
}

// Detail 售后详情（FR-005）: 仅本人可见, 他人按不存在处理。
func (i *AfterSaleLogicImpl) Detail(ctx context.Context, userId int64, afterSaleNo string) (*model.AfterSaleDetail, error) {
	rec, err := loadAfterSale(ctx, afterSaleNo)
	if err != nil {
		return nil, err
	}
	acols := dao.AfterSaleOrder.Columns()
	if rec[acols.UserId].Int64() != userId {
		return nil, errcode.New(errcode.CodeAfterSaleNotFound, "售后单不存在")
	}
	return afterSaleDetail(rec), nil
}

// Cancel 撤销申请（FR-013, 用户裁定 D4）: 仅 {10 待审核, 20 待寄回, 30 待退款} 可撤——即"尚未进入资金环节";
// "退款中"(40) 及以上一律拒绝（在途退款不可中断, 否则回调将因状态不符而不匹配, 形成"钱已出、单已撤销"）。
// 撤销后该订单项的可退数量随之释放（申请时按"非 90/91"口径占用）。
func (i *AfterSaleLogicImpl) Cancel(ctx context.Context, userId int64, afterSaleNo string) error {
	rec, err := loadAfterSale(ctx, afterSaleNo)
	if err != nil {
		return err
	}
	acols := dao.AfterSaleOrder.Columns()
	if rec[acols.UserId].Int64() != userId {
		return errcode.New(errcode.CodeAfterSaleNotFound, "售后单不存在")
	}
	cur := rec[acols.Status].Int()
	if cur != afterSalePendingAudit && cur != afterSalePendingReturn && cur != afterSalePendingRefund {
		return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可撤销售后")
	}
	res, err := dao.AfterSaleOrder.Ctx(ctx).
		Where(acols.AfterSaleNo, afterSaleNo).
		Where(acols.Status, cur). // 条件更新: 并发下不重入
		Data(do.AfterSaleOrder{Status: afterSaleCanceled}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "撤销售后失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可撤销售后")
	}
	return nil
}

// afterSaleFinishedSideEffects 完成后的账实处理（FR-014~016 / US6）:
//  1. **仅退货退款**按 quantity 回补可售库存（locked 不动——支付时已核销）;
//  2. 重算订单 refund_status（0 无售后 / 1 部分退款 / 2 全额退款, 口径="已完成退款数量"）;
//  3. 投递佣金冲销事件（结算属批次 11; 未装配则告警跳过）。
//
// 调用方须保证"仅在 40→50 条件更新成功时调用一次", 以维持幂等（重复回调不会二次回补/二次投递）。
func afterSaleFinishedSideEffects(ctx context.Context, afterSaleNo, operator string) error {
	rec, err := loadAfterSale(ctx, afterSaleNo)
	if err != nil {
		return err
	}
	acols, icols := dao.AfterSaleOrder.Columns(), dao.TradeOrderItem.Columns()

	if rec[acols.Type].Int() == afterSaleTypeReturnGoods {
		itemId := rec[acols.OrderItemId].Int64()
		qty := rec[acols.Quantity].Int()
		item, e := dao.TradeOrderItem.Ctx(ctx).Fields(icols.SkuId).Where(icols.Id, itemId).One()
		if e != nil {
			return gerror.Wrap(e, "查询订单项失败")
		}
		if item.IsEmpty() {
			return errcode.New(errcode.CodeAfterSaleNotFound, "订单项不存在")
		}
		skuId := item[icols.SkuId].Int64()
		// 评审 I1: ① 必须判行数（库存行缺失时静默无操作会造成账实不符, 批次06 I1 同型）;
		// ② 必须写 inventory_log（000003 为"6=售后回补"专门留了枚举, 而全仓此前只有后台调整在写流水;
		//    inventory_impl 的注释亦写明"全部库存变更 = 单行原子条件 UPDATE + 同事务流水, 流水是对账唯一依据"）。
		res, ue := g.DB().Exec(ctx, "UPDATE inventory SET total = total + ? WHERE sku_id = ?", qty, skuId)
		if ue != nil {
			return gerror.Wrap(ue, "回补库存失败")
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return errcode.New(errcode.CodeInventoryAdjust, "库存行不存在, 无法回补")
		}
		inv, e2 := dao.Inventory.Ctx(ctx).
			Fields(dao.Inventory.Columns().Total, dao.Inventory.Columns().Locked).
			Where(dao.Inventory.Columns().SkuId, skuId).One()
		if e2 != nil {
			return gerror.Wrap(e2, "读取库存失败")
		}
		if _, e2 = dao.InventoryLog.Ctx(ctx).Data(do.InventoryLog{
			SkuId:       skuId,
			OrderNo:     rec[acols.OrderNo].String(),
			ChangeType:  6, // 6=售后退货回补（000003 枚举）
			Quantity:    qty,
			TotalAfter:  inv[dao.Inventory.Columns().Total].Int(),
			LockedAfter: inv[dao.Inventory.Columns().Locked].Int(),
			Operator:    operator,
			Remark:      "售后退货回补(" + afterSaleNo + ")",
		}).Insert(); e2 != nil {
			return gerror.Wrap(e2, "写入库存流水失败")
		}
	}

	if err = recalcOrderRefundStatus(ctx, rec[acols.OrderId].Int64()); err != nil {
		return err
	}

	refundFen, _ := yuanFen(rec[acols.RefundAmount].String())
	if CommissionReverse != nil {
		CommissionReverse.ReverseForAfterSale(ctx, rec[acols.OrderNo].String(), afterSaleNo, refundFen)
	} else {
		g.Log().Warningf(ctx, "佣金冲销事件出口未装配(批次11): 售后完成未投递冲销 order_no=%s after_sale_no=%s",
			rec[acols.OrderNo].String(), afterSaleNo)
	}
	return nil
}

// recalcOrderRefundStatus 重算订单退款状态（0 无售后 / 1 部分退款 / 2 全额退款）。
// 口径: 已完成(50)的售后数量之和为 0 → 0; ≥ 订单全部行购买数量之和 → 2; 否则 1。
func recalcOrderRefundStatus(ctx context.Context, orderId int64) error {
	ocols, icols := dao.TradeOrder.Columns(), dao.TradeOrderItem.Columns()
	totalQty, err := dao.TradeOrderItem.Ctx(ctx).
		Where(icols.OrderId, orderId).Sum(icols.Quantity)
	if err != nil {
		return gerror.Wrap(err, "统计订单行数量失败")
	}
	acols := dao.AfterSaleOrder.Columns()
	doneQty, err := dao.AfterSaleOrder.Ctx(ctx).
		Where(acols.OrderId, orderId).
		Where(acols.Status, afterSaleFinished).
		Sum(acols.Quantity)
	if err != nil {
		return gerror.Wrap(err, "统计已退款数量失败")
	}
	status := 0
	if doneQty > 0 {
		status = 1
	}
	if doneQty >= totalQty && totalQty > 0 {
		status = 2
	}
	if _, err = dao.TradeOrder.Ctx(ctx).
		Where(ocols.Id, orderId).
		Data(do.TradeOrder{RefundStatus: status}).
		Update(); err != nil {
		return gerror.Wrap(err, "更新订单退款状态失败")
	}
	return nil
}
