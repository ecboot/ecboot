// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PromotionActivityScopeDao is the data access object for the table promotion_activity_scope.
type PromotionActivityScopeDao struct {
	table    string                        // table is the underlying table name of the DAO.
	group    string                        // group is the database configuration group name of the current DAO.
	columns  PromotionActivityScopeColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler            // handlers for customized model modification.
}

// PromotionActivityScopeColumns defines and stores column names for the table promotion_activity_scope.
type PromotionActivityScopeColumns struct {
	Id           string // 范围ID
	ActivityId   string // 活动ID
	ScopeType    string // 范围类型:1全场 2分类 3商品
	TargetId     string // 目标ID(scope_type=2分类ID/3商品SPU_ID;全场为NULL,每活动至多一条全场行)
	TargetIdNorm string // 唯一键载体:全场(NULL)归一为0,堵住NULL≠NULL多行漏洞
	CreatedAt    string // 创建时间
}

// promotionActivityScopeColumns holds the columns for the table promotion_activity_scope.
var promotionActivityScopeColumns = PromotionActivityScopeColumns{
	Id:           "id",
	ActivityId:   "activity_id",
	ScopeType:    "scope_type",
	TargetId:     "target_id",
	TargetIdNorm: "target_id_norm",
	CreatedAt:    "created_at",
}

// NewPromotionActivityScopeDao creates and returns a new DAO object for table data access.
func NewPromotionActivityScopeDao(handlers ...gdb.ModelHandler) *PromotionActivityScopeDao {
	return &PromotionActivityScopeDao{
		group:    "default",
		table:    "promotion_activity_scope",
		columns:  promotionActivityScopeColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PromotionActivityScopeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PromotionActivityScopeDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PromotionActivityScopeDao) Columns() PromotionActivityScopeColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PromotionActivityScopeDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PromotionActivityScopeDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PromotionActivityScopeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
