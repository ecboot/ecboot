// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserFavorite is the golang structure for table user_favorite.
type UserFavorite struct {
	Id        uint64      `json:"id"        orm:"id"         ` // 收藏ID
	UserId    uint64      `json:"userId"    orm:"user_id"    ` // 用户ID
	SpuId     uint64      `json:"spuId"     orm:"spu_id"     ` // SPU ID(实时联查现价与可售状态,不快照)
	Deleted   int         `json:"deleted"   orm:"deleted"    ` // 软删除:0否 1是
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 收藏时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` // 更新时间
}
