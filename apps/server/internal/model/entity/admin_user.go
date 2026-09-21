// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminUser is the golang structure for table admin_user.
type AdminUser struct {
	Id            uint64      `json:"id"            orm:"id"              ` // 后台账号ID
	Username      string      `json:"username"      orm:"username"        ` // 登录名(唯一)
	PasswordHash  string      `json:"passwordHash"  orm:"password_hash"   ` // 密码哈希(bcrypt)
	RealName      string      `json:"realName"      orm:"real_name"       ` // 姓名
	IsSuper       int         `json:"isSuper"       orm:"is_super"        ` // 超级管理员:1是(跳过权限校验) 0否
	Status        int         `json:"status"        orm:"status"          ` // 状态:1正常 2禁用
	LastLoginTime *gtime.Time `json:"lastLoginTime" orm:"last_login_time" ` // 最后登录时间
	Deleted       int         `json:"deleted"       orm:"deleted"         ` // 软删除:0否 1是
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      ` // 创建时间
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"      ` // 更新时间
}
