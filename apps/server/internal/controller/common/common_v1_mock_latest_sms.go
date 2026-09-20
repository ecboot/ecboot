package common

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/common/v1"
)

// MockLatestSms 联调取码（仅 mock 模式注册路由; 生产不存在该端点）
func (c *ControllerV1) MockLatestSms(ctx context.Context, req *v1.MockLatestSmsReq) (res *v1.MockLatestSmsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
