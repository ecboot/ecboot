// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserFootprint is the golang structure for table user_footprint.
type UserFootprint struct {
	Id         uint64      `json:"id"         orm:"id"           ` // 足迹ID
	UserId     uint64      `json:"userId"     orm:"user_id"      ` // 用户ID
	SpuId      uint64      `json:"spuId"      orm:"spu_id"       ` // SPU ID
	ViewCount  uint        `json:"viewCount"  orm:"view_count"   ` // 累计浏览次数(重复浏览累加)
	LastViewAt *gtime.Time `json:"lastViewAt" orm:"last_view_at" ` // 最近浏览时间(重复浏览更新此列,清理任务扫描)
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"   ` // 首次浏览时间
}
