package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	IndexReq struct {
		g.Meta `path:"/index" tags:"Shop" method:"GET" summary:"首页接口"`
	}

	IndexRes struct {
		Msg string `json:"msg" v:"required" dc:"首页接口"`
	}
)