package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 风控规则列表
	AdminRiskRuleListReq struct {
		g.Meta  `path:"/risk-rules" method:"GET" summary:"风控规则列表"`
		Status  int `json:"status" dc:"状态筛选"`
		PageReq
	}
	AdminRiskRuleItem struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		RuleType      int    `json:"ruleType" dc:"1黑名单 2高频下单 3异常领券 4佣金套利 5休眠分级"`
		ConditionExpr string `json:"conditionExpr" dc:"条件描述"`
		Action        int    `json:"action" dc:"1拦截 2标记"`
		Status        int    `json:"status"`
	}
	AdminRiskRuleListRes struct {
		PageRes
		List []AdminRiskRuleItem `json:"list"`
	}

	AdminRiskRuleCreateReq struct {
		g.Meta        `path:"/risk-rules" method:"POST" summary:"新增风控规则"`
		Name          string `json:"name" v:"required" dc:"名称"`
		RuleType      int    `json:"ruleType" v:"required|in:1,2,3,4,5" dc:"类型"`
		ConditionExpr string `json:"conditionExpr" v:"required" dc:"条件描述"`
		Action        int    `json:"action" dc:"处置"`
	}
	AdminRiskRuleCreateRes struct {
		Id string `json:"id"`
	}

	AdminRiskRuleUpdateReq struct {
		g.Meta        `path:"/risk-rules/{id}" method:"PUT" summary:"修改风控规则"`
		Id            string `json:"id" v:"required" dc:"规则ID"`
		Name          string `json:"name" dc:"名称"`
		ConditionExpr string `json:"conditionExpr" dc:"条件描述"`
		Action        int    `json:"action" dc:"处置"`
		Status        int    `json:"status" dc:"状态"`
	}
	AdminRiskRuleUpdateRes struct {
		Success bool `json:"success"`
	}

	AdminRiskRuleDeleteReq struct {
		g.Meta `path:"/risk-rules/{id}" method:"DELETE" summary:"删除风控规则(软删)"`
		Id     string `json:"id" v:"required" dc:"规则ID"`
	}
	AdminRiskRuleDeleteRes struct {
		Success bool `json:"success"`
	}

	// 风控事件查询
	AdminRiskRecordListReq struct {
		g.Meta       `path:"/risk-records" method:"GET" summary:"风控事件列表"`
		UserId       string `json:"userId" dc:"命中用户"`
		AppealStatus int    `json:"appealStatus" dc:"申诉状态筛选"`
		PageReq
	}
	AdminRiskRecordItem struct {
		Id           string `json:"id"`
		UserId       string `json:"userId"`
		RuleName     string `json:"ruleName" dc:"命中规则"`
		ObjectType   int    `json:"objectType" dc:"1订单 2券 3提现 4售后"`
		ObjectNo     string `json:"objectNo" dc:"对象单号"`
		Action       int    `json:"action" dc:"1拦截 2标记"`
		AppealStatus int    `json:"appealStatus" dc:"0无 1申诉中 2通过 3驳回"`
		Remark       string `json:"remark"`
		CreatedAt    string `json:"createdAt"`
	}
	AdminRiskRecordListRes struct {
		PageRes
		List []AdminRiskRecordItem `json:"list"`
	}

	// 申诉处理
	AdminRiskAppealReq struct {
		g.Meta `path:"/risk-records/{id}/appeal" method:"POST" summary:"申诉处理"`
		Id     string `json:"id" v:"required" dc:"事件ID"`
		Pass   bool   `json:"pass" dc:"通过/驳回"`
		Remark string `json:"remark" dc:"处理说明"`
	}
	AdminRiskAppealRes struct {
		Success bool `json:"success"`
	}
)
