// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FlashSaleItemDao is the data access object for the table flash_sale_item.
type FlashSaleItemDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  FlashSaleItemColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// FlashSaleItemColumns defines and stores column names for the table flash_sale_item.
type FlashSaleItemColumns struct {
	Id         string // 场次商品ID
	ActivityId string // 活动ID
	SkuId      string // SKU ID
	FlashPrice string // 秒杀价(快照进订单项price)
	StockCount string // 活动限量(与inventory分账,契约R4)
	SoldCount  string // 已售(取消回补本列;可售=stock_count-sold_count)
	PerLimit   string // 每人限购
	CreatedAt  string // 创建时间
	UpdatedAt  string // 更新时间
}

// flashSaleItemColumns holds the columns for the table flash_sale_item.
var flashSaleItemColumns = FlashSaleItemColumns{
	Id:         "id",
	ActivityId: "activity_id",
	SkuId:      "sku_id",
	FlashPrice: "flash_price",
	StockCount: "stock_count",
	SoldCount:  "sold_count",
	PerLimit:   "per_limit",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewFlashSaleItemDao creates and returns a new DAO object for table data access.
func NewFlashSaleItemDao(handlers ...gdb.ModelHandler) *FlashSaleItemDao {
	return &FlashSaleItemDao{
		group:    "default",
		table:    "flash_sale_item",
		columns:  flashSaleItemColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *FlashSaleItemDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *FlashSaleItemDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *FlashSaleItemDao) Columns() FlashSaleItemColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *FlashSaleItemDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *FlashSaleItemDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *FlashSaleItemDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
