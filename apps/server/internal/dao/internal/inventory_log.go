// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// InventoryLogDao is the data access object for the table inventory_log.
type InventoryLogDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  InventoryLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// InventoryLogColumns defines and stores column names for the table inventory_log.
type InventoryLogColumns struct {
	Id          string // 流水ID
	SkuId       string // SKU ID
	OrderNo     string // 关联订单号(后台调整等无单操作为空)
	ChangeType  string // 变动类型:1下单锁定 2支付核销 3取消释放 4超时释放 5后台调整 6售后回补
	Quantity    string // 变动数量(绝对值,方向由change_type决定)
	TotalAfter  string // 变动后total快照(对账用)
	LockedAfter string // 变动后locked快照(对账用)
	Operator    string // 操作者:system/user:{id}/admin:{id}
	Remark      string // 备注
	CreatedAt   string // 创建时间
}

// inventoryLogColumns holds the columns for the table inventory_log.
var inventoryLogColumns = InventoryLogColumns{
	Id:          "id",
	SkuId:       "sku_id",
	OrderNo:     "order_no",
	ChangeType:  "change_type",
	Quantity:    "quantity",
	TotalAfter:  "total_after",
	LockedAfter: "locked_after",
	Operator:    "operator",
	Remark:      "remark",
	CreatedAt:   "created_at",
}

// NewInventoryLogDao creates and returns a new DAO object for table data access.
func NewInventoryLogDao(handlers ...gdb.ModelHandler) *InventoryLogDao {
	return &InventoryLogDao{
		group:    "default",
		table:    "inventory_log",
		columns:  inventoryLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *InventoryLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *InventoryLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *InventoryLogDao) Columns() InventoryLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *InventoryLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *InventoryLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *InventoryLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
