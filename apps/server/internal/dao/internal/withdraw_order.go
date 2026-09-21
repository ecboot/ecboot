// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// WithdrawOrderDao is the data access object for the table withdraw_order.
type WithdrawOrderDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  WithdrawOrderColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// WithdrawOrderColumns defines and stores column names for the table withdraw_order.
type WithdrawOrderColumns struct {
	Id              string // 提现单ID
	WithdrawNo      string // 提现单号(全局唯一)
	UserId          string // 用户ID
	Amount          string // 提现金额
	WithdrawChannel string // 渠道:1微信商家转账(V1单渠道)
	ChannelOrderNo  string // 渠道打款单号(回填;与渠道联合唯一,防重复打款)
	Status          string // 状态:10待审核 20审核通过 30打款中 40成功 50审核拒绝 60打款失败已回退
	AuditTime       string // 审核时间
	PayTime         string // 打款成功时间
	FailReason      string // 拒绝/失败原因
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
}

// withdrawOrderColumns holds the columns for the table withdraw_order.
var withdrawOrderColumns = WithdrawOrderColumns{
	Id:              "id",
	WithdrawNo:      "withdraw_no",
	UserId:          "user_id",
	Amount:          "amount",
	WithdrawChannel: "withdraw_channel",
	ChannelOrderNo:  "channel_order_no",
	Status:          "status",
	AuditTime:       "audit_time",
	PayTime:         "pay_time",
	FailReason:      "fail_reason",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewWithdrawOrderDao creates and returns a new DAO object for table data access.
func NewWithdrawOrderDao(handlers ...gdb.ModelHandler) *WithdrawOrderDao {
	return &WithdrawOrderDao{
		group:    "default",
		table:    "withdraw_order",
		columns:  withdrawOrderColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *WithdrawOrderDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *WithdrawOrderDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *WithdrawOrderDao) Columns() WithdrawOrderColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *WithdrawOrderDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *WithdrawOrderDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *WithdrawOrderDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
