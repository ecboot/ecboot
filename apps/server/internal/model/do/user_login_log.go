// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserLoginLog is the golang structure of table user_login_log for DAO operations like Where/Data.
type UserLoginLog struct {
	g.Meta       `orm:"table:user_login_log, do:true"`
	Id           any         // 日志ID
	UserId       any         // 用户ID
	LoginChannel any         // 登录渠道:1微信小程序 2H5
	LoginStatus  any         // 结果:1成功 2失败
	Ip           any         // 来源IP(IPv6最长45字符)
	UserAgent    any         // 浏览器/客户端UA
	CreatedAt    *gtime.Time // 登录时间
}
