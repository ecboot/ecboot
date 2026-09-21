// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskRecord is the golang structure of table risk_record for DAO operations like Where/Data.
type RiskRecord struct {
	g.Meta       `orm:"table:risk_record, do:true"`
	Id           any         // 事件ID
	UserId       any         // 命中用户ID
	RuleId       any         // 命中规则ID
	ObjectType   any         // 关联对象类型:1订单 2优惠券 3提现 4售后
	ObjectNo     any         // 关联对象单号(订单号/券ID/提现单号等)
	Action       any         // 处置结果:1拦截 2标记
	AppealStatus any         // 申诉状态:0无 1申诉中 2申诉通过(解除拦截) 3申诉驳回
	Remark       any         // 备注(命中明细/申诉结论)
	CreatedAt    *gtime.Time // 命中时间
}
