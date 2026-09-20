// Package routes 路由注册：四渠道分组（契约见 specs/004-api-surface/contracts/）。
// 鉴权级别按组承载——Auth 中间件当前为占位（认证实现在 003/后续特性），
// 其白名单语义：公开路径直通，会员组校验 Bearer，管理组校验凭证+权限点。
package routes

import (
	"ecboot/internal/consts"
	"ecboot/internal/controller/admin"
	"ecboot/internal/controller/common"
	"ecboot/internal/controller/shop"
	"ecboot/internal/controller/user"
	"ecboot/internal/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func RouterGroup(s *ghttp.Server) {
	// 根组：全局中间件 + 健康探针
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.AccessLog)
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.GET("/health", healthHandler)

		// 公共渠道（全部公开：验证码/门店/探针/分享上报）
		group.Group("/common", func(commonGroup *ghttp.RouterGroup) {
			commonGroup.Bind(common.NewV1())
		})

		// 会员中心渠道（login/refresh 公开, 其余会员鉴权——白名单在 Auth 内）
		group.Group("/user", func(userGroup *ghttp.RouterGroup) {
			userGroup.Middleware(middleware.Auth)
			userGroup.Bind(user.NewV1())
		})

		// 商城渠道（浏览公开、交易会员鉴权——白名单在 Auth 内）
		group.Group("/shop", func(shopGroup *ghttp.RouterGroup) {
			shopGroup.Middleware(middleware.Auth)
			shopGroup.Bind(shop.NewV1())
		})

		// 管理后台渠道（login 公开, 其余管理员鉴权+权限点——白名单在 Auth 内）
		group.Group("/admin", func(adminGroup *ghttp.RouterGroup) {
			adminGroup.Middleware(middleware.Auth)
			adminGroup.Bind(admin.NewV1())
		})
	})
}

func healthHandler(r *ghttp.Request) {
	r.Response.WriteJsonExit(g.Map{
		"status":  "up",
		"version": consts.Version,
	})
}
