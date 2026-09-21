package system

import (
	"context"
	"errors"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/library/captcha"
)

// errCode 提取错误中的契约码（兼容 gerror 包装, 与 user 域测试同式）。
func errCode(err error) int {
	var ge *gerror.Error
	if errors.As(err, &ge) {
		return ge.Code().Code()
	}
	return -1
}

// lastLoginLog 取某用户名最近一条登录审计行。
func lastLoginLog(ctx context.Context, t *gtest.T, username string) gdb.Record {
	rec, err := g.DB().GetOne(ctx,
		"SELECT * FROM admin_login_log WHERE username=? ORDER BY id DESC LIMIT 1", username)
	t.AssertNil(err)
	return rec
}

const (
	tSeedPassword = "Ecboot@Admin2026"
	tSeedHash     = "$2a$10$AX.WGxKFiEoaIxKDzcbdGOIwtl555zjc1zLPFLDoz8mHu80DW1xg2"
)

// TestAdminLogin 登录行为与审计（spec FR-001~003, US1 验收 1~3）。
func TestAdminLogin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		id := seedAdmin(ctx, t, "t_login_a", 0)
		defer cleanupAdmin(ctx, t, "t_login_a")

		// ① 成功: 双凭证 + isSuper=false + 审计成功行 + last_login_time
		res, err := AdminLogin(ctx, "t_login_a", tSeedPassword, "", "", "10.0.0.1", "UA-test")
		t.AssertNil(err)
		t.Assert(res.Token != "", true)
		t.Assert(res.RefreshToken != "", true)
		t.Assert(res.RealName, "t_login_a")
		t.Assert(res.IsSuper, false)
		log1 := lastLoginLog(ctx, t, "t_login_a")
		t.Assert(log1["login_status"].Int(), 1)
		t.Assert(log1["admin_id"].Int64(), id)
		t.Assert(log1["ip"].String(), "10.0.0.1")
		t.Assert(log1["user_agent"].String(), "UA-test")
		rec, err := g.DB().GetOne(ctx, "SELECT last_login_time FROM admin_user WHERE id=?", id)
		t.AssertNil(err)
		t.Assert(rec["last_login_time"].IsNil(), false)

		// ② 密码错误: 80001 + 审计 login_status=2
		_, err = AdminLogin(ctx, "t_login_a", "wrong-pass", "", "", "10.0.0.1", "UA-test")
		t.Assert(errCode(err), errcode.CodeAdminBadCredential)
		t.Assert(lastLoginLog(ctx, t, "t_login_a")["login_status"].Int(), 2)

		// ③ 账号不存在: 80001 + 审计 login_status=3（无 admin_id 回填）
		_, err = AdminLogin(ctx, "t_ghost", tSeedPassword, "", "", "10.0.0.1", "UA-test")
		t.Assert(errCode(err), errcode.CodeAdminBadCredential)
		ghost := lastLoginLog(ctx, t, "t_ghost")
		t.Assert(ghost["login_status"].Int(), 3)
		t.Assert(ghost["admin_id"].IsNil(), true)

		// ④ 禁用账号: 80009 + 审计 login_status=3
		_, err = g.DB().Exec(ctx, "UPDATE admin_user SET status=2 WHERE id=?", id)
		t.AssertNil(err)
		_, err = AdminLogin(ctx, "t_login_a", tSeedPassword, "", "", "10.0.0.1", "UA-test")
		t.Assert(errCode(err), errcode.CodeAdminDisabled)
		t.Assert(lastLoginLog(ctx, t, "t_login_a")["login_status"].Int(), 3)
	})
}

// TestAdminLoginCaptcha 验证码携带即强校验（FR-003, research D7）。
func TestAdminLoginCaptcha(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		seedAdmin(ctx, t, "t_login_c", 0)
		defer cleanupAdmin(ctx, t, "t_login_c")

		key, _, _, err := captcha.Generate(ctx)
		t.AssertNil(err)

		// 错误验证码 → 20001, 不发会话
		_, err = AdminLogin(ctx, "t_login_c", tSeedPassword, key, "zzzz", "10.0.0.1", "UA")
		t.Assert(errCode(err), errcode.CodeCaptchaError)

		// 正确验证码 → 成功
		key2, _, answer2, err := captcha.Generate(ctx)
		t.AssertNil(err)
		res, err := AdminLogin(ctx, "t_login_c", tSeedPassword, key2, answer2, "10.0.0.1", "UA")
		t.AssertNil(err)
		t.Assert(res.Token != "", true)
	})
}

// TestAdminChangePassword 改密（FR-007, US1 验收 7）。
func TestAdminChangePassword(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		id := seedAdmin(ctx, t, "t_pwd_a", 0)
		defer cleanupAdmin(ctx, t, "t_pwd_a")

		// 旧密码错误 → 80002
		err := ChangePassword(ctx, id, "bad-old", "NewPass@2026")
		t.Assert(errCode(err), errcode.CodeOldPasswordWrong)

		// 正确改密 → 旧密码登录被拒, 新密码成功
		t.AssertNil(ChangePassword(ctx, id, tSeedPassword, "NewPass@2026"))
		_, err = AdminLogin(ctx, "t_pwd_a", tSeedPassword, "", "", "10.0.0.1", "UA")
		t.Assert(errCode(err), errcode.CodeAdminBadCredential)
		res, err := AdminLogin(ctx, "t_pwd_a", "NewPass@2026", "", "", "10.0.0.1", "UA")
		t.AssertNil(err)
		t.Assert(res.Token != "", true)
	})
}

// TestAdminProfile 个人信息（FR-006, research D6）：含角色编码列表。
func TestAdminProfile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		id := seedAdmin(ctx, t, "t_prof_a", 0)
		defer cleanupAdmin(ctx, t, "t_prof_a")
		_, err := g.DB().Exec(ctx,
			"INSERT INTO admin_role(name,code,description,status) VALUES('测试角色','t_role_prof','',1)")
		t.AssertNil(err)
		_, err = g.DB().Exec(ctx,
			"DELETE FROM admin_user_role WHERE admin_id=?", id)
		t.AssertNil(err)
		_, err = g.DB().Exec(ctx,
			"INSERT INTO admin_user_role(admin_id,role_id) SELECT ?, id FROM admin_role WHERE code='t_role_prof'", id)
		t.AssertNil(err)
		defer func() {
			_, _ = g.DB().Exec(ctx, "DELETE FROM admin_role WHERE code='t_role_prof'")
			_, _ = g.DB().Exec(ctx, "DELETE FROM admin_user_role WHERE admin_id=?", id)
		}()

		p, err := Profile(ctx, id)
		t.AssertNil(err)
		t.Assert(p.Username, "t_prof_a")
		t.Assert(p.RealName, "t_prof_a")
		t.Assert(len(p.Roles), 1)
		t.Assert(p.Roles[0], "t_role_prof")
	})
}

// TestActiveAdmin 后台账号态校验（spec FR-008）：禁用/软删/不存在 → false。
func TestActiveAdmin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		id := seedAdmin(ctx, t, "t_act_admin", 0)
		defer cleanupAdmin(ctx, t, "t_act_admin")

		ok, err := ActiveAdmin(ctx, id)
		t.AssertNil(err)
		t.Assert(ok, true)

		// 禁用 → false
		_, err = g.DB().Exec(ctx, "UPDATE `admin_user` SET status=2 WHERE id=?", id)
		t.AssertNil(err)
		ok, err = ActiveAdmin(ctx, id)
		t.AssertNil(err)
		t.Assert(ok, false)

		// 恢复启用后软删 → false
		_, err = g.DB().Exec(ctx, "UPDATE `admin_user` SET status=1, deleted=1 WHERE id=?", id)
		t.AssertNil(err)
		ok, err = ActiveAdmin(ctx, id)
		t.AssertNil(err)
		t.Assert(ok, false)

		// 不存在 → false
		ok, err = ActiveAdmin(ctx, 999999999)
		t.AssertNil(err)
		t.Assert(ok, false)
	})
}
