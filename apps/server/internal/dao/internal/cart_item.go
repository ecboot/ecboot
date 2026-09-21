// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CartItemDao is the data access object for the table cart_item.
type CartItemDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  CartItemColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// CartItemColumns defines and stores column names for the table cart_item.
type CartItemColumns struct {
	Id        string // 购物车项ID
	UserId    string // 用户ID
	SkuId     string // SKU ID
	Quantity  string // 数量(应用层限制上限,如99)
	Checked   string // 结算勾选:0未选 1选中
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// cartItemColumns holds the columns for the table cart_item.
var cartItemColumns = CartItemColumns{
	Id:        "id",
	UserId:    "user_id",
	SkuId:     "sku_id",
	Quantity:  "quantity",
	Checked:   "checked",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewCartItemDao creates and returns a new DAO object for table data access.
func NewCartItemDao(handlers ...gdb.ModelHandler) *CartItemDao {
	return &CartItemDao{
		group:    "default",
		table:    "cart_item",
		columns:  cartItemColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CartItemDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CartItemDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CartItemDao) Columns() CartItemColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CartItemDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CartItemDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CartItemDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
