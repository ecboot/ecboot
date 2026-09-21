// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TradeOrderItemDao is the data access object for the table trade_order_item.
type TradeOrderItemDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  TradeOrderItemColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// TradeOrderItemColumns defines and stores column names for the table trade_order_item.
type TradeOrderItemColumns struct {
	Id                  string // 订单项ID
	OrderId             string // 订单ID
	OrderNo             string // 订单号(冗余,免联查)
	SpuId               string // SPU ID(溯源用)
	SkuId               string // SKU ID(售后/库存回补定位)
	SkuNo               string // SKU编码(快照)
	SpuName             string // SPU名称(下单时快照)
	SkuName             string // SKU名称=名称+规格串(快照)
	SkuImage            string // SKU主图URL(下单时快照)
	SkuSpecs            string // 规格组合(快照):{"颜色":"黑","尺码":"M"}
	FlashSaleItemId     string // 秒杀场次商品ID(NULL=非秒杀;限购校验与秒杀订单识别依据)
	Quantity            string // 购买数量
	OriginalPrice       string // 下单时页面价(快照)
	Price               string // 成交单价(优惠分摊前)
	PromotionAmount     string // 分摊到本行的优惠(按行金额比例,尾差记末行)
	PointAmount         string // 积分抵扣行分摊(合计=订单头point_amount,尾差记末行)
	AccountAmount       string // 余额抵扣行分摊(合计=订单头account_amount,售后按行原路退回依据)
	CouponAmount        string // 优惠券抵扣行分摊
	FullReductionAmount string // 满减优惠行分摊(两项合计=订单头对应明细,尾差记末行)
	PayAmount           string // 本行实付=price*quantity-promotion_amount(币种继承订单头表)
	CreatedAt           string // 创建时间
	UpdatedAt           string // 更新时间(与created_at相等,订单项不可变)
}

// tradeOrderItemColumns holds the columns for the table trade_order_item.
var tradeOrderItemColumns = TradeOrderItemColumns{
	Id:                  "id",
	OrderId:             "order_id",
	OrderNo:             "order_no",
	SpuId:               "spu_id",
	SkuId:               "sku_id",
	SkuNo:               "sku_no",
	SpuName:             "spu_name",
	SkuName:             "sku_name",
	SkuImage:            "sku_image",
	SkuSpecs:            "sku_specs",
	FlashSaleItemId:     "flash_sale_item_id",
	Quantity:            "quantity",
	OriginalPrice:       "original_price",
	Price:               "price",
	PromotionAmount:     "promotion_amount",
	PointAmount:         "point_amount",
	AccountAmount:       "account_amount",
	CouponAmount:        "coupon_amount",
	FullReductionAmount: "full_reduction_amount",
	PayAmount:           "pay_amount",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
}

// NewTradeOrderItemDao creates and returns a new DAO object for table data access.
func NewTradeOrderItemDao(handlers ...gdb.ModelHandler) *TradeOrderItemDao {
	return &TradeOrderItemDao{
		group:    "default",
		table:    "trade_order_item",
		columns:  tradeOrderItemColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TradeOrderItemDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TradeOrderItemDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TradeOrderItemDao) Columns() TradeOrderItemColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TradeOrderItemDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TradeOrderItemDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TradeOrderItemDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
