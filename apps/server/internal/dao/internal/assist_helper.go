// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AssistHelperDao is the data access object for the table assist_helper.
type AssistHelperDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  AssistHelperColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// AssistHelperColumns defines and stores column names for the table assist_helper.
type AssistHelperColumns struct {
	Id           string // 助力记录ID
	RecordId     string // 参与记录ID
	HelperUserId string // 助力人用户ID
	CreatedAt    string // 助力时间
}

// assistHelperColumns holds the columns for the table assist_helper.
var assistHelperColumns = AssistHelperColumns{
	Id:           "id",
	RecordId:     "record_id",
	HelperUserId: "helper_user_id",
	CreatedAt:    "created_at",
}

// NewAssistHelperDao creates and returns a new DAO object for table data access.
func NewAssistHelperDao(handlers ...gdb.ModelHandler) *AssistHelperDao {
	return &AssistHelperDao{
		group:    "default",
		table:    "assist_helper",
		columns:  assistHelperColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AssistHelperDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AssistHelperDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AssistHelperDao) Columns() AssistHelperColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AssistHelperDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AssistHelperDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AssistHelperDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
