package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	PingReq struct {
		g.Meta `path:"/ping" method:"GET" summary:"健康探针"`
	}
	PingRes struct {
		Pong string `json:"pong" dc:"固定pong"`
	}
)
