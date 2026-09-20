package middleware

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Recovery panic 恢复：转统一系统错误响应（10002 + 追踪号），不泄漏堆栈（spec FR-003）。
// 必须挂在业务分组最外层（TraceId 之后，保证追踪号已写入）。
func Recovery(r *ghttp.Request) {
	defer func() {
		if e := recover(); e != nil {
			g.Log().Errorf(r.Context(), "panic recovered: %v", e)
			traceId := r.Response.Header().Get("X-Trace-Id")
			r.Response.WriteJson(g.Map{
				"code":    10002,
				"message": "系统繁忙,请稍后重试(追踪号: " + traceId + ")",
				"data":    nil,
			})
		}
	}()
	r.Middleware.Next()
}
