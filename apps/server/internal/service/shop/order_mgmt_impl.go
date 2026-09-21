// order_mgmt_impl.go 后台订单管理面（012-trade-wiring; 接口契约见 order.go IOrderLogic 后台段）。
// 语义: 发货仅待发货(20)→待收货(30) 且**物流公司须存在且启用**（批次 03 字典）;
// 后台取消与 C 端同语义（订单 90 + 资源释放, 复用既有取消路径）;
// 内部备注写 seller_remark（C 端 Detail 查询口径不含该列——买家不可见）;
// CancelTimeout/AutoConfirm 为定时入口（配置化时限由调度侧控制, 本层按入参/默认推进）。
package shop

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// AdminList 后台订单列表（FR-009）: 状态/订单号/用户关键词/时间区间筛选 + 分页。
func (i *OrderLogicImpl) AdminList(ctx context.Context, q model.AdminOrderQuery) (*model.PageResult[model.AdminOrderSummary], error) {
	page := q.Normalized()
	cols := dao.TradeOrder.Columns()
	m := dao.TradeOrder.Ctx(ctx)
	if q.Status > 0 {
		m = m.Where(cols.Status, q.Status)
	}
	if q.OrderNo != "" {
		m = m.Where(cols.OrderNo, q.OrderNo)
	}
	if q.UserKeyword != "" {
		// 用户关键词（ID 或手机号精确）→ 先解析目标用户集合
		uids, err := matchUserIds(ctx, q.UserKeyword)
		if err != nil {
			return nil, err
		}
		if len(uids) == 0 {
			return &model.PageResult[model.AdminOrderSummary]{List: []model.AdminOrderSummary{}, Total: 0}, nil
		}
		m = m.WhereIn(cols.UserId, uids)
	}
	if q.StartTime != "" {
		m = m.WhereGTE(cols.CreatedAt, q.StartTime)
	}
	if q.EndTime != "" {
		m = m.WhereLTE(cols.CreatedAt, q.EndTime)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计订单失败")
	}
	recs, err := m.OrderDesc(cols.Id).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单失败")
	}
	list := make([]model.AdminOrderSummary, 0, len(recs))
	for _, r := range recs {
		orderNo := r[cols.OrderNo].String()
		briefs, err := orderItemBriefs(ctx, orderNo)
		if err != nil {
			return nil, err
		}
		list = append(list, model.AdminOrderSummary{
			OrderNo:   orderNo,
			UserId:    fmt.Sprintf("%d", r[cols.UserId].Int64()),
			Status:    r[cols.Status].Int(),
			PayAmount: r[cols.PayAmount].String(),
			Items:     briefs,
			CreatedAt: r[cols.CreatedAt].String(),
		})
	}
	return &model.PageResult[model.AdminOrderSummary]{List: list, Total: int64(total)}, nil
}

// orderItemBriefs 订单项摘要（后台无归属限制——运营视角）。
func orderItemBriefs(ctx context.Context, orderNo string) ([]model.OrderItemBrief, error) {
	icols := dao.TradeOrderItem.Columns()
	recs, err := dao.TradeOrderItem.Ctx(ctx).Where(icols.OrderNo, orderNo).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单项失败")
	}
	out := make([]model.OrderItemBrief, 0, len(recs))
	for _, r := range recs {
		out = append(out, model.OrderItemBrief{
			SpuName:  r[icols.SpuName].String(),
			SkuSpecs: specsMap(r[icols.SkuSpecs].String()),
			Image:    r[icols.SkuImage].String(),
			Quantity: r[icols.Quantity].Int(),
			Price:    r[icols.Price].String(),
		})
	}
	return out, nil
}

// matchUserIds 用户关键词 → 用户 ID 集合。
// 注: 本批仅支持**用户 ID 精确匹配**; 手机号匹配需 PhoneCipher（其单例在 user 域私有）——
// 跨域直接 import 违反分层（兄弟域互禁）, 已在 PROGRESS 记账待后续以共享库方式补齐。
func matchUserIds(ctx context.Context, keyword string) ([]int64, error) {
	var id int64
	if _, err := fmt.Sscan(keyword, &id); err != nil || id <= 0 {
		return nil, nil
	}
	ucols := dao.User.Columns()
	rec, err := dao.User.Ctx(ctx).Fields(ucols.Id).Where(ucols.Id, id).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询用户失败")
	}
	if rec.IsEmpty() {
		return nil, nil
	}
	return []int64{rec[ucols.Id].Int64()}, nil
}

// AdminDetail 后台订单详情（无归属限制——复用 C 端详情, userId=0 表后台视角）。
func (i *OrderLogicImpl) AdminDetail(ctx context.Context, orderNo string) (*model.OrderDetailView, error) {
	return i.OrderDetail(ctx, 0, orderNo)
}

// Deliver 发货（FR-010）: 仅待发货(20) → 待收货(30); 物流公司须存在且启用。
func (i *OrderLogicImpl) Deliver(ctx context.Context, orderNo, logisticsCode, deliverNo, operator string) error {
	// 物流公司校验（批次 03 字典: 存在且启用）
	lcols := dao.LogisticsCompany.Columns()
	cnt, err := dao.LogisticsCompany.Ctx(ctx).
		Where(lcols.Code, logisticsCode).
		Where(lcols.Status, 1).
		Where(lcols.Deleted, 0).
		Count()
	if err != nil {
		return gerror.Wrap(err, "查询物流公司失败")
	}
	if cnt == 0 {
		return errcode.New(errcode.CodeInvalidParam, "物流公司不存在或已停用")
	}

	cols := dao.TradeOrder.Columns()
	rec, err := dao.TradeOrder.Ctx(ctx).Where(cols.OrderNo, orderNo).One()
	if err != nil {
		return gerror.Wrap(err, "查询订单失败")
	}
	if rec.IsEmpty() {
		return errcode.New(errcode.CodeOrderNotFound, "订单不存在")
	}
	if rec[cols.Status].Int() != 20 {
		return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可发货")
	}
	res, err := dao.TradeOrder.Ctx(ctx).
		Where(cols.OrderNo, orderNo).
		Where(cols.Status, 20).
		Data(g.Map{
			cols.Status:         30,
			cols.DeliverCompany: logisticsCode,
			cols.DeliverNo:      deliverNo,
			cols.DeliverTime:    gtime.Now(),
		}).Update()
	if err != nil {
		return gerror.Wrap(err, "发货失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可发货")
	}
	_, _ = dao.TradeOrderLog.Ctx(ctx).Data(g.Map{
		"order_no":      orderNo,
		"order_id":      rec[cols.Id].Int64(),
		"from_status":   20,
		"to_status":     30,
		"remark":        "发货: " + logisticsCode + " " + deliverNo,
		"operator_type": 3, // 3=管理员
		"operator_id":   operator,
	}).Insert()
	return nil
}

// AdminCancel 后台取消（FR-011）: 与 C 端同语义（状态 90 + 资源释放）, 但审计口径为管理员。
// 012 修复轮 I6 补漏: 原先丢弃 operator 且走 Cancel(0) → cancel_type 记"用户"、operator_type 也用错枚举,
// 审计上无法区分"用户自己取消"与"后台取消"。现在 operator（controller 侧 admin:{id}）原样落到 operator_id。
func (i *OrderLogicImpl) AdminCancel(ctx context.Context, orderNo, reason, operator string) error {
	opId := operator
	if opId == "" {
		opId = "admin"
	}
	return i.cancelBy(ctx, 0, orderNo, reason, cancelTypeAdmin, opTypeAdmin, opId)
}

// SellerRemark 内部备注（FR-012）: 买家不可见（C 端查询口径不含该列）。
func (i *OrderLogicImpl) SellerRemark(ctx context.Context, orderNo, remark, operator string) error {
	cols := dao.TradeOrder.Columns()
	res, err := dao.TradeOrder.Ctx(ctx).
		Where(cols.OrderNo, orderNo).
		Data(g.Map{
			cols.SellerRemark: remark,
			cols.UpdatedAt:    gtime.Now(),
		}).Update()
	if err != nil {
		return gerror.Wrap(err, "写备注失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		// 值未变或订单不存在: 判存在性
		cnt, cerr := dao.TradeOrder.Ctx(ctx).Where(cols.OrderNo, orderNo).Count()
		if cerr != nil {
			return gerror.Wrap(cerr, "查询订单失败")
		}
		if cnt == 0 {
			return errcode.New(errcode.CodeOrderNotFound, "订单不存在")
		}
	}
	_ = operator
	return nil
}

// CancelTimeout 超时取消扫描（内部方法; 时限由配置/调度控制, 本层按传入分钟数扫描）。
func (i *OrderLogicImpl) CancelTimeout(ctx context.Context) (int, error) {
	cols := dao.TradeOrder.Columns()
	recs, err := dao.TradeOrder.Ctx(ctx).
		Fields(cols.OrderNo).
		Where(cols.Status, 10).
		Where("created_at < DATE_SUB(NOW(), INTERVAL 30 MINUTE)").
		All()
	if err != nil {
		return 0, gerror.Wrap(err, "扫描超时订单失败")
	}
	n := 0
	for _, r := range recs {
		// 012 修复轮 I6 补漏: 超时取消的取消方是"系统", 不是用户也不是管理员——
		// 原走 Cancel(0) 致 cancel_type/operator_type 双双记错。
		if err := i.cancelBy(ctx, 0, r[cols.OrderNo].String(), "超时未支付自动取消",
			cancelTypeSystem, opTypeSystem, "system:timeout"); err != nil {
			continue
		}
		n++
	}
	return n, nil
}

// AutoConfirm 自动确认收货（内部方法; 发货后 7 天）。
func (i *OrderLogicImpl) AutoConfirm(ctx context.Context) (int, error) {
	cols := dao.TradeOrder.Columns()
	res, err := dao.TradeOrder.Ctx(ctx).
		Where(cols.Status, 30).
		Where("deliver_time IS NOT NULL AND deliver_time < DATE_SUB(NOW(), INTERVAL 7 DAY)").
		Data(g.Map{cols.Status: 40, cols.FinishTime: gtime.Now()}).
		Update()
	if err != nil {
		return 0, gerror.Wrap(err, "自动确认收货失败")
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}
