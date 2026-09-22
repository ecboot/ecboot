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
	"/shop/index", // 首页聚合对游客开放（015 批次 09; 由端点级可达性测试发现漏配）
	// M11（015 评审）: 本行是**前缀**形态而 /shop/index 并无子路由——将来若出现
	// `/shop/index/*` 的会员端点会被同型吞掉（批次 08/09 两次 C1 同型）; 届时应改用 publicGetOnlyPrefixes。
	// 另一条已知防线边界: gf 支持 `X-Url-Path` 头覆盖路由路径, 而本白名单只查 r.URL.Path——
	// "白名单路径 + 该头"可让请求不注入 userId 直达任意 handler。今天无害的条件是
	// **所有会员 controller 都调用了 requireMember**（47 个 shop controller 已逐一核对, 015 评审）,
	// 新增会员端点时必须保持该惯例, 否则此头即成绕过面。
	"/shop/categories", "/shop/brands", "/shop/products", "/shop/search",
	"/shop/activities/", "/shop/full-reductions", "/shop/banners", "/shop/floors",
	// 注（015 批次 09 全端点扫描）: `/shop/bargains/` 与 `/shop/assists/` **不再**列在此处——
	// 这两个前缀下既有公开的 GET（进度查询）, 又有**会员的 POST**（帮砍 /{id}/cut、助力 /{id}/helpers）;
	// 前缀放行会把会员动作一并直通（Auth 不注入 ctx userId）→ 控制器 requireMember 恒判未登录。
	// 它们改列 publicGetOnlyPrefixes（只放行 GET）。
	"/shop/pay/notify", "/shop/refund/notify",
	// 注（014 评审 C1）: 此处**不得**加 `/shop/reviews`——它是前缀匹配, 会把同一前缀下的
	// `/shop/reviews/mine`（我的评价）与 `/shop/reviews/{id}/extra`（追评）一并放行,
	// 而 Auth 对白名单路径**直通不注入 ctx userId** → 控制器 requireMember 恒判未登录(10003)。
	// 公开的"商品评价列表"是 `/shop/products/{spuId}/reviews`, 已由上面的 `/shop/products` 覆盖。
	// （该行原先在此, 并配有一处 `POST /shop/reviews` 的精确例外; 端点未连线时不可观测, 连线后即成故障。）
}

// publicGetOnlyPrefixes 只放行 **GET** 的公开前缀。
// 场景（015 批次 09 全端点扫描得出）: 同一前缀下既有公开读（如砍价/助力进度查询）, 又有会员写
// （帮砍、助力）——前缀级放行会把"写"也当公开路径直通, 控制器拿不到 userId（批次 08 的 C1 同型）。
var publicGetOnlyPrefixes = []string{"/shop/bargains/", "/shop/assists/"}

func isPublicPath(r *ghttp.Request) bool {
	path := r.URL.Path
	if r.Method == "GET" {
		for _, p := range publicGetOnlyPrefixes {
			if path == p || strings.HasPrefix(path, p+"/") || strings.HasPrefix(path, strings.TrimSuffix(p, "/")+"/") {
				return true
			}
		}
	}
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
