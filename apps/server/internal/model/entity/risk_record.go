// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskRecord is the golang structure for table risk_record.
type RiskRecord struct {
	Id           uint64      `json:"id"           orm:"id"            ` // 事件ID
	UserId       uint64      `json:"userId"       orm:"user_id"       ` // 命中用户ID
	RuleId       uint64      `json:"ruleId"       orm:"rule_id"       ` // 命中规则ID
	ObjectType   int         `json:"objectType"   orm:"object_type"   ` // 关联对象类型:1订单 2优惠券 3提现 4售后
	ObjectNo     string      `json:"objectNo"     orm:"object_no"     ` // 关联对象单号(订单号/券ID/提现单号等)
	Action       int         `json:"action"       orm:"action"        ` // 处置结果:1拦截 2标记
	AppealStatus int         `json:"appealStatus" orm:"appeal_status" ` // 申诉状态:0无 1申诉中 2申诉通过(解除拦截) 3申诉驳回
	Remark       string      `json:"remark"       orm:"remark"        ` // 备注(命中明细/申诉结论)
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` // 命中时间
}
