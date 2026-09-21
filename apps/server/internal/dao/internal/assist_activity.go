// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AssistActivityDao is the data access object for the table assist_activity.
type AssistActivityDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  AssistActivityColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// AssistActivityColumns defines and stores column names for the table assist_activity.
type AssistActivityColumns struct {
	Id            string // 助力活动ID
	Name          string // 活动名称
	RewardType    string // 奖励类型:1优惠券 2积分
	RewardRef     string // 奖励载体ID(如coupon.id;积分为点数存config)
	RequiredCount string // 所需助力人数
	PerLimit      string // 每人可发起次数
	StartTime     string // 开始时间(含)
	EndTime       string // 结束时间(不含)
	Config        string // 扩展参数(积分数/阶梯奖励等)
	Status        string // 状态:1启用 0停用
	Deleted       string // 软删除:0否 1是
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
}

// assistActivityColumns holds the columns for the table assist_activity.
var assistActivityColumns = AssistActivityColumns{
	Id:            "id",
	Name:          "name",
	RewardType:    "reward_type",
	RewardRef:     "reward_ref",
	RequiredCount: "required_count",
	PerLimit:      "per_limit",
	StartTime:     "start_time",
	EndTime:       "end_time",
	Config:        "config",
	Status:        "status",
	Deleted:       "deleted",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewAssistActivityDao creates and returns a new DAO object for table data access.
func NewAssistActivityDao(handlers ...gdb.ModelHandler) *AssistActivityDao {
	return &AssistActivityDao{
		group:    "default",
		table:    "assist_activity",
		columns:  assistActivityColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AssistActivityDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AssistActivityDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AssistActivityDao) Columns() AssistActivityColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AssistActivityDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AssistActivityDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AssistActivityDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
