package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 发起助力（会员; 受活动次数限制）
	AssistLaunchReq struct {
		g.Meta     `path:"/assists" method:"POST" summary:"发起助力"`
		ActivityId string `json:"activityId" v:"required" dc:"助力活动ID"`
	}
	AssistLaunchRes struct {
		RecordId string `json:"recordId" dc:"参与记录ID"`
	}

	// 助力进度（公开）
	AssistProgressReq struct {
		g.Meta   `path:"/assists/{recordId}" method:"GET" summary:"助力进度"`
		RecordId string `json:"recordId" v:"required" dc:"参与记录ID"`
	}
	AssistProgressRes struct {
		RecordId      string         `json:"recordId"`
		ActivityId    string         `json:"activityId"`
		HelperCount   int            `json:"helperCount" dc:"已助力人数"`
		RequiredCount int            `json:"requiredCount" dc:"所需人数"`
		Status        int            `json:"status" dc:"1进行中 2已完成发奖 3已过期"`
		FinishTime    string         `json:"finishTime"`
		Helpers       []AssistHelper `json:"helpers" dc:"助力人列表(脱敏)"`
	}
	AssistHelper struct {
		UserId    string `json:"userId"`
		Nickname  string `json:"nickname" dc:"助力人(脱敏)"`
		CreatedAt string `json:"createdAt"`
	}

	// 助力（会员; 一人一助力; 风控挂载位: risk_rule 2/3/4）
	AssistHelpReq struct {
		g.Meta   `path:"/assists/{recordId}/helpers" method:"POST" summary:"助力"`
		RecordId string `json:"recordId" v:"required" dc:"参与记录ID"`
	}
	AssistHelpRes struct {
		Success bool `json:"success"`
		Done    bool `json:"done" dc:"本次助力后是否达成发奖"`
	}
)
