package admin

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"ecboot/api/admin/v1"
	"ecboot/internal/service/system"
)

// AdminLogin 后台登录（含登录审计）
func (c *ControllerV1) AdminLogin(ctx context.Context, req *v1.AdminLoginReq) (res *v1.AdminLoginRes, err error) {
	r := g.RequestFromCtx(ctx)
	out, err := system.AdminLogin(ctx, req.Username, req.Password,
		req.CaptchaKey, req.CaptchaCode, r.GetClientIp(), r.Header.Get("User-Agent"))
	if err != nil {
		return nil, err
	}
	return &v1.AdminLoginRes{
		Token:        out.Token,
		RefreshToken: out.RefreshToken,
		RealName:     out.RealName,
		IsSuper:      out.IsSuper,
	}, nil
}
