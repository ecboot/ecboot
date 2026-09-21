// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PayOrder is the golang structure of table pay_order for DAO operations like Where/Data.
type PayOrder struct {
	g.Meta         `orm:"table:pay_order, do:true"`
	Id             any         // 支付单ID
	PayNo          any         // 支付单号(全局唯一,传给支付渠道的商户单号)
	OrderId        any         // 订单ID
	OrderNo        any         // 订单号
	UserId         any         // 付款用户ID
	PayChannel     any         // 支付渠道:1微信支付 2支付宝(预留)
	Amount         any         // 应付金额(必须等于订单pay_amount,回调时校验)
	Currency       any         // 币种(ISO 4217,须与订单一致)
	Status         any         // 支付状态:10待支付 20支付成功 30支付失败 90已关闭
	ChannelTradeNo any         // 第三方交易号(微信/支付宝,成功回填;唯一防重放)
	SuccessTime    *gtime.Time // 支付成功时间(毫秒精度,渠道回调为准)
	ExpireTime     *gtime.Time // 支付截止时间(超时未付则关单)
	ClosedTime     *gtime.Time // 关闭时间
	FailReason     any         // 失败原因(渠道返回)
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
}
