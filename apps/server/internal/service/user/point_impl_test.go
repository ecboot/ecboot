package user

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/model"
)

// TestPointAccountAndLogs 积分账户与流水（FR-012/013）: 余额可负如实 / 类型筛选 / 分页。
func TestPointAccountAndLogs(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900006001"
		defer cleanupMember(ctx, t, phone)
		uid := seedMember(ctx, t, phone, "积分测试", 0)
		defer cleanupPoint(ctx, t, uid)

		days := 10
		seedPointAccount(ctx, t, uid, 120, &days)
		seedPointLog(ctx, t, uid, 1, 10, 10, "")
		seedPointLog(ctx, t, uid, 3, -5, 5, "ORD001")
		seedPointLog(ctx, t, uid, 3, -20, -15, "ORD002") // 可致负（欠款语义）

		// 账户: 余额如实（含负）
		acc, err := PointAccount(ctx, uid)
		t.AssertNil(err)
		t.Assert(acc.Balance, 120)
		t.Assert(acc.LastEarnedAt != "", true)

		// 流水: 全部
		res, err := PointLogs(ctx, uid, 0, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.Total, 3)

		// 类型筛选（3=下单消耗 → 2 条）
		res, err = PointLogs(ctx, uid, 3, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.Total, 2)
		for _, it := range res.List {
			t.Assert(it.BizType, 3)
		}
		// 负数流水如实返回
		hasNeg := false
		for _, it := range res.List {
			if it.Points < 0 {
				hasNeg = true
			}
		}
		t.Assert(hasNeg, true)

		// 分页
		res, err = PointLogs(ctx, uid, 0, model.PageReq{Page: 1, PageSize: 2})
		t.AssertNil(err)
		t.Assert(len(res.List), 2)
		t.Assert(res.Total, 3)

		// 余额为负的账户如实展示
		neg := "13900006002"
		defer cleanupMember(ctx, t, neg)
		uidNeg := seedMember(ctx, t, neg, "负余额", 0)
		defer cleanupPoint(ctx, t, uidNeg)
		seedPointAccount(ctx, t, uidNeg, -30, nil)
		accNeg, err := PointAccount(ctx, uidNeg)
		t.AssertNil(err)
		t.Assert(accNeg.Balance, -30)
	})
}

// TestPointEarnConsume 积分变动内部方法（D6）: 变动 + 流水 + 余额快照。
func TestPointEarnConsume(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900006003"
		defer cleanupMember(ctx, t, phone)
		uid := seedMember(ctx, t, phone, "变动测试", 0)
		defer cleanupPoint(ctx, t, uid)
		seedPointAccount(ctx, t, uid, 0, nil)

		t.AssertNil(PointEarn(ctx, uid, 1, 50, ""))   // 签到 +50
		t.AssertNil(PointConsume(ctx, uid, 30, "O1")) // 下单消耗 -30
		t.AssertNil(PointRefund(ctx, uid, 20, "O1"))  // 退款回退 +20

		acc, err := PointAccount(ctx, uid)
		t.AssertNil(err)
		t.Assert(acc.Balance, 40) // 0+50-30+20

		res, err := PointLogs(ctx, uid, 0, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.Total, 3)
	})
}

// TestInviteRecords 邀请记录（FR-014）: 仅本人 + 分页 + 被邀请人脱敏。
func TestInviteRecords(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			phoneA = "13900006010"
			phoneB = "13900006011"
			phoneC = "13900006012"
			phoneD = "13900006013"
		)
		defer cleanupMember(ctx, t, phoneA)
		defer cleanupMember(ctx, t, phoneB)
		defer cleanupMember(ctx, t, phoneC)
		defer cleanupMember(ctx, t, phoneD)
		uidA := seedMember(ctx, t, phoneA, "邀请人", 0)
		uidB := seedMember(ctx, t, phoneB, "被邀请人B", 0)
		uidC := seedMember(ctx, t, phoneC, "被邀请人C", 0)
		uidD := seedMember(ctx, t, phoneD, "被邀请人D", 0)
		defer cleanupInvite(ctx, t, uidA)
		defer cleanupInvite(ctx, t, uidB)

		// 注意: invite_record 有 uk(new_user_id)——一人只能被邀请一次
		seedInviteRecord(ctx, t, uidA, uidB, 1)
		seedInviteRecord(ctx, t, uidA, uidC, 2) // 同 inviter 两条
		seedInviteRecord(ctx, t, uidB, uidD, 1) // 他人记录（不得出现在 A 的列表）

		res, err := InviteRecords(ctx, uidA, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.Total, 2) // 仅 A 自己的
		// 被邀请人昵称脱敏（首字符 + 掩码）
		for _, it := range res.List {
			t.Assert(it.NewUser != "", true)
			t.Assert(it.NewUser != "被邀请人B", true) // 非明文
		}
		// 分页
		res, err = InviteRecords(ctx, uidA, model.PageReq{Page: 1, PageSize: 1})
		t.AssertNil(err)
		t.Assert(len(res.List), 1)
		t.Assert(res.Total, 2)
	})
}
