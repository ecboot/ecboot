// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CommissionRecordDao is the data access object for the table commission_record.
type CommissionRecordDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  CommissionRecordColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// CommissionRecordColumns defines and stores column names for the table commission_record.
type CommissionRecordColumns struct {
	Id                string // 佣金记录ID
	OrderNo           string // 订单号
	OrderItemId       string // 订单项ID
	BeneficiaryUserId string // 受益人用户ID
	Level             string // 层级:1直接邀请 2间接邀请
	BaseAmount        string // 计佣基数(=订单项实付pay_amount)
	Rate              string // 命中比例%
	Amount            string // 佣金金额(冲销记录为负值)
	Status            string // 状态:1待结算 2已结算 3已失效(保护期退款) 4欠款冲销中(结算后退款)
	SettleTime        string // 结算时间(确认收货+7天保护期满)
	ReversalRecordId  string // 冲销关联:本记录被哪条负额冲销记录回指(原记录侧)
	ReversalOfId      string // 冲销关联:本冲销记录冲销的原记录ID(冲销记录侧)
	CreatedAt         string // 创建时间
	UpdatedAt         string // 更新时间
}

// commissionRecordColumns holds the columns for the table commission_record.
var commissionRecordColumns = CommissionRecordColumns{
	Id:                "id",
	OrderNo:           "order_no",
	OrderItemId:       "order_item_id",
	BeneficiaryUserId: "beneficiary_user_id",
	Level:             "level",
	BaseAmount:        "base_amount",
	Rate:              "rate",
	Amount:            "amount",
	Status:            "status",
	SettleTime:        "settle_time",
	ReversalRecordId:  "reversal_record_id",
	ReversalOfId:      "reversal_of_id",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
}

// NewCommissionRecordDao creates and returns a new DAO object for table data access.
func NewCommissionRecordDao(handlers ...gdb.ModelHandler) *CommissionRecordDao {
	return &CommissionRecordDao{
		group:    "default",
		table:    "commission_record",
		columns:  commissionRecordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CommissionRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CommissionRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CommissionRecordDao) Columns() CommissionRecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CommissionRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CommissionRecordDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CommissionRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
