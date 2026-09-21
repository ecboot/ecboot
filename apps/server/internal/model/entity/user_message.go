// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserMessage is the golang structure for table user_message.
type UserMessage struct {
	Id        uint64      `json:"id"        orm:"id"         ` // 消息ID
	UserId    uint64      `json:"userId"    orm:"user_id"    ` // 接收用户ID
	Title     string      `json:"title"     orm:"title"      ` // 标题
	Content   string      `json:"content"   orm:"content"    ` // 内容
	BizType   int         `json:"bizType"   orm:"biz_type"   ` // 业务类型:1订单 2营销 3售后
	BizNo     string      `json:"bizNo"     orm:"biz_no"     ` // 关联业务单号(点击跳转)
	IsRead    int         `json:"isRead"    orm:"is_read"    ` // 已读:0否 1是
	ReadTime  *gtime.Time `json:"readTime"  orm:"read_time"  ` // 阅读时间
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 创建时间
}
