package routes

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/os/gfile"
)

func detectPublic(s *ghttp.Server) {
	if gfile.IsDir("public") {
		s.SetServerRoot("public")
		s.SetOpenApiPath("/api.json")
		oai := s.GetOpenApi()
		oai.Servers = &goai.Servers{
			{
				URL:         "/",
				Description: "开发环境",
			}, {
				URL:         "/test/",
				Description: "测试环境",
			},
		}
		// 1. 定义安全方案 (SecurityScheme)
		oai.Components = goai.Components{
			SecuritySchemes: goai.SecuritySchemes{
				"BearerAuth": goai.SecuritySchemeRef{
					Value: &goai.SecurityScheme{
						Type:         "http",   // 认证类型为 HTTP
						Scheme:       "bearer", // 具体方案为 Bearer
						BearerFormat: "JWT",    // 可选：指定 Token 格式，如 JWT
					},
				},
			},
		}
		// 2. 设置全局安全要求 (SecurityRequirements)
		oai.Security = &goai.SecurityRequirements{
			goai.SecurityRequirement{
				"BearerAuth": []string{},
			},
		}
		oai.Config.CommonResponse = ghttp.DefaultHandlerResponse{}
		oai.Config.CommonResponseDataField = "Data"
	}
}
