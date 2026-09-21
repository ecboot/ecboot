// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TradeOrderLog is the golang structure for table trade_order_log.
type TradeOrderLog struct {
	Id           uint64      `json:"id"           orm:"id"            ` // 流水ID
	OrderId      uint64      `json:"orderId"      orm:"order_id"      ` // 订单ID
	OrderNo      string      `json:"orderNo"      orm:"order_no"      ` // 订单号
	FromStatus   int         `json:"fromStatus"   orm:"from_status"   ` // 迁移前状态(建单时为NULL)
	ToStatus     int         `json:"toStatus"     orm:"to_status"     ` // 迁移后状态
	OperatorType int         `json:"operatorType" orm:"operator_type" ` // 操作者类型:1系统 2用户 3管理员
	OperatorId   string      `json:"operatorId"   orm:"operator_id"   ` // 操作者标识(system/user:{id}/admin:{id})
	Remark       string      `json:"remark"       orm:"remark"        ` // 备注(如超时取消/确认收货)
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` // 创建时间
}
