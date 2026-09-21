package user

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/model"
)

// TestProfileDetailUser 资料查询（FR-001）: 手机号脱敏 + 等级名按成长值匹配 + 成长值。
func TestProfileDetailUser(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900002001"
		defer cleanupMember(ctx, t, phone)
		cleanupLevelRule(ctx, t, "T青铜", "T白银", "T黄金")
		defer cleanupLevelRule(ctx, t, "T青铜", "T白银", "T黄金")

		uid := seedMember(ctx, t, phone, "昵称A", 150)
		seedLevelRule(ctx, t, "T青铜", 0)
		seedLevelRule(ctx, t, "T白银", 100)
		seedLevelRule(ctx, t, "T黄金", 500)

		d, err := ProfileDetail(ctx, uid)
		t.AssertNil(err)
		t.Assert(d.Nickname, "昵称A")
		t.Assert(d.GrowthValue, 150)
		// 手机号脱敏: 前3+****+后4（明文绝不出参）
		t.Assert(d.PhoneMasked, "139****2001")
		t.Assert(d.PhoneMasked != phone, true)
		// 等级: 150 落在 [100,500) → 白银（取最大满足门槛者）
		t.Assert(d.LevelName, "T白银")
	})
}

// TestProfileUpdateUser 资料修改（FR-002）: 改动生效; 未传字段保持不变。
func TestProfileUpdateUser(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900002002"
		defer cleanupMember(ctx, t, phone)
		uid := seedMember(ctx, t, phone, "旧昵称", 0)

		// 预置头像与性别, 改名时不传它们 → 须保持不变
		t.AssertNil(ProfileUpdate(ctx, uid, model.ProfileUpdateInput{
			Nickname: "新昵称", Avatar: "http://img/a.png", Gender: 1,
		}))
		d, err := ProfileDetail(ctx, uid)
		t.AssertNil(err)
		t.Assert(d.Nickname, "新昵称")
		t.Assert(d.Avatar, "http://img/a.png")
		t.Assert(d.Gender, 1)

		// 仅改性别 → 昵称/头像不变
		t.AssertNil(ProfileUpdate(ctx, uid, model.ProfileUpdateInput{Gender: 2}))
		d, err = ProfileDetail(ctx, uid)
		t.AssertNil(err)
		t.Assert(d.Nickname, "新昵称")
		t.Assert(d.Avatar, "http://img/a.png")
		t.Assert(d.Gender, 2)
	})
}

// TestLoginLogsUser 登录记录（FR-003）: 近 30 天窗口 + 时间倒序 + 分页。
func TestLoginLogsUser(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900002003"
		defer cleanupMember(ctx, t, phone)
		uid := seedMember(ctx, t, phone, "登录测试", 0)

		seedLoginLog(ctx, t, uid, 1, 1, 1)  // 1 天前（窗口内）
		seedLoginLog(ctx, t, uid, 5, 1, 1)  // 5 天前（窗口内）
		seedLoginLog(ctx, t, uid, 40, 1, 1) // 40 天前（窗口外, 不得出现）
		// 他人日志（不得出现）
		other := seedMember(ctx, t, "13900002999", "他人", 0)
		defer cleanupMember(ctx, t, "13900002999")
		seedLoginLog(ctx, t, other, 1, 1, 1)

		res, err := LoginLogs(ctx, uid, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.Total, 2) // 30 天窗口内仅 2 条（他人与超窗口均排除）
		t.Assert(len(res.List), 2)
		// 倒序: 最近的在最前
		t.Assert(res.List[0].CreatedAt > res.List[1].CreatedAt, true)
		// 字段映射
		t.Assert(res.List[0].Ip, "127.0.*.*") // IP 脱敏

		// 分页
		res, err = LoginLogs(ctx, uid, model.PageReq{Page: 1, PageSize: 1})
		t.AssertNil(err)
		t.Assert(len(res.List), 1)
		t.Assert(res.Total, 2)
	})
}
