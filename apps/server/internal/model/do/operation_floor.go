// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OperationFloor is the golang structure of table operation_floor for DAO operations like Where/Data.
type OperationFloor struct {
	g.Meta    `orm:"table:operation_floor, do:true"`
	Id        any         // 楼层ID
	FloorType any         // 类型:1金刚区 2商品楼层 3专题
	Title     any         // 楼层标题
	Config    any         // 楼层配置(金刚区=图标入口数组/商品楼层=商品ID列表与参数;类型化schema由应用层约定)
	Sort      any         // 排序,越小越靠前
	Status    any         // 状态:1启用 0停用
	Deleted   any         // 软删除:0否 1是
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
