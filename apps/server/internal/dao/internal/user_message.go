// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserMessageDao is the data access object for the table user_message.
type UserMessageDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserMessageColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserMessageColumns defines and stores column names for the table user_message.
type UserMessageColumns struct {
	Id        string // 消息ID
	UserId    string // 接收用户ID
	Title     string // 标题
	Content   string // 内容
	BizType   string // 业务类型:1订单 2营销 3售后
	BizNo     string // 关联业务单号(点击跳转)
	IsRead    string // 已读:0否 1是
	ReadTime  string // 阅读时间
	CreatedAt string // 创建时间
}

// userMessageColumns holds the columns for the table user_message.
var userMessageColumns = UserMessageColumns{
	Id:        "id",
	UserId:    "user_id",
	Title:     "title",
	Content:   "content",
	BizType:   "biz_type",
	BizNo:     "biz_no",
	IsRead:    "is_read",
	ReadTime:  "read_time",
	CreatedAt: "created_at",
}

// NewUserMessageDao creates and returns a new DAO object for table data access.
func NewUserMessageDao(handlers ...gdb.ModelHandler) *UserMessageDao {
	return &UserMessageDao{
		group:    "default",
		table:    "user_message",
		columns:  userMessageColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserMessageDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserMessageDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserMessageDao) Columns() UserMessageColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserMessageDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserMessageDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserMessageDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
