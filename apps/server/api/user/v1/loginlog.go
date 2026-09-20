package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 近 30 天登录记录（安全中心; FR-023）
	LoginLogListReq struct {
		g.Meta `path:"/login-logs" method:"GET" summary:"登录记录"`
		PageReq
	}
	LoginLogItem struct {
		Channel   int    `json:"channel" dc:"1小程序 2H5"`
		Status    int    `json:"status" dc:"1成功 2失败"`
		Ip        string `json:"ip" dc:"来源IP(脱敏)"`
		UserAgent string `json:"userAgent" dc:"客户端摘要"`
		CreatedAt string `json:"createdAt" dc:"登录时间"`
	}
	LoginLogListRes struct {
		PageRes
		List []LoginLogItem `json:"list"`
	}
)
