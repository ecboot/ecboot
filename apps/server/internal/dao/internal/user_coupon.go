// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserCouponDao is the data access object for the table user_coupon.
type UserCouponDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserCouponColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserCouponColumns defines and stores column names for the table user_coupon.
type UserCouponColumns struct {
	Id         string // 用户券ID(订单user_coupon_id引用此表)
	UserId     string // 持券用户ID
	CouponId   string // 优惠券模板ID
	Status     string // 状态:1未使用 2已使用 3已过期 4已退回(订单取消,可再用)
	ExpireTime string // 过期时间(领取时按模板规则计算落库)
	OrderNo    string // 核销订单号
	UsedTime   string // 核销时间
	CreatedAt  string // 领取时间
	UpdatedAt  string // 更新时间
}

// userCouponColumns holds the columns for the table user_coupon.
var userCouponColumns = UserCouponColumns{
	Id:         "id",
	UserId:     "user_id",
	CouponId:   "coupon_id",
	Status:     "status",
	ExpireTime: "expire_time",
	OrderNo:    "order_no",
	UsedTime:   "used_time",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewUserCouponDao creates and returns a new DAO object for table data access.
func NewUserCouponDao(handlers ...gdb.ModelHandler) *UserCouponDao {
	return &UserCouponDao{
		group:    "default",
		table:    "user_coupon",
		columns:  userCouponColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserCouponDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserCouponDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserCouponDao) Columns() UserCouponColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserCouponDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserCouponDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserCouponDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
