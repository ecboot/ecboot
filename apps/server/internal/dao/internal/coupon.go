// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CouponDao is the data access object for the table coupon.
type CouponDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  CouponColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// CouponColumns defines and stores column names for the table coupon.
type CouponColumns struct {
	Id              string // 优惠券模板ID
	Name            string // 券名称(如:满100减20)
	Type            string // 券类型:1满减券 2无门槛券 3折扣券(预留)
	ThresholdAmount string // 使用门槛(满X元可用,0=无门槛)
	DiscountAmount  string // 抵扣金额(满减/无门槛券)
	DiscountRate    string // 折扣率0.01-0.99(折扣券预留,如0.85=85折)
	TotalCount      string // 发放总量,0=不限量
	ReceivedCount   string // 已领取数量(领券事务内条件更新,防超发)
	PerLimit        string // 每人限领数量
	ValidType       string // 有效期方式:1固定区间 2领取后N天
	ValidStartAt    string // 固定区间开始(valid_type=1)
	ValidEndAt      string // 固定区间结束(valid_type=1)
	ValidDays       string // 领取后N天有效(valid_type=2)
	Status          string // 状态:1启用 0停用(停用不影响已领取的券)
	Deleted         string // 软删除:0否 1是
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
}

// couponColumns holds the columns for the table coupon.
var couponColumns = CouponColumns{
	Id:              "id",
	Name:            "name",
	Type:            "type",
	ThresholdAmount: "threshold_amount",
	DiscountAmount:  "discount_amount",
	DiscountRate:    "discount_rate",
	TotalCount:      "total_count",
	ReceivedCount:   "received_count",
	PerLimit:        "per_limit",
	ValidType:       "valid_type",
	ValidStartAt:    "valid_start_at",
	ValidEndAt:      "valid_end_at",
	ValidDays:       "valid_days",
	Status:          "status",
	Deleted:         "deleted",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewCouponDao creates and returns a new DAO object for table data access.
func NewCouponDao(handlers ...gdb.ModelHandler) *CouponDao {
	return &CouponDao{
		group:    "default",
		table:    "coupon",
		columns:  couponColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CouponDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CouponDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CouponDao) Columns() CouponColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CouponDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CouponDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CouponDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
