// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminPermission is the golang structure of table admin_permission for DAO operations like Where/Data.
type AdminPermission struct {
	g.Meta    `orm:"table:admin_permission, do:true"`
	Id        any         // 权限ID
	ParentId  any         // 父权限ID,0为根(菜单树)
	Name      any         // 权限名称(如:商品管理/SPU上架)
	Code      any         // 权限编码(如 product:spu:create;接口权限与API路径对应)
	Type      any         // 类型:1菜单 2按钮/操作 3接口
	Path      any         // 前端路由地址(菜单)或API路径(接口)
	Sort      any         // 同级排序
	Status    any         // 状态:1启用 0禁用
	Deleted   any         // 软删除:0否 1是
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
