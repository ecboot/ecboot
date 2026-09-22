package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"ecboot/internal/consts"
	"ecboot/internal/library/security"
	"ecboot/internal/service/system"
)

// CtxUserId/CtxToken ctx 键已下沉 consts（service 层可无环读取），此处别名保持既有引用。
const (
	CtxUserId = consts.CtxUserId
	CtxToken  = consts.CtxToken
)

// CtxUserIdFrom 从上下文读取当前会员 ID（0=未登录）。
func CtxUserIdFrom(ctx context.Context) int64 {
	if v, ok := ctx.Value(CtxUserId).(int64); ok {
		return v
	}
	return 0
}

// publicPrefixes 公开路径白名单（与路由注释语义一致, 直通不校验）。
// 回调路径永不走 Bearer（渠道签名验证语义）。
var publicPrefixes = []string{
	"/common/",
	"/user/login", "/user/token/refresh",
	"/admin/login", "/admin/token/refresh",
	"/shop/categories", "/shop/brands", "/shop/products", "/shop/search",
	"/shop/activities/", "/shop/full-reductions", "/shop/banners", "/shop/floors",
	"/shop/bargains/", "/shop/assists/",
	"/shop/pay/notify", "/shop/refund/notify",
	// 注（014 评审 C1）: 此处**不得**加 `/shop/reviews`——它是前缀匹配, 会把同一前缀下的
	// `/shop/reviews/mine`（我的评价）与 `/shop/reviews/{id}/extra`（追评）一并放行,
	// 而 Auth 对白名单路径**直通不注入 ctx userId** → 控制器 requireMember 恒判未登录(10003)。
	// 公开的"商品评价列表"是 `/shop/products/{spuId}/reviews`, 已由上面的 `/shop/products` 覆盖。
	// （该行原先在此, 并配有一处 `POST /shop/reviews` 的精确例外; 端点未连线时不可观测, 连线后即成故障。）
}

func isPublicPath(r *ghttp.Request) bool {
	path := r.URL.Path
	for _, p := range publicPrefixes {
		// 段边界匹配: path==p 或 path 以 p+"/" 开头（防 /user/login* 误匹配）
		if path == p || strings.HasPrefix(path, p+"/") || strings.HasPrefix(path, strings.TrimSuffix(p, "/")+"/") {
			return true
		}
	}
	return false
}

// Auth 鉴权中间件：公开路径直通；其余要求有效 Bearer（滑动续期）。
// 会话按渠道隔离（research D1）：/admin/ 前缀走 admin 会话并校验账号态（FR-008:
// 禁用/软删账号的既有会话立即不可用）；其余走 user 会话。
// 管理组权限点校验由 controller 显式 RequirePerm 挂接（research D2）。
func Auth(r *ghttp.Request) {
	if isPublicPath(r) {
		r.Middleware.Next()
		return
	}
	aud := "user"
	if strings.HasPrefix(r.URL.Path, "/admin/") {
		aud = "admin"
	}
	token := strings.TrimPrefix(r.GetHeader("Authorization", ""), "Bearer ")
	if token == "" {
		unauthorized(r)
		return
	}
	sm := security.NewSessionManager(aud, 7) // TTL 实际由会话创建时的值控制, 校验续期用同默认
	userId, ok, err := sm.Validate(r.Context(), token)
	if err != nil || !ok {
		unauthorized(r)
		return
	}
	// 管理渠道加验账号态（禁用/软删即拒）
	if aud == "admin" {
		active, aErr := system.ActiveAdmin(r.Context(), userId)
		if aErr != nil || !active {
			unauthorized(r)
			return
		}
	}
	// userId 注入 ctx（控制器经 CtxUserIdFrom 读取）
	ctx := context.WithValue(r.Context(), CtxUserId, userId)
	ctx = context.WithValue(ctx, CtxToken, token)
	r.SetCtx(ctx)

	r.Middleware.Next()
}

func unauthorized(r *ghttp.Request) {
	r.Response.WriteJson(g.Map{
		"code":    10003,
		"message": "未登录或凭证失效",
		"data":    nil,
	})
	r.Exit()
}

var _ = http.StatusOK
