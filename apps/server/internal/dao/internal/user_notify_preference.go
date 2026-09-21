// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserNotifyPreferenceDao is the data access object for the table user_notify_preference.
type UserNotifyPreferenceDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  UserNotifyPreferenceColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// UserNotifyPreferenceColumns defines and stores column names for the table user_notify_preference.
type UserNotifyPreferenceColumns struct {
	Id        string // 偏好ID
	UserId    string // 会员ID
	Channel   string // 渠道:1小程序订阅 2短信
	Enabled   string // 是否接收:1接收 0关闭
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// userNotifyPreferenceColumns holds the columns for the table user_notify_preference.
var userNotifyPreferenceColumns = UserNotifyPreferenceColumns{
	Id:        "id",
	UserId:    "user_id",
	Channel:   "channel",
	Enabled:   "enabled",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewUserNotifyPreferenceDao creates and returns a new DAO object for table data access.
func NewUserNotifyPreferenceDao(handlers ...gdb.ModelHandler) *UserNotifyPreferenceDao {
	return &UserNotifyPreferenceDao{
		group:    "default",
		table:    "user_notify_preference",
		columns:  userNotifyPreferenceColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserNotifyPreferenceDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserNotifyPreferenceDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserNotifyPreferenceDao) Columns() UserNotifyPreferenceColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserNotifyPreferenceDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserNotifyPreferenceDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserNotifyPreferenceDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
