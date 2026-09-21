// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CommissionRecord is the golang structure of table commission_record for DAO operations like Where/Data.
type CommissionRecord struct {
	g.Meta            `orm:"table:commission_record, do:true"`
	Id                any         // 佣金记录ID
	OrderNo           any         // 订单号
	OrderItemId       any         // 订单项ID
	BeneficiaryUserId any         // 受益人用户ID
	Level             any         // 层级:1直接邀请 2间接邀请
	BaseAmount        any         // 计佣基数(=订单项实付pay_amount)
	Rate              any         // 命中比例%
	Amount            any         // 佣金金额(冲销记录为负值)
	Status            any         // 状态:1待结算 2已结算 3已失效(保护期退款) 4欠款冲销中(结算后退款)
	SettleTime        *gtime.Time // 结算时间(确认收货+7天保护期满)
	ReversalRecordId  any         // 冲销关联:本记录被哪条负额冲销记录回指(原记录侧)
	ReversalOfId      any         // 冲销关联:本冲销记录冲销的原记录ID(冲销记录侧)
	CreatedAt         *gtime.Time // 创建时间
	UpdatedAt         *gtime.Time // 更新时间
}
