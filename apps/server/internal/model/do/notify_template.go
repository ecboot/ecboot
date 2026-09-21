// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// NotifyTemplate is the golang structure of table notify_template for DAO operations like Where/Data.
type NotifyTemplate struct {
	g.Meta             `orm:"table:notify_template, do:true"`
	Id                 any         // 模板ID
	Code               any         // 模板编码(notify_task.template_code引用)
	Channel            any         // 渠道:1小程序订阅消息 2短信 3站内信
	ExternalTemplateId any         // 平台侧模板ID(微信订阅消息/短信平台模板,平台审核产物;站内信为空)
	Title              any         // 标题(站内信自控;渠道侧为文案备份)
	ContentTemplate    any         // 内容模板({{变量}}占位;站内信全文自控,渠道侧为文案备份)
	ParamsSchema       any         // 变量契约:[{"name":"orderNo","type":"string","example":"O1"}]开发与运营的参数约定
	Status             any         // 状态:1启用 0停用(渠道级开关,短信成本闸门)
	Deleted            any         // 软删除:0否 1是
	CreatedAt          *gtime.Time // 创建时间
	UpdatedAt          *gtime.Time // 更新时间
}
