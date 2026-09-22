// dto_shop.go shop 域 DTO（自 internal/service/shop 迁入 2026-09-21, 契约出入参统一归属 model）。
package model

type AfterSaleApplyInput struct {
	OrderItemId   int64
	Type          int // 1仅退款 2退货退款
	Quantity      int
	Reason        string
	Description   string
	VoucherImages []string
}

type AfterSaleSummary struct {
	AfterSaleNo  string `json:"afterSaleNo"`
	OrderNo      string `json:"orderNo"`
	UserId       int64  `json:"userId" dc:"申请人（后台展示用; C 端不映射）"`
	Type         int    `json:"type"`
	Quantity     int    `json:"quantity"`
	RefundAmount string `json:"refundAmount"`
	Status       int    `json:"status"`
	CreatedAt    string `json:"createdAt"`
}

type AfterSaleDetail struct {
	AfterSaleNo       string   `json:"afterSaleNo"`
	OrderNo           string   `json:"orderNo"`
	OrderItemId       int64    `json:"orderItemId"`
	UserId            int64    `json:"userId" dc:"申请人（后台展示用; C 端不映射）"`
	Type              int      `json:"type"`
	Quantity          int      `json:"quantity"`
	Reason            string   `json:"reason"`
	Description       string   `json:"description"`
	VoucherImages     []string `json:"voucherImages"`
	RefundAmount      string   `json:"refundAmount"`
	ReturnLogisticsNo string   `json:"returnLogisticsNo"`
	RefundNo          string   `json:"refundNo" dc:"渠道退款单号"`
	RejectReason      string   `json:"rejectReason"`
	AuditTime         string   `json:"auditTime"`
	RefundTime        string   `json:"refundTime"`
	OperatorId        string   `json:"operatorId" dc:"最后操作人（后台可见; C 端不映射）"`
	FailReason        string   `json:"failReason" dc:"渠道退款失败原因（后台可见; 可重试时的排障依据）"`
	Status            int      `json:"status"`
}

type CartView struct {
	Items []CartLine `json:"items"`
}

type CartLine struct {
	ItemId   int64             `json:"itemId"`
	SkuId    int64             `json:"skuId"`
	SpuName  string            `json:"spuName"`
	Specs    map[string]string `json:"specs"`
	Image    string            `json:"image"`
	Price    string            `json:"price"`
	Sellable bool              `json:"sellable"`
	Quantity int               `json:"quantity"`
	Checked  bool              `json:"checked"`
}

type CheckoutQuery struct {
	AddressId  int64
	CouponId   int64
	UsePoint   bool
	UseAccount bool
}

type CheckoutResult struct {
	Items         []CartLine          `json:"items"`
	Amount        AmountBook          `json:"amount"`
	UsableCoupons []UsableCouponBrief `json:"usableCoupons"`
	Errors        []string            `json:"errors" dc:"失效商品提示"`
}

type AmountBook struct {
	TotalAmount         string `json:"totalAmount"`
	CouponAmount        string `json:"couponAmount"`
	FullReductionAmount string `json:"fullReductionAmount"`
	PointAmount         string `json:"pointAmount"`
	PointUsed           int    `json:"pointUsed"`
	AccountAmount       string `json:"accountAmount" dc:"余额抵扣(现金应付=payAmount-accountAmount)"`
	FreightAmount       string `json:"freightAmount"`
	PayAmount           string `json:"payAmount" dc:"实付(含余额抵扣前口径=total-promotion+freight)"`
}

type UsableCouponBrief struct {
	UserCouponId int64  `json:"userCouponId"`
	Name         string `json:"name"`
	Discount     string `json:"discount"`
}

type FreightTemplateItem struct {
	Id               int64   `json:"id"`
	Name             string  `json:"name"`
	ChargeType       int     `json:"chargeType" dc:"1按件 2按重"`
	FreeThreshold    *string `json:"freeThreshold" dc:"满额包邮,空=不包邮"`
	FreeExcludeCodes string  `json:"freeExcludeCodes" dc:"不参与包邮的省级代码(JSON)"`
	Status           int     `json:"status"`
}

type FreightTemplateInput struct {
	Name             string
	ChargeType       int
	FreeThreshold    *string
	FreeExcludeCodes []string
	Rules            []FreightRuleInput
}

type FreightRuleInput struct {
	RegionCodes  []string `json:"regionCodes" dc:"省级代码; 空=全国兜底"`
	FirstUnit    int
	FirstFee     string
	ContinueUnit int
	ContinueFee  string
}

type FreightCalcItem struct {
	SkuId    int64
	Weight   string // 克（按重计费用）
	Quantity int
}

type InventoryItem struct {
	SkuId     int64  `json:"skuId"`
	SkuNo     string `json:"skuNo"`
	SkuName   string `json:"skuName"`
	Total     int    `json:"total"`
	Locked    int    `json:"locked"`
	Available int    `json:"available" dc:"total-locked"`
	WarnCount int    `json:"warnCount"`
}

type LogisticsCompany struct {
	Id           int64  `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	TrackingRule string `json:"trackingRule"`
	Status       int    `json:"status"`
}

type LogisticsCompanyInput struct {
	Code         string
	Name         string
	TrackingRule string
	Status       int
}

type StoreQuery struct {
	DistrictCode string
	Longitude    float64
	Latitude     float64
	RadiusKm     int
	PageReq
}

type StoreItem struct {
	Id            int64   `json:"id"`
	StoreNo       string  `json:"storeNo"`
	Name          string  `json:"name"`
	ProvinceCode  string  `json:"provinceCode"`
	CityCode      string  `json:"cityCode"`
	DistrictCode  string  `json:"districtCode"`
	DetailAddress string  `json:"detailAddress"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
	BusinessHours string  `json:"businessHours"`
	ContactPhone  string  `json:"contactPhone"`
	PickupEnabled bool    `json:"pickupEnabled"`
	Status        int     `json:"status"`
	DistanceM     int64   `json:"distanceM" dc:"附近检索返回"`
}

type StoreInput struct {
	Name          string
	ProvinceCode  string
	CityCode      string
	DistrictCode  string
	DetailAddress string
	Longitude     float64
	Latitude      float64
	BusinessHours string
	ContactPhone  string
	PickupEnabled bool
	Status        int
}

type RiskRuleItem struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	RuleType      int    `json:"ruleType" dc:"1黑名单 2高频下单 3异常领券 4佣金套利 5休眠分级"`
	ConditionExpr string `json:"conditionExpr"`
	Action        int    `json:"action" dc:"1拦截 2标记"`
	Status        int    `json:"status"`
}

type RiskRuleInput struct {
	Name          string
	RuleType      int
	ConditionExpr string
	Action        int
	Status        int
}

type RiskRecordItem struct {
	Id           int64  `json:"id"`
	UserId       int64  `json:"userId"`
	RuleName     string `json:"ruleName"`
	ObjectType   int    `json:"objectType" dc:"1订单 2券 3提现 4售后"`
	ObjectNo     string `json:"objectNo"`
	Action       int    `json:"action"`
	AppealStatus int    `json:"appealStatus" dc:"0无 1申诉中 2通过 3驳回"`
	CreatedAt    string `json:"createdAt"`
}

type OrderCreateInput struct {
	RequestToken string
	AddressId    int64
	UserCouponId int64
	UsePoint     bool
	UseAccount   bool
	UserRemark   string
	Channel      int
	// 玩法上下文（互斥）:
	CartItemIds     []int64 // 普通
	GroupBuyTeamId  int64   // 拼团
	FlashSaleItemId int64   // 秒杀
	BargainRecordId int64   // 砍价成交
	SkuId           int64   // 秒杀/直购 SKU
	Quantity        int
}

type OrderCreated struct {
	OrderNo   string `json:"orderNo"`
	PayAmount string `json:"payAmount"`
}

type OrderSummary struct {
	OrderNo   string           `json:"orderNo"`
	Status    int              `json:"status"`
	Amount    AmountBook       `json:"amount"`
	Items     []OrderItemBrief `json:"items"`
	CreatedAt string           `json:"createdAt"`
}

type OrderItemBrief struct {
	SpuName  string            `json:"spuName"`
	SkuSpecs map[string]string `json:"skuSpecs"`
	Image    string            `json:"image"`
	Quantity int               `json:"quantity"`
	Price    string            `json:"price"`
}

type OrderDetail struct {
	OrderNo         string            `json:"orderNo"`
	Status          int               `json:"status"`
	RefundStatus    int               `json:"refundStatus"`
	Amount          AmountBook        `json:"amount"`
	Items           []OrderItemDetail `json:"items"`
	Receiver        map[string]string `json:"receiver" dc:"收货快照"`
	UserRemark      string            `json:"userRemark"`
	SellerRemark    string            `json:"sellerRemark"`
	Pay             map[string]string `json:"pay"`
	Deliver         map[string]string `json:"deliver"`
	Cancel          map[string]string `json:"cancel"`
	StatusLogs      []OrderStatusLog  `json:"statusLogs"`
	GroupBuyTeamId  int64             `json:"groupBuyTeamId"`
	BargainRecordId int64             `json:"bargainRecordId"`
	CreatedAt       string            `json:"createdAt"`
}

type OrderItemDetail struct {
	Id                  int64             `json:"id"`
	SpuId               int64             `json:"spuId"`
	SkuId               int64             `json:"skuId"`
	SkuNo               string            `json:"skuNo"`
	SpuName             string            `json:"spuName" dc:"快照"`
	SkuName             string            `json:"skuName" dc:"快照"`
	SkuImage            string            `json:"skuImage" dc:"快照"`
	SkuSpecs            map[string]string `json:"skuSpecs" dc:"快照"`
	Quantity            int               `json:"quantity"`
	OriginalPrice       string            `json:"originalPrice" dc:"快照"`
	Price               string            `json:"price" dc:"成交价快照"`
	PromotionAmount     string            `json:"promotionAmount"`
	CouponAmount        string            `json:"couponAmount"`
	FullReductionAmount string            `json:"fullReductionAmount"`
	PointAmount         string            `json:"pointAmount"`
	AccountAmount       string            `json:"accountAmount"`
	PayAmount           string            `json:"payAmount"`
	CanReview           bool              `json:"canReview"`
}

type OrderStatusLog struct {
	FromStatus int    `json:"fromStatus"`
	ToStatus   int    `json:"toStatus"`
	Remark     string `json:"remark"`
	CreatedAt  string `json:"createdAt"`
}

type AdminOrderQuery struct {
	Status      int
	OrderNo     string
	UserKeyword string
	StartTime   string
	EndTime     string
	PageReq
}

type AdminOrderSummary struct {
	OrderNo   string           `json:"orderNo"`
	UserId    string           `json:"userId"`
	Status    int              `json:"status"`
	PayAmount string           `json:"payAmount"`
	Items     []OrderItemBrief `json:"items"`
	CreatedAt string           `json:"createdAt"`
}

type PayCreated struct {
	PayNo         string         `json:"payNo"`
	ChannelParams map[string]any `json:"channelParams" dc:"渠道唤起参数(mock 为占位)"`
}

type PayStatus struct {
	Status int    `json:"status" dc:"10待支付 20成功 30失败 90关闭"`
	PaidAt string `json:"paidAt"`
}

type CategoryNode struct {
	Id       int64          `json:"id"`
	Name     string         `json:"name"`
	Icon     string         `json:"icon"`
	Children []CategoryNode `json:"children"`
}

type CategoryInput struct {
	ParentId int64
	Name     string
	Icon     string
	Level    int
	Sort     int
	Status   int // 修改时生效（创建忽略——连线适配 010-product-admin）
}

type BrandItem struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	Description string `json:"description" dc:"简介(管理面)"`
	Sort        int    `json:"sort" dc:"排序(管理面)"`
	Status      int    `json:"status" dc:"状态(管理面)"`
}

type BrandInput struct {
	Name        string
	Logo        string
	Description string
	Sort        int
	Status      int
}

type ProductQuery struct {
	CategoryId int64
	BrandId    int64
	Sort       int // 0综合 1销量 2价格 3上新
	PriceMin   string
	PriceMax   string
	PageReq
}

type ProductCard struct {
	SpuId      int64  `json:"spuId"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	PriceRange string `json:"priceRange"`
	SaleCount  int    `json:"saleCount"`
}

type ProductDetail struct {
	SpuId           int64             `json:"spuId"`
	SpuNo           string            `json:"spuNo"`
	Name            string            `json:"name"`
	SubTitle        string            `json:"subTitle"`
	Images          []string          `json:"images"`
	VideoUrl        string            `json:"videoUrl"`
	Description     string            `json:"description"`
	SpecDefinitions []map[string]any  `json:"specDefinitions"`
	Attributes      map[string]string `json:"attributes"`
	FreightSummary  string            `json:"freightSummary"`
	Skus            []SkuCard         `json:"skus"`
	ReviewSummary   ReviewSummary     `json:"reviewSummary"`
}

type SkuCard struct {
	SkuId     int64             `json:"skuId"`
	SkuNo     string            `json:"skuNo"`
	Specs     map[string]string `json:"specs"`
	Price     string            `json:"price"`
	LinePrice string            `json:"linePrice"`
	Sellable  bool              `json:"sellable"`
}

type ReviewSummary struct {
	Avg          string         `json:"avg"`
	Distribution map[string]int `json:"distribution"`
	Total        int            `json:"total"`
}

type ReviewCard struct {
	ReviewId  int64             `json:"reviewId"`
	User      string            `json:"user" dc:"匿名脱敏"`
	Score     int               `json:"score"`
	Content   string            `json:"content"`
	Images    []string          `json:"images"`
	Specs     map[string]string `json:"specs"`
	Reply     string            `json:"reply"`
	Extra     string            `json:"extra"`
	CreatedAt string            `json:"createdAt"`
}

type AdminProductQuery struct {
	Status     int
	CategoryId int64
	Keyword    string
	PageReq
}

type AdminProductItem struct {
	SpuId      int64  `json:"spuId"`
	SpuNo      string `json:"spuNo"`
	Name       string `json:"name"`
	CategoryId int64  `json:"categoryId"`
	BrandId    int64  `json:"brandId" dc:"品牌(连线适配 010)"`
	Status     int    `json:"status"`
	SaleCount  int    `json:"saleCount"`
	CreatedAt  string `json:"createdAt"`
}

type SpuInput struct {
	Name              string
	SubTitle          string
	CategoryId        int64
	BrandId           int64
	FreightTemplateId int64
	Images            []string
	VideoUrl          string
	Description       string
	SpecDefinitions   []map[string]any
	Attributes        map[string]string
}

type AdminProductDetailView struct {
	SpuId             int64             `json:"spuId"`
	SpuNo             string            `json:"spuNo"`
	Name              string            `json:"name"`
	SubTitle          string            `json:"subTitle"`
	VideoUrl          string            `json:"videoUrl"`
	CategoryId        int64             `json:"categoryId"`
	BrandId           int64             `json:"brandId"`
	FreightTemplateId int64             `json:"freightTemplateId"`
	Images            []string          `json:"images"`
	Description       string            `json:"description"`
	SpecDefinitions   []map[string]any  `json:"specDefinitions"`
	Attributes        map[string]string `json:"attributes"`
	SaleRestrictCodes []string          `json:"saleRestrictCodes"`
	Status            int               `json:"status"`
	Skus              []AdminSkuDetail  `json:"skus"`
}

type AdminSkuDetail struct {
	SkuId     int64             `json:"skuId"`
	SkuNo     string            `json:"skuNo"`
	Specs     map[string]string `json:"specs"`
	Price     string            `json:"price"`
	LinePrice string            `json:"linePrice"`
	CostPrice string            `json:"costPrice" dc:"管理可见"`
	Weight    string            `json:"weight" dc:"重量克(连线适配 010)"`
	Barcode   string            `json:"barcode" dc:"条码(连线适配 010)"`
	Status    int               `json:"status"`
}

type SkuInput struct {
	Specs     map[string]string
	Price     string
	LinePrice string
	CostPrice string
	Image     string
	Weight    string
	Barcode   string
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
	ValidType  int    `json:"validType" dc:"1固定区间 2领取后N天(016 契约微扩 D3-②)"`
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
	Status       int `dc:"启停(仅 AdminUpdate 消费; 016 契约微扩 D3-③: api Update Req 有 status 而 DTO 漏)"`
}

type CouponRecordItem struct {
	UserCouponId int64  `json:"userCouponId"`
	UserId       int64  `json:"userId"`
	Status       int    `json:"status"`
	OrderNo      string `json:"orderNo"`
	CreatedAt    string `json:"createdAt"`
}

type PromotionActivityItem struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	SpuId     int64  `json:"spuId" dc:"拼团/砍价有"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Status    int    `json:"status"`
	// 016 契约微扩 D3-⑧: 各玩法列表项出参要求（拼团 groupSize/perLimit; 助力 rewardType/
	// requiredCount/perLimit/rewardDesc）——缺列的玩法类型读行为零值, 不影响既有消费
	GroupSize     int    `json:"groupSize" dc:"拼团成团人数"`
	PerLimit      int    `json:"perLimit" dc:"拼团限购/助力每人可发起"`
	RewardType    int    `json:"rewardType" dc:"助力:1券 2积分"`
	RequiredCount int    `json:"requiredCount" dc:"助力所需人数"`
	RewardDesc    string `json:"rewardDesc" dc:"助力奖励说明"`
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
	Status    int `dc:"启停(仅 Update 消费; 016 契约微扩 D3-④: api Update Req 有 status 而 DTO 漏)"`
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
	Status    int `dc:"启停(仅 Update 消费; 016 契约微扩 D3-⑤)"`
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
	Status    int `dc:"启停(仅 Update 消费; 016 契约微扩 D3-⑥)"`
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
	Status        int `dc:"启停(仅 Update 消费; 016 契约微扩 D3-⑦)"`
}

type ReviewCreateInput struct {
	OrderItemId int64
	Score       int
	Content     string
	Images      []string
	IsAnonymous bool
}

type MyReviewItem struct {
	ReviewId    int64  `json:"reviewId"`
	SpuName     string `json:"spuName"`
	Score       int    `json:"score"`
	Content     string `json:"content"`
	Extra       string `json:"extra"`
	Reply       string `json:"reply"`
	AuditStatus int    `json:"auditStatus" dc:"0待审 1通过 2驳回"`
	CreatedAt   string `json:"createdAt"`
}

type AdminReviewItem struct {
	ReviewId    int64  `json:"reviewId"`
	OrderNo     string `json:"orderNo"`
	UserMasked  string `json:"user" dc:"脱敏"`
	SpuName     string `json:"spuName" dc:"快照"`
	Score       int    `json:"score"`
	Content     string `json:"content"`
	AuditStatus int    `json:"auditStatus"`
	CreatedAt   string `json:"createdAt"`
}

// AdminCategoryNode 管理端分类树节点（含禁用/ParentId）。
type AdminCategoryNode struct {
	Id       int64               `json:"id"`
	ParentId int64               `json:"parentId"`
	Name     string              `json:"name"`
	Icon     string              `json:"icon"`
	Level    int                 `json:"level"`
	Sort     int                 `json:"sort"`
	Status   int                 `json:"status"`
	Children []AdminCategoryNode `json:"children"`
}

// ProductDetailView C 端商品详情（FR-004; 不含成本价; 含限售标记）。
type ProductDetailView struct {
	SpuId           int64             `json:"spuId"`
	SpuNo           string            `json:"spuNo"`
	Name            string            `json:"name"`
	SubTitle        string            `json:"subTitle"`
	Images          []string          `json:"images"`
	VideoUrl        string            `json:"videoUrl"`
	Description     string            `json:"description"`
	SpecDefinitions []map[string]any  `json:"specDefinitions"`
	Attributes      map[string]string `json:"attributes"`
	FreightSummary  string            `json:"freightSummary"`
	SaleRestricted  bool              `json:"saleRestricted"`
	Skus            []SkuCard         `json:"skus"`
	ReviewSummary   ReviewSummary     `json:"reviewSummary"`
}

// OrderDetailView 订单详情（含金额账本/商品行/状态时间线）。
type OrderDetailView struct {
	OrderNo      string            `json:"orderNo"`
	Status       int               `json:"status"`
	RefundStatus int               `json:"refundStatus"`
	Amount       AmountBook        `json:"amount"`
	Items        []OrderItemBrief  `json:"items"`
	Receiver     map[string]string `json:"receiver" dc:"收货快照"`
	UserRemark   string            `json:"userRemark"`
	SellerRemark string            `json:"sellerRemark"`
	Pay          map[string]string `json:"pay" dc:"支付摘要"`
	Deliver      map[string]string `json:"deliver" dc:"物流信息"`
	Cancel       map[string]string `json:"cancel" dc:"取消信息"`
	StatusLogs   []OrderStatusLog  `json:"statusLogs"`
	CreatedAt    string            `json:"createdAt"`
}

// ---------- 运营装修（009-logistics-ops, research D1） ----------

// OperBannerItem 轮播/弹窗管理项（admin）。
type OperBannerItem struct {
	Id        int64  `json:"id"`
	Position  int    `json:"position" dc:"1首页轮播 2首页弹窗"`
	ImageUrl  string `json:"imageUrl"`
	LinkUrl   string `json:"linkUrl"`
	Sort      int    `json:"sort"`
	StartTime string `json:"startTime" dc:"投放起(RFC3339; 空=立即)"`
	EndTime   string `json:"endTime" dc:"投放止(RFC3339; 空=长期)"`
	Status    int    `json:"status"`
}

// OperBannerInput 轮播创建/修改入参（全量覆盖；空时段=立即/长期）。
type OperBannerInput struct {
	Position  int
	ImageUrl  string
	LinkUrl   string
	Sort      int
	StartTime string // RFC3339 或空
	EndTime   string
	Status    int
}

// OperFloorItem 楼层管理项（admin）。
type OperFloorItem struct {
	Id        int64          `json:"id"`
	FloorType int            `json:"floorType" dc:"1金刚区 2商品楼层 3专题"`
	Title     string         `json:"title"`
	Config    map[string]any `json:"config"`
	Sort      int            `json:"sort"`
	Status    int            `json:"status"`
}

// OperFloorInput 楼层创建/修改入参（config 为不透明 JSON 对象）。
type OperFloorInput struct {
	FloorType int
	Title     string
	Config    map[string]any
	Sort      int
	Status    int
}

// PublicBannerItem C 端轮播项（精简字段）。
type PublicBannerItem struct {
	Id       int64  `json:"id"`
	ImageUrl string `json:"imageUrl"`
	LinkUrl  string `json:"linkUrl"`
}

// FloorProductSummary 商品楼层装配的商品摘要（research D3）。
type FloorProductSummary struct {
	SpuId int64  `json:"spuId"`
	Name  string `json:"name"`
	Image string `json:"image" dc:"首图"`
	Price string `json:"price" dc:"价格(元, 取 price_min 与商品列表口径一致)"`
}

// PublicFloorItem C 端楼层项（商品楼层含装配摘要）。
type PublicFloorItem struct {
	FloorId   int64                 `json:"floorId"`
	FloorType int                   `json:"floorType"`
	Title     string                `json:"title"`
	Config    map[string]any        `json:"config"`
	Products  []FloorProductSummary `json:"products" dc:"仅商品楼层有值"`
}

// ---------- 营销 C 端 DTO（015-marketing-c 批次 09 新建） ----------
// 说明: C 端营销在批次 09 之前**无任何接口与 DTO**（IActivityLogic 全为管理面）,
// 故按批次 03（IOperationLogic）先例在本批新建; 管理面 DTO 一律不动。

// ActivitySkuBrief 场次商品摘要（五类玩法共用: 秒杀价/成团价/起始价等）。
type ActivitySkuBrief struct {
	ItemId        int64  `json:"itemId" dc:"场次商品ID(砍价发起必传; 其余玩法为 0)"`
	SkuId         int64  `json:"skuId"`
	Price         string `json:"price" dc:"活动价(元; 秒杀价/成团价/砍价起始价)"`
	OriginalPrice string `json:"originalPrice" dc:"原价(元; 砍价用)"`
	FloorPrice    string `json:"floorPrice" dc:"底价(元; 砍价用)"`
	MaxCutCount   int    `json:"maxCutCount" dc:"最大刀数(砍价用)"`
	StockRemain   int    `json:"stockRemain" dc:"活动剩余量(秒杀用)"`
	PerLimit      int    `json:"perLimit" dc:"每人限购(秒杀用)"`
}

// PublicGroupBuyItem 拼团活动（公开列表项）。
type PublicGroupBuyItem struct {
	ActivityId int64              `json:"activityId"`
	Name       string             `json:"name"`
	SpuId      int64              `json:"spuId"`
	SpuName    string             `json:"spuName"`
	Image      string             `json:"image"`
	GroupSize  int                `json:"groupSize" dc:"成团人数"`
	Items      []ActivitySkuBrief `json:"items"`
	EndTime    string             `json:"endTime"`
}

// PublicFlashSaleItem 秒杀场次（公开列表项；含预告）。
type PublicFlashSaleItem struct {
	ActivityId int64              `json:"activityId"`
	Name       string             `json:"name"`
	StartTime  string             `json:"startTime"`
	EndTime    string             `json:"endTime"`
	Upcoming   bool               `json:"upcoming" dc:"是否预告(未开始)"`
	Items      []ActivitySkuBrief `json:"items"`
}

// PublicBargainItem 砍价活动（公开列表项）。
type PublicBargainItem struct {
	ActivityId int64              `json:"activityId"`
	Name       string             `json:"name"`
	SpuId      int64              `json:"spuId"`
	SpuName    string             `json:"spuName"`
	Image      string             `json:"image"`
	Items      []ActivitySkuBrief `json:"items"`
	EndTime    string             `json:"endTime"`
}

// PublicAssistItem 助力活动（公开列表项）。
type PublicAssistItem struct {
	ActivityId    int64  `json:"activityId"`
	Name          string `json:"name"`
	RequiredCount int    `json:"requiredCount" dc:"所需助力人数"`
	PerLimit      int    `json:"perLimit" dc:"每人可发起次数"`
	RewardType    int    `json:"rewardType" dc:"奖励类型:1优惠券 2积分"`
	EndTime       string `json:"endTime"`
}

// LadderBrief 满减档位摘要。
type LadderBrief struct {
	Threshold string `json:"threshold" dc:"满(元)"`
	Discount  string `json:"discount" dc:"减(元)"`
}

// PublicFullReductionItem 满减活动（公开列表项）。
type PublicFullReductionItem struct {
	ActivityId int64         `json:"activityId"`
	Name       string        `json:"name"`
	Ladders    []LadderBrief `json:"ladders"`
	ScopeDesc  string        `json:"scopeDesc" dc:"适用范围摘要(全场/分类/商品)"`
	EndTime    string        `json:"endTime"`
}

// PlayHelperItem 帮砍/助力条目（昵称脱敏）。
type PlayHelperItem struct {
	UserId    int64  `json:"userId"`
	Nickname  string `json:"nickname" dc:"脱敏"`
	Amount    string `json:"amount" dc:"本刀金额(砍价用,元)"`
	CreatedAt string `json:"createdAt"`
}

// BargainLaunchResult 发起砍价出参。
type BargainLaunchResult struct {
	RecordId     int64  `json:"recordId"`
	CurrentPrice string `json:"currentPrice" dc:"发起即首刀后的当前价(元)"`
}

// BargainProgressView 砍价进度。
type BargainProgressView struct {
	RecordId      int64            `json:"recordId"`
	SkuId         int64            `json:"skuId"`
	OriginalPrice string           `json:"originalPrice" dc:"起始价(元)"`
	CurrentPrice  string           `json:"currentPrice" dc:"当前价(元)"`
	FloorPrice    string           `json:"floorPrice" dc:"底价(元)"`
	CutCount      int              `json:"cutCount" dc:"已砍刀数"`
	Status        int              `json:"status" dc:"1砍价中 2到底价 3已下单 4超时 5取消"`
	ExpireTime    string           `json:"expireTime"`
	OrderNo       string           `json:"orderNo"`
	Helpers       []PlayHelperItem `json:"helpers"`
}

// BargainCutResult 帮砍出参。
type BargainCutResult struct {
	CutAmount    string `json:"cutAmount"`
	CurrentPrice string `json:"currentPrice"`
	FloorReached bool   `json:"floorReached" dc:"是否已到底价(可下单)"`
}

// AssistLaunchResult 发起助力出参。
type AssistLaunchResult struct {
	RecordId int64 `json:"recordId"`
}

// AssistProgressView 助力进度。
type AssistProgressView struct {
	RecordId      int64            `json:"recordId"`
	ActivityId    int64            `json:"activityId"`
	HelperCount   int              `json:"helperCount"`
	RequiredCount int              `json:"requiredCount"`
	Status        int              `json:"status" dc:"1进行中 2已完成发奖 3已过期"`
	FinishTime    string           `json:"finishTime"`
	Helpers       []PlayHelperItem `json:"helpers"`
}

// AssistHelpResult 助力出参。
type AssistHelpResult struct {
	Done bool `json:"done" dc:"本次助力后是否达成发奖"`
}

// IndexAggregate 首页聚合（轮播 + 楼层 + 五类活动入口 + 可领券; 任一块可为空数组）。
type IndexAggregate struct {
	Banners        []PublicBannerItem        `json:"banners"`
	Floors         []PublicFloorItem         `json:"floors"`
	FlashSales     []PublicFlashSaleItem     `json:"flashSales"`
	GroupBuys      []PublicGroupBuyItem      `json:"groupBuys"`
	Bargains       []PublicBargainItem       `json:"bargains"`
	Assists        []PublicAssistItem        `json:"assists"`
	FullReductions []PublicFullReductionItem `json:"fullReductions"`
	Coupons        []CouponTemplate          `json:"coupons" dc:"可领券"`
}
