// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// NotifyTask is the golang structure for table notify_task.
type NotifyTask struct {
	Id            uint64      `json:"id"            orm:"id"              ` // 任务ID
	UserId        uint64      `json:"userId"        orm:"user_id"         ` // 接收用户ID
	Channel       int         `json:"channel"       orm:"channel"         ` // 渠道:1小程序订阅消息 2短信 3站内信
	BizType       int         `json:"bizType"       orm:"biz_type"        ` // 业务类型:1订单 2营销 3售后
	BizNo         string      `json:"bizNo"         orm:"biz_no"          ` // 关联业务单号(如订单号)
	TemplateCode  string      `json:"templateCode"  orm:"template_code"   ` // 消息模板编码
	Params        string      `json:"params"        orm:"params"          ` // 模板参数(键值对)
	Status        int         `json:"status"        orm:"status"          ` // 状态:10待发送 20已发送 30失败 40达重试上限 50已跳过(渠道关闭)
	RetryCount    uint        `json:"retryCount"    orm:"retry_count"     ` // 已重试次数
	NextRetryTime *gtime.Time `json:"nextRetryTime" orm:"next_retry_time" ` // 下次重试时间(失败退避后回填)
	SentTime      *gtime.Time `json:"sentTime"      orm:"sent_time"       ` // 发送成功时间
	FailReason    string      `json:"failReason"    orm:"fail_reason"     ` // 最后失败原因
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      ` // 创建时间
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"      ` // 更新时间
}
