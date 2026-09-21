// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserRelationDao is the data access object for the table user_relation.
type UserRelationDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  UserRelationColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// UserRelationColumns defines and stores column names for the table user_relation.
type UserRelationColumns struct {
	Id          string // 关系ID
	UserId      string // 用户ID(一人至多一条关系)
	InviterId   string // 直接上级(邀请人);二级=上级的上级,两次单列查询解析;无祖父列/路径列(ADR-0003)
	BindChannel string // 绑定方式:1分享链接 2邀请码
	BindTime    string // 绑定时间
	Locked      string // 锁定:0保护期内(可换绑) 1已锁定(绑定窗口期满,防抢人纠纷)
	LockTime    string // 锁定时间
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
}

// userRelationColumns holds the columns for the table user_relation.
var userRelationColumns = UserRelationColumns{
	Id:          "id",
	UserId:      "user_id",
	InviterId:   "inviter_id",
	BindChannel: "bind_channel",
	BindTime:    "bind_time",
	Locked:      "locked",
	LockTime:    "lock_time",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewUserRelationDao creates and returns a new DAO object for table data access.
func NewUserRelationDao(handlers ...gdb.ModelHandler) *UserRelationDao {
	return &UserRelationDao{
		group:    "default",
		table:    "user_relation",
		columns:  userRelationColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserRelationDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserRelationDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserRelationDao) Columns() UserRelationColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserRelationDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserRelationDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserRelationDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
