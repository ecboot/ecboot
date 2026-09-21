// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ShareRecord is the golang structure of table share_record for DAO operations like Where/Data.
type ShareRecord struct {
	g.Meta       `orm:"table:share_record, do:true"`
	Id           any         // 分享记录ID
	SharerUserId any         // 分享人用户ID
	SpuId        any         // 分享对象SPU(NULL=整店/页面分享)
	ShareChannel any         // 渠道:1小程序卡片 2海报图片 3复制链接/口令 4朋友圈/社群
	SceneValue   any         // 小程序场景值(scene参数,落地追踪)
	ShareCode    any         // 使用的推广码
	CreatedAt    *gtime.Time // 分享时间(归因窗口起点)
}
