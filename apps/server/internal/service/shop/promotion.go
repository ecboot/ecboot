// promotion.go 促销域——券模板(coupon V9) / 满减(promotion_* V20) /
// 拼团(group_buy_* V17/V24) / 秒杀(flash_sale_* V18) / 砍价(bargain_* V29) / 助力(assist_* V29)。
// 规则: 领券防超发(条件更新); 满减自动命中最优档+先满减后券; 秒杀活动分账库存(条件更新防超卖);
// 拼团人齐自动成团/超时解散退款; 砍价当前价条件更新防超砍+超时失败; 帮砍/助力挂风控(rule 2/3/4)。
package shop

import (
	"context"

	"ecboot/internal/model"
)

// ICouponLogic 优惠券模板（shop 侧; 用户持有侧在 user 域 IUserCouponLogic）。
type ICouponLogic interface {
	AdminList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.CouponTemplate], error)
	AdminCreate(ctx context.Context, in model.CouponInput) (int64, error)
	AdminUpdate(ctx context.Context, id int64, in model.CouponInput) error
	AdminDelete(ctx context.Context, id int64) error
	AdminRecords(ctx context.Context, couponId int64, page model.PageReq) (*model.PageResult[model.CouponRecordItem], error)
	// PublicList 公开可领列表（过滤停发/领完, 标记 canReceive 由会员态补充）。
	PublicList(ctx context.Context, page model.PageReq) (*model.PageResult[model.CouponTemplate], error)
}

// IActivityLogic 满减/拼团/秒杀/砍价/助力活动（管理+浏览; 各活动同构: 活动→场次商品→参与）。
type IActivityLogic interface {
	// ---- 满减 ----
	FullReductionList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.PromotionActivityItem], error)
	FullReductionCreate(ctx context.Context, in model.PromotionActivityInput) (int64, error) // 档位+范围嵌套, 门槛唯一 50008
	FullReductionDetail(ctx context.Context, id int64) (*model.PromotionActivityDetail, error)
	FullReductionUpdate(ctx context.Context, id int64, in model.PromotionActivityInput) error
	FullReductionDelete(ctx context.Context, id int64) error

	// ---- 拼团 ----
	GroupBuyList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.PromotionActivityItem], error)
	GroupBuyCreate(ctx context.Context, in model.GroupBuyInput) (int64, error)
	GroupBuyDetail(ctx context.Context, id int64) (*model.PromotionActivityDetail, error)
	GroupBuyUpdate(ctx context.Context, id int64, in model.GroupBuyInput) error
	GroupBuyDelete(ctx context.Context, id int64) error
	// GroupBuySetItems 场次商品（SKU 级成团价, 全量替换）。
	GroupBuySetItems(ctx context.Context, id int64, items []model.ActivitySkuInput) error

	// ---- 秒杀 ----
	FlashSaleList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.PromotionActivityItem], error)
	FlashSaleCreate(ctx context.Context, in model.ActivityTimeInput) (int64, error)
	FlashSaleDetail(ctx context.Context, id int64) (*model.PromotionActivityDetail, error)
	FlashSaleUpdate(ctx context.Context, id int64, in model.ActivityTimeInput) error
	FlashSaleDelete(ctx context.Context, id int64) error
	// FlashSaleSetItems 场次商品（秒杀价/限量/限购; 活动分账库存）。
	FlashSaleSetItems(ctx context.Context, id int64, items []model.ActivitySkuInput) error

	// ---- 砍价 ----
	BargainList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.PromotionActivityItem], error)
	BargainCreate(ctx context.Context, in model.BargainActivityInput) (int64, error)
	BargainDetail(ctx context.Context, id int64) (*model.PromotionActivityDetail, error)
	BargainUpdate(ctx context.Context, id int64, in model.BargainActivityInput) error
	BargainDelete(ctx context.Context, id int64) error
	BargainSetItems(ctx context.Context, id int64, items []model.ActivitySkuInput) error

	// ---- 助力 ----
	AssistList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.PromotionActivityItem], error)
	AssistCreate(ctx context.Context, in model.AssistActivityInput) (int64, error)
	AssistDetail(ctx context.Context, id int64) (*model.PromotionActivityDetail, error)
	AssistUpdate(ctx context.Context, id int64, in model.AssistActivityInput) error
	AssistDelete(ctx context.Context, id int64) error
}
