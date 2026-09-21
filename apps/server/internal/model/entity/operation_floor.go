// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OperationFloor is the golang structure for table operation_floor.
type OperationFloor struct {
	Id        uint64      `json:"id"        orm:"id"         ` // 楼层ID
	FloorType int         `json:"floorType" orm:"floor_type" ` // 类型:1金刚区 2商品楼层 3专题
	Title     string      `json:"title"     orm:"title"      ` // 楼层标题
	Config    string      `json:"config"    orm:"config"     ` // 楼层配置(金刚区=图标入口数组/商品楼层=商品ID列表与参数;类型化schema由应用层约定)
	Sort      int         `json:"sort"      orm:"sort"       ` // 排序,越小越靠前
	Status    int         `json:"status"    orm:"status"     ` // 状态:1启用 0停用
	Deleted   int         `json:"deleted"   orm:"deleted"    ` // 软删除:0否 1是
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` // 更新时间
}
