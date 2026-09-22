// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TradeOrderDao is the data access object for the table trade_order.
type TradeOrderDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TradeOrderColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TradeOrderColumns defines and stores column names for the table trade_order.
type TradeOrderColumns struct {
	Id                  string // 订单ID
	OrderNo             string // 订单号(业务号:日期+雪花/随机,全局唯一,分片友好)
	UserId              string // 买家用户ID
	OrderChannel        string // 下单渠道:1微信小程序 2H5
	Status              string // 订单状态:10待付款 20待发货 30待收货 40已完成 90已取消
	RefundStatus        string // 退款状态:0无售后 1部分退款 2全额退款(不影响主状态机)
	Currency            string // 币种(ISO 4217,订单项继承此字段)
	TotalAmount         string // 商品总额(原价*数量合计)
	PromotionAmount     string // 优惠总额(优惠券等)
	FreightAmount       string // 运费(按运费模板计算,下单时快照)
	FreightTemplateId   string // 运费模板ID(下单时快照,NULL=包邮)
	PointAmount         string // 积分抵扣金额(promotion_amount构成之一)
	PointUsed           string // 消耗积分点数
	AccountAmount       string // 佣金余额抵扣金额(用户资产消耗,非营销优惠;现金实付=pay_amount-account_amount,渠道单金额口径)
	CouponAmount        string // 优惠券抵扣金额(promotion_amount构成之一)
	FullReductionAmount string // 满减优惠金额(promotion_amount构成之一)
	PromotionActivityId string // 命中的满减活动ID(NULL=未参与)
	PayAmount           string // 实付金额=total_amount-promotion_amount+freight_amount
	UserCouponId        string // 使用的用户优惠券ID(与订单同事务核销)
	AttributedUserId    string // 订单归因人(分享带来本单的用户;佣金按 分享归因>关系链兜底>自然流量 三级判定)
	AttributionType     string // 归因类型:1分享归因 2关系链兜底 3自然流量
	GroupBuyTeamId      string // 拼团团ID(NULL=普通订单)
	BargainRecordId     string // 砍价单ID(NULL=非砍价订单)
	ReceiverName        string // 收货人姓名(下单时快照)
	ReceiverPhone       string // 收货人手机号(快照)
	ReceiverProvince    string // 省(快照,运费区域匹配依据)
	ReceiverCity        string // 市(快照)
	ReceiverDistrict    string // 区/县(快照)
	ReceiverDetail      string // 详细地址(快照)
	UserRemark          string // 买家留言
	SellerRemark        string // 卖家备注(客服/仓库内部使用,买家不可见)
	RequestToken        string // 下单幂等token(确认页发放,Redis抢占+唯一索引兜底;NULL不参与唯一)
	FlashSaleItemId     string // 秒杀场次商品ID(NULL=普通单; 取消时据此回补活动库存)
	PayTime             string // 支付完成时间
	DeliverCompany      string // 物流公司(发货预留)
	DeliverNo           string // 物流单号(发货预留)
	DeliverTime         string // 发货时间(预留)
	FinishTime          string // 订单完成时间(确认收货)
	CancelType          string // 取消方:1用户 2系统超时 3管理员
	CancelReason        string // 取消原因
	CancelTime          string // 取消时间
	CreatedAt           string // 创建时间
	UpdatedAt           string // 更新时间
}

// tradeOrderColumns holds the columns for the table trade_order.
var tradeOrderColumns = TradeOrderColumns{
	Id:                  "id",
	OrderNo:             "order_no",
	UserId:              "user_id",
	OrderChannel:        "order_channel",
	Status:              "status",
	RefundStatus:        "refund_status",
	Currency:            "currency",
	TotalAmount:         "total_amount",
	PromotionAmount:     "promotion_amount",
	FreightAmount:       "freight_amount",
	FreightTemplateId:   "freight_template_id",
	PointAmount:         "point_amount",
	PointUsed:           "point_used",
	AccountAmount:       "account_amount",
	CouponAmount:        "coupon_amount",
	FullReductionAmount: "full_reduction_amount",
	PromotionActivityId: "promotion_activity_id",
	PayAmount:           "pay_amount",
	UserCouponId:        "user_coupon_id",
	AttributedUserId:    "attributed_user_id",
	AttributionType:     "attribution_type",
	GroupBuyTeamId:      "group_buy_team_id",
	BargainRecordId:     "bargain_record_id",
	ReceiverName:        "receiver_name",
	ReceiverPhone:       "receiver_phone",
	ReceiverProvince:    "receiver_province",
	ReceiverCity:        "receiver_city",
	ReceiverDistrict:    "receiver_district",
	ReceiverDetail:      "receiver_detail",
	UserRemark:          "user_remark",
	SellerRemark:        "seller_remark",
	RequestToken:        "request_token",
	FlashSaleItemId:     "flash_sale_item_id",
	PayTime:             "pay_time",
	DeliverCompany:      "deliver_company",
	DeliverNo:           "deliver_no",
	DeliverTime:         "deliver_time",
	FinishTime:          "finish_time",
	CancelType:          "cancel_type",
	CancelReason:        "cancel_reason",
	CancelTime:          "cancel_time",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
}

// NewTradeOrderDao creates and returns a new DAO object for table data access.
func NewTradeOrderDao(handlers ...gdb.ModelHandler) *TradeOrderDao {
	return &TradeOrderDao{
		group:    "default",
		table:    "trade_order",
		columns:  tradeOrderColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TradeOrderDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TradeOrderDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TradeOrderDao) Columns() TradeOrderColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TradeOrderDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TradeOrderDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TradeOrderDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
