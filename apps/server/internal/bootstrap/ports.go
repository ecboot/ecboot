// ports.go 跨域端口装配（012 评审 I5）。
// 形态: 依赖倒置——shop 域在 service/shop/ports.go 声明接口, 本层（装配层）把 user 域实现注入,
// 两个域兄弟因此互不 import（分层契约: service/user ↔ service/shop 禁止互相引用）。
// 已接: ICouponQuery（结算试算的可用券列表——本批 22 端点中 /shop/cart/checkout 的出参义务）。
// 未接: ICouponTrade/IPointTrade/IAccountTrade——均属下单事务三段联动, 随账户域（批次 11）一并接,
// 见 PROGRESS §五 2026-09-21 批次 06 记账。
package bootstrap

import (
	"context"

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
