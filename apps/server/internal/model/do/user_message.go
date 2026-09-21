// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserMessage is the golang structure of table user_message for DAO operations like Where/Data.
type UserMessage struct {
	g.Meta    `orm:"table:user_message, do:true"`
	Id        any         // 消息ID
	UserId    any         // 接收用户ID
	Title     any         // 标题
	Content   any         // 内容
	BizType   any         // 业务类型:1订单 2营销 3售后
	BizNo     any         // 关联业务单号(点击跳转)
	IsRead    any         // 已读:0否 1是
	ReadTime  *gtime.Time // 阅读时间
	CreatedAt *gtime.Time // 创建时间
}
