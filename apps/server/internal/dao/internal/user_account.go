// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserAccountDao is the data access object for the table user_account.
type UserAccountDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserAccountColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserAccountColumns defines and stores column names for the table user_account.
type UserAccountColumns struct {
	Id        string // 账户ID
	UserId    string // 用户ID
	Balance   string // 可用余额(有符号,欠款为负,后续佣金入账抵扣;欠款不可用于余额消费)
	Frozen    string // 冻结金额(提现审核/打款中)
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// userAccountColumns holds the columns for the table user_account.
var userAccountColumns = UserAccountColumns{
	Id:        "id",
	UserId:    "user_id",
	Balance:   "balance",
	Frozen:    "frozen",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewUserAccountDao creates and returns a new DAO object for table data access.
func NewUserAccountDao(handlers ...gdb.ModelHandler) *UserAccountDao {
	return &UserAccountDao{
		group:    "default",
		table:    "user_account",
		columns:  userAccountColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserAccountDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserAccountDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserAccountDao) Columns() UserAccountColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserAccountDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserAccountDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserAccountDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
