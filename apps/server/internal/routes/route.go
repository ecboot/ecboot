// Package routes 路由注册：四渠道分组（契约见 specs/004-api-surface/contracts/）。
// 中间件链：TraceId → Recovery → AccessLog → Response（统一三段式包装与错误映射）。
// 鉴权级别按组承载——Auth 中间件白名单直通公开路径；会员组校验 Bearer（003 实现真实校验）；
// 管理组校验凭证+权限点（后台特性落地）。
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
		group.Middleware(middleware.TraceId)
		group.Middleware(middleware.Recovery)
		group.Middleware(middleware.AccessLog)
		group.Middleware(middleware.Response)
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
