// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminRole is the golang structure for table admin_role.
type AdminRole struct {
	Id          uint64      `json:"id"          orm:"id"          ` // 角色ID
	Name        string      `json:"name"        orm:"name"        ` // 角色名称(如:运营/客服/财务)
	Code        string      `json:"code"        orm:"code"        ` // 角色编码(唯一,如 ops/service/finance)
	Description string      `json:"description" orm:"description" ` // 角色描述
	Status      int         `json:"status"      orm:"status"      ` // 状态:1启用 0停用
	Deleted     int         `json:"deleted"     orm:"deleted"     ` // 软删除:0否 1是
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  ` // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  ` // 更新时间
}
