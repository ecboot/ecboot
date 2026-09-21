// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TradeOrder is the golang structure of table trade_order for DAO operations like Where/Data.
type TradeOrder struct {
	g.Meta              `orm:"table:trade_order, do:true"`
	Id                  any         // 订单ID
	OrderNo             any         // 订单号(业务号:日期+雪花/随机,全局唯一,分片友好)
	UserId              any         // 买家用户ID
	SellerId            any         // 订单归属商家(1=自营;V1不拆单,多商家上线时购物车按商家分组结算)
	OrderChannel        any         // 下单渠道:1微信小程序 2H5
	Status              any         // 订单状态:10待付款 20待发货 30待收货 40已完成 90已取消
	RefundStatus        any         // 退款状态:0无售后 1部分退款 2全额退款(不影响主状态机)
	Currency            any         // 币种(ISO 4217,订单项继承此字段)
	TotalAmount         any         // 商品总额(原价*数量合计)
	PromotionAmount     any         // 优惠总额(优惠券等)
	FreightAmount       any         // 运费(按运费模板计算,下单时快照)
	FreightTemplateId   any         // 运费模板ID(下单时快照,NULL=包邮)
	PointAmount         any         // 积分抵扣金额(promotion_amount构成之一)
	PointUsed           any         // 消耗积分点数
	AccountAmount       any         // 佣金余额抵扣金额(用户资产消耗,非营销优惠;现金实付=pay_amount-account_amount,渠道单金额口径)
	CouponAmount        any         // 优惠券抵扣金额(promotion_amount构成之一)
	FullReductionAmount any         // 满减优惠金额(promotion_amount构成之一)
	PromotionActivityId any         // 命中的满减活动ID(NULL=未参与)
	PayAmount           any         // 实付金额=total_amount-promotion_amount+freight_amount
	UserCouponId        any         // 使用的用户优惠券ID(与订单同事务核销)
	AttributedUserId    any         // 订单归因人(分享带来本单的用户;佣金按 分享归因>关系链兜底>自然流量 三级判定)
	AttributionType     any         // 归因类型:1分享归因 2关系链兜底 3自然流量
	GroupBuyTeamId      any         // 拼团团ID(NULL=普通订单)
	BargainRecordId     any         // 砍价单ID(NULL=非砍价订单)
	ReceiverName        any         // 收货人姓名(下单时快照)
	ReceiverPhone       any         // 收货人手机号(快照)
	ReceiverProvince    any         // 省(快照,运费区域匹配依据)
	ReceiverCity        any         // 市(快照)
	ReceiverDistrict    any         // 区/县(快照)
	ReceiverDetail      any         // 详细地址(快照)
	UserRemark          any         // 买家留言
	SellerRemark        any         // 卖家备注(客服/仓库内部使用,买家不可见)
	RequestToken        any         // 下单幂等token(确认页发放,Redis抢占+唯一索引兜底;NULL不参与唯一)
	PayTime             *gtime.Time // 支付完成时间
	DeliverCompany      any         // 物流公司(发货预留)
	DeliverNo           any         // 物流单号(发货预留)
	DeliverTime         *gtime.Time // 发货时间(预留)
	FinishTime          *gtime.Time // 订单完成时间(确认收货)
	CancelType          any         // 取消方:1用户 2系统超时 3管理员
	CancelReason        any         // 取消原因
	CancelTime          *gtime.Time // 取消时间
	CreatedAt           *gtime.Time // 创建时间
	UpdatedAt           *gtime.Time // 更新时间
}
