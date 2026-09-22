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
func (i *AfterSaleLogicImpl) Apply(
	ctx context.Context, userId int64, in model.AfterSaleApplyInput,
) (string, error) {
	if in.Type != afterSaleTypeRefundOnly && in.Type != afterSaleTypeReturnGoods {
		return "", errcode.New(errcode.CodeInvalidParam, "售后类型非法")
	}
	if in.Quantity <= 0 {
		return "", errcode.New(errcode.CodeInvalidParam, "退货数量需大于 0")
	}

	item, order, err := afterSaleItemAndOrder(ctx, in.OrderItemId, userId)
	if err != nil {
		return "", err
	}
	ocols := dao.TradeOrder.Columns()
	if order[ocols.Status].Int() != orderStatusCompleted {
		// 契约: 仅**已完成**订单可申请（已发货未收货是否开放属产品决策, 本批不开放）
		return "", errcode.New(errcode.CodeAfterSaleDenied, "当前订单状态不可申请售后")
	}

	icols := dao.TradeOrderItem.Columns()
	lineQty := item[icols.Quantity].Int()
	if lineQty <= 0 {
		return "", errcode.New(errcode.CodeAfterSaleDenied, "订单项数量异常")
	}
	payFen, err := yuanFen(item[icols.PayAmount].String())
	if err != nil {
		return "", err
	}

	// 累计额度: 占用 = 非「已拒绝/已撤销」的申请数量（D5）
	usedQty, usedRefundFen, err := afterSaleUsed(ctx, in.OrderItemId)
	if err != nil {
		return "", err
	}
	if usedQty+in.Quantity > lineQty {
		return "", errcode.New(errcode.CodeAfterSaleDenied, "已达可退数量上限")
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

	no, err := nextAfterSaleNo()
	if err != nil {
		return "", err
	}
	_, err = dao.AfterSaleOrder.Ctx(ctx).Data(do.AfterSaleOrder{
		AfterSaleNo:   no,
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
	}).Insert()
	if err != nil {
		return "", gerror.Wrap(err, "创建售后单失败")
	}
	return no, nil
}

// afterSaleItemAndOrder 取订单项与其订单, 并校验归属（他人资源一律按"不存在"处理, 不泄露存在性）。
func afterSaleItemAndOrder(
	ctx context.Context, orderItemId, userId int64,
) (gdb.Record, gdb.Record, error) {
	icols, ocols := dao.TradeOrderItem.Columns(), dao.TradeOrder.Columns()
	item, err := dao.TradeOrderItem.Ctx(ctx).Where(icols.Id, orderItemId).One()
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

// tryRefund 发起渠道退款（FR-011/012）: 成功 → 40 退款中; 渠道失败 → **停 30** 并写 fail_reason（可由后台重试）;
// 退款额为 0（全优惠行）→ 不调渠道, 直接 50 完成（D6）。
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
			return afterSaleFinishedSideEffects(ctx, afterSaleNo) // 与状态推进同事务
		})
		if e != nil {
			return e
		}
		if !advanced {
			return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可发起退款")
		}
		return nil
	}

	if e := afterSaleRefundChannel.Refund(afterSaleNo, amountFen); e != nil {
		// 渠道失败: 停在待退款并把原因留痕（可重试）; 不回滚状态机（30 是失败后的正确落点）
		_, _ = dao.AfterSaleOrder.Ctx(ctx).
			Where(acols.AfterSaleNo, afterSaleNo).
			Where(acols.Status, afterSalePendingRefund).
			Data(do.AfterSaleOrder{FailReason: e.Error(), OperatorId: operator}).Update()
		g.Log().Errorf(ctx, "[售后] 渠道退款失败, 待重试: after_sale_no=%s amount_fen=%d err=%v",
			afterSaleNo, amountFen, e)
		return nil
	}
	res, err := dao.AfterSaleOrder.Ctx(ctx).
		Where(acols.AfterSaleNo, afterSaleNo).
		Where(acols.Status, afterSalePendingRefund).
		Data(do.AfterSaleOrder{Status: afterSaleRefunding, FailReason: "", OperatorId: operator}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "推进退款中失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可发起退款")
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
	if rec[acols.Status].Int() != afterSalePendingRefund {
		return errcode.New(errcode.CodeStatusNotAllowed, "仅待退款可重试退款")
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
func afterSaleFinishedSideEffects(ctx context.Context, afterSaleNo string) error {
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
		if _, e = g.DB().Exec(ctx,
			"UPDATE inventory SET total = total + ? WHERE sku_id = ?", qty, skuId); e != nil {
			return gerror.Wrap(e, "回补库存失败")
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
