// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AfterSaleOrder is the golang structure of table after_sale_order for DAO operations like Where/Data.
type AfterSaleOrder struct {
	g.Meta            `orm:"table:after_sale_order, do:true"`
	Id                any         // 售后单ID
	AfterSaleNo       any         // 售后单号(全局唯一,兼作渠道退款out_refund_no幂等键)
	OrderId           any         // 订单ID
	OrderNo           any         // 订单号
	OrderItemId       any         // 订单项ID(售后粒度=订单项)
	UserId            any         // 申请人用户ID
	Type              any         // 售后类型:1仅退款 2退货退款 3换货(预留未启用)
	Status            any         // 状态:10待审核 20待买家寄回 30待退款 40退款中 50已完成 90已拒绝 91已撤销
	Currency          any         // 币种(ISO 4217,随订单)
	Quantity          any         // 退货/退款数量(≤订单项quantity,应用校验;部分退货的金额与回补依据)
	Reason            any         // 售后原因(不想要了/质量问题/发错货等)
	Description       any         // 问题描述
	VoucherImages     any         // 凭证图片URL数组
	RefundAmount      any         // 退款金额(≤订单项pay_amount)
	ReturnLogisticsNo any         // 买家寄回物流单号(type=2)
	RefundNo          any         // 渠道退款单号(微信退款ID,回填)
	RejectReason      any         // 拒绝原因
	OperatorId        any         // 最后操作人标识:admin:{id}(审核/确认收货/退款重试)
	FailReason        any         // 渠道退款失败原因(重试前保留,成功后清空)
	AuditTime         *gtime.Time // 审核时间
	RefundTime        *gtime.Time // 退款完成时间
	CreatedAt         *gtime.Time // 创建时间
	UpdatedAt         *gtime.Time // 更新时间
}
