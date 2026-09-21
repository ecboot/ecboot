// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BargainItemDao is the data access object for the table bargain_item.
type BargainItemDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  BargainItemColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// BargainItemColumns defines and stores column names for the table bargain_item.
type BargainItemColumns struct {
	Id            string // 砍价场次商品ID
	ActivityId    string // 活动ID
	SkuId         string // SKU ID
	OriginalPrice string // 起始价(=发起时售价口径)
	FloorPrice    string // 底价(砍到底价即可下单;floor<=original应用校验)
	MaxCutCount   string // 最大帮砍刀数(0=不限,金额收敛到底价)
	Config        string // 玩法参数(随机递减算法边界/首刀上限等,应用层约定)
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
}

// bargainItemColumns holds the columns for the table bargain_item.
var bargainItemColumns = BargainItemColumns{
	Id:            "id",
	ActivityId:    "activity_id",
	SkuId:         "sku_id",
	OriginalPrice: "original_price",
	FloorPrice:    "floor_price",
	MaxCutCount:   "max_cut_count",
	Config:        "config",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewBargainItemDao creates and returns a new DAO object for table data access.
func NewBargainItemDao(handlers ...gdb.ModelHandler) *BargainItemDao {
	return &BargainItemDao{
		group:    "default",
		table:    "bargain_item",
		columns:  bargainItemColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BargainItemDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BargainItemDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BargainItemDao) Columns() BargainItemColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BargainItemDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BargainItemDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *BargainItemDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
