// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GroupBuyItemDao is the data access object for the table group_buy_item.
type GroupBuyItemDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  GroupBuyItemColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// GroupBuyItemColumns defines and stores column names for the table group_buy_item.
type GroupBuyItemColumns struct {
	Id         string // 拼团场次商品ID
	ActivityId string // 拼团活动ID
	SkuId      string // SKU ID
	GroupPrice string // 该SKU成团价(下单快照进订单项price)
	CreatedAt  string // 创建时间
	UpdatedAt  string // 更新时间
}

// groupBuyItemColumns holds the columns for the table group_buy_item.
var groupBuyItemColumns = GroupBuyItemColumns{
	Id:         "id",
	ActivityId: "activity_id",
	SkuId:      "sku_id",
	GroupPrice: "group_price",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewGroupBuyItemDao creates and returns a new DAO object for table data access.
func NewGroupBuyItemDao(handlers ...gdb.ModelHandler) *GroupBuyItemDao {
	return &GroupBuyItemDao{
		group:    "default",
		table:    "group_buy_item",
		columns:  groupBuyItemColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *GroupBuyItemDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *GroupBuyItemDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *GroupBuyItemDao) Columns() GroupBuyItemColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *GroupBuyItemDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *GroupBuyItemDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *GroupBuyItemDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
