// promotion.go 促销域——券模板(coupon V9) / 满减(promotion_* V20) /
// 拼团(group_buy_* V17/V24) / 秒杀(flash_sale_* V18) / 砍价(bargain_* V29) / 助力(assist_* V29)。
// 规则: 领券防超发(条件更新); 满减自动命中最优档+先满减后券; 秒杀活动分账库存(条件更新防超卖);
// 拼团人齐自动成团/超时解散退款; 砍价当前价条件更新防超砍+超时失败; 帮砍/助力挂风控(rule 2/3/4)。
package shop

import "context"

// ICouponLogic 优惠券模板（shop 侧; 用户持有侧在 user 域 IUserCouponLogic）。
type ICouponLogic interface {
	AdminList(ctx context.Context, status int, page PageQuery) (*PageResult[CouponTemplate], error)
	AdminCreate(ctx context.Context, in CouponInput) (int64, error)
	AdminUpdate(ctx context.Context, id int64, in CouponInput) error
	AdminDelete(ctx context.Context, id int64) error
	AdminRecords(ctx context.Context, couponId int64, page PageQuery) (*PageResult[CouponRecordItem], error)
	// PublicList 公开可领列表（过滤停发/领完, 标记 canReceive 由会员态补充）。
	PublicList(ctx context.Context, page PageQuery) (*PageResult[CouponTemplate], error)
}

type CouponTemplate struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Type       int    `json:"type" dc:"1满减 2无门槛"`
	Threshold  string `json:"threshold"`
	Discount   string `json:"discount"`
	TotalCount int    `json:"totalCount"`
	Received   int    `json:"receivedCount"`
	PerLimit   int    `json:"perLimit"`
	ValidDesc  string `json:"validDesc"`
	Status     int    `json:"status"`
}

type CouponInput struct {
	Name         string
	Type         int
	Threshold    string
	Discount     string
	TotalCount   int
	PerLimit     int
	ValidType    int
	ValidStartAt string
	ValidEndAt   string
	ValidDays    int
}

type CouponRecordItem struct {
	UserCouponId int64  `json:"userCouponId"`
	UserId       int64  `json:"userId"`
	Status       int    `json:"status"`
	OrderNo      string `json:"orderNo"`
	CreatedAt    string `json:"createdAt"`
}

// IActivityLogic 满减/拼团/秒杀/砍价/助力活动（管理+浏览; 各活动同构: 活动→场次商品→参与）。
type IActivityLogic interface {
	// ---- 满减 ----
	FullReductionList(ctx context.Context, status int, page PageQuery) (*PageResult[PromotionActivityItem], error)
	FullReductionCreate(ctx context.Context, in PromotionActivityInput) (int64, error) // 档位+范围嵌套, 门槛唯一 50008
	FullReductionDetail(ctx context.Context, id int64) (*PromotionActivityDetail, error)
	FullReductionUpdate(ctx context.Context, id int64, in PromotionActivityInput) error
	FullReductionDelete(ctx context.Context, id int64) error

	// ---- 拼团 ----
	GroupBuyList(ctx context.Context, status int, page PageQuery) (*PageResult[PromotionActivityItem], error)
	GroupBuyCreate(ctx context.Context, in GroupBuyInput) (int64, error)
	GroupBuyDetail(ctx context.Context, id int64) (*PromotionActivityDetail, error)
	GroupBuyUpdate(ctx context.Context, id int64, in GroupBuyInput) error
	GroupBuyDelete(ctx context.Context, id int64) error
	// GroupBuySetItems 场次商品（SKU 级成团价, 全量替换）。
	GroupBuySetItems(ctx context.Context, id int64, items []ActivitySkuInput) error

	// ---- 秒杀 ----
	FlashSaleList(ctx context.Context, status int, page PageQuery) (*PageResult[PromotionActivityItem], error)
	FlashSaleCreate(ctx context.Context, in ActivityTimeInput) (int64, error)
	FlashSaleDetail(ctx context.Context, id int64) (*PromotionActivityDetail, error)
	FlashSaleUpdate(ctx context.Context, id int64, in ActivityTimeInput) error
	FlashSaleDelete(ctx context.Context, id int64) error
	// FlashSaleSetItems 场次商品（秒杀价/限量/限购; 活动分账库存）。
	FlashSaleSetItems(ctx context.Context, id int64, items []ActivitySkuInput) error

	// ---- 砍价 ----
	BargainList(ctx context.Context, status int, page PageQuery) (*PageResult[PromotionActivityItem], error)
	BargainCreate(ctx context.Context, in BargainActivityInput) (int64, error)
	BargainDetail(ctx context.Context, id int64) (*PromotionActivityDetail, error)
	BargainUpdate(ctx context.Context, id int64, in BargainActivityInput) error
	BargainDelete(ctx context.Context, id int64) error
	BargainSetItems(ctx context.Context, id int64, items []ActivitySkuInput) error

	// ---- 助力 ----
	AssistList(ctx context.Context, status int, page PageQuery) (*PageResult[PromotionActivityItem], error)
	AssistCreate(ctx context.Context, in AssistActivityInput) (int64, error)
	AssistDetail(ctx context.Context, id int64) (*PromotionActivityDetail, error)
	AssistUpdate(ctx context.Context, id int64, in AssistActivityInput) error
	AssistDelete(ctx context.Context, id int64) error
}

type PromotionActivityItem struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	SpuId     int64  `json:"spuId" dc:"拼团/砍价有"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Status    int    `json:"status"`
}

type PromotionActivityDetail struct {
	Id         int64                  `json:"id"`
	Name       string                 `json:"name"`
	StartTime  string                 `json:"startTime"`
	EndTime    string                 `json:"endTime"`
	Status     int                    `json:"status"`
	Ladders    []PromotionLadder      `json:"ladders" dc:"满减档位"`
	Scopes     []PromotionScope       `json:"scopes" dc:"满减范围"`
	Items      []PromotionActivitySku `json:"items" dc:"场次商品"`
	GroupSize  int                    `json:"groupSize" dc:"拼团人数"`
	RewardDesc string                 `json:"rewardDesc" dc:"助力奖励"`
}

type PromotionLadder struct {
	Threshold string `json:"threshold"`
	Discount  string `json:"discount"`
}

type PromotionScope struct {
	ScopeType int   `json:"scopeType" dc:"1全场 2分类 3商品"`
	TargetId  int64 `json:"targetId" dc:"全场为0"`
}

type PromotionActivitySku struct {
	SkuId         int64  `json:"skuId"`
	GroupPrice    string `json:"groupPrice" dc:"拼团"`
	FlashPrice    string `json:"flashPrice" dc:"秒杀"`
	StockCount    int    `json:"stockCount" dc:"秒杀限量"`
	PerLimit      int    `json:"perLimit" dc:"秒杀限购"`
	OriginalPrice string `json:"originalPrice" dc:"砍价起始价"`
	FloorPrice    string `json:"floorPrice" dc:"砍价底价"`
	MaxCutCount   int    `json:"maxCutCount" dc:"砍价最大刀数"`
}

type PromotionActivityInput struct {
	Name      string
	StartTime string
	EndTime   string
	Ladders   []PromotionLadder
	Scopes    []PromotionScope
}

type GroupBuyInput struct {
	Name      string
	SpuId     int64
	GroupSize int
	PerLimit  int
	StartTime string
	EndTime   string
}

type ActivityTimeInput struct {
	Name      string
	StartTime string
	EndTime   string
	Status    int
}

type BargainActivityInput struct {
	Name      string
	SpuId     int64
	StartTime string
	EndTime   string
}

type ActivitySkuInput struct {
	SkuId         int64
	GroupPrice    string
	FlashPrice    string
	StockCount    int
	PerLimit      int
	OriginalPrice string
	FloorPrice    string
	MaxCutCount   int
	Config        map[string]any
}

type AssistActivityInput struct {
	Name          string
	RewardType    int
	RewardRef     int64
	PointAmount   int
	RequiredCount int
	PerLimit      int
	StartTime     string
	EndTime       string
}
