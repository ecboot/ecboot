package middleware

import (
	"errors"

	bizerr "ecboot/internal/library/err"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gvalid"
)

// Response 统一响应中间件（替代 ghttp.MiddlewareHandlerResponse）：
// handler 无错 → {code:0, message:"ok", data:res}；
// 业务异常 → {code:业务码, message:文案}；
// 校验失败 → {code:10001, message:首错含字段}；
// 其他错误 → {code:10002, message:系统错误+追踪号}（不泄漏内部细节, spec FR-003）。
func Response(r *ghttp.Request) {
	r.Middleware.Next()

	var (
		code    int
		message string
		data    any
	)

	if err := r.GetError(); err != nil {
		code, message = mapError(err, r)
		// 清除错误避免框架再按默认形态输出
		r.SetError(nil)
	} else {
		code, message = 0, "ok"
		data = r.GetHandlerResponse()
	}

	r.Response.WriteJson(g.Map{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func mapError(err error, r *ghttp.Request) (int, string) {
	// 业务异常：携带契约错误码
	var biz *bizerr.Business
	if errors.As(err, &biz) {
		return biz.Code(), biz.Message()
	}

	// 参数校验失败：首错含字段与原因（FR-003）
	var verr gvalid.Error
	if errors.As(err, &verr) {
		msg := "参数校验失败"
		if fe := verr.FirstError(); fe != nil {
			msg = "参数校验失败: " + fe.Error()
		}
		return 10001, msg
	}

	// 系统错误：追踪号回传便于排查，不暴露内部细节
	traceId := r.Response.Header().Get("X-Trace-Id")
	g.Log().Error(r.Context(), err)
	return 10002, "系统繁忙,请稍后重试(追踪号: " + traceId + ")"
}
