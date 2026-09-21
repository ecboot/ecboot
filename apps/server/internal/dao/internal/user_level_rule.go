// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserLevelRuleDao is the data access object for the table user_level_rule.
type UserLevelRuleDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  UserLevelRuleColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// UserLevelRuleColumns defines and stores column names for the table user_level_rule.
type UserLevelRuleColumns struct {
	Id              string // 等级规则ID
	Name            string // 等级名称(如 普通会员/白银/黄金)
	GrowthThreshold string // 成长值门槛(唯一,等级判定=满足的最高门槛)
	Benefits        string // 等级权益配置(折扣/优先客服等,键值对)
	Status          string // 状态:1启用 0停用
	Deleted         string // 软删除:0否 1是
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
}

// userLevelRuleColumns holds the columns for the table user_level_rule.
var userLevelRuleColumns = UserLevelRuleColumns{
	Id:              "id",
	Name:            "name",
	GrowthThreshold: "growth_threshold",
	Benefits:        "benefits",
	Status:          "status",
	Deleted:         "deleted",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewUserLevelRuleDao creates and returns a new DAO object for table data access.
func NewUserLevelRuleDao(handlers ...gdb.ModelHandler) *UserLevelRuleDao {
	return &UserLevelRuleDao{
		group:    "default",
		table:    "user_level_rule",
		columns:  userLevelRuleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserLevelRuleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserLevelRuleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserLevelRuleDao) Columns() UserLevelRuleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserLevelRuleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserLevelRuleDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserLevelRuleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
