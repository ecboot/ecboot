package admin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/admin/v1"
)

// AdminMemberRebindPhone 改绑手机号（旧号解占; 审计留痕）
func (c *ControllerV1) AdminMemberRebindPhone(ctx context.Context, req *v1.AdminMemberRebindPhoneReq) (res *v1.AdminMemberRebindPhoneRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
