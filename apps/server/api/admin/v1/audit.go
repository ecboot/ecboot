package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 操作审计日志
	AdminOperationLogListReq struct {
		g.Meta    `path:"/operation-logs" method:"GET" summary:"操作审计日志"`
		AdminId   string `json:"adminId" dc:"操作人"`
		Module    string `json:"module" dc:"业务模块"`
		StartTime string `json:"startTime" dc:"起 RFC3339"`
		EndTime   string `json:"endTime" dc:"止"`
		PageReq
	}
	AdminOperationLogItem struct {
		Id            string `json:"id"`
		Username      string `json:"username" dc:"操作人"`
		Module        string `json:"module" dc:"模块"`
		Operation     string `json:"operation" dc:"操作"`
		Method        string `json:"method" dc:"HTTP方法"`
		RequestUri    string `json:"requestUri" dc:"请求路径"`
		ResultStatus  int    `json:"resultStatus" dc:"1成功 0失败"`
		Ip            string `json:"ip"`
		CostMs        int    `json:"costMs" dc:"耗时毫秒"`
		CreatedAt     string `json:"createdAt"`
	}
	AdminOperationLogListRes struct {
		PageRes
		List []AdminOperationLogItem `json:"list"`
	}

	// 后台登录审计
	AdminLoginLogListReq struct {
		g.Meta  `path:"/admin-login-logs" method:"GET" summary:"后台登录审计"`
		Username string `json:"username" dc:"用户名"`
		PageReq
	}
	AdminLoginLogItem struct {
		Username    string `json:"username"`
		AdminId     string `json:"adminId" dc:"成功时回填"`
		LoginStatus int    `json:"loginStatus" dc:"1成功 2失败"`
		Ip          string `json:"ip"`
		CreatedAt   string `json:"createdAt"`
	}
	AdminLoginLogListRes struct {
		PageRes
		List []AdminLoginLogItem `json:"list"`
	}
)
