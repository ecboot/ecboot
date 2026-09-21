// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PointAccountDao is the data access object for the table point_account.
type PointAccountDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  PointAccountColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// PointAccountColumns defines and stores column names for the table point_account.
type PointAccountColumns struct {
	Id           string // 积分账户ID
	UserId       string // 用户ID
	Balance      string // 积分余额(有符号,回退欠款为负,禁止UNSIGNED)
	LastEarnedAt string // 最后获得时间(滚动有效期口径:自最后获得日起12个月内有效;过期任务清零并记流水)
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
}

// pointAccountColumns holds the columns for the table point_account.
var pointAccountColumns = PointAccountColumns{
	Id:           "id",
	UserId:       "user_id",
	Balance:      "balance",
	LastEarnedAt: "last_earned_at",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewPointAccountDao creates and returns a new DAO object for table data access.
func NewPointAccountDao(handlers ...gdb.ModelHandler) *PointAccountDao {
	return &PointAccountDao{
		group:    "default",
		table:    "point_account",
		columns:  pointAccountColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PointAccountDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PointAccountDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PointAccountDao) Columns() PointAccountColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PointAccountDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PointAccountDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PointAccountDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
