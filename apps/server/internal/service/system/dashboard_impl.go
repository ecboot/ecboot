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
	// 已支付口径（C1 评审修复: **显式枚举 20/30/40**, 原 `>= 20` 把 90=已取消计入——
	// cancelBy 只改状态不清零 pay_amount, 取消单的应付额被计成销售额（评审探针: +888）;
	// 10=待付款同样不计入）
	m := g.DB().Model("trade_order").Ctx(ctx).WhereIn("status", []int{20, 30, 40})
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

	// 活跃（I4 评审修复——口径统一）: last_active_at 在窗口内; **无窗口默认近 30 天**
	// （spec/api dc 已同步为"近 30 天内活跃"; 原三处表述不一: spec"全量"/api"周期内"/实现"近30天"）
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

	// 休眠（I5 二次收口——**读法与权威实现同源**）: 阈值取 `system_config` 表（dormant.tier1.days,
	// 000033 种子）, 与 user/wx.go 的 ensureNotDormant→cfgInt 完全同源——原修复误用 `g.Cfg()`
	// （配置文件/环境变量），而该键**只存在于 system_config 表**, 故恒取兜底 90 = 与原硬编码等价,
	// 且运营改阈值后登录门禁与看板统计会分叉（复审决定性探针: 表值改 2 后 g.Cfg 仍 90）。
	// 休眠定义在 **000027**（原注释误引 000029=bargain_assist）; 口径: 最后活跃 ≥ 阈值天 + 未删 + 正常态
	dormantDays := cfgIntSystem(ctx, "dormant.tier1.days", 90)
	dc, err := g.DB().Model("user").Ctx(ctx).
		Where("deleted", 0).
		Where("status", 1).
		Where("last_active_at IS NOT NULL").
		Where("last_active_at <= DATE_SUB(NOW(), INTERVAL ? DAY)", dormantDays).Count()
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

// cfgIntSystem 读 system_config 整型配置（停用/缺失回退默认; 与 user 域的 cfgInt 同源同表——
// 两域各自持有一份读法以免兄弟域互引, 但**数据源唯一**=system_config）。
func cfgIntSystem(ctx context.Context, code string, fallback int) int {
	v, err := g.DB().GetOne(ctx,
		"SELECT value FROM system_config WHERE code=? AND status=1 AND deleted=0", code)
	if err != nil || v.IsEmpty() {
		return fallback
	}
	if n := v["value"].Int(); n > 0 {
		return n
	}
	return fallback
}

// mustFen 元字符串 → 分（解析失败归 0——统计面不因单值异常中断）。
func mustFen(yuan string) int64 {
	fen, err := money.FromYuanString(yuan)
	if err != nil {
		return 0
	}
	return fen
}
