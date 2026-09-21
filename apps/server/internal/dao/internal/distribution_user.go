// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DistributionUserDao is the data access object for the table distribution_user.
type DistributionUserDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  DistributionUserColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// DistributionUserColumns defines and stores column names for the table distribution_user.
type DistributionUserColumns struct {
	Id        string // 推广员ID
	UserId    string // 用户ID
	Status    string // 状态:1待审核 2通过 3冻结(冻结不产生新佣金,存量可提现)
	Level     string // 推广员等级(V1单一等级=1;多级启用阈值见system_config: distribution.level.threshold)
	ApplyTime string // 申请时间
	AuditTime string // 审核时间
	Deleted   string // 软删除:0否 1是
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// distributionUserColumns holds the columns for the table distribution_user.
var distributionUserColumns = DistributionUserColumns{
	Id:        "id",
	UserId:    "user_id",
	Status:    "status",
	Level:     "level",
	ApplyTime: "apply_time",
	AuditTime: "audit_time",
	Deleted:   "deleted",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewDistributionUserDao creates and returns a new DAO object for table data access.
func NewDistributionUserDao(handlers ...gdb.ModelHandler) *DistributionUserDao {
	return &DistributionUserDao{
		group:    "default",
		table:    "distribution_user",
		columns:  distributionUserColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DistributionUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DistributionUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DistributionUserDao) Columns() DistributionUserColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DistributionUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DistributionUserDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DistributionUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
