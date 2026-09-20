package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 创建订单（幂等; 玩法上下文二选一见字段）
	OrderCreateReq struct {
		g.Meta        `path:"/orders" method:"POST" summary:"创建订单"`
		RequestToken  string `json:"requestToken" v:"required" dc:"下单幂等凭证"`
		AddressId     string `json:"addressId" v:"required" dc:"收货地址ID"`
		UserCouponId  string `json:"userCouponId" dc:"使用的用户券"`
		UsePoint      bool   `json:"usePoint" dc:"积分抵扣"`
		UseAccount    bool   `json:"useAccount" dc:"佣金余额抵扣"`
		UserRemark    string `json:"userRemark" dc:"买家留言"`
		CartItemIds   []string `json:"cartItemIds" dc:"普通下单:购物车项ID列表"`
		GroupBuyTeamId string `json:"groupBuyTeamId" dc:"拼团下单:团ID"`
		SkuId         string `json:"skuId" dc:"秒杀/直购:SKU ID"`
		Quantity      int    `json:"quantity" dc:"秒杀/直购数量"`
		FlashSaleItemId string `json:"flashSaleItemId" dc:"秒杀下单:场次商品ID"`
		BargainRecordId string `json:"bargainRecordId" dc:"砍价下单:砍价单ID"`
	}
	OrderCreateRes struct {
		OrderNo   string `json:"orderNo" dc:"订单号"`
		PayAmount string `json:"payAmount" dc:"应付(元)"`
	}

	OrderListItem struct {
		OrderNo   string           `json:"orderNo"`
		Status    int              `json:"status" dc:"10待付款 20待发货 30待收货 40已完成 90已取消"`
		Amount    OrderAmountBrief `json:"amount" dc:"金额摘要"`
		Items     []OrderItemBrief `json:"items" dc:"商品摘要"`
		CreatedAt string           `json:"createdAt"`
	}
	OrderAmountBrief struct {
		TotalAmount     string `json:"totalAmount"`
		PromotionAmount string `json:"promotionAmount"`
		FreightAmount   string `json:"freightAmount"`
		AccountAmount   string `json:"accountAmount" dc:"余额抵扣"`
		PayAmount       string `json:"payAmount" dc:"应付"`
	}
	OrderItemBrief struct {
		SpuName  string            `json:"spuName"`
		SkuSpecs map[string]string `json:"skuSpecs"`
		Image    string            `json:"image"`
		Quantity int               `json:"quantity"`
		Price    string            `json:"price"`
	}
	OrderStatusLog struct {
		FromStatus int    `json:"fromStatus" dc:"前状态(0=建单)"`
		ToStatus   int    `json:"toStatus"`
		Remark     string `json:"remark"`
		CreatedAt  string `json:"createdAt"`
	}

	OrderListReq struct {
		g.Meta `path:"/orders" method:"GET" summary:"我的订单列表"`
		Status int `json:"status" dc:"状态筛选"`
		PageReq
	}
	OrderListRes struct {
		PageRes
		List []OrderListItem `json:"list"`
	}

	OrderDetailReq struct {
		g.Meta  `path:"/orders/{orderNo}" method:"GET" summary:"订单详情"`
		OrderNo string `json:"orderNo" v:"required" dc:"订单号"`
	}
	OrderDetailRes struct {
		OrderNo        string           `json:"orderNo"`
		Status         int              `json:"status"`
		RefundStatus   int              `json:"refundStatus" dc:"0无 1部分退款 2全额退款"`
		Amount         OrderAmountBrief `json:"amount"`
		Items          []OrderItemBrief `json:"items"`
		Receiver       map[string]string `json:"receiver" dc:"收货快照"`
		PayTime        string           `json:"payTime" dc:"支付时间"`
		DeliverInfo    map[string]string `json:"deliverInfo" dc:"物流信息(发货后)"`
		CancelInfo     map[string]string `json:"cancelInfo" dc:"取消信息(取消后)"`
		StatusLogs     []OrderStatusLog `json:"statusLogs" dc:"状态时间线"`
		CreatedAt      string           `json:"createdAt"`
	}

	OrderCancelReq struct {
		g.Meta  `path:"/orders/{orderNo}/cancel" method:"POST" summary:"取消订单(仅待付款)"`
		OrderNo string `json:"orderNo" v:"required" dc:"订单号"`
		Reason  string `json:"reason" dc:"取消原因"`
	}
	OrderCancelRes struct {
		Success bool `json:"success"`
	}

	OrderConfirmReq struct {
		g.Meta  `path:"/orders/{orderNo}/confirm" method:"POST" summary:"确认收货(仅待收货)"`
		OrderNo string `json:"orderNo" v:"required" dc:"订单号"`
	}
	OrderConfirmRes struct {
		Success bool `json:"success"`
	}
)
