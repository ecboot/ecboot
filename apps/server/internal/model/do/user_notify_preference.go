// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserNotifyPreference is the golang structure of table user_notify_preference for DAO operations like Where/Data.
type UserNotifyPreference struct {
	g.Meta    `orm:"table:user_notify_preference, do:true"`
	Id        any         // 偏好ID
	UserId    any         // 会员ID
	Channel   any         // 渠道:1小程序订阅 2短信
	Enabled   any         // 是否接收:1接收 0关闭
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
