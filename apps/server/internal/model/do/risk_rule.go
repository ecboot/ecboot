// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskRule is the golang structure of table risk_rule for DAO operations like Where/Data.
type RiskRule struct {
	g.Meta        `orm:"table:risk_rule, do:true"`
	Id            any         // 规则ID
	Name          any         // 规则名称(如 黑名单/高频下单/佣金套利特征)
	RuleType      any         // 规则类型:1黑名单 2高频下单 3异常领券 4佣金套利特征 5休眠账户分级
	ConditionExpr any         // 规则条件描述(如 1分钟内下单>10次)
	Action        any         // 处置:1拦截 2标记(放行但留痕)
	Status        any         // 状态:1启用 0停用
	Deleted       any         // 软删除:0否 1是
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
}
