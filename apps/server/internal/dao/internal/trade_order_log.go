// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TradeOrderLogDao is the data access object for the table trade_order_log.
type TradeOrderLogDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  TradeOrderLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// TradeOrderLogColumns defines and stores column names for the table trade_order_log.
type TradeOrderLogColumns struct {
	Id           string // 流水ID
	OrderId      string // 订单ID
	OrderNo      string // 订单号
	FromStatus   string // 迁移前状态(建单时为NULL)
	ToStatus     string // 迁移后状态
	OperatorType string // 操作者类型:1系统 2用户 3管理员
	OperatorId   string // 操作者标识(system/user:{id}/admin:{id})
	Remark       string // 备注(如超时取消/确认收货)
	CreatedAt    string // 创建时间
}

// tradeOrderLogColumns holds the columns for the table trade_order_log.
var tradeOrderLogColumns = TradeOrderLogColumns{
	Id:           "id",
	OrderId:      "order_id",
	OrderNo:      "order_no",
	FromStatus:   "from_status",
	ToStatus:     "to_status",
	OperatorType: "operator_type",
	OperatorId:   "operator_id",
	Remark:       "remark",
	CreatedAt:    "created_at",
}

// NewTradeOrderLogDao creates and returns a new DAO object for table data access.
func NewTradeOrderLogDao(handlers ...gdb.ModelHandler) *TradeOrderLogDao {
	return &TradeOrderLogDao{
		group:    "default",
		table:    "trade_order_log",
		columns:  tradeOrderLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TradeOrderLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TradeOrderLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TradeOrderLogDao) Columns() TradeOrderLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TradeOrderLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TradeOrderLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TradeOrderLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
