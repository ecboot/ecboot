// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AfterSaleOrder is the golang structure for table after_sale_order.
type AfterSaleOrder struct {
	Id                uint64      `json:"id"                orm:"id"                  ` // 售后单ID
	AfterSaleNo       string      `json:"afterSaleNo"       orm:"after_sale_no"       ` // 售后单号(全局唯一,兼作渠道退款out_refund_no幂等键)
	OrderId           uint64      `json:"orderId"           orm:"order_id"            ` // 订单ID
	OrderNo           string      `json:"orderNo"           orm:"order_no"            ` // 订单号
	OrderItemId       uint64      `json:"orderItemId"       orm:"order_item_id"       ` // 订单项ID(售后粒度=订单项)
	UserId            uint64      `json:"userId"            orm:"user_id"             ` // 申请人用户ID
	SellerId          uint64      `json:"sellerId"          orm:"seller_id"           ` // 售后处理商家(1=自营;商家后台按此筛选)
	Type              int         `json:"type"              orm:"type"                ` // 售后类型:1仅退款 2退货退款 3换货(预留未启用)
	Status            int         `json:"status"            orm:"status"              ` // 状态:10待审核 20待买家寄回 30待退款 40退款中 50已完成 90已拒绝 91已撤销
	Currency          string      `json:"currency"          orm:"currency"            ` // 币种(ISO 4217,随订单)
	Quantity          uint        `json:"quantity"          orm:"quantity"            ` // 退货/退款数量(≤订单项quantity,应用校验;部分退货的金额与回补依据)
	Reason            string      `json:"reason"            orm:"reason"              ` // 售后原因(不想要了/质量问题/发错货等)
	Description       string      `json:"description"       orm:"description"         ` // 问题描述
	VoucherImages     string      `json:"voucherImages"     orm:"voucher_images"      ` // 凭证图片URL数组
	RefundAmount      float64     `json:"refundAmount"      orm:"refund_amount"       ` // 退款金额(≤订单项pay_amount)
	ReturnLogisticsNo string      `json:"returnLogisticsNo" orm:"return_logistics_no" ` // 买家寄回物流单号(type=2)
	RefundNo          string      `json:"refundNo"          orm:"refund_no"           ` // 渠道退款单号(微信退款ID,回填)
	RejectReason      string      `json:"rejectReason"      orm:"reject_reason"       ` // 拒绝原因
	AuditTime         *gtime.Time `json:"auditTime"         orm:"audit_time"          ` // 审核时间
	RefundTime        *gtime.Time `json:"refundTime"        orm:"refund_time"         ` // 退款完成时间
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"          ` // 创建时间
	UpdatedAt         *gtime.Time `json:"updatedAt"         orm:"updated_at"          ` // 更新时间
}
