package shop

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

// seedPointUser 建测试会员（仅需满足 user 表非空列; 本测试不涉解密）。
func seedPointUser(ctx context.Context, t *gtest.T, hash string) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE phone_hash=?", hash)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO `user`(nickname,phone,phone_hash,growth_value,status) VALUES('积分测试','x',?,0,1)", hash)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

func cleanupPointUser(ctx context.Context, t *gtest.T, hash string) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM point_account WHERE user_id IN (SELECT id FROM `user` WHERE phone_hash=?)", hash)
	_, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE phone_hash=?", hash)
}

// TestCalcPointDeductFen 积分抵扣（FR-001 前置修复）:
// 有余额且勾选 → min(余额, 上限); 未勾选/上限<=0/无账户行/余额<=0 → 0（且不报错）。
func TestCalcPointDeductFen(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "PT-FIX-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		_, err := g.DB().Exec(ctx, "INSERT INTO point_account(user_id,balance) VALUES(?,500)", uid)
		t.AssertNil(err)

		// 余额 500、上限 300 → 300（修复前: 查询报错 → 恒 0）
		t.Assert(calcPointDeductFen(ctx, uid, true, 300), int64(300))
		// 上限大于余额 → 受余额约束
		t.Assert(calcPointDeductFen(ctx, uid, true, 800), int64(500))
		// 未勾选 → 0
		t.Assert(calcPointDeductFen(ctx, uid, false, 300), int64(0))
		// 上限 <= 0 → 0
		t.Assert(calcPointDeductFen(ctx, uid, true, 0), int64(0))

		// 无账户行 → 0（无行不是错误）
		const h2 = "PT-FIX-2"
		defer cleanupPointUser(ctx, t, h2)
		uid2 := seedPointUser(ctx, t, h2)
		t.Assert(calcPointDeductFen(ctx, uid2, true, 300), int64(0))

		// 余额 0（含负）→ 0
		_, err = g.DB().Exec(ctx, "UPDATE point_account SET balance=0 WHERE user_id=?", uid)
		t.AssertNil(err)
		t.Assert(calcPointDeductFen(ctx, uid, true, 300), int64(0))
		_, err = g.DB().Exec(ctx, "UPDATE point_account SET balance=-50 WHERE user_id=?", uid)
		t.AssertNil(err)
		t.Assert(calcPointDeductFen(ctx, uid, true, 300), int64(0))
	})
}
