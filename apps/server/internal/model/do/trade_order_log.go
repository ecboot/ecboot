// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TradeOrderLog is the golang structure of table trade_order_log for DAO operations like Where/Data.
type TradeOrderLog struct {
	g.Meta       `orm:"table:trade_order_log, do:true"`
	Id           any         // 流水ID
	OrderId      any         // 订单ID
	OrderNo      any         // 订单号
	FromStatus   any         // 迁移前状态(建单时为NULL)
	ToStatus     any         // 迁移后状态
	OperatorType any         // 操作者类型:1系统 2用户 3管理员
	OperatorId   any         // 操作者标识(system/user:{id}/admin:{id})
	Remark       any         // 备注(如超时取消/确认收货)
	CreatedAt    *gtime.Time // 创建时间
}
