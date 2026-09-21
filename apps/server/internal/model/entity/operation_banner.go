// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OperationBanner is the golang structure for table operation_banner.
type OperationBanner struct {
	Id        uint64      `json:"id"        orm:"id"         ` // bannerID
	Position  int         `json:"position"  orm:"position"   ` // 位置:1首页轮播 2首页弹窗(枚举可扩展)
	ImageUrl  string      `json:"imageUrl"  orm:"image_url"  ` // 图片URL(OSS/CDN)
	LinkUrl   string      `json:"linkUrl"   orm:"link_url"   ` // 跳转链接(小程序页面路径或H5;空=纯展示)
	Sort      int         `json:"sort"      orm:"sort"       ` // 排序,越小越靠前
	StartTime *gtime.Time `json:"startTime" orm:"start_time" ` // 投放开始(NULL=立即生效)
	EndTime   *gtime.Time `json:"endTime"   orm:"end_time"   ` // 投放结束(NULL=长期有效)
	Status    int         `json:"status"    orm:"status"     ` // 状态:1启用 0停用
	Deleted   int         `json:"deleted"   orm:"deleted"    ` // 软删除:0否 1是
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` // 更新时间
}
