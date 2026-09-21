// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NotifyTaskDao is the data access object for the table notify_task.
type NotifyTaskDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  NotifyTaskColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// NotifyTaskColumns defines and stores column names for the table notify_task.
type NotifyTaskColumns struct {
	Id            string // 任务ID
	UserId        string // 接收用户ID
	Channel       string // 渠道:1小程序订阅消息 2短信 3站内信
	BizType       string // 业务类型:1订单 2营销 3售后
	BizNo         string // 关联业务单号(如订单号)
	TemplateCode  string // 消息模板编码
	Params        string // 模板参数(键值对)
	Status        string // 状态:10待发送 20已发送 30失败 40达重试上限 50已跳过(渠道关闭)
	RetryCount    string // 已重试次数
	NextRetryTime string // 下次重试时间(失败退避后回填)
	SentTime      string // 发送成功时间
	FailReason    string // 最后失败原因
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
}

// notifyTaskColumns holds the columns for the table notify_task.
var notifyTaskColumns = NotifyTaskColumns{
	Id:            "id",
	UserId:        "user_id",
	Channel:       "channel",
	BizType:       "biz_type",
	BizNo:         "biz_no",
	TemplateCode:  "template_code",
	Params:        "params",
	Status:        "status",
	RetryCount:    "retry_count",
	NextRetryTime: "next_retry_time",
	SentTime:      "sent_time",
	FailReason:    "fail_reason",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewNotifyTaskDao creates and returns a new DAO object for table data access.
func NewNotifyTaskDao(handlers ...gdb.ModelHandler) *NotifyTaskDao {
	return &NotifyTaskDao{
		group:    "default",
		table:    "notify_task",
		columns:  notifyTaskColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *NotifyTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *NotifyTaskDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *NotifyTaskDao) Columns() NotifyTaskColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *NotifyTaskDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *NotifyTaskDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *NotifyTaskDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
