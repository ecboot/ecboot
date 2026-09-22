// ports.go 跨域端口装配（012 评审 I5）。
// 形态: 依赖倒置——shop 域在 service/shop/ports.go 声明接口, 本层（装配层）把 user 域实现注入,
// 两个域兄弟因此互不 import（分层契约: service/user ↔ service/shop 禁止互相引用）。
// 已接: ICouponQuery（结算试算的可用券列表——本批 22 端点中 /shop/cart/checkout 的出参义务）。
// 未接: ICouponTrade/IPointTrade/IAccountTrade——均属下单事务三段联动, 随账户域（批次 11）一并接,
// 见 PROGRESS §五 2026-09-21 批次 06 记账。
package bootstrap

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/model"
	"ecboot/internal/service/shop"
	"ecboot/internal/service/user"
)

func init() {
	shop.CouponQuery = couponQueryAdapter{}
}

// couponQueryAdapter 券查询口适配: user 域 DTO → shop 域 DTO（各域只认自己的出参形态）。
type couponQueryAdapter struct{}

func (couponQueryAdapter) UsableForOrder(
	ctx context.Context, userId int64, goodsAmountYuan string,
) ([]model.UsableCouponBrief, error) {
	items, err := user.UsableForOrder(ctx, userId, goodsAmountYuan)
	if err != nil {
		return nil, err
	}
	// 两个 DTO 字段完全一致 → 直接类型转换（加字段时编译器会在此报错, 比手写映射更安全）
	out := make([]model.UsableCouponBrief, 0, len(items))
	for _, it := range items {
		out = append(out, model.UsableCouponBrief(it))
	}
	return out, nil
}

// ---------- 分销资金域装配（017 批次 11） ----------

func init() {
	// 正向: 订单确认收货 → 佣金计提（SettleOrder, 归因+规则命中在 user 域）
	shop.CommissionSettle = commissionSettleAdapter{}
	// 反向: 售后完成 → 佣金冲销（批次 07 投递侧的消费端落地; 适配器把 afterSaleNo 还原为订单项集合）
	shop.CommissionReverse = commissionReverseAdapter{}
}

// commissionSettleAdapter 正向计提适配。
type commissionSettleAdapter struct{}

func (commissionSettleAdapter) OnOrderConfirmed(ctx context.Context, orderNo string) {
	if err := user.NewDistributionLogic().SettleOrder(ctx, orderNo); err != nil {
		g.Log().Errorf(ctx, "[佣金计提] 确认收货计提失败(事件已投递, 需人工核对): order_no=%s err=%v", orderNo, err)
	}
}

// commissionReverseAdapter 冲销适配: user 域冲销按 order_item_id 粒度。
// I3（评审修复）: 按售后单的 order_item_id **精确冲销**（spec FR-4"该项退款→该项冲销"）——
// 原实现展开订单全部订单项、忽略 refundFen, 多行订单部分退款会冲掉整单佣金（零记账偏离）。
type commissionReverseAdapter struct{}

func (commissionReverseAdapter) ReverseForAfterSale(ctx context.Context, orderNo, afterSaleNo string, refundFen int64) {
	itemId, err := g.DB().Model("after_sale_order").Ctx(ctx).
		Fields("order_item_id").Where("after_sale_no", afterSaleNo).Value("order_item_id")
	if err != nil {
		g.Log().Errorf(ctx, "[佣金冲销] 售后单查询失败(需人工核对): after_sale_no=%s err=%v", afterSaleNo, err)
		return
	}
	if itemId == nil || itemId.Int64() <= 0 {
		g.Log().Warningf(ctx, "[佣金冲销] 售后单无订单项(仅退款整单?): after_sale_no=%s — 按订单全项冲销", afterSaleNo)
		ids, e := g.DB().Model("trade_order_item").Ctx(ctx).Fields("id").Where("order_no", orderNo).All()
		if e != nil {
			g.Log().Errorf(ctx, "[佣金冲销] 订单项查询失败(需人工核对): order_no=%s err=%v", orderNo, e)
			return
		}
		logic := user.NewDistributionLogic()
		for _, r := range ids {
			if err = logic.ReverseOnRefund(ctx, r["id"].Int64()); err != nil {
				g.Log().Errorf(ctx, "[佣金冲销] 冲销失败(需人工核对): order_no=%s order_item_id=%d err=%v", orderNo, r["id"].Int64(), err)
			}
		}
		return
	}
	if err = user.NewDistributionLogic().ReverseOnRefund(ctx, itemId.Int64()); err != nil {
		g.Log().Errorf(ctx, "[佣金冲销] 冲销失败(需人工核对): order_no=%s order_item_id=%d after_sale_no=%s err=%v",
			orderNo, itemId.Int64(), afterSaleNo, err)
	}
}
