// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PayCallbackLog is the golang structure of table pay_callback_log for DAO operations like Where/Data.
type PayCallbackLog struct {
	g.Meta        `orm:"table:pay_callback_log, do:true"`
	Id            any         // 日志ID
	PayNo         any         // 关联支付单号
	PayChannel    any         // 支付渠道:1微信 2支付宝
	NotifyType    any         // 通知类型:1支付结果 2退款结果
	VerifyStatus  any         // 验签结果:0失败 1通过
	ProcessStatus any         // 处理结果:0未处理/重复忽略 1已处理
	RawBody       any         // 回调原文(全量留档,可重放审计)
	CreatedAt     *gtime.Time // 创建时间
}
