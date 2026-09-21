// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NotifyTemplateDao is the data access object for the table notify_template.
type NotifyTemplateDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  NotifyTemplateColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// NotifyTemplateColumns defines and stores column names for the table notify_template.
type NotifyTemplateColumns struct {
	Id                 string // 模板ID
	Code               string // 模板编码(notify_task.template_code引用)
	Channel            string // 渠道:1小程序订阅消息 2短信 3站内信
	ExternalTemplateId string // 平台侧模板ID(微信订阅消息/短信平台模板,平台审核产物;站内信为空)
	Title              string // 标题(站内信自控;渠道侧为文案备份)
	ContentTemplate    string // 内容模板({{变量}}占位;站内信全文自控,渠道侧为文案备份)
	ParamsSchema       string // 变量契约:[{"name":"orderNo","type":"string","example":"O1"}]开发与运营的参数约定
	Status             string // 状态:1启用 0停用(渠道级开关,短信成本闸门)
	Deleted            string // 软删除:0否 1是
	CreatedAt          string // 创建时间
	UpdatedAt          string // 更新时间
}

// notifyTemplateColumns holds the columns for the table notify_template.
var notifyTemplateColumns = NotifyTemplateColumns{
	Id:                 "id",
	Code:               "code",
	Channel:            "channel",
	ExternalTemplateId: "external_template_id",
	Title:              "title",
	ContentTemplate:    "content_template",
	ParamsSchema:       "params_schema",
	Status:             "status",
	Deleted:            "deleted",
	CreatedAt:          "created_at",
	UpdatedAt:          "updated_at",
}

// NewNotifyTemplateDao creates and returns a new DAO object for table data access.
func NewNotifyTemplateDao(handlers ...gdb.ModelHandler) *NotifyTemplateDao {
	return &NotifyTemplateDao{
		group:    "default",
		table:    "notify_template",
		columns:  notifyTemplateColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *NotifyTemplateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *NotifyTemplateDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *NotifyTemplateDao) Columns() NotifyTemplateColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *NotifyTemplateDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *NotifyTemplateDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *NotifyTemplateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
