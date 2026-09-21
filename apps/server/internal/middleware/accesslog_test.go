package middleware

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// TestMaskBody 访问日志脱敏（评审 C1, FR-007）: 凭证/验证码/凭证对一律遮蔽。
func TestMaskBody(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 登录体: 密码遮蔽, 用户名保留
		out := maskBody(`{"username":"admin","password":"Ecboot@Admin2026"}`)
		t.Assert(out, `{"username":"admin","password":"***"}`)

		// 改密体: 双密码遮蔽（大小写不敏感）
		out = maskBody(`{"oldPassword":"Aa1!","NEWPASSWORD":"Bb2@"}`)
		t.Assert(out, `{"oldPassword":"***","NEWPASSWORD":"***"}`) // 键大小写保留, 值遮蔽

		// 验证码与凭证对
		out = maskBody(`{"captchaKey":"k1","captchaCode":"ab12","smsCode":"654321","refreshToken":"r1","token":"t1"}`)
		t.Assert(out, `{"captchaKey":"***","captchaCode":"***","smsCode":"***","refreshToken":"***","token":"***"}`)

		// 无敏感字段/空体原样
		t.Assert(maskBody(`{"username":"u"}`), `{"username":"u"}`)
		t.Assert(maskBody(``), ``)
	})
}

// TestMaskHeaders 头部脱敏: 授权/会话头只留存在性。
func TestMaskHeaders(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		out := maskHeaders(map[string]string{
			"Authorization": "Bearer abc123",
			"Content-Type":  "application/json",
			"Cookie":        "sid=xyz",
		})
		t.Assert(out["Authorization"], "***")
		t.Assert(out["Cookie"], "***")
		t.Assert(out["Content-Type"], "application/json")
	})
}
