package middleware

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// publicPrefixes 公开路径白名单（评审 C1：与路由注释语义一致，直通不校验）。
// 完整清单权威：specs/004-api-surface/contracts/conventions.md §鉴权三级。
// 回调路径（/shop/pay/notify 等）必须始终在此列——渠道签名验证语义，永不走 Bearer。
var publicPrefixes = []string{
	"/common/",           // 公共渠道全部公开
	"/user/login",        // 登录两端点
	"/user/token/refresh", // 刷新凭证
	"/admin/login",       // 后台登录
	"/admin/token/refresh",
	// shop 渠道公开浏览（contracts/shop-api.md 标注公开的端点）
	"/shop/categories", "/shop/brands", "/shop/products", "/shop/search",
	"/shop/activities/", "/shop/full-reductions", "/shop/banners", "/shop/floors",
	"/shop/bargains/", "/shop/assists/", // 进度查看（帮砍/助力动作的鉴权属实现层细化）
	"/shop/pay/notify", "/shop/refund/notify", // 渠道回调（签名验证）
	"/shop/reviews", // 商品评价列表（注意与会员提交 /shop/reviews POST 的区分，见下方方法判定）
}

// isPublicPath 判定请求是否公开（前缀 + POST /shop/reviews 的特例处理：
// GET 同路径为公开评价列表，POST 为会员提交——占位阶段对 POST 不放行）。
func isPublicPath(r *ghttp.Request) bool {
	path := r.URL.Path
	for _, p := range publicPrefixes {
		if strings.HasPrefix(path, p) {
			// 特例：POST /shop/reviews 是会员提交评价，非公开
			if path == "/shop/reviews" && r.Method == "POST" {
				return false
			}
			return true
		}
	}
	return false
}

// Auth 鉴权占位中间件：公开路径直通；其余要求 Bearer 存在（真实校验属认证特性）。
// 响应遵循统一三段式契约（code=10003），不使用裸 401。
func Auth(r *ghttp.Request) {
	if isPublicPath(r) {
		r.Middleware.Next()
		return
	}
	token := strings.TrimPrefix(r.GetHeader("Authorization", ""), "Bearer ")
	if token == "" {
		r.Response.WriteJson(g.Map{
			"code":    10003,
			"message": "未登录或凭证失效",
			"data":    nil,
		})
		return
	}

	// TODO(认证特性): 校验 token 有效性（会话/续期）、管理组权限点校验

	r.Middleware.Next()
}
