package middleware

import (
	"errors"
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
)

// TestMapErrorCode 契约码映射（评审 I4）: 业务码直出; 框架码归一为 10001/10002。
func TestMapErrorCode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 契约业务码直出（errcode 工厂）
		code, msg := mapErrorCode(errcode.New(errcode.CodeAdminNameTaken, "登录名已存在"))
		t.Assert(code, errcode.CodeAdminNameTaken)
		t.Assert(msg, "登录名已存在")

		// 通用业务码（10001 自身, 契约域边界）
		code, _ = mapErrorCode(errcode.New(errcode.CodeInvalidParam, "区划码非法"))
		t.Assert(code, errcode.CodeInvalidParam)

		// gf 框架：参数校验失败 → 10001（非 51）
		code, msg = mapErrorCode(gerror.NewCode(gcode.CodeValidationFailed, "Validation Failed"))
		t.Assert(code, errcode.CodeInvalidParam)
		t.Assert(msg, "参数校验失败")

		// gf 框架：DB 操作错误 → 10002（非 52）
		code, _ = mapErrorCode(gerror.NewCode(gcode.CodeDbOperationError, "Db Operation Error"))
		t.Assert(code, errcode.CodeSystemError)

		// 普通错误 → 10002
		code, _ = mapErrorCode(errors.New("boom"))
		t.Assert(code, errcode.CodeSystemError)
	})
}
