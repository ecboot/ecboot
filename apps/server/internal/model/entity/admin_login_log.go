// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminLoginLog is the golang structure for table admin_login_log.
type AdminLoginLog struct {
	Id          uint64      `json:"id"          orm:"id"           ` // 日志ID
	Username    string      `json:"username"    orm:"username"     ` // 尝试登录的用户名(含失败,账号可能不存在)
	AdminId     uint64      `json:"adminId"     orm:"admin_id"     ` // 成功时回填后台账号ID
	LoginStatus int         `json:"loginStatus" orm:"login_status" ` // 结果:1成功 2失败(密码错误) 3失败(账号禁用/不存在)
	Ip          string      `json:"ip"          orm:"ip"           ` // 来源IP(IPv6最长45字符)
	UserAgent   string      `json:"userAgent"   orm:"user_agent"   ` // 浏览器User-Agent
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   ` // 登录时间
}
