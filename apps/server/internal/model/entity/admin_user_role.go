// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminUserRole is the golang structure for table admin_user_role.
type AdminUserRole struct {
	Id        uint64      `json:"id"        orm:"id"         ` // 关联ID
	AdminId   uint64      `json:"adminId"   orm:"admin_id"   ` // 后台账号ID
	RoleId    uint64      `json:"roleId"    orm:"role_id"    ` // 角色ID
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 创建时间
}
