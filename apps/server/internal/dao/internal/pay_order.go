// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PayOrderDao is the data access object for the table pay_order.
type PayOrderDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  PayOrderColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// PayOrderColumns defines and stores column names for the table pay_order.
type PayOrderColumns struct {
	Id             string // 支付单ID
	PayNo          string // 支付单号(全局唯一,传给支付渠道的商户单号)
	OrderId        string // 订单ID
	OrderNo        string // 订单号
	UserId         string // 付款用户ID
	PayChannel     string // 支付渠道:1微信支付 2支付宝(预留)
	Amount         string // 应付金额(必须等于订单pay_amount,回调时校验)
	Currency       string // 币种(ISO 4217,须与订单一致)
	Status         string // 支付状态:10待支付 20支付成功 30支付失败 90已关闭
	ChannelTradeNo string // 第三方交易号(微信/支付宝,成功回填;唯一防重放)
	SuccessTime    string // 支付成功时间(毫秒精度,渠道回调为准)
	ExpireTime     string // 支付截止时间(超时未付则关单)
	ClosedTime     string // 关闭时间
	FailReason     string // 失败原因(渠道返回)
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
}

// payOrderColumns holds the columns for the table pay_order.
var payOrderColumns = PayOrderColumns{
	Id:             "id",
	PayNo:          "pay_no",
	OrderId:        "order_id",
	OrderNo:        "order_no",
	UserId:         "user_id",
	PayChannel:     "pay_channel",
	Amount:         "amount",
	Currency:       "currency",
	Status:         "status",
	ChannelTradeNo: "channel_trade_no",
	SuccessTime:    "success_time",
	ExpireTime:     "expire_time",
	ClosedTime:     "closed_time",
	FailReason:     "fail_reason",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewPayOrderDao creates and returns a new DAO object for table data access.
func NewPayOrderDao(handlers ...gdb.ModelHandler) *PayOrderDao {
	return &PayOrderDao{
		group:    "default",
		table:    "pay_order",
		columns:  payOrderColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PayOrderDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PayOrderDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PayOrderDao) Columns() PayOrderColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PayOrderDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PayOrderDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PayOrderDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
