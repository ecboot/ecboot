package routes

import (
	"ecboot/internal/consts"
	"ecboot/internal/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func RouterGroup(s *ghttp.Server) {
	// 注册路由
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/health", healthHandler)
		group.Middleware(middleware.AccessLog)
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(
			// hello.NewV1(),
		)
	})
}

func healthHandler(r *ghttp.Request) {
	r.Response.WriteJsonExit(g.Map{
		"status":  "up",
		"version": consts.Version,
	})
}