// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GroupBuyTeamMemberDao is the data access object for the table group_buy_team_member.
type GroupBuyTeamMemberDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  GroupBuyTeamMemberColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// GroupBuyTeamMemberColumns defines and stores column names for the table group_buy_team_member.
type GroupBuyTeamMemberColumns struct {
	Id        string // 团员ID
	TeamId    string // 团ID
	UserId    string // 成员用户ID
	OrderNo   string // 成员订单号
	JoinTime  string // 参团时间
	CreatedAt string // 创建时间
}

// groupBuyTeamMemberColumns holds the columns for the table group_buy_team_member.
var groupBuyTeamMemberColumns = GroupBuyTeamMemberColumns{
	Id:        "id",
	TeamId:    "team_id",
	UserId:    "user_id",
	OrderNo:   "order_no",
	JoinTime:  "join_time",
	CreatedAt: "created_at",
}

// NewGroupBuyTeamMemberDao creates and returns a new DAO object for table data access.
func NewGroupBuyTeamMemberDao(handlers ...gdb.ModelHandler) *GroupBuyTeamMemberDao {
	return &GroupBuyTeamMemberDao{
		group:    "default",
		table:    "group_buy_team_member",
		columns:  groupBuyTeamMemberColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *GroupBuyTeamMemberDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *GroupBuyTeamMemberDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *GroupBuyTeamMemberDao) Columns() GroupBuyTeamMemberColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *GroupBuyTeamMemberDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *GroupBuyTeamMemberDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *GroupBuyTeamMemberDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
