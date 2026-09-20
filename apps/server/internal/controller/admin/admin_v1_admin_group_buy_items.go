package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminGroupBuyItems 拼团场次商品（SKU 级成团价; 全量替换）
func (c *ControllerV1) AdminGroupBuyItems(ctx context.Context, req *v1.AdminGroupBuyItemsReq) (res *v1.AdminGroupBuyItemsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
