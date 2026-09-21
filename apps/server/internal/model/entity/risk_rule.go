// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskRule is the golang structure for table risk_rule.
type RiskRule struct {
	Id            uint64      `json:"id"            orm:"id"             ` // 规则ID
	Name          string      `json:"name"          orm:"name"           ` // 规则名称(如 黑名单/高频下单/佣金套利特征)
	RuleType      int         `json:"ruleType"      orm:"rule_type"      ` // 规则类型:1黑名单 2高频下单 3异常领券 4佣金套利特征 5休眠账户分级
	ConditionExpr string      `json:"conditionExpr" orm:"condition_expr" ` // 规则条件描述(如 1分钟内下单>10次)
	Action        int         `json:"action"        orm:"action"         ` // 处置:1拦截 2标记(放行但留痕)
	Status        int         `json:"status"        orm:"status"         ` // 状态:1启用 0停用
	Deleted       int         `json:"deleted"       orm:"deleted"        ` // 软删除:0否 1是
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     ` // 创建时间
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     ` // 更新时间
}
