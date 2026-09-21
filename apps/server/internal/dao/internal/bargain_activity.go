// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BargainActivityDao is the data access object for the table bargain_activity.
type BargainActivityDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  BargainActivityColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// BargainActivityColumns defines and stores column names for the table bargain_activity.
type BargainActivityColumns struct {
	Id        string // 砍价活动ID
	Name      string // 活动名称
	SpuId     string // SPU ID
	StartTime string // 开始时间(含)
	EndTime   string // 结束时间(不含)
	Status    string // 状态:1启用 0停用
	Deleted   string // 软删除:0否 1是
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// bargainActivityColumns holds the columns for the table bargain_activity.
var bargainActivityColumns = BargainActivityColumns{
	Id:        "id",
	Name:      "name",
	SpuId:     "spu_id",
	StartTime: "start_time",
	EndTime:   "end_time",
	Status:    "status",
	Deleted:   "deleted",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewBargainActivityDao creates and returns a new DAO object for table data access.
func NewBargainActivityDao(handlers ...gdb.ModelHandler) *BargainActivityDao {
	return &BargainActivityDao{
		group:    "default",
		table:    "bargain_activity",
		columns:  bargainActivityColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BargainActivityDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BargainActivityDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BargainActivityDao) Columns() BargainActivityColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BargainActivityDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BargainActivityDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *BargainActivityDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
