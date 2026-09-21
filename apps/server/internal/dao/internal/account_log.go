// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AccountLogDao is the data access object for the table account_log.
type AccountLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  AccountLogColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// AccountLogColumns defines and stores column names for the table account_log.
type AccountLogColumns struct {
	Id           string // 流水ID
	UserId       string // 用户ID
	BizType      string // 业务类型:1佣金入账 2提现冻结 3提现完成 4提现失败回退 5冲销扣回 6余额消费冻结 7余额消费完成 8余额消费退回
	Amount       string // 变动金额(有符号:入账正,冻结/扣回负)
	BalanceAfter string // 变动后余额快照(对账用)
	FrozenAfter  string // 变动后冻结快照(对账用)
	BizNo        string // 关联业务单号(commission_record.id或withdraw_order.withdraw_no)
	CreatedAt    string // 创建时间
}

// accountLogColumns holds the columns for the table account_log.
var accountLogColumns = AccountLogColumns{
	Id:           "id",
	UserId:       "user_id",
	BizType:      "biz_type",
	Amount:       "amount",
	BalanceAfter: "balance_after",
	FrozenAfter:  "frozen_after",
	BizNo:        "biz_no",
	CreatedAt:    "created_at",
}

// NewAccountLogDao creates and returns a new DAO object for table data access.
func NewAccountLogDao(handlers ...gdb.ModelHandler) *AccountLogDao {
	return &AccountLogDao{
		group:    "default",
		table:    "account_log",
		columns:  accountLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AccountLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AccountLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AccountLogDao) Columns() AccountLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AccountLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AccountLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AccountLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
