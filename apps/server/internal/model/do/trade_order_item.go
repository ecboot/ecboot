// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TradeOrderItem is the golang structure of table trade_order_item for DAO operations like Where/Data.
type TradeOrderItem struct {
	g.Meta              `orm:"table:trade_order_item, do:true"`
	Id                  any         // 订单项ID
	OrderId             any         // 订单ID
	OrderNo             any         // 订单号(冗余,免联查)
	SpuId               any         // SPU ID(溯源用)
	SkuId               any         // SKU ID(售后/库存回补定位)
	SkuNo               any         // SKU编码(快照)
	SpuName             any         // SPU名称(下单时快照)
	SkuName             any         // SKU名称=名称+规格串(快照)
	SkuImage            any         // SKU主图URL(下单时快照)
	SkuSpecs            any         // 规格组合(快照):{"颜色":"黑","尺码":"M"}
	FlashSaleItemId     any         // 秒杀场次商品ID(NULL=非秒杀;限购校验与秒杀订单识别依据)
	Quantity            any         // 购买数量
	OriginalPrice       any         // 下单时页面价(快照)
	Price               any         // 成交单价(优惠分摊前)
	PromotionAmount     any         // 分摊到本行的优惠(按行金额比例,尾差记末行)
	PointAmount         any         // 积分抵扣行分摊(合计=订单头point_amount,尾差记末行)
	AccountAmount       any         // 余额抵扣行分摊(合计=订单头account_amount,售后按行原路退回依据)
	CouponAmount        any         // 优惠券抵扣行分摊
	FullReductionAmount any         // 满减优惠行分摊(两项合计=订单头对应明细,尾差记末行)
	PayAmount           any         // 本行实付=price*quantity-promotion_amount(币种继承订单头表)
	CreatedAt           *gtime.Time // 创建时间
	UpdatedAt           *gtime.Time // 更新时间(与created_at相等,订单项不可变)
}
