// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OperationBanner is the golang structure of table operation_banner for DAO operations like Where/Data.
type OperationBanner struct {
	g.Meta    `orm:"table:operation_banner, do:true"`
	Id        any         // bannerID
	Position  any         // 位置:1首页轮播 2首页弹窗(枚举可扩展)
	ImageUrl  any         // 图片URL(OSS/CDN)
	LinkUrl   any         // 跳转链接(小程序页面路径或H5;空=纯展示)
	Sort      any         // 排序,越小越靠前
	StartTime *gtime.Time // 投放开始(NULL=立即生效)
	EndTime   *gtime.Time // 投放结束(NULL=长期有效)
	Status    any         // 状态:1启用 0停用
	Deleted   any         // 软删除:0否 1是
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
