// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminLoginLog is the golang structure of table admin_login_log for DAO operations like Where/Data.
type AdminLoginLog struct {
	g.Meta      `orm:"table:admin_login_log, do:true"`
	Id          any         // 日志ID
	Username    any         // 尝试登录的用户名(含失败,账号可能不存在)
	AdminId     any         // 成功时回填后台账号ID
	LoginStatus any         // 结果:1成功 2失败(密码错误) 3失败(账号禁用/不存在)
	Ip          any         // 来源IP(IPv6最长45字符)
	UserAgent   any         // 浏览器User-Agent
	CreatedAt   *gtime.Time // 登录时间
}
