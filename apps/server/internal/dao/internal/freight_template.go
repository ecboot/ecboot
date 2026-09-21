// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FreightTemplateDao is the data access object for the table freight_template.
type FreightTemplateDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  FreightTemplateColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// FreightTemplateColumns defines and stores column names for the table freight_template.
type FreightTemplateColumns struct {
	Id               string // 模板ID
	Name             string // 模板名称(如:华东包邮-全国2件10元)
	ChargeType       string // 计费方式:1按件数 2按重量
	FreeThreshold    string // 满额包邮阈值(订单商品金额≥此值免运费,NULL=不包邮)
	FreeExcludeCodes string // 不参与满额包邮的省级区划代码列表(如["650000","540000"]);NULL=全国均可包邮
	Status           string // 状态:1启用 0停用(停用后新订单不可用,已下单不受影响)
	Deleted          string // 软删除:0否 1是
	CreatedAt        string // 创建时间
	UpdatedAt        string // 更新时间
}

// freightTemplateColumns holds the columns for the table freight_template.
var freightTemplateColumns = FreightTemplateColumns{
	Id:               "id",
	Name:             "name",
	ChargeType:       "charge_type",
	FreeThreshold:    "free_threshold",
	FreeExcludeCodes: "free_exclude_codes",
	Status:           "status",
	Deleted:          "deleted",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
}

// NewFreightTemplateDao creates and returns a new DAO object for table data access.
func NewFreightTemplateDao(handlers ...gdb.ModelHandler) *FreightTemplateDao {
	return &FreightTemplateDao{
		group:    "default",
		table:    "freight_template",
		columns:  freightTemplateColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *FreightTemplateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *FreightTemplateDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *FreightTemplateDao) Columns() FreightTemplateColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *FreightTemplateDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *FreightTemplateDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *FreightTemplateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
