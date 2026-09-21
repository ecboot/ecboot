// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LogisticsCompanyDao is the data access object for the table logistics_company.
type LogisticsCompanyDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  LogisticsCompanyColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// LogisticsCompanyColumns defines and stores column names for the table logistics_company.
type LogisticsCompanyColumns struct {
	Id           string // 物流公司ID
	Code         string // 编码(唯一,订单deliver_company存此值)
	Name         string // 公司名称(如 顺丰速运)
	TrackingRule string // 运单号校验规则描述(如 SF+12位数字)
	Sort         string // 排序
	Status       string // 状态:1启用 0停用
	Deleted      string // 软删除:0否 1是
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
}

// logisticsCompanyColumns holds the columns for the table logistics_company.
var logisticsCompanyColumns = LogisticsCompanyColumns{
	Id:           "id",
	Code:         "code",
	Name:         "name",
	TrackingRule: "tracking_rule",
	Sort:         "sort",
	Status:       "status",
	Deleted:      "deleted",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewLogisticsCompanyDao creates and returns a new DAO object for table data access.
func NewLogisticsCompanyDao(handlers ...gdb.ModelHandler) *LogisticsCompanyDao {
	return &LogisticsCompanyDao{
		group:    "default",
		table:    "logistics_company",
		columns:  logisticsCompanyColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LogisticsCompanyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LogisticsCompanyDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LogisticsCompanyDao) Columns() LogisticsCompanyColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LogisticsCompanyDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LogisticsCompanyDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *LogisticsCompanyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
