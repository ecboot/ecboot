// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminRolePermission is the golang structure for table admin_role_permission.
type AdminRolePermission struct {
	Id           uint64      `json:"id"           orm:"id"            ` // 关联ID
	RoleId       uint64      `json:"roleId"       orm:"role_id"       ` // 角色ID
	PermissionId uint64      `json:"permissionId" orm:"permission_id" ` // 权限ID
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` // 创建时间
}
