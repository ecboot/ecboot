// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserFootprintDao is the data access object for the table user_footprint.
type UserFootprintDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  UserFootprintColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// UserFootprintColumns defines and stores column names for the table user_footprint.
type UserFootprintColumns struct {
	Id         string // 足迹ID
	UserId     string // 用户ID
	SpuId      string // SPU ID
	ViewCount  string // 累计浏览次数(重复浏览累加)
	LastViewAt string // 最近浏览时间(重复浏览更新此列,清理任务扫描)
	CreatedAt  string // 首次浏览时间
}

// userFootprintColumns holds the columns for the table user_footprint.
var userFootprintColumns = UserFootprintColumns{
	Id:         "id",
	UserId:     "user_id",
	SpuId:      "spu_id",
	ViewCount:  "view_count",
	LastViewAt: "last_view_at",
	CreatedAt:  "created_at",
}

// NewUserFootprintDao creates and returns a new DAO object for table data access.
func NewUserFootprintDao(handlers ...gdb.ModelHandler) *UserFootprintDao {
	return &UserFootprintDao{
		group:    "default",
		table:    "user_footprint",
		columns:  userFootprintColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserFootprintDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserFootprintDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserFootprintDao) Columns() UserFootprintColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserFootprintDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserFootprintDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserFootprintDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
