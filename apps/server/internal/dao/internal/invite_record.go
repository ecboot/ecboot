// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// InviteRecordDao is the data access object for the table invite_record.
type InviteRecordDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  InviteRecordColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// InviteRecordColumns defines and stores column names for the table invite_record.
type InviteRecordColumns struct {
	Id            string // 激励记录ID
	NewUserId     string // 新用户ID(仅可被激励一次)
	InviterId     string // 邀请人用户ID
	RewardType    string // 奖励类型:1优惠券
	RewardTrigger string // 奖励触发时机:1注册即发 2首单后发(防刷,需订单回调触发)
	RewardRef     string // 奖励载体ID(如user_coupon.id)
	Status        string // 状态:1已发放(发放与订单同事务,幂等)
	CreatedAt     string // 创建时间
}

// inviteRecordColumns holds the columns for the table invite_record.
var inviteRecordColumns = InviteRecordColumns{
	Id:            "id",
	NewUserId:     "new_user_id",
	InviterId:     "inviter_id",
	RewardType:    "reward_type",
	RewardTrigger: "reward_trigger",
	RewardRef:     "reward_ref",
	Status:        "status",
	CreatedAt:     "created_at",
}

// NewInviteRecordDao creates and returns a new DAO object for table data access.
func NewInviteRecordDao(handlers ...gdb.ModelHandler) *InviteRecordDao {
	return &InviteRecordDao{
		group:    "default",
		table:    "invite_record",
		columns:  inviteRecordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *InviteRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *InviteRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *InviteRecordDao) Columns() InviteRecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *InviteRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *InviteRecordDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *InviteRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
