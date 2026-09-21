// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GroupBuyActivity is the golang structure for table group_buy_activity.
type GroupBuyActivity struct {
	Id           uint64      `json:"id"           orm:"id"             ` // 活动ID
	SpuId        uint64      `json:"spuId"        orm:"spu_id"         ` // SPU ID(成团价对全SKU生效)
	Name         string      `json:"name"         orm:"name"           ` // 活动名称
	GroupSize    uint        `json:"groupSize"    orm:"group_size"     ` // 成团人数
	ValidStartAt *gtime.Time `json:"validStartAt" orm:"valid_start_at" ` // 活动开始时间
	ValidEndAt   *gtime.Time `json:"validEndAt"   orm:"valid_end_at"   ` // 活动结束时间
	PerLimit     uint        `json:"perLimit"     orm:"per_limit"      ` // 每人限购(整个活动期)
	Status       int         `json:"status"       orm:"status"         ` // 状态:1启用 0停用
	Deleted      int         `json:"deleted"      orm:"deleted"        ` // 软删除:0否 1是
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"     ` // 创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"     ` // 更新时间
}
