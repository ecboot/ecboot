package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminFlashSaleList ---------- 秒杀 ----------
func (c *ControllerV1) AdminFlashSaleList(ctx context.Context, req *v1.AdminFlashSaleListReq) (res *v1.AdminFlashSaleListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
