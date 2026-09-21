// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminUserRole is the golang structure of table admin_user_role for DAO operations like Where/Data.
type AdminUserRole struct {
	g.Meta    `orm:"table:admin_user_role, do:true"`
	Id        any         // 关联ID
	AdminId   any         // 后台账号ID
	RoleId    any         // 角色ID
	CreatedAt *gtime.Time // 创建时间
}
