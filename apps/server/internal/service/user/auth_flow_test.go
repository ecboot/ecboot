package user

import (
	"context"
	"errors"
	"testing"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/library/captcha"
	"ecboot/internal/library/security"
)

func init() {
	// 单测 cwd 为包目录: 显式指向 manifest/config 使 g.Redis()/g.DB() 可用
	if a, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile); ok {
		_ = a.SetPath("../../manifest/config")
	}
}

// TestAuthFlow 认证链路端到端（依赖 compose 的 MySQL/Redis 与已应用迁移）。
func TestAuthFlow(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13800001111"
		// 清理上次运行遗留的用户与 Redis 状态键
		_, _ = g.DB().Exec(ctx, "DELETE FROM user WHERE phone_hash=?", phoneCipher().Hash(phone))
		// 清理上次运行遗留的 Redis 状态键（重发锁/失败计数/旧码）
		_, _ = g.Redis().Do(ctx, "DEL",
			"captcha:sms:interval:"+phone, "captcha:fail:"+phone,
			"captcha:sms:"+phone, "mock:sms:"+phone)

		// ---- ① 图形验证码: 获取(直调组件可拿答案) → 一次性断言 ----
		key, _, answer, err := captcha.Generate(ctx)
		t.AssertNil(err)
		ok, err := captcha.Verify(ctx, key, answer)
		t.AssertNil(err)
		t.Assert(ok, true)
		ok, err = captcha.Verify(ctx, key, answer) // 同凭证复用必须失败
		t.AssertNil(err)
		t.Assert(ok, false)

		// ---- ② 发短信码（重新取图形码） ----
		key2, _, answer2, err := captcha.Generate(ctx)
		t.AssertNil(err)
		_, err = SendSmsCode(ctx, phone, key2, answer2)
		t.AssertNil(err)

		// ---- ③ 注册即登录 ----
		v, err := g.Redis().Do(ctx, "GET", "mock:sms:"+phone)
		t.AssertNil(err)
		code := v.String()
		t.Assert(code != "", true)
		out, err := SmsLogin(ctx, phone, code, 1)
		t.AssertNil(err)
		t.Assert(out.IsNew, true)
		t.Assert(out.UserId > 0, true)

		// ---- ④ 库内零明文（SC-007） ----
		rec, err := g.DB().GetOne(ctx, "SELECT phone, phone_hash FROM `user` WHERE id=?", out.UserId)
		t.AssertNil(err)
		t.Assert(rec["phone_hash"].String() != "", true)
		phoneStored := rec["phone"].String()
		if len(phoneStored) == 11 && phoneStored[0] == '1' {
			t.Error("user.phone 疑似明文存储")
		}

		// ---- ⑤ 会话有效性 ----
		sm := security.NewSessionManager(7)
		uid, ok, err := sm.Validate(ctx, out.Token)
		t.AssertNil(err)
		t.Assert(ok, true)
		t.Assert(uid, out.UserId)

		// ---- ⑥ 微信归并: 同手机号 → 同一账号 ----
		out2, err := WxLogin(ctx, "dev001", phone, "", 1)
		t.AssertNil(err)
		t.Assert(out2.UserId, out.UserId)

		// ---- ⑦ 绑定冲突: 另一微信身份携同号 → 20004 ----
		_, err = WxLogin(ctx, "dev002", phone, "", 1)
		t.Assert(err != nil, true)
		t.Assert(errCode(err), errcode.CodeWxBindConflict)

		// ---- ⑧ 重复登录（已注册路径, 幂等） ----
		out3, err := smsLoginReplay(ctx, phone)
		t.AssertNil(err)
		t.Assert(out3.IsNew, false)
		t.Assert(out3.UserId, out.UserId)

		// 清理测试数据
		_, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE id=?", out.UserId)
		_, _ = g.DB().Exec(ctx, "DELETE FROM user_login_log WHERE user_id=?", out.UserId)
	})
}

// smsLoginReplay 重新发码并登录（已注册账号路径）。
func smsLoginReplay(ctx context.Context, phone string) (*LoginOutcome, error) {
	_, _ = g.Redis().Do(ctx, "DEL", "captcha:sms:interval:"+phone) // 测试内重发绕开 60s 窗口
	key, _, answer, err := captcha.Generate(ctx)
	if err != nil {
		return nil, err
	}
	if _, err = SendSmsCode(ctx, phone, key, answer); err != nil {
		return nil, err
	}
	v, err := g.Redis().Do(ctx, "GET", "mock:sms:"+phone)
	if err != nil {
		return nil, err
	}
	return SmsLogin(ctx, phone, v.String(), 1)
}

// errCode 提取错误中的契约码（兼容 gerror 包装与裸 Business）。
func errCode(err error) int {
	var ge *gerror.Error
	if errors.As(err, &ge) {
		return ge.Code().Code()
	}
	type coder interface{ Code() int }
	var c coder
	if errors.As(err, &c) {
		return c.Code()
	}
	return -1
}
