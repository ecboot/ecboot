// permission.go 管理端权限点挂接 helper（research D2: controller 显式调用, 调用点即文档）。
// 用法: 受权限点保护的 controller 方法首行 `if err := middleware.RequirePerm(ctx, "system:xxx:yyy"); err != nil { return nil, err }`。
// 权限码以 api 层注释与 000032 种子为准绳（internal/consts/permission.go 同源）。
package middleware

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/internal/errcode"
	"ecboot/internal/service/system"
)

// RequirePerm 权限点校验: 未登录 10003; 无权(非超管) 10005; 超管直通（FR-017）。
func RequirePerm(ctx context.Context, code string) error {
	adminId := CtxUserIdFrom(ctx)
	if adminId == 0 {
		return gerror.NewCode(gcode.New(errcode.CodeUnauthorized, "未登录或凭证失效", nil))
	}
	ok, err := system.HasPermission(ctx, adminId, code)
	if err != nil {
		return err
	}
	if !ok {
		return errcode.New(errcode.CodeForbidden, "权限不足")
	}
	return nil
}
