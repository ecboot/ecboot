// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PromotionActivityLadderDao is the data access object for the table promotion_activity_ladder.
type PromotionActivityLadderDao struct {
	table    string                         // table is the underlying table name of the DAO.
	group    string                         // group is the database configuration group name of the current DAO.
	columns  PromotionActivityLadderColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler             // handlers for customized model modification.
}

// PromotionActivityLadderColumns defines and stores column names for the table promotion_activity_ladder.
type PromotionActivityLadderColumns struct {
	Id              string // 档位ID
	ActivityId      string // 活动ID
	ThresholdAmount string // 门槛金额(满X元,范围内商品合计)
	DiscountAmount  string // 优惠金额(减Y元)
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
}

// promotionActivityLadderColumns holds the columns for the table promotion_activity_ladder.
var promotionActivityLadderColumns = PromotionActivityLadderColumns{
	Id:              "id",
	ActivityId:      "activity_id",
	ThresholdAmount: "threshold_amount",
	DiscountAmount:  "discount_amount",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewPromotionActivityLadderDao creates and returns a new DAO object for table data access.
func NewPromotionActivityLadderDao(handlers ...gdb.ModelHandler) *PromotionActivityLadderDao {
	return &PromotionActivityLadderDao{
		group:    "default",
		table:    "promotion_activity_ladder",
		columns:  promotionActivityLadderColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PromotionActivityLadderDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PromotionActivityLadderDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PromotionActivityLadderDao) Columns() PromotionActivityLadderColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PromotionActivityLadderDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PromotionActivityLadderDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PromotionActivityLadderDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
