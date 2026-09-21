// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OperationFloorDao is the data access object for the table operation_floor.
type OperationFloorDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  OperationFloorColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// OperationFloorColumns defines and stores column names for the table operation_floor.
type OperationFloorColumns struct {
	Id        string // 楼层ID
	FloorType string // 类型:1金刚区 2商品楼层 3专题
	Title     string // 楼层标题
	Config    string // 楼层配置(金刚区=图标入口数组/商品楼层=商品ID列表与参数;类型化schema由应用层约定)
	Sort      string // 排序,越小越靠前
	Status    string // 状态:1启用 0停用
	Deleted   string // 软删除:0否 1是
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// operationFloorColumns holds the columns for the table operation_floor.
var operationFloorColumns = OperationFloorColumns{
	Id:        "id",
	FloorType: "floor_type",
	Title:     "title",
	Config:    "config",
	Sort:      "sort",
	Status:    "status",
	Deleted:   "deleted",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewOperationFloorDao creates and returns a new DAO object for table data access.
func NewOperationFloorDao(handlers ...gdb.ModelHandler) *OperationFloorDao {
	return &OperationFloorDao{
		group:    "default",
		table:    "operation_floor",
		columns:  operationFloorColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *OperationFloorDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *OperationFloorDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *OperationFloorDao) Columns() OperationFloorColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *OperationFloorDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *OperationFloorDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *OperationFloorDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
