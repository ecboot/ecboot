package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/gogf/gf/v2/net/ghttp"
)

// TraceId 追踪号中间件：每请求生成（或透传上游）X-Trace-Id，
// 回写响应头——贯穿日志与响应（spec FR-004）。
func TraceId(r *ghttp.Request) {
	traceId := r.Header.Get("X-Trace-Id")
	if traceId == "" {
		b := make([]byte, 8)
		_, _ = rand.Read(b)
		traceId = hex.EncodeToString(b)
	}
	r.Response.Header().Set("X-Trace-Id", traceId)
	log.Println("[trace] traceId=", traceId, "uri=", r.URL.Path)

	r.Middleware.Next()
}
