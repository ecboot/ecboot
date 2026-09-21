// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserFavoriteDao is the data access object for the table user_favorite.
type UserFavoriteDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  UserFavoriteColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// UserFavoriteColumns defines and stores column names for the table user_favorite.
type UserFavoriteColumns struct {
	Id        string // 收藏ID
	UserId    string // 用户ID
	SpuId     string // SPU ID(实时联查现价与可售状态,不快照)
	Deleted   string // 软删除:0否 1是
	CreatedAt string // 收藏时间
	UpdatedAt string // 更新时间
}

// userFavoriteColumns holds the columns for the table user_favorite.
var userFavoriteColumns = UserFavoriteColumns{
	Id:        "id",
	UserId:    "user_id",
	SpuId:     "spu_id",
	Deleted:   "deleted",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewUserFavoriteDao creates and returns a new DAO object for table data access.
func NewUserFavoriteDao(handlers ...gdb.ModelHandler) *UserFavoriteDao {
	return &UserFavoriteDao{
		group:    "default",
		table:    "user_favorite",
		columns:  userFavoriteColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserFavoriteDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserFavoriteDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserFavoriteDao) Columns() UserFavoriteColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserFavoriteDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserFavoriteDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserFavoriteDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
