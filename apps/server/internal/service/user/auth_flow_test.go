package user

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/database/gdb"
	gredis "github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/library/captcha"
	"ecboot/internal/library/security"
	"ecboot/internal/library/sms"
)

func init() {
	// 时区口径统一（009 评审 C2）: 与库内 UTC 墙钟一致（见 main.go）
	time.Local = time.UTC
	// 测试显式开启 mock（fail-closed 语义下的白盒开关）
	os.Setenv("ECBOOT_MOCK", "true")
	// 确定性测试配置（不依赖 manifest/config 的本地差异——compose 基线环境）。
	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			{
				Link: "mysql:myuser:secret@tcp(127.0.0.1:13306)/mydatabase",
			},
		},
	})
	gredis.SetConfig(&gredis.Config{
		Address: "127.0.0.1:6379",
		Db:      0,
	})
}

// issueCode 清理频控状态并走真实发码链路, 返回 mock 短信码。
func issueCode(ctx context.Context, t *gtest.T, phone string) string {
	_, _ = g.Redis().Do(ctx, "DEL", "captcha:sms:interval:"+phone)
	key, _, answer, err := captcha.Generate(ctx)
	t.AssertNil(err)
	_, err = sms.SendSmsCode(ctx, phone, key, answer)
	t.AssertNil(err)
	v, err := g.Redis().Do(ctx, "GET", "mock:sms:"+phone)
	t.AssertNil(err)
	return v.String()
}

// cleanState 清理测试遗留（user 行 + Redis 状态键）。
func cleanState(ctx context.Context, t *gtest.T, phone string) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE phone_hash=?", phoneCipher().Hash(phone))
	_, _ = g.Redis().Do(ctx, "DEL",
		"captcha:sms:interval:"+phone, "captcha:fail:"+phone,
		"captcha:sms:"+phone, "mock:sms:"+phone)
}

// errCode 提取错误中的契约码（兼容 gerror 包装）。
func errCode(err error) int {
	var ge *gerror.Error
	if errors.As(err, &ge) {
		return ge.Code().Code()
	}
	return -1
}

// TestAuthFlow 认证链路端到端。
func TestAuthFlow(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13800001111"
		cleanState(ctx, t, phone)

		// ① 图形码一次性
		key, _, answer, err := captcha.Generate(ctx)
		t.AssertNil(err)
		ok, err := captcha.Verify(ctx, key, answer)
		t.AssertNil(err)
		t.Assert(ok, true)
		ok, err = captcha.Verify(ctx, key, answer)
		t.AssertNil(err)
		t.Assert(ok, false)

		// ②③ 注册即登录
		code := issueCode(ctx, t, phone)
		out, err := SmsLogin(ctx, phone, code, 1)
		t.AssertNil(err)
		t.Assert(out.IsNew, true)

		// ④ 库内零明文
		rec, err := g.DB().GetOne(ctx, "SELECT phone, phone_hash FROM `user` WHERE id=?", out.UserId)
		t.AssertNil(err)
		t.Assert(rec["phone_hash"].String() != "", true)
		phoneStored := rec["phone"].String()
		if len(phoneStored) == 11 && phoneStored[0] == '1' {
			t.Error("user.phone 疑似明文存储")
		}

		// ⑤ 会话有效
		sm := security.NewSessionManager("user", 7)
		uid, ok, err := sm.Validate(ctx, out.Token)
		t.AssertNil(err)
		t.Assert(ok, true)
		t.Assert(uid, out.UserId)

		// ⑥ 微信归并: 同号+有效码 → 同一账号
		mergeCode := issueCode(ctx, t, phone)
		out2, err := WxLogin(ctx, "dev001", phone, mergeCode, 1)
		t.AssertNil(err)
		t.Assert(out2.UserId, out.UserId)

		// ⑦ 绑定冲突 → 20004
		conflictCode := issueCode(ctx, t, phone)
		_, err = WxLogin(ctx, "dev002", phone, conflictCode, 1)
		t.Assert(errCode(err), errcode.CodeWxBindConflict)

		// ⑧ 重复登录幂等
		code3 := issueCode(ctx, t, phone)
		out3, err := SmsLogin(ctx, phone, code3, 1)
		t.AssertNil(err)
		t.Assert(out3.IsNew, false)
		t.Assert(out3.UserId, out.UserId)

		// ⑨ 休眠分级: 微信静默被拒; 短信通道放行
		_, _ = g.DB().Exec(ctx, "UPDATE `user` SET last_active_at=DATE_SUB(NOW(), INTERVAL 91 DAY) WHERE id=?", out.UserId)
		_, err = WxLogin(ctx, "dev001", "", "", 1)
		t.Assert(err != nil, true)
		replayCode := issueCode(ctx, t, phone)
		out4, err := SmsLogin(ctx, phone, replayCode, 1)
		t.AssertNil(err)
		t.Assert(out4.UserId, out.UserId)

		_, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE id=?", out.UserId)
		_, _ = g.DB().Exec(ctx, "DELETE FROM user_login_log WHERE user_id=?", out.UserId)
	})
}

// TestSessionLifecycle 会话三态与 refresh 轮换（评审 I12）。
func TestSessionLifecycle(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sm := security.NewSessionManager("user", 7)

		token, refresh, err := sm.Create(ctx, 999)
		t.AssertNil(err)

		nt, nr, ruid, err := sm.Refresh(ctx, refresh)
		t.AssertNil(err)
		t.Assert(ruid, 999)
		t.Assert(nt != token, true)

		_, _, _, err = sm.Refresh(ctx, refresh) // 旧 refresh 已消费
		t.Assert(err != nil, true)

		t.AssertNil(sm.Destroy(ctx, nt, nr))
		_, ok, err := sm.Validate(ctx, nt)
		t.AssertNil(err)
		t.Assert(ok, false)
		_, _, _, err = sm.Refresh(ctx, nr)
		t.Assert(err != nil, true)
	})
}

// TestAntiAbuse 防刷规则（60s 重发窗口 / 5 次失败锁定——评审 I12）。
func TestAntiAbuse(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13800002222"
		cleanState(ctx, t, phone)

		_ = issueCode(ctx, t, phone)
		key, _, answer, err := captcha.Generate(ctx)
		t.AssertNil(err)
		_, err = sms.SendSmsCode(ctx, phone, key, answer)
		t.Assert(errCode(err), errcode.CodeTooFrequent)

		for i := 0; i < 5; i++ {
			dbg, _ := g.Redis().Do(ctx, "EXISTS", "captcha:fail:"+phone)
			t.Log("DEBUG i=", i, "failKeyExists=", dbg.Int())
			_, err = SmsLogin(ctx, phone, "000000", 1)
			t.Log("DEBUG i=", i, "errCode=", errCode(err))
			t.Assert(errCode(err), errcode.CodeCaptchaError) // 1~5 次: 20001
		}
		// 第 6 次: 计数达上限 → 锁定 20002
		_, err = SmsLogin(ctx, phone, "000000", 1)
		t.Assert(errCode(err), errcode.CodeLocked)

		cleanState(ctx, t, phone)
	})
}
