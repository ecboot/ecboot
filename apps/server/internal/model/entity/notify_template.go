// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// NotifyTemplate is the golang structure for table notify_template.
type NotifyTemplate struct {
	Id                 uint64      `json:"id"                 orm:"id"                   ` // 模板ID
	Code               string      `json:"code"               orm:"code"                 ` // 模板编码(notify_task.template_code引用)
	Channel            int         `json:"channel"            orm:"channel"              ` // 渠道:1小程序订阅消息 2短信 3站内信
	ExternalTemplateId string      `json:"externalTemplateId" orm:"external_template_id" ` // 平台侧模板ID(微信订阅消息/短信平台模板,平台审核产物;站内信为空)
	Title              string      `json:"title"              orm:"title"                ` // 标题(站内信自控;渠道侧为文案备份)
	ContentTemplate    string      `json:"contentTemplate"    orm:"content_template"     ` // 内容模板({{变量}}占位;站内信全文自控,渠道侧为文案备份)
	ParamsSchema       string      `json:"paramsSchema"       orm:"params_schema"        ` // 变量契约:[{"name":"orderNo","type":"string","example":"O1"}]开发与运营的参数约定
	Status             int         `json:"status"             orm:"status"               ` // 状态:1启用 0停用(渠道级开关,短信成本闸门)
	Deleted            int         `json:"deleted"            orm:"deleted"              ` // 软删除:0否 1是
	CreatedAt          *gtime.Time `json:"createdAt"          orm:"created_at"           ` // 创建时间
	UpdatedAt          *gtime.Time `json:"updatedAt"          orm:"updated_at"           ` // 更新时间
}
