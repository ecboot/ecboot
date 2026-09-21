// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GroupBuyTeamDao is the data access object for the table group_buy_team.
type GroupBuyTeamDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  GroupBuyTeamColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// GroupBuyTeamColumns defines and stores column names for the table group_buy_team.
type GroupBuyTeamColumns struct {
	Id           string // 团ID
	ActivityId   string // 活动ID
	LeaderUserId string // 团长用户ID
	Status       string // 状态:1拼团中 2已成团 3已解散(超时/取消)
	MemberCount  string // 当前成员数(应用层与member表同事务维护)
	ExpireTime   string // 成团截止时间(超时扫描依据)
	SuccessTime  string // 成团时间
	CreatedAt    string // 开团时间
	UpdatedAt    string // 更新时间
}

// groupBuyTeamColumns holds the columns for the table group_buy_team.
var groupBuyTeamColumns = GroupBuyTeamColumns{
	Id:           "id",
	ActivityId:   "activity_id",
	LeaderUserId: "leader_user_id",
	Status:       "status",
	MemberCount:  "member_count",
	ExpireTime:   "expire_time",
	SuccessTime:  "success_time",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewGroupBuyTeamDao creates and returns a new DAO object for table data access.
func NewGroupBuyTeamDao(handlers ...gdb.ModelHandler) *GroupBuyTeamDao {
	return &GroupBuyTeamDao{
		group:    "default",
		table:    "group_buy_team",
		columns:  groupBuyTeamColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *GroupBuyTeamDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *GroupBuyTeamDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *GroupBuyTeamDao) Columns() GroupBuyTeamColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *GroupBuyTeamDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *GroupBuyTeamDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *GroupBuyTeamDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
