// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserNotifyPreference is the golang structure for table user_notify_preference.
type UserNotifyPreference struct {
	Id        uint64      `json:"id"        orm:"id"         ` // 偏好ID
	UserId    uint64      `json:"userId"    orm:"user_id"    ` // 会员ID
	Channel   int         `json:"channel"   orm:"channel"    ` // 渠道:1小程序订阅 2短信
	Enabled   int         `json:"enabled"   orm:"enabled"    ` // 是否接收:1接收 0关闭
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` // 更新时间
}
