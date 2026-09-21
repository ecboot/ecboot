// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProductBrandDao is the data access object for the table product_brand.
type ProductBrandDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  ProductBrandColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// ProductBrandColumns defines and stores column names for the table product_brand.
type ProductBrandColumns struct {
	Id          string // 品牌ID
	Name        string // 品牌名称(存活期间唯一,软删后仍占用)
	Logo        string // 品牌Logo URL
	Description string // 品牌简介
	Sort        string // 排序,越小越靠前
	Status      string // 状态:1启用 0禁用
	Deleted     string // 软删除:0否 1是
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
}

// productBrandColumns holds the columns for the table product_brand.
var productBrandColumns = ProductBrandColumns{
	Id:          "id",
	Name:        "name",
	Logo:        "logo",
	Description: "description",
	Sort:        "sort",
	Status:      "status",
	Deleted:     "deleted",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewProductBrandDao creates and returns a new DAO object for table data access.
func NewProductBrandDao(handlers ...gdb.ModelHandler) *ProductBrandDao {
	return &ProductBrandDao{
		group:    "default",
		table:    "product_brand",
		columns:  productBrandColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ProductBrandDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ProductBrandDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ProductBrandDao) Columns() ProductBrandColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ProductBrandDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ProductBrandDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ProductBrandDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
