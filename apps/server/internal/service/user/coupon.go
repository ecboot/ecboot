// coupon.go 用户券持有侧——表: user_coupon（V9）+ 联查 coupon 模板（V9, 模板归 shop 域）。
// 规则: 领券防超发（received_count 条件更新, FR-021 契约）+ per_limit 校验（同事务计数）;
// 核销/退回与订单同事务（下单成功置已使用、取消置退回, 可再用）; 过期由使用时惰性判定+过期任务。
package user

import (
	"context"

	"ecboot/internal/model"
)

// IUserCouponLogic 用户优惠券。
type IUserCouponLogic interface {
	// AvailableTemplates 可领模板列表（过滤已领完/超个人限领/停发）。
	AvailableTemplates(ctx context.Context, userId int64, page model.PageReq) (*model.PageResult[model.AvailableCoupon], error)
	// Receive 领取（防超发+限领校验同事务）。
	Receive(ctx context.Context, userId, couponId int64) (int64, error)
	// Mine 我的券（status: 1未使用 2已使用 3已过期 4已退回; 过期惰性判定）。
	Mine(ctx context.Context, userId int64, status int, page model.PageReq) (*model.PageResult[model.MyCouponItem], error)
	// UsableForOrder 结算可用券匹配（门槛<=商品金额, 按抵扣降序）。
	UsableForOrder(ctx context.Context, userId int64, goodsAmount string) ([]model.UsableCouponItem, error)
	// Consume 核销（下单事务内由 order 域调用: 置已使用+绑单号）。
	Consume(ctx context.Context, userId, userCouponId int64, orderNo string) error
	// ReturnBack 订单取消退回（置可再用, 过期时间不变）。
	ReturnBack(ctx context.Context, userCouponId int64) error
}
