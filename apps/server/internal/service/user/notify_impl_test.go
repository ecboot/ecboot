package user

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// TestMessages 站内信（FR-009）: 未读计数（全量非当前页）+ 已读态筛选 + 分页。
func TestMessages(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900005001"
		defer cleanupMember(ctx, t, phone)
		uid := seedMember(ctx, t, phone, "消息测试", 0)
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM user_message WHERE user_id=?", uid) }()

		seedMessage(ctx, t, uid, "未读1", false)
		seedMessage(ctx, t, uid, "未读2", false)
		seedMessage(ctx, t, uid, "已读1", true)

		// 全部: total=3, unread=2
		res, err := Messages(ctx, uid, -1, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.Total, 3)
		t.Assert(res.UnreadCount, 2)

		// 仅未读: total=2
		res, err = Messages(ctx, uid, 0, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.Total, 2)
		for _, it := range res.List {
			t.Assert(it.IsRead, false)
		}
		// 未读计数仍是全量的 2（不受筛选影响——"全量非当前页"）
		t.Assert(res.UnreadCount, 2)

		// 仅已读: total=1
		res, err = Messages(ctx, uid, 1, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.Total, 1)

		// 分页
		res, err = Messages(ctx, uid, -1, model.PageReq{Page: 1, PageSize: 2})
		t.AssertNil(err)
		t.Assert(len(res.List), 2)
		t.Assert(res.Total, 3)
	})
}

// TestMarkRead 标记已读（FR-009）: 单条（归属校验）/全部。
func TestMarkRead(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900005002"
		defer cleanupMember(ctx, t, phone)
		uid := seedMember(ctx, t, phone, "已读测试", 0)
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM user_message WHERE user_id=?", uid) }()

		m1 := seedMessage(ctx, t, uid, "M1", false)
		m2 := seedMessage(ctx, t, uid, "M2", false)

		// 单条已读 → 计数 1
		t.AssertNil(MarkRead(ctx, uid, m1))
		res, err := Messages(ctx, uid, -1, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.UnreadCount, 1)

		// 他人消息 → 10006
		other := seedMember(ctx, t, "13900005999", "他人M", 0)
		defer cleanupMember(ctx, t, "13900005999")
		err = MarkRead(ctx, other, m2)
		t.Assert(errCode(err), errcode.CodeNotFound)

		// 全部已读 → 计数 0
		t.AssertNil(MarkAllRead(ctx, uid))
		res, err = Messages(ctx, uid, -1, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.UnreadCount, 0)
	})
}

// TestNotifyPreferences 通知偏好（FR-010/011）: 未设置=全开; 设置幂等; 两渠道。
func TestNotifyPreferences(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900005003"
		defer cleanupMember(ctx, t, phone)
		uid := seedMember(ctx, t, phone, "偏好测试", 0)

		// 未设置 → 默认全开（两渠道）
		prefs, err := Preferences(ctx, uid)
		t.AssertNil(err)
		t.Assert(len(prefs), 2)
		for _, p := range prefs {
			t.Assert(p.Enabled, true)
		}

		// 关闭短信（channel=2）→ 再查反映
		t.AssertNil(SetPreferences(ctx, uid, []model.NotifyPreference{
			{Channel: 2, Enabled: false},
		}))
		prefs, err = Preferences(ctx, uid)
		t.AssertNil(err)
		for _, p := range prefs {
			if p.Channel == 2 {
				t.Assert(p.Enabled, false)
			}
			if p.Channel == 1 {
				t.Assert(p.Enabled, true) // 未设置渠道保持默认开
			}
		}

		// 幂等: 重复设置同值不产生重复行
		t.AssertNil(SetPreferences(ctx, uid, []model.NotifyPreference{{Channel: 2, Enabled: false}}))
		prefs, err = Preferences(ctx, uid)
		t.AssertNil(err)
		t.Assert(len(prefs), 2) // 仍两行（每渠道至多一行）
		n, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM user_notify_preference WHERE user_id=? AND channel=2", uid)
		t.AssertNil(err)
		t.Assert(n.Int(), 1)
	})
}
