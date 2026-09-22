package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/api/user/v1"
	"ecboot/internal/library/security"
	"ecboot/internal/service/user"
)

// TokenRefresh 刷新访问凭证（双凭证会话）。
// N3 收口: 刷新前校验账号态——原实现**不校验**（会员禁用后 refresh 仍换发新双凭证,
// 且 refresh 每次轮换长寿 4×TTL → 近乎无限续期, 禁用可被完全绕过）。
// 注: 须先从旧 refresh token 解析出 userId（Refresh 内部 GETDEL 即作废）, 故改两段式校验;
// 简化实现: 先解码 refresh 得 userId（不消费）, 校验账号态后再执行 Refresh。
func (c *ControllerV1) TokenRefresh(ctx context.Context, req *v1.TokenRefreshReq) (res *v1.TokenRefreshRes, err error) {
	sm := security.NewSessionManager("user", 7)
	uid, err := sm.PeekUserId(ctx, req.RefreshToken)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(10003, "刷新凭证失效", nil))
	}
	if active, aErr := user.ActiveUser(ctx, uid); aErr != nil || !active {
		return nil, gerror.NewCode(gcode.New(10003, "账号不可用", nil)) // 禁用/软删 → 拒绝续期
	}
	token, refresh, _, err := sm.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(10003, "刷新凭证失效", nil))
	}
	return &v1.TokenRefreshRes{Token: token, RefreshToken: refresh}, nil
}
