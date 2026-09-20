package v1

import "github.com/gogf/gf/v2/frame/g"

// 管理端订单摘要结构（渠道隔离: 不跨渠道引用 shop 包）
type AdminOrderItemBrief struct {
	SpuName  string            `json:"spuName"`
	SkuSpecs map[string]string `json:"skuSpecs"`
	Image    string            `json:"image"`
	Quantity int               `json:"quantity"`
	Price    string            `json:"price"`
}
type AdminOrderAmountBrief struct {
	TotalAmount     string `json:"totalAmount"`
	PromotionAmount string `json:"promotionAmount"`
	AccountAmount   string `json:"accountAmount"`
	FreightAmount   string `json:"freightAmount"`
	PayAmount       string `json:"payAmount"`
}
type AdminOrderStatusLog struct {
	FromStatus int    `json:"fromStatus"`
	ToStatus   int    `json:"toStatus"`
	Remark     string `json:"remark"`
	CreatedAt  string `json:"createdAt"`
}

type (
	// 订单列表（多条件筛选）
	AdminOrderListReq struct {
		g.Meta  `path:"/orders" method:"GET" summary:"订单列表"`
		Status  int    `json:"status" dc:"状态筛选"`
		OrderNo string `json:"orderNo" dc:"订单号"`
		UserKeyword string `json:"userKeyword" dc:"用户ID/手机号(精确)"`
		StartTime  string `json:"startTime" dc:"下单起 RFC3339"`
		EndTime    string `json:"endTime" dc:"下单止"`
		PageReq
	}
	AdminOrderItem struct {
		OrderNo   string           `json:"orderNo"`
		UserId    string           `json:"userId"`
		Status    int              `json:"status"`
		PayAmount string           `json:"payAmount"`
		Items     []AdminOrderItemBrief `json:"items" dc:"商品摘要"`
		CreatedAt string           `json:"createdAt"`
	}
	AdminOrderListRes struct {
		PageRes
		List []AdminOrderItem `json:"list"`
	}

	// 订单详情（含收货快照/状态流水）
	AdminOrderDetailReq struct {
		g.Meta  `path:"/orders/{orderNo}" method:"GET" summary:"订单详情"`
		OrderNo string `json:"orderNo" v:"required" dc:"订单号"`
	}
	AdminOrderDetailRes struct {
		OrderNo      string            `json:"orderNo"`
		UserId       string            `json:"userId"`
		Status       int               `json:"status"`
		RefundStatus int               `json:"refundStatus"`
		Amount       AdminOrderAmountBrief `json:"amount"`
		Items        []AdminOrderItemBrief `json:"items"`
		Receiver     map[string]string `json:"receiver" dc:"收货快照"`
		UserRemark   string            `json:"userRemark" dc:"买家留言"`
		SellerRemark string            `json:"sellerRemark" dc:"卖家备注"`
		Pay          map[string]string `json:"pay" dc:"支付摘要(payNo/channel/paidAt)"`
		Deliver      map[string]string `json:"deliver" dc:"物流信息"`
		StatusLogs   []AdminOrderStatusLog `json:"statusLogs"`
		CreatedAt    string            `json:"createdAt"`
	}

	// 发货
	AdminOrderDeliverReq struct {
		g.Meta   `path:"/orders/{orderNo}/deliver" method:"POST" summary:"订单发货"`
		OrderNo  string `json:"orderNo" v:"required" dc:"订单号"`
		LogisticsCode string `json:"logisticsCode" v:"required" dc:"物流公司编码"`
		DeliverNo     string `json:"deliverNo" v:"required" dc:"运单号"`
	}
	AdminOrderDeliverRes struct {
		Success bool `json:"success"`
	}

	// 管理员取消（待付款）
	AdminOrderCancelReq struct {
		g.Meta  `path:"/orders/{orderNo}/cancel" method:"POST" summary:"管理员取消订单"`
		OrderNo string `json:"orderNo" v:"required" dc:"订单号"`
		Reason  string `json:"reason" v:"required" dc:"取消原因"`
	}
	AdminOrderCancelRes struct {
		Success bool `json:"success"`
	}

	// 卖家备注
	AdminOrderRemarkReq struct {
		g.Meta  `path:"/orders/{orderNo}/seller-remark" method:"PUT" summary:"卖家备注"`
		OrderNo string `json:"orderNo" v:"required" dc:"订单号"`
		SellerRemark string `json:"sellerRemark" v:"required" dc:"备注内容"`
	}
	AdminOrderRemarkRes struct {
		Success bool `json:"success"`
	}
)
