// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminRolePermission is the golang structure of table admin_role_permission for DAO operations like Where/Data.
type AdminRolePermission struct {
	g.Meta       `orm:"table:admin_role_permission, do:true"`
	Id           any         // 关联ID
	RoleId       any         // 角色ID
	PermissionId any         // 权限ID
	CreatedAt    *gtime.Time // 创建时间
}
