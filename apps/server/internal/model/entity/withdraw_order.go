// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// WithdrawOrder is the golang structure for table withdraw_order.
type WithdrawOrder struct {
	Id              uint64      `json:"id"              orm:"id"               ` // 提现单ID
	WithdrawNo      string      `json:"withdrawNo"      orm:"withdraw_no"      ` // 提现单号(全局唯一)
	UserId          uint64      `json:"userId"          orm:"user_id"          ` // 用户ID
	Amount          float64     `json:"amount"          orm:"amount"           ` // 提现金额
	WithdrawChannel int         `json:"withdrawChannel" orm:"withdraw_channel" ` // 渠道:1微信商家转账(V1单渠道)
	ChannelOrderNo  string      `json:"channelOrderNo"  orm:"channel_order_no" ` // 渠道打款单号(回填;与渠道联合唯一,防重复打款)
	Status          int         `json:"status"          orm:"status"           ` // 状态:10待审核 20审核通过 30打款中 40成功 50审核拒绝 60打款失败已回退
	AuditTime       *gtime.Time `json:"auditTime"       orm:"audit_time"       ` // 审核时间
	PayTime         *gtime.Time `json:"payTime"         orm:"pay_time"         ` // 打款成功时间
	FailReason      string      `json:"failReason"      orm:"fail_reason"      ` // 拒绝/失败原因
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       ` // 创建时间
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       ` // 更新时间
}
