// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// BargainActivity is the golang structure for table bargain_activity.
type BargainActivity struct {
	Id        uint64      `json:"id"        orm:"id"         ` // 砍价活动ID
	Name      string      `json:"name"      orm:"name"       ` // 活动名称
	SpuId     uint64      `json:"spuId"     orm:"spu_id"     ` // SPU ID
	StartTime *gtime.Time `json:"startTime" orm:"start_time" ` // 开始时间(含)
	EndTime   *gtime.Time `json:"endTime"   orm:"end_time"   ` // 结束时间(不含)
	Status    int         `json:"status"    orm:"status"     ` // 状态:1启用 0停用
	Deleted   int         `json:"deleted"   orm:"deleted"    ` // 软删除:0否 1是
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` // 更新时间
}
