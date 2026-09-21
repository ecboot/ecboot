// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ShareRecord is the golang structure for table share_record.
type ShareRecord struct {
	Id           uint64      `json:"id"           orm:"id"             ` // 分享记录ID
	SharerUserId uint64      `json:"sharerUserId" orm:"sharer_user_id" ` // 分享人用户ID
	SpuId        uint64      `json:"spuId"        orm:"spu_id"         ` // 分享对象SPU(NULL=整店/页面分享)
	ShareChannel int         `json:"shareChannel" orm:"share_channel"  ` // 渠道:1小程序卡片 2海报图片 3复制链接/口令 4朋友圈/社群
	SceneValue   string      `json:"sceneValue"   orm:"scene_value"    ` // 小程序场景值(scene参数,落地追踪)
	ShareCode    string      `json:"shareCode"    orm:"share_code"     ` // 使用的推广码
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"     ` // 分享时间(归因窗口起点)
}
