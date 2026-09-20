package middleware

import (
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

func Auth(r *ghttp.Request) {
	token := r.GetHeader("Authorization", "")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    http.StatusUnauthorized,
			Message: "未授权",
		})
		return
	}

	// TODO: 验证 JWT 令牌

	// TODO: 提取 JWT 中的用户信息

	// TODO: 验证用户权限

	r.Middleware.Next()
}
