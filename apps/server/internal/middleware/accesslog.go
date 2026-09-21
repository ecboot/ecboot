package middleware

import (
	"regexp"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

// maskPattern 敏感字段遮蔽（007-admin-base 评审 C1, FR-007: 凭证/验证码/凭证对不得入日志）。
// 覆盖 JSON 形态的 "key":"value"（大小写不敏感）。
var maskPattern = regexp.MustCompile(
	`(?i)("(?:password|oldpassword|newpassword|smscode|captchacode|captchakey|refreshtoken|token|authorization)"\s*:\s*")[^"]*(")`)

// maskBody 请求体敏感字段脱敏；非 JSON/无敏感字段原样返回。
func maskBody(body string) string {
	if body == "" {
		return body
	}
	return maskPattern.ReplaceAllString(body, `${1}***${2}`)
}

// maskHeaders 头部脱敏：授权/凭证类头仅保留存在性。
func maskHeaders(headers map[string]string) map[string]string {
	out := make(map[string]string, len(headers))
	for k, v := range headers {
		switch k {
		case "Authorization", "Cookie", "Set-Cookie":
			if v != "" {
				out[k] = "***"
			}
		default:
			out[k] = v
		}
	}
	return out
}

func AccessLog(r *ghttp.Request) {
	startTime := time.Now()
	r.Middleware.Next()
	endTime := time.Now()

	// 计算处理时间
	duration := endTime.Sub(startTime).Milliseconds()

	// 获取响应内容
	resp := r.Response.BufferString()

	// 获取请求头信息（授权/凭证类头脱敏）
	headers := map[string]string{}
	for k, v := range r.Header {
		headers[k] = gconv.String(v)
	}

	// 记录请求信息（请求体敏感字段脱敏, FR-007）
	g.Log().Info(r.Context(), g.Map{
		"method":   r.Method,
		"url":      r.URL.String(),
		"headers":  maskHeaders(headers),
		"request":  maskBody(r.GetBodyString()),
		"response": resp,
		"status":   r.Response.Status,
		"duration": duration,
		"logTime":  time.Now(),
	})
}
