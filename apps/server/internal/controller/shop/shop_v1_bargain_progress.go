package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/shop/v1"
)

// BargainProgress 砍价进度（公开; 含帮砍列表）
func (c *ControllerV1) BargainProgress(ctx context.Context, req *v1.BargainProgressReq) (res *v1.BargainProgressRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
