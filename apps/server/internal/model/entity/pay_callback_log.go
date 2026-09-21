// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PayCallbackLog is the golang structure for table pay_callback_log.
type PayCallbackLog struct {
	Id            uint64      `json:"id"            orm:"id"             ` // 日志ID
	PayNo         string      `json:"payNo"         orm:"pay_no"         ` // 关联支付单号
	PayChannel    int         `json:"payChannel"    orm:"pay_channel"    ` // 支付渠道:1微信 2支付宝
	NotifyType    int         `json:"notifyType"    orm:"notify_type"    ` // 通知类型:1支付结果 2退款结果
	VerifyStatus  int         `json:"verifyStatus"  orm:"verify_status"  ` // 验签结果:0失败 1通过
	ProcessStatus int         `json:"processStatus" orm:"process_status" ` // 处理结果:0未处理/重复忽略 1已处理
	RawBody       string      `json:"rawBody"       orm:"raw_body"       ` // 回调原文(全量留档,可重放审计)
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     ` // 创建时间
}
