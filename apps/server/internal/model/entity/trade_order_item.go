// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TradeOrderItem is the golang structure for table trade_order_item.
type TradeOrderItem struct {
	Id                  uint64      `json:"id"                  orm:"id"                    ` // 订单项ID
	OrderId             uint64      `json:"orderId"             orm:"order_id"              ` // 订单ID
	OrderNo             string      `json:"orderNo"             orm:"order_no"              ` // 订单号(冗余,免联查)
	SpuId               uint64      `json:"spuId"               orm:"spu_id"                ` // SPU ID(溯源用)
	SkuId               uint64      `json:"skuId"               orm:"sku_id"                ` // SKU ID(售后/库存回补定位)
	SkuNo               string      `json:"skuNo"               orm:"sku_no"                ` // SKU编码(快照)
	SpuName             string      `json:"spuName"             orm:"spu_name"              ` // SPU名称(下单时快照)
	SkuName             string      `json:"skuName"             orm:"sku_name"              ` // SKU名称=名称+规格串(快照)
	SkuImage            string      `json:"skuImage"            orm:"sku_image"             ` // SKU主图URL(下单时快照)
	SkuSpecs            string      `json:"skuSpecs"            orm:"sku_specs"             ` // 规格组合(快照):{"颜色":"黑","尺码":"M"}
	FlashSaleItemId     uint64      `json:"flashSaleItemId"     orm:"flash_sale_item_id"    ` // 秒杀场次商品ID(NULL=非秒杀;限购校验与秒杀订单识别依据)
	Quantity            uint        `json:"quantity"            orm:"quantity"              ` // 购买数量
	OriginalPrice       float64     `json:"originalPrice"       orm:"original_price"        ` // 下单时页面价(快照)
	Price               float64     `json:"price"               orm:"price"                 ` // 成交单价(优惠分摊前)
	PromotionAmount     float64     `json:"promotionAmount"     orm:"promotion_amount"      ` // 分摊到本行的优惠(按行金额比例,尾差记末行)
	PointAmount         float64     `json:"pointAmount"         orm:"point_amount"          ` // 积分抵扣行分摊(合计=订单头point_amount,尾差记末行)
	AccountAmount       float64     `json:"accountAmount"       orm:"account_amount"        ` // 余额抵扣行分摊(合计=订单头account_amount,售后按行原路退回依据)
	CouponAmount        float64     `json:"couponAmount"        orm:"coupon_amount"         ` // 优惠券抵扣行分摊
	FullReductionAmount float64     `json:"fullReductionAmount" orm:"full_reduction_amount" ` // 满减优惠行分摊(两项合计=订单头对应明细,尾差记末行)
	PayAmount           float64     `json:"payAmount"           orm:"pay_amount"            ` // 本行实付=price*quantity-promotion_amount(币种继承订单头表)
	CreatedAt           *gtime.Time `json:"createdAt"           orm:"created_at"            ` // 创建时间
	UpdatedAt           *gtime.Time `json:"updatedAt"           orm:"updated_at"            ` // 更新时间(与created_at相等,订单项不可变)
}
