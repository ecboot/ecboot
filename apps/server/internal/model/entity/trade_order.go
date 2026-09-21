// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TradeOrder is the golang structure for table trade_order.
type TradeOrder struct {
	Id                  uint64      `json:"id"                  orm:"id"                    ` // 订单ID
	OrderNo             string      `json:"orderNo"             orm:"order_no"              ` // 订单号(业务号:日期+雪花/随机,全局唯一,分片友好)
	UserId              uint64      `json:"userId"              orm:"user_id"               ` // 买家用户ID
	OrderChannel        int         `json:"orderChannel"        orm:"order_channel"         ` // 下单渠道:1微信小程序 2H5
	Status              int         `json:"status"              orm:"status"                ` // 订单状态:10待付款 20待发货 30待收货 40已完成 90已取消
	RefundStatus        int         `json:"refundStatus"        orm:"refund_status"         ` // 退款状态:0无售后 1部分退款 2全额退款(不影响主状态机)
	Currency            string      `json:"currency"            orm:"currency"              ` // 币种(ISO 4217,订单项继承此字段)
	TotalAmount         float64     `json:"totalAmount"         orm:"total_amount"          ` // 商品总额(原价*数量合计)
	PromotionAmount     float64     `json:"promotionAmount"     orm:"promotion_amount"      ` // 优惠总额(优惠券等)
	FreightAmount       float64     `json:"freightAmount"       orm:"freight_amount"        ` // 运费(按运费模板计算,下单时快照)
	FreightTemplateId   uint64      `json:"freightTemplateId"   orm:"freight_template_id"   ` // 运费模板ID(下单时快照,NULL=包邮)
	PointAmount         float64     `json:"pointAmount"         orm:"point_amount"          ` // 积分抵扣金额(promotion_amount构成之一)
	PointUsed           uint        `json:"pointUsed"           orm:"point_used"            ` // 消耗积分点数
	AccountAmount       float64     `json:"accountAmount"       orm:"account_amount"        ` // 佣金余额抵扣金额(用户资产消耗,非营销优惠;现金实付=pay_amount-account_amount,渠道单金额口径)
	CouponAmount        float64     `json:"couponAmount"        orm:"coupon_amount"         ` // 优惠券抵扣金额(promotion_amount构成之一)
	FullReductionAmount float64     `json:"fullReductionAmount" orm:"full_reduction_amount" ` // 满减优惠金额(promotion_amount构成之一)
	PromotionActivityId uint64      `json:"promotionActivityId" orm:"promotion_activity_id" ` // 命中的满减活动ID(NULL=未参与)
	PayAmount           float64     `json:"payAmount"           orm:"pay_amount"            ` // 实付金额=total_amount-promotion_amount+freight_amount
	UserCouponId        uint64      `json:"userCouponId"        orm:"user_coupon_id"        ` // 使用的用户优惠券ID(与订单同事务核销)
	AttributedUserId    uint64      `json:"attributedUserId"    orm:"attributed_user_id"    ` // 订单归因人(分享带来本单的用户;佣金按 分享归因>关系链兜底>自然流量 三级判定)
	AttributionType     int         `json:"attributionType"     orm:"attribution_type"      ` // 归因类型:1分享归因 2关系链兜底 3自然流量
	GroupBuyTeamId      uint64      `json:"groupBuyTeamId"      orm:"group_buy_team_id"     ` // 拼团团ID(NULL=普通订单)
	BargainRecordId     uint64      `json:"bargainRecordId"     orm:"bargain_record_id"     ` // 砍价单ID(NULL=非砍价订单)
	ReceiverName        string      `json:"receiverName"        orm:"receiver_name"         ` // 收货人姓名(下单时快照)
	ReceiverPhone       string      `json:"receiverPhone"       orm:"receiver_phone"        ` // 收货人手机号(快照)
	ReceiverProvince    string      `json:"receiverProvince"    orm:"receiver_province"     ` // 省(快照,运费区域匹配依据)
	ReceiverCity        string      `json:"receiverCity"        orm:"receiver_city"         ` // 市(快照)
	ReceiverDistrict    string      `json:"receiverDistrict"    orm:"receiver_district"     ` // 区/县(快照)
	ReceiverDetail      string      `json:"receiverDetail"      orm:"receiver_detail"       ` // 详细地址(快照)
	UserRemark          string      `json:"userRemark"          orm:"user_remark"           ` // 买家留言
	SellerRemark        string      `json:"sellerRemark"        orm:"seller_remark"         ` // 卖家备注(客服/仓库内部使用,买家不可见)
	RequestToken        string      `json:"requestToken"        orm:"request_token"         ` // 下单幂等token(确认页发放,Redis抢占+唯一索引兜底;NULL不参与唯一)
	PayTime             *gtime.Time `json:"payTime"             orm:"pay_time"              ` // 支付完成时间
	DeliverCompany      string      `json:"deliverCompany"      orm:"deliver_company"       ` // 物流公司(发货预留)
	DeliverNo           string      `json:"deliverNo"           orm:"deliver_no"            ` // 物流单号(发货预留)
	DeliverTime         *gtime.Time `json:"deliverTime"         orm:"deliver_time"          ` // 发货时间(预留)
	FinishTime          *gtime.Time `json:"finishTime"          orm:"finish_time"           ` // 订单完成时间(确认收货)
	CancelType          int         `json:"cancelType"          orm:"cancel_type"           ` // 取消方:1用户 2系统超时 3管理员
	CancelReason        string      `json:"cancelReason"        orm:"cancel_reason"         ` // 取消原因
	CancelTime          *gtime.Time `json:"cancelTime"          orm:"cancel_time"           ` // 取消时间
	CreatedAt           *gtime.Time `json:"createdAt"           orm:"created_at"            ` // 创建时间
	UpdatedAt           *gtime.Time `json:"updatedAt"           orm:"updated_at"            ` // 更新时间
}
