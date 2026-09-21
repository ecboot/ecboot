// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FreightRuleDao is the data access object for the table freight_rule.
type FreightRuleDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  FreightRuleColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// FreightRuleColumns defines and stores column names for the table freight_rule.
type FreightRuleColumns struct {
	Id           string // 规则ID
	TemplateId   string // 所属模板ID
	RegionCodes  string // 适用省级行政区代码列表(如["310000","330000"];空数组=全国兜底规则,每模板至多一条)
	FirstUnit    string // 首段额度:按件=首件数;按重=首重克数
	FirstFee     string // 首段运费
	ContinueUnit string // 续段步长:按件=每续N件;按重=每续N克
	ContinueFee  string // 续段单价(每步长追加费用)
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
}

// freightRuleColumns holds the columns for the table freight_rule.
var freightRuleColumns = FreightRuleColumns{
	Id:           "id",
	TemplateId:   "template_id",
	RegionCodes:  "region_codes",
	FirstUnit:    "first_unit",
	FirstFee:     "first_fee",
	ContinueUnit: "continue_unit",
	ContinueFee:  "continue_fee",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewFreightRuleDao creates and returns a new DAO object for table data access.
func NewFreightRuleDao(handlers ...gdb.ModelHandler) *FreightRuleDao {
	return &FreightRuleDao{
		group:    "default",
		table:    "freight_rule",
		columns:  freightRuleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *FreightRuleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *FreightRuleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *FreightRuleDao) Columns() FreightRuleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *FreightRuleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *FreightRuleDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *FreightRuleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
