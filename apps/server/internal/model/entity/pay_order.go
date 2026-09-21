// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PayOrder is the golang structure for table pay_order.
type PayOrder struct {
	Id             uint64      `json:"id"             orm:"id"               ` // 支付单ID
	PayNo          string      `json:"payNo"          orm:"pay_no"           ` // 支付单号(全局唯一,传给支付渠道的商户单号)
	OrderId        uint64      `json:"orderId"        orm:"order_id"         ` // 订单ID
	OrderNo        string      `json:"orderNo"        orm:"order_no"         ` // 订单号
	UserId         uint64      `json:"userId"         orm:"user_id"          ` // 付款用户ID
	PayChannel     int         `json:"payChannel"     orm:"pay_channel"      ` // 支付渠道:1微信支付 2支付宝(预留)
	Amount         float64     `json:"amount"         orm:"amount"           ` // 应付金额(必须等于订单pay_amount,回调时校验)
	Currency       string      `json:"currency"       orm:"currency"         ` // 币种(ISO 4217,须与订单一致)
	Status         int         `json:"status"         orm:"status"           ` // 支付状态:10待支付 20支付成功 30支付失败 90已关闭
	ChannelTradeNo string      `json:"channelTradeNo" orm:"channel_trade_no" ` // 第三方交易号(微信/支付宝,成功回填;唯一防重放)
	SuccessTime    *gtime.Time `json:"successTime"    orm:"success_time"     ` // 支付成功时间(毫秒精度,渠道回调为准)
	ExpireTime     *gtime.Time `json:"expireTime"     orm:"expire_time"      ` // 支付截止时间(超时未付则关单)
	ClosedTime     *gtime.Time `json:"closedTime"     orm:"closed_time"      ` // 关闭时间
	FailReason     string      `json:"failReason"     orm:"fail_reason"      ` // 失败原因(渠道返回)
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       ` // 创建时间
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"       ` // 更新时间
}
