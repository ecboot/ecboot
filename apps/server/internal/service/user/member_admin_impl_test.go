// member_admin_impl_test.go 会员管理（018 批次 12）——治理动作与唯一键防线。
package user

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// TestMemberAdminLifecycle 列表脱敏/禁用启用/详情资产/改绑唯一键（FR-1）。
func TestMemberAdminLifecycle(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		uid := distCleanupUser(ctx, t, "TF-MA-1")
		defer distCleanupAll(ctx, uid)
		// fixture 补真实密文: phone 列存 Encrypt 密文（distCleanupUser 插的是明文占位,
		// 解密失败 → 列表 Phone 为空——脱敏断言需要真实密文走通 memberPhone 链路）
		enc, err := phoneCipher().Encrypt("13900001111")
		t.AssertNil(err)
		_, err = g.DB().Exec(ctx, "UPDATE `user` SET phone=?, phone_hash=? WHERE id=?", enc, phoneCipher().Hash("13900001111"), uid)
		t.AssertNil(err)
		logic := NewMemberAdminLogic()

		// 列表可见 + 手机号脱敏（前3后4）
		list, err := logic.AdminList(ctx, "", 0, "", model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		var phone string
		found := false
		for _, it := range list.List {
			if it.UserId == uid {
				found = true
				phone = it.Phone
			}
		}
		t.Assert(found, true)
		t.Assert(len(phone) == 11 && phone[3:7] == "****", true)

		// 详情: 资产与订单统计字段齐
		d, err := logic.AdminDetail(ctx, uid)
		t.AssertNil(err)
		t.Assert(d.Assets != nil && d.OrderStats != nil, true)
		_, err = logic.AdminDetail(ctx, 999999999)
		t.Assert(errCode(err), errcode.CodeNotFound)

		// 禁用 → 详情 status=2（api 口径）→ 启用恢复
		t.AssertNil(logic.Disable(ctx, uid, true, "测试禁用"))
		d, _ = logic.AdminDetail(ctx, uid)
		t.Assert(d.Status, 2)
		t.AssertNil(logic.Disable(ctx, uid, false, ""))
		d, _ = logic.AdminDetail(ctx, uid)
		t.Assert(d.Status, 1)

		// 改绑: 非法号 → 10001; 合法 → phone/hash 更新; 已占号 → 业务码
		t.Assert(errCode(logic.RebindPhone(ctx, uid, "20000000000")), errcode.CodeInvalidParam)
		t.AssertNil(logic.RebindPhone(ctx, uid, "13900001111"))
		d, _ = logic.AdminDetail(ctx, uid)
		t.Assert(d.Phone, "139****1111")
		// 同号再绑到另一用户 → uk_phone_hash 1062 → 业务码
		uid2 := distCleanupUser(ctx, t, "TF-MA-2")
		defer distCleanupAll(ctx, uid2)
		t.Assert(errCode(logic.RebindPhone(ctx, uid2, "13900001111")), errcode.CodeInvalidParam)
		// 不存在用户 → 404 类
		t.Assert(errCode(logic.RebindPhone(ctx, 999999999, "13900002222")), errcode.CodeNotFound)
	})
}
