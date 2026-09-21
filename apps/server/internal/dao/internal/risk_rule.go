// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RiskRuleDao is the data access object for the table risk_rule.
type RiskRuleDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  RiskRuleColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// RiskRuleColumns defines and stores column names for the table risk_rule.
type RiskRuleColumns struct {
	Id            string // 规则ID
	Name          string // 规则名称(如 黑名单/高频下单/佣金套利特征)
	RuleType      string // 规则类型:1黑名单 2高频下单 3异常领券 4佣金套利特征 5休眠账户分级
	ConditionExpr string // 规则条件描述(如 1分钟内下单>10次)
	Action        string // 处置:1拦截 2标记(放行但留痕)
	Status        string // 状态:1启用 0停用
	Deleted       string // 软删除:0否 1是
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
}

// riskRuleColumns holds the columns for the table risk_rule.
var riskRuleColumns = RiskRuleColumns{
	Id:            "id",
	Name:          "name",
	RuleType:      "rule_type",
	ConditionExpr: "condition_expr",
	Action:        "action",
	Status:        "status",
	Deleted:       "deleted",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewRiskRuleDao creates and returns a new DAO object for table data access.
func NewRiskRuleDao(handlers ...gdb.ModelHandler) *RiskRuleDao {
	return &RiskRuleDao{
		group:    "default",
		table:    "risk_rule",
		columns:  riskRuleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RiskRuleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RiskRuleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RiskRuleDao) Columns() RiskRuleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RiskRuleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RiskRuleDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *RiskRuleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
