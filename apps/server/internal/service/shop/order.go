// order.go 订单域——表: trade_order / trade_order_item / trade_order_log（V6/V17/V23~V25/V28/V29）。
// 规则:
//
//	创建=单一本地事务: 幂等token(UNIQUE兜底) → 锁库存(ADR-0001条件更新) → 锁券/积分/余额(各域核销) →
//	快照(商品/收货/运费/规格) → 优惠三构成分摊(尾差记末行) → 状态流水; 四玩法上下文: 普通/拼团(teamId)/秒杀(itemId)/砍价(recordId)。
//	状态机: 10待付款→20待发货→30待收货→40已完成; 10→90取消(用户/30min超时/管理员)——条件UPDATE, 非法迁移 affected=0。
//	取消副作用: 释放库存+退回券/积分/余额; 支付成功: 核销库存+分销计提事件。
//	快照不可变(ADR-0002); 归因字段(attributed_user_id/attribution_type)下单时判定(V28)。
package shop

import "context"

// IOrderLogic 订单。
type IOrderLogic interface {
	// ---- C 端 ----
	// Create 创建订单（幂等: request_token; 玩法上下文互斥; 返回订单号与应付）。
	Create(ctx context.Context, userId int64, in OrderCreateInput) (*OrderCreated, error)
	List(ctx context.Context, userId int64, status int, page PageQuery) (*PageResult[OrderSummary], error)
	Detail(ctx context.Context, userId int64, orderNo string) (*OrderDetail, error) // 他人资源按不存在
	Cancel(ctx context.Context, userId int64, orderNo string, reason string) error  // 释放库存+退券/积分/余额
	Confirm(ctx context.Context, userId int64, orderNo string) error                // 确认收货→已完成（触发分销结算事件）

	// ---- 后台 ----
	AdminList(ctx context.Context, q AdminOrderQuery) (*PageResult[AdminOrderSummary], error)
	AdminDetail(ctx context.Context, orderNo string) (*OrderDetail, error)
	// Deliver 发货（状态机 20→30; 记物流+流水+通知事件）。
	Deliver(ctx context.Context, orderNo, logisticsCode, deliverNo, operator string) error
	AdminCancel(ctx context.Context, orderNo, reason, operator string) error
	SellerRemark(ctx context.Context, orderNo, remark, operator string) error

	// ---- 定时 ----
	// CancelTimeout 超时取消扫描（30 分钟配置化; 释放库存+退优惠+退余额）。
	CancelTimeout(ctx context.Context) (int, error)
	// AutoConfirm 自动确认收货（发货后 7 天配置化）。
	AutoConfirm(ctx context.Context) (int, error)
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
	Amount    OrderAmountBook  `json:"amount"`
	Items     []OrderItemBrief `json:"items"`
	CreatedAt string           `json:"createdAt"`
}

type OrderAmountBook struct {
	TotalAmount         string `json:"totalAmount"`
	PromotionAmount     string `json:"promotionAmount"`
	CouponAmount        string `json:"couponAmount"`
	FullReductionAmount string `json:"fullReductionAmount"`
	PointAmount         string `json:"pointAmount"`
	PointUsed           int    `json:"pointUsed"`
	AccountAmount       string `json:"accountAmount"`
	FreightAmount       string `json:"freightAmount"`
	PayAmount           string `json:"payAmount"`
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
	Amount          OrderAmountBook   `json:"amount"`
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
	PageQuery
}

type AdminOrderSummary struct {
	OrderNo   string           `json:"orderNo"`
	UserId    string           `json:"userId"`
	Status    int              `json:"status"`
	PayAmount string           `json:"payAmount"`
	Items     []OrderItemBrief `json:"items"`
	CreatedAt string           `json:"createdAt"`
}
