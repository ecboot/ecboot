package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"ecboot/internal/library/security"
)

// CtxUserId ctx 中当前会员 ID 的键（控制器经 middleware.CtxUserIdFrom(ctx) 读取）。
const CtxUserId = "currentUserId"

// CtxToken ctx 中当前访问凭证的键（登出销毁用）。
const CtxToken = "currentToken"

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
	"/shop/reviews",
}

func isPublicPath(r *ghttp.Request) bool {
	path := r.URL.Path
	for _, p := range publicPrefixes {
		// 段边界匹配: path==p 或 path 以 p+"/" 开头（防 /user/login* 误匹配）
		if path == p || strings.HasPrefix(path, p+"/") || strings.HasPrefix(path, strings.TrimSuffix(p, "/")+"/") {
			if path == "/shop/reviews" && r.Method == "POST" {
				return false // POST 提交评价是会员行为
			}
			return true
		}
	}
	return false
}

// Auth 鉴权中间件：公开路径直通；其余要求有效 Bearer（会员会话, 滑动续期）。
// 管理组校验（凭证+权限点）属后台特性, 当前管理路径同走会员会话占位。
func Auth(r *ghttp.Request) {
	if isPublicPath(r) {
		r.Middleware.Next()
		return
	}
	token := strings.TrimPrefix(r.GetHeader("Authorization", ""), "Bearer ")
	if token == "" {
		unauthorized(r)
		return
	}
	sm := security.NewSessionManager(7) // TTL 实际由会话创建时的值控制, 校验续期用同默认
	userId, ok, err := sm.Validate(r.Context(), token)
	if err != nil || !ok {
		unauthorized(r)
		return
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
