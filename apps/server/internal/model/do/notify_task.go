// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// NotifyTask is the golang structure of table notify_task for DAO operations like Where/Data.
type NotifyTask struct {
	g.Meta        `orm:"table:notify_task, do:true"`
	Id            any         // 任务ID
	UserId        any         // 接收用户ID
	Channel       any         // 渠道:1小程序订阅消息 2短信 3站内信
	BizType       any         // 业务类型:1订单 2营销 3售后
	BizNo         any         // 关联业务单号(如订单号)
	TemplateCode  any         // 消息模板编码
	Params        any         // 模板参数(键值对)
	Status        any         // 状态:10待发送 20已发送 30失败 40达重试上限 50已跳过(渠道关闭)
	RetryCount    any         // 已重试次数
	NextRetryTime *gtime.Time // 下次重试时间(失败退避后回填)
	SentTime      *gtime.Time // 发送成功时间
	FailReason    any         // 最后失败原因
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
}
