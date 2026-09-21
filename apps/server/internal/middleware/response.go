package middleware

import (
	"errors"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gvalid"

	"ecboot/internal/errcode"
)

// Response 统一响应中间件（替代 ghttp.MiddlewareHandlerResponse）：
//   - handler 无错且未自行写出 → {code:0, message:"ok", data:res}
//   - 携带 gcode 的业务错误（errcode.New 工厂）→ {code:契约码, message:文案}
//   - 校验失败 → {code:10001, message:首错含字段}
//   - 其他错误 → {code:10002, message:系统错误+追踪号}（不泄漏内部细节, FR-003）
//
// 若下游已写出响应体（鉴权直写/渠道回调/Exit 提前返回），本中间件不再覆盖（评审 C4）。
func Response(r *ghttp.Request) {
	r.Middleware.Next()

	if r.Response.BufferLength() > 0 {
		// 下游已显式写出——保留原响应
		r.SetError(nil)
		return
	}

	var (
		code    int
		message string
		data    any
	)

	if err := r.GetError(); err != nil {
		code, message = mapError(err, r)
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
	code, message := mapErrorCode(err)
	if code == errcode.CodeSystemError {
		// 系统错误：服务端留日志 + 追踪号回传（不暴露内部细节, FR-003）
		traceId := r.Response.Header().Get("X-Trace-Id")
		g.Log().Error(r.Context(), err)
		message = "系统繁忙,请稍后重试(追踪号: " + traceId + ")"
	}
	return code, message
}

// mapErrorCode 错误 → 契约码与文案（纯函数便于单测; 系统错误的追踪号与日志由调用方补）。
//   - 契约业务码（>= 10001, errcode.New 工厂产物）→ 直出（评审 C1）
//   - 参数校验失败（gvalid 错误或 gf 校验码）→ 10001, 首错含字段（FR-003）
//   - 其余（**含 gf 框架码 51/52 等**）→ 10002 系统错误
//
// 007/008 评审 I4: 初版判据 `gc.Code() > 0` 使 gf 框架码（51 Validation Failed /
// 52 Db Operation Error）直接透出，契约的 10001/10002 语义失效——现以契约码域
// （>= CodeInvalidParam）为界，框架码统一归入 10002 并恢复日志与追踪号。
func mapErrorCode(err error) (int, string) {
	if gc := gerror.Code(err); gc.Code() >= errcode.CodeInvalidParam {
		return gc.Code(), gc.Message()
	}

	// 参数校验失败：首错含字段与原因
	var verr gvalid.Error
	if errors.As(err, &verr) {
		msg := "参数校验失败"
		if fe := verr.FirstError(); fe != nil {
			msg = "参数校验失败: " + fe.Error()
		}
		return errcode.CodeInvalidParam, msg
	}
	if gerror.Code(err).Code() == gcode.CodeValidationFailed.Code() {
		return errcode.CodeInvalidParam, "参数校验失败"
	}

	return errcode.CodeSystemError, ""
}
