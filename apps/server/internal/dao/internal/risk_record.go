// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RiskRecordDao is the data access object for the table risk_record.
type RiskRecordDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  RiskRecordColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// RiskRecordColumns defines and stores column names for the table risk_record.
type RiskRecordColumns struct {
	Id           string // 事件ID
	UserId       string // 命中用户ID
	RuleId       string // 命中规则ID
	ObjectType   string // 关联对象类型:1订单 2优惠券 3提现 4售后
	ObjectNo     string // 关联对象单号(订单号/券ID/提现单号等)
	Action       string // 处置结果:1拦截 2标记
	AppealStatus string // 申诉状态:0无 1申诉中 2申诉通过(解除拦截) 3申诉驳回
	Remark       string // 备注(命中明细/申诉结论)
	CreatedAt    string // 命中时间
}

// riskRecordColumns holds the columns for the table risk_record.
var riskRecordColumns = RiskRecordColumns{
	Id:           "id",
	UserId:       "user_id",
	RuleId:       "rule_id",
	ObjectType:   "object_type",
	ObjectNo:     "object_no",
	Action:       "action",
	AppealStatus: "appeal_status",
	Remark:       "remark",
	CreatedAt:    "created_at",
}

// NewRiskRecordDao creates and returns a new DAO object for table data access.
func NewRiskRecordDao(handlers ...gdb.ModelHandler) *RiskRecordDao {
	return &RiskRecordDao{
		group:    "default",
		table:    "risk_record",
		columns:  riskRecordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RiskRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RiskRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RiskRecordDao) Columns() RiskRecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RiskRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RiskRecordDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *RiskRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
