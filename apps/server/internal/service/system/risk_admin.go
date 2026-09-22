// risk_admin.go 风控管理（admin, 018 批次 12）——接口定义 + 评估器。
// 规则: 规则/事件全部条件更新判行数; 申诉四态（0无 1申诉中 2通过 3驳回）已结论拒绝重放;
// 评估器实现 shop.IRiskHit 端口（批次 09/11 降级放行点就此闭合为真实判定）,
// 评估失败不阻断主流程（降级放行+告警——批次 09 语义延续）。
package system

import (
	"context"

	"ecboot/internal/model"
)

// AdminRiskRuleItem 后台风控规则项。
type AdminRiskRuleItem struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	RuleType      int    `json:"ruleType" dc:"1黑名单 2高频下单 3异常领券 4佣金套利 5休眠分级"`
	ConditionExpr string `json:"conditionExpr"`
	Action        int    `json:"action" dc:"1拦截 2标记"`
	Status        int    `json:"status"`
}

// AdminRiskRecordItem 后台风控事件项。
type AdminRiskRecordItem struct {
	Id           int64  `json:"id"`
	UserId       int64  `json:"userId"`
	RuleName     string `json:"ruleName" dc:"命中规则"`
	ObjectType   int    `json:"objectType"`
	ObjectNo     string `json:"objectNo"`
	Action       int    `json:"action"`
	AppealStatus int    `json:"appealStatus" dc:"0无 1申诉中 2通过 3驳回"`
	Remark       string `json:"remark"`
	CreatedAt    string `json:"createdAt"`
}

// IRiskAdminLogic 风控管理（admin）。
type IRiskAdminLogic interface {
	AdminRuleList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[AdminRiskRuleItem], error)
	AdminRuleCreate(ctx context.Context, name string, ruleType int, conditionExpr string, action int) (int64, error)
	AdminRuleUpdate(ctx context.Context, id int64, name, conditionExpr string, action int, status *int) error
	AdminRuleDelete(ctx context.Context, id int64) error
	AdminRuleDetail(ctx context.Context, id int64) (*AdminRiskRuleItem, error)
	AdminRecordList(ctx context.Context, userId int64, appealStatus int, page model.PageReq) (*model.PageResult[AdminRiskRecordItem], error)
	// Appeal 申诉处理（2通过/3驳回; 已有结论拒绝——条件更新）。
	AdminAppeal(ctx context.Context, id int64, pass bool, remark string) error
}
