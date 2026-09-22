// dashboard_impl.go IDashboardLogic 实现（019 批次 13）——只读统计, 口径全部锚定既有权威查询。
// 交易: 已支付订单口径（status>=20）/销售额=pay_amount 合计/退款额=售后完成实退合计/待发货=status 20。
// 会员: 窗口新增（created_at）/活跃（last_active_at 窗口内）/休眠（last_active_at ≥90 天, 与 000029 同源）。
// 商品: 在售 SPU/低库存（inventory available<=warn_count 既有口径）/待审核评价（audit_status=0）。
package system

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/library/money"
	"ecboot/internal/model"
)

// DashboardLogicImpl IDashboardLogic 实现。
type DashboardLogicImpl struct{}

func NewDashboardLogic() *DashboardLogicImpl { return &DashboardLogicImpl{} }

var _ IDashboardLogic = (*DashboardLogicImpl)(nil)

// Trade 交易看板（FR-1）。
func (i *DashboardLogicImpl) Trade(ctx context.Context, startTime, endTime string) (*model.TradeDashboard, error) {
	out := &model.TradeDashboard{}
	// 已支付口径（status>=20 含发货/完成; 10=待付款不计入销售额）
	m := g.DB().Model("trade_order").Ctx(ctx).Where("status >= 20")
	if startTime != "" {
		m = m.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		m = m.Where("created_at <= ?", endTime)
	}
	cnt, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计订单数失败")
	}
	out.OrderCount = int64(cnt)
	// 金额走 Value（字符串 decimal）+ money 规范化——Sum 返回 float64 违反"禁 float 存算金额"铁律
	sales, err := m.Value("COALESCE(SUM(pay_amount),0)")
	if err != nil {
		return nil, gerror.Wrap(err, "统计销售额失败")
	}
	out.SalesAmount = money.ToYuanString(mustFen(sales.String()))

	// 待发货（当前值, 不受窗口影响——运营实时口径）
	pd, err := g.DB().Model("trade_order").Ctx(ctx).Where("status", 20).Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计待发货失败")
	}
	out.PendingDeliver = int64(pd)

	// 退款额: 售后完成（status=50）实退合计; 口径与售后域状态枚举对齐
	rm := g.DB().Model("after_sale_order").Ctx(ctx).Where("status", 50)
	if startTime != "" {
		rm = rm.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		rm = rm.Where("created_at <= ?", endTime)
	}
	refund, err := rm.Value("COALESCE(SUM(refund_amount),0)")
	if err != nil {
		return nil, gerror.Wrap(err, "统计退款额失败")
	}
	out.RefundAmount = money.ToYuanString(mustFen(refund.String()))
	return out, nil
}

// Member 会员看板（FR-2）。
func (i *DashboardLogicImpl) Member(ctx context.Context, startTime, endTime string) (*model.MemberDashboard, error) {
	out := &model.MemberDashboard{}
	nm := g.DB().Model("user").Ctx(ctx).Where("deleted", 0)
	if startTime != "" {
		nm = nm.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		nm = nm.Where("created_at <= ?", endTime)
	}
	nc, err := nm.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计新增会员失败")
	}
	out.NewCount = int64(nc)

	// 活跃: last_active_at 在窗口内（无窗口则近 30 天）
	activeFrom := startTime
	if activeFrom == "" {
		activeFrom = time.Now().AddDate(0, 0, -30).UTC().Format("2006-01-02 15:04:05")
	}
	am := g.DB().Model("user").Ctx(ctx).
		Where("deleted", 0).Where("last_active_at >= ?", activeFrom)
	if endTime != "" {
		am = am.Where("last_active_at <= ?", endTime)
	}
	ac, err := am.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计活跃会员失败")
	}
	out.ActiveCount = int64(ac)

	// 休眠: 最后活跃 ≥90 天（与 000029 休眠分级同源）
	dc, err := g.DB().Model("user").Ctx(ctx).
		Where("deleted", 0).
		Where("status", 1).
		Where("last_active_at IS NOT NULL").
		Where("last_active_at <= DATE_SUB(NOW(), INTERVAL 90 DAY)").Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计休眠会员失败")
	}
	out.DormantCount = int64(dc)
	return out, nil
}

// Product 商品看板（FR-3）。
func (i *DashboardLogicImpl) Product(ctx context.Context) (*model.ProductDashboard, error) {
	out := &model.ProductDashboard{}
	onSale, err := g.DB().Model("product_spu").Ctx(ctx).
		Where("status", 1).Where("deleted", 0).Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计在售商品失败")
	}
	out.OnSaleCount = int64(onSale)

	// 低库存: 复用 inventory 既有预警口径（available <= warn_count, 与 inventory_impl 一致）
	low, err := g.DB().Model("inventory").Ctx(ctx).
		Where("(total - locked) <= warn_count").Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计低库存失败")
	}
	out.LowStockCount = int64(low)

	// 待审核评价（C 端 4 端点批次 08 的 audit_status 口径; V1 自动通过时该数恒小）
	pr, err := g.DB().Model("product_review").Ctx(ctx).
		Where("audit_status", 0).Where("deleted", 0).Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计待审核评价失败")
	}
	out.PendingReview = int64(pr)
	return out, nil
}

// mustFen 元字符串 → 分（解析失败归 0——统计面不因单值异常中断）。
func mustFen(yuan string) int64 {
	fen, err := money.FromYuanString(yuan)
	if err != nil {
		return 0
	}
	return fen
}
