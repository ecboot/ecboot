// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AfterSaleOrderDao is the data access object for the table after_sale_order.
type AfterSaleOrderDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  AfterSaleOrderColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// AfterSaleOrderColumns defines and stores column names for the table after_sale_order.
type AfterSaleOrderColumns struct {
	Id                string // 售后单ID
	AfterSaleNo       string // 售后单号(全局唯一,兼作渠道退款out_refund_no幂等键)
	OrderId           string // 订单ID
	OrderNo           string // 订单号
	OrderItemId       string // 订单项ID(售后粒度=订单项)
	UserId            string // 申请人用户ID
	Type              string // 售后类型:1仅退款 2退货退款 3换货(预留未启用)
	Status            string // 状态:10待审核 20待买家寄回 30待退款 40退款中 50已完成 90已拒绝 91已撤销
	Currency          string // 币种(ISO 4217,随订单)
	Quantity          string // 退货/退款数量(≤订单项quantity,应用校验;部分退货的金额与回补依据)
	Reason            string // 售后原因(不想要了/质量问题/发错货等)
	Description       string // 问题描述
	VoucherImages     string // 凭证图片URL数组
	RefundAmount      string // 退款金额(≤订单项pay_amount)
	ReturnLogisticsNo string // 买家寄回物流单号(type=2)
	RefundNo          string // 渠道退款单号(微信退款ID,回填)
	RejectReason      string // 拒绝原因
	OperatorId        string // 最后操作人标识:admin:{id}(审核/确认收货/退款重试)
	FailReason        string // 渠道退款失败原因(重试前保留,成功后清空)
	AuditTime         string // 审核时间
	RefundTime        string // 退款完成时间
	CreatedAt         string // 创建时间
	UpdatedAt         string // 更新时间
}

// afterSaleOrderColumns holds the columns for the table after_sale_order.
var afterSaleOrderColumns = AfterSaleOrderColumns{
	Id:                "id",
	AfterSaleNo:       "after_sale_no",
	OrderId:           "order_id",
	OrderNo:           "order_no",
	OrderItemId:       "order_item_id",
	UserId:            "user_id",
	Type:              "type",
	Status:            "status",
	Currency:          "currency",
	Quantity:          "quantity",
	Reason:            "reason",
	Description:       "description",
	VoucherImages:     "voucher_images",
	RefundAmount:      "refund_amount",
	ReturnLogisticsNo: "return_logistics_no",
	RefundNo:          "refund_no",
	RejectReason:      "reject_reason",
	OperatorId:        "operator_id",
	FailReason:        "fail_reason",
	AuditTime:         "audit_time",
	RefundTime:        "refund_time",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
}

// NewAfterSaleOrderDao creates and returns a new DAO object for table data access.
func NewAfterSaleOrderDao(handlers ...gdb.ModelHandler) *AfterSaleOrderDao {
	return &AfterSaleOrderDao{
		group:    "default",
		table:    "after_sale_order",
		columns:  afterSaleOrderColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AfterSaleOrderDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AfterSaleOrderDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AfterSaleOrderDao) Columns() AfterSaleOrderColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AfterSaleOrderDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AfterSaleOrderDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AfterSaleOrderDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
