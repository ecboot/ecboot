// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProductCategoryDao is the data access object for the table product_category.
type ProductCategoryDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  ProductCategoryColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// ProductCategoryColumns defines and stores column names for the table product_category.
type ProductCategoryColumns struct {
	Id        string // 分类ID
	ParentId  string // 父分类ID,0为根(固定三级:1/2/3)
	Name      string // 分类名称
	Icon      string // 分类图标URL
	Level     string // 层级:1一级 2二级 3三级
	Sort      string // 同级排序,越小越靠前
	Status    string // 状态:1启用 0禁用
	Deleted   string // 软删除:0否 1是
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// productCategoryColumns holds the columns for the table product_category.
var productCategoryColumns = ProductCategoryColumns{
	Id:        "id",
	ParentId:  "parent_id",
	Name:      "name",
	Icon:      "icon",
	Level:     "level",
	Sort:      "sort",
	Status:    "status",
	Deleted:   "deleted",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewProductCategoryDao creates and returns a new DAO object for table data access.
func NewProductCategoryDao(handlers ...gdb.ModelHandler) *ProductCategoryDao {
	return &ProductCategoryDao{
		group:    "default",
		table:    "product_category",
		columns:  productCategoryColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ProductCategoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ProductCategoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ProductCategoryDao) Columns() ProductCategoryColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ProductCategoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ProductCategoryDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ProductCategoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
