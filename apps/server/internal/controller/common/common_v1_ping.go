package common

import (
	"context"

	"ecboot/api/common/v1"
)

// Ping 健康探针（公开, FR-022）
func (c *ControllerV1) Ping(ctx context.Context, req *v1.PingReq) (res *v1.PingRes, err error) {
	return &v1.PingRes{Pong: "pong"}, nil
}
