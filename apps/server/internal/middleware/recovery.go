package middleware

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Recovery panic 恢复：真实 panic 转统一系统错误响应（10002 + 追踪号）。
// gf 的 Request.Exit() 以内部哨兵 panic 实现提前返回——此类"正常中止"
// （响应已写出）必须原样放行，否则会产生双 JSON 拼接（评审 C4）。
func Recovery(r *ghttp.Request) {
	defer func() {
		if e := recover(); e != nil {
			// 响应已有内容 = 正常提前返回（Exit），静默放行
			if r.Response.BufferLength() > 0 {
				return
			}
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
