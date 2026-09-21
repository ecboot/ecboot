// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PointLogDao is the data access object for the table point_log.
type PointLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  PointLogColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// PointLogColumns defines and stores column names for the table point_log.
type PointLogColumns struct {
	Id           string // 流水ID
	UserId       string // 用户ID
	BizType      string // 业务类型:1签到 2消费获得 3下单消耗 4退款回退 5分享获得 6评价获得 7注册赠送 8邀请奖励 9过期扣减
	Points       string // 积分变动(有符号:获得正,消耗/回退负)
	BalanceAfter string // 变动后余额快照(对账用)
	OrderNo      string // 关联订单号
	CreatedAt    string // 创建时间
}

// pointLogColumns holds the columns for the table point_log.
var pointLogColumns = PointLogColumns{
	Id:           "id",
	UserId:       "user_id",
	BizType:      "biz_type",
	Points:       "points",
	BalanceAfter: "balance_after",
	OrderNo:      "order_no",
	CreatedAt:    "created_at",
}

// NewPointLogDao creates and returns a new DAO object for table data access.
func NewPointLogDao(handlers ...gdb.ModelHandler) *PointLogDao {
	return &PointLogDao{
		group:    "default",
		table:    "point_log",
		columns:  pointLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PointLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PointLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PointLogDao) Columns() PointLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PointLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PointLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PointLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
