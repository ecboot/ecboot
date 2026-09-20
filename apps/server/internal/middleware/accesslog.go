package middleware

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

func AccessLog(r *ghttp.Request) {
	startTime := time.Now()
	r.Middleware.Next()
	endTime := time.Now()

	// 计算处理时间
	duration := endTime.Sub(startTime).Milliseconds()

	// 获取响应内容
	resp := r.Response.BufferString()

	// 获取请求头信息
	headers := map[string]string{}
	for k, v := range r.Header {
		headers[k] = gconv.String(v)
	}

	// 记录请求信息
	g.Log().Info(r.Context(), g.Map{
		"method":   r.Method,
		"url":      r.URL.String(),
		"headers":  headers,
		"request":  r.GetBodyString(),
		"response": resp,
		"status":   r.Response.Status,
		"duration": duration,
		"logTime":  time.Now(),
	})
}
