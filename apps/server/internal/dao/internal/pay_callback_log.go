// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PayCallbackLogDao is the data access object for the table pay_callback_log.
type PayCallbackLogDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  PayCallbackLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// PayCallbackLogColumns defines and stores column names for the table pay_callback_log.
type PayCallbackLogColumns struct {
	Id            string // 日志ID
	PayNo         string // 关联支付单号
	PayChannel    string // 支付渠道:1微信 2支付宝
	NotifyType    string // 通知类型:1支付结果 2退款结果
	VerifyStatus  string // 验签结果:0失败 1通过
	ProcessStatus string // 处理结果:0未处理/重复忽略 1已处理
	RawBody       string // 回调原文(全量留档,可重放审计)
	CreatedAt     string // 创建时间
}

// payCallbackLogColumns holds the columns for the table pay_callback_log.
var payCallbackLogColumns = PayCallbackLogColumns{
	Id:            "id",
	PayNo:         "pay_no",
	PayChannel:    "pay_channel",
	NotifyType:    "notify_type",
	VerifyStatus:  "verify_status",
	ProcessStatus: "process_status",
	RawBody:       "raw_body",
	CreatedAt:     "created_at",
}

// NewPayCallbackLogDao creates and returns a new DAO object for table data access.
func NewPayCallbackLogDao(handlers ...gdb.ModelHandler) *PayCallbackLogDao {
	return &PayCallbackLogDao{
		group:    "default",
		table:    "pay_callback_log",
		columns:  payCallbackLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PayCallbackLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PayCallbackLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PayCallbackLogDao) Columns() PayCallbackLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PayCallbackLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PayCallbackLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PayCallbackLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
