// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminPermission is the golang structure for table admin_permission.
type AdminPermission struct {
	Id        uint64      `json:"id"        orm:"id"         ` // 权限ID
	ParentId  uint64      `json:"parentId"  orm:"parent_id"  ` // 父权限ID,0为根(菜单树)
	Name      string      `json:"name"      orm:"name"       ` // 权限名称(如:商品管理/SPU上架)
	Code      string      `json:"code"      orm:"code"       ` // 权限编码(如 product:spu:create;接口权限与API路径对应)
	Type      int         `json:"type"      orm:"type"       ` // 类型:1菜单 2按钮/操作 3接口
	Path      string      `json:"path"      orm:"path"       ` // 前端路由地址(菜单)或API路径(接口)
	Sort      int         `json:"sort"      orm:"sort"       ` // 同级排序
	Status    int         `json:"status"    orm:"status"     ` // 状态:1启用 0禁用
	Deleted   int         `json:"deleted"   orm:"deleted"    ` // 软删除:0否 1是
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` // 更新时间
}
