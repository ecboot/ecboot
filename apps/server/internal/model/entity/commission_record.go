// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CommissionRecord is the golang structure for table commission_record.
type CommissionRecord struct {
	Id                uint64      `json:"id"                orm:"id"                  ` // 佣金记录ID
	OrderNo           string      `json:"orderNo"           orm:"order_no"            ` // 订单号
	OrderItemId       uint64      `json:"orderItemId"       orm:"order_item_id"       ` // 订单项ID
	BeneficiaryUserId uint64      `json:"beneficiaryUserId" orm:"beneficiary_user_id" ` // 受益人用户ID
	Level             int         `json:"level"             orm:"level"               ` // 层级:1直接邀请 2间接邀请
	BaseAmount        float64     `json:"baseAmount"        orm:"base_amount"         ` // 计佣基数(=订单项实付pay_amount)
	Rate              float64     `json:"rate"              orm:"rate"                ` // 命中比例%
	Amount            float64     `json:"amount"            orm:"amount"              ` // 佣金金额(冲销记录为负值)
	Status            int         `json:"status"            orm:"status"              ` // 状态:1待结算 2已结算 3已失效(保护期退款) 4欠款冲销中(结算后退款)
	SettleTime        *gtime.Time `json:"settleTime"        orm:"settle_time"         ` // 结算时间(确认收货+7天保护期满)
	ReversalRecordId  uint64      `json:"reversalRecordId"  orm:"reversal_record_id"  ` // 冲销关联:本记录被哪条负额冲销记录回指(原记录侧)
	ReversalOfId      uint64      `json:"reversalOfId"      orm:"reversal_of_id"      ` // 冲销关联:本冲销记录冲销的原记录ID(冲销记录侧)
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"          ` // 创建时间
	UpdatedAt         *gtime.Time `json:"updatedAt"         orm:"updated_at"          ` // 更新时间
}
