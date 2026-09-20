// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserLoginLog is the golang structure for table user_login_log.
type UserLoginLog struct {
	Id           uint64      `json:"id"           orm:"id"            ` // 日志ID
	UserId       uint64      `json:"userId"       orm:"user_id"       ` // 用户ID
	LoginChannel int         `json:"loginChannel" orm:"login_channel" ` // 登录渠道:1微信小程序 2H5
	LoginStatus  int         `json:"loginStatus"  orm:"login_status"  ` // 结果:1成功 2失败
	Ip           string      `json:"ip"           orm:"ip"            ` // 来源IP(IPv6最长45字符)
	UserAgent    string      `json:"userAgent"    orm:"user_agent"    ` // 浏览器/客户端UA
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` // 登录时间
}
