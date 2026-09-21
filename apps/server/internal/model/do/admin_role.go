// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminRole is the golang structure of table admin_role for DAO operations like Where/Data.
type AdminRole struct {
	g.Meta      `orm:"table:admin_role, do:true"`
	Id          any         // 角色ID
	Name        any         // 角色名称(如:运营/客服/财务)
	Code        any         // 角色编码(唯一,如 ops/service/finance)
	Description any         // 角色描述
	Status      any         // 状态:1启用 0停用
	Deleted     any         // 软删除:0否 1是
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
}
