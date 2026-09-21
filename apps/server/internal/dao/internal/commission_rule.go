// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CommissionRuleDao is the data access object for the table commission_rule.
type CommissionRuleDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  CommissionRuleColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// CommissionRuleColumns defines and stores column names for the table commission_rule.
type CommissionRuleColumns struct {
	Id         string // 规则ID
	ScopeType  string // 作用域:1分类(默认) 2商品(覆盖)
	ScopeId    string // 作用域目标ID(分类ID或SPU ID)
	Level1Rate string // 一级佣金比例%(直接邀请人,0-100)
	Level2Rate string // 二级佣金比例%(邀请人的邀请人,0-100)
	Status     string // 状态:1启用 0停用
	Deleted    string // 软删除:0否 1是
	CreatedAt  string // 创建时间
	UpdatedAt  string // 更新时间
}

// commissionRuleColumns holds the columns for the table commission_rule.
var commissionRuleColumns = CommissionRuleColumns{
	Id:         "id",
	ScopeType:  "scope_type",
	ScopeId:    "scope_id",
	Level1Rate: "level1_rate",
	Level2Rate: "level2_rate",
	Status:     "status",
	Deleted:    "deleted",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewCommissionRuleDao creates and returns a new DAO object for table data access.
func NewCommissionRuleDao(handlers ...gdb.ModelHandler) *CommissionRuleDao {
	return &CommissionRuleDao{
		group:    "default",
		table:    "commission_rule",
		columns:  commissionRuleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CommissionRuleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CommissionRuleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CommissionRuleDao) Columns() CommissionRuleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CommissionRuleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CommissionRuleDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CommissionRuleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
