package user

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/model"
)

// ---- 内部方法测试（011 评审 I4: 此前 6 个内部方法零覆盖） ----

// cleanupNotifyTasks 清理通知任务。
func cleanupNotifyTasks(ctx context.Context, t *gtest.T, userId int64) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM notify_task WHERE user_id=?", userId)
}

// TestGrowthAndLevel 成长值累加与等级重算（内部方法）: 跨门槛升级; 无规则时不动。
func TestGrowthAndLevel(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900007001"
		defer cleanupMember(ctx, t, phone)
		cleanupLevelRule(ctx, t, "I青", "I银", "I金")
		defer cleanupLevelRule(ctx, t, "I青", "I银", "I金")
		uid := seedMember(ctx, t, phone, "成长测试", 50)

		seedLevelRule(ctx, t, "I青", 0)
		seedLevelRule(ctx, t, "I银", 100)
		seedLevelRule(ctx, t, "I金", 500)

		// 成长值 50 → 跨过 100 → 升到银
		t.AssertNil(GrowthAdd(ctx, uid, 60))
		d, err := ProfileDetail(ctx, uid)
		t.AssertNil(err)
		t.Assert(d.GrowthValue, 110)
		t.Assert(d.LevelName, "I银")

		// 再跨 500 → 金
		t.AssertNil(GrowthAdd(ctx, uid, 400))
		d, err = ProfileDetail(ctx, uid)
		t.AssertNil(err)
		t.Assert(d.GrowthValue, 510)
		t.Assert(d.LevelName, "I金")

		// 非法增量（<=0）静默无操作
		t.AssertNil(GrowthAdd(ctx, uid, 0))
		d, err = ProfileDetail(ctx, uid)
		t.AssertNil(err)
		t.Assert(d.GrowthValue, 510)
	})
}

// TestFootprintRecordAndClean 足迹 UPSERT 与过期清理（内部方法）。
func TestFootprintRecordAndClean(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			phone = "13900007002"
			sfx   = "intfp"
		)
		defer cleanupMember(ctx, t, phone)
		defer cleanupSpuForMember(ctx, t, sfx)
		uid := seedMember(ctx, t, phone, "足迹内部", 0)
		defer cleanupMemberCollections(ctx, t, uid)
		spuId := seedSpuForMember(ctx, t, sfx, "5.00", `["http://img/i.png"]`, 1)

		// UPSERT: 三次浏览 → 一行, view_count=3
		for i := 0; i < 3; i++ {
			t.AssertNil(FootprintRecord(ctx, uid, spuId))
		}
		n, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM user_footprint WHERE user_id=? AND spu_id=?", uid, spuId)
		t.AssertNil(err)
		t.Assert(n.Int(), 1)
		vc, err := g.DB().GetValue(ctx,
			"SELECT view_count FROM user_footprint WHERE user_id=? AND spu_id=?", uid, spuId)
		t.AssertNil(err)
		t.Assert(vc.Int(), 3)

		// 过期清理: 把 last_view_at 改到 100 天前 → 被清
		_, err = g.DB().Exec(ctx,
			"UPDATE user_footprint SET last_view_at=DATE_SUB(NOW(), INTERVAL 100 DAY) WHERE user_id=?", uid)
		t.AssertNil(err)
		cleaned, err := FootprintCleanExpired(ctx)
		t.AssertNil(err)
		t.AssertGE(cleaned, int64(1))
		n, err = g.DB().GetValue(ctx, "SELECT COUNT(*) FROM user_footprint WHERE user_id=?", uid)
		t.AssertNil(err)
		t.Assert(n.Int(), 0)
	})
}

// TestEnqueueAndDispatch 通知入队与投递（内部方法）: 按偏好拆渠道; 站内信落库。
func TestEnqueueAndDispatch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900007003"
		defer cleanupMember(ctx, t, phone)
		uid := seedMember(ctx, t, phone, "通知内部", 0)
		defer cleanupNotifyTasks(ctx, t, uid)
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM user_message WHERE user_id=?", uid) }()

		// 关闭短信（channel=2）→ Enqueue 只建 站内信(3) + 小程序(1)
		t.AssertNil(SetPreferences(ctx, uid, []model.NotifyPreference{{Channel: 2, Enabled: false}}))
		t.AssertNil(Enqueue(ctx, uid, 1, "ORD-INT-1", "ORDER_PAID", map[string]any{"orderNo": "ORD-INT-1"}))
		chans, err := g.DB().GetArray(ctx,
			"SELECT channel FROM notify_task WHERE user_id=? ORDER BY channel", uid)
		t.AssertNil(err)
		t.Assert(len(chans), 2) // 1 与 3（2 被偏好关闭, 不建行）

		// 投递站内信任务 → user_message 落库
		taskId, err := g.DB().GetValue(ctx,
			"SELECT id FROM notify_task WHERE user_id=? AND channel=3 ORDER BY id DESC LIMIT 1", uid)
		t.AssertNil(err)
		t.AssertNil(DispatchTask(ctx, taskId.Int64()))
		mc, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM user_message WHERE user_id=?", uid)
		t.AssertNil(err)
		t.Assert(mc.Int(), 1)
		// 任务状态 → 20 已发送
		st, err := g.DB().GetValue(ctx, "SELECT status FROM notify_task WHERE id=?", taskId.Int64())
		t.AssertNil(err)
		t.Assert(st.Int(), 20)
		// 重复投递幂等（已非 10 状态 → 无操作, 不重复落库）
		t.AssertNil(DispatchTask(ctx, taskId.Int64()))
		mc, err = g.DB().GetValue(ctx, "SELECT COUNT(*) FROM user_message WHERE user_id=?", uid)
		t.AssertNil(err)
		t.Assert(mc.Int(), 1)
	})
}

// TestPointExpireDormant 积分滚动过期（内部方法）: 超 12 个月未获得且余额>0 → 清零 + bizType=9 流水。
func TestPointExpireDormant(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			phoneOld = "13900007004"
			phoneNew = "13900007005"
		)
		defer cleanupMember(ctx, t, phoneOld)
		defer cleanupMember(ctx, t, phoneNew)
		uidOld := seedMember(ctx, t, phoneOld, "过期", 0)
		uidNew := seedMember(ctx, t, phoneNew, "未过期", 0)
		defer cleanupPoint(ctx, t, uidOld)
		defer cleanupPoint(ctx, t, uidNew)

		old := 400 // 天
		seedPointAccount(ctx, t, uidOld, 88, &old)
		recent := 10
		seedPointAccount(ctx, t, uidNew, 66, &recent)

		n, err := PointExpireDormant(ctx)
		t.AssertNil(err)
		t.AssertGE(n, int64(1))

		// 过期账户清零 + 流水 bizType=9
		acc, err := PointAccount(ctx, uidOld)
		t.AssertNil(err)
		t.Assert(acc.Balance, 0)
		lg, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM point_log WHERE user_id=? AND biz_type=9", uidOld)
		t.AssertNil(err)
		t.Assert(lg.Int(), 1)

		// 未过期账户不动
		acc2, err := PointAccount(ctx, uidNew)
		t.AssertNil(err)
		t.Assert(acc2.Balance, 66)
	})
}
