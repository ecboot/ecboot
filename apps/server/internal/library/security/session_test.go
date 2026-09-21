package security

import (
	"context"
	"os"
	"testing"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	gredis "github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/test/gtest"
)

func init() {
	// 确定性测试配置（compose 基线, 与 service 层测试同源）
	_ = os.Setenv("ECBOOT_MOCK", "true")
	gredis.SetConfig(&gredis.Config{
		Address: "127.0.0.1:6379",
		Db:      0,
	})
}

// TestSessionAudienceIsolation 会话渠道隔离（007-admin-base research D1 安全属性）：
// user 会话对 admin 渠道不可见, 反之亦然——修复跨渠道越权（同号 userId 冒充）。
func TestSessionAudienceIsolation(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		userMgr := NewSessionManager("user", 7)
		adminMgr := NewSessionManager("admin", 7)

		token, refresh, err := userMgr.Create(ctx, 7)
		t.AssertNil(err)
		defer func() { _ = userMgr.Destroy(ctx, token, refresh) }()

		// 本渠道可验
		uid, ok, err := userMgr.Validate(ctx, token)
		t.AssertNil(err)
		t.Assert(ok, true)
		t.Assert(uid, 7)

		// 跨渠道拒绝（核心安全属性: 修复前同号 userId 可冒充管理员）
		_, ok, err = adminMgr.Validate(ctx, token)
		t.AssertNil(err)
		t.Assert(ok, false)

		// 跨渠道 refresh 拒绝
		_, _, _, err = adminMgr.Refresh(ctx, refresh)
		t.AssertNE(err, nil)

		// 跨渠道登出不可销毁本渠道会话
		_ = adminMgr.Destroy(ctx, token, refresh)
		_, ok, err = userMgr.Validate(ctx, token)
		t.AssertNil(err)
		t.Assert(ok, true)
	})
}

// TestSessionRefreshOneShot 刷新凭证一次性：refresh 后旧凭证全失效（既有语义, audience 后保持）。
func TestSessionRefreshOneShot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		mgr := NewSessionManager("user", 7)
		token, refresh, err := mgr.Create(ctx, 42)
		t.AssertNil(err)
		defer func() { _ = mgr.Destroy(ctx, token, refresh) }()

		newToken, newRefresh, uid, err := mgr.Refresh(ctx, refresh)
		t.AssertNil(err)
		t.Assert(uid, 42)
		defer func() { _ = mgr.Destroy(ctx, newToken, newRefresh) }()

		// 旧 refreshToken 已被 GETDEL, 重放拒绝
		_, _, _, err = mgr.Refresh(ctx, refresh)
		t.AssertNE(err, nil)
	})
}
