// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FlashSaleActivityDao is the data access object for the table flash_sale_activity.
type FlashSaleActivityDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  FlashSaleActivityColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// FlashSaleActivityColumns defines and stores column names for the table flash_sale_activity.
type FlashSaleActivityColumns struct {
	Id        string // 活动ID
	Name      string // 活动名称
	StartTime string // 开始时间(含)
	EndTime   string // 结束时间(不含)
	Status    string // 状态:1启用 0停用
	Deleted   string // 软删除:0否 1是
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// flashSaleActivityColumns holds the columns for the table flash_sale_activity.
var flashSaleActivityColumns = FlashSaleActivityColumns{
	Id:        "id",
	Name:      "name",
	StartTime: "start_time",
	EndTime:   "end_time",
	Status:    "status",
	Deleted:   "deleted",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewFlashSaleActivityDao creates and returns a new DAO object for table data access.
func NewFlashSaleActivityDao(handlers ...gdb.ModelHandler) *FlashSaleActivityDao {
	return &FlashSaleActivityDao{
		group:    "default",
		table:    "flash_sale_activity",
		columns:  flashSaleActivityColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *FlashSaleActivityDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *FlashSaleActivityDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *FlashSaleActivityDao) Columns() FlashSaleActivityColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *FlashSaleActivityDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *FlashSaleActivityDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *FlashSaleActivityDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
