// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminUser is the golang structure of table admin_user for DAO operations like Where/Data.
type AdminUser struct {
	g.Meta        `orm:"table:admin_user, do:true"`
	Id            any         // 后台账号ID
	Username      any         // 登录名(唯一)
	PasswordHash  any         // 密码哈希(bcrypt)
	RealName      any         // 姓名
	IsSuper       any         // 超级管理员:1是(跳过权限校验) 0否
	Status        any         // 状态:1正常 2禁用
	LastLoginTime *gtime.Time // 最后登录时间
	Deleted       any         // 软删除:0否 1是
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
}
