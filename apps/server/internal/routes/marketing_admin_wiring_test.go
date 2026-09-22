// marketing_admin_wiring_test.go 营销后台 30 端点的**可达性与权限挂载**（016 批次 10）。
// 为什么必须走 HTTP: RequirePerm 是 controller 首行显式挂接（research D2）——漏挂/挂错权限点
// 只有真实路由 + 真实凭证才能暴露（批次 08/09 的 C1 同源教训: service 层测试看不到中间件跳）。
// 守卫规则: 未登录 → 10003; 超管（is_super 直通, 000035 种子语义）合法凭证 → 不得回 10003/10005。
package routes

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/library/security"
)

// seedRouteAdmin 造超管（is_super=1 直通权限判定, 不依赖 RBAC 种子数据）。
func seedRouteAdmin(ctx context.Context, t *gtest.T, username string) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM admin_user WHERE username=?", username)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO admin_user(username,password_hash,real_name,is_super,status) VALUES(?,?,?,?,1)",
		username, "x", "端点探活超管", 1)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

// TestMarketingAdminEndpointsReachable 30 端点可达性（每类抽全 4 动作形态: GET/POST/PUT/DELETE）。
func TestMarketingAdminEndpointsReachable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		s := startTestServer(t, 38877)
		defer func() { _ = s.Shutdown() }()
		base := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())

		adminId := seedRouteAdmin(ctx, t, "ROUTE-MKT-ADMIN")
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM admin_user WHERE id=?", adminId) }()
		token, _, err := security.NewSessionManager("admin", 7).Create(ctx, adminId)
		t.AssertNil(err)

		// 30 端点全量探活（带合法超管凭证 → 不得 10003/10005; 业务错如 50003"不存在"是正常回包）
		endpoints := []struct{ method, path, body string }{
			{"GET", "/admin/coupons", ""},
			{"POST", "/admin/coupons", `{"name":"TF-WIRE","type":2,"discount":"5.00","validType":2,"validDays":7}`},
			{"GET", "/admin/coupons/999999999", ""},
			{"PUT", "/admin/coupons/999999999", `{"status":0}`},
			{"DELETE", "/admin/coupons/999999999", ""},
			{"GET", "/admin/coupons/999999999/records", ""},
			{"GET", "/admin/full-reductions", ""},
			{"POST", "/admin/full-reductions", `{"name":"TF-WIRE-FR","startTime":"2026-01-01 00:00:00","endTime":"2026-12-31 00:00:00","ladders":[{"threshold":"100.00","discount":"10.00"}]}`},
			{"GET", "/admin/full-reductions/999999999", ""},
			{"PUT", "/admin/full-reductions/999999999", `{"name":"TF-WIRE-FR","startTime":"2026-01-01 00:00:00","endTime":"2026-12-31 00:00:00","status":1}`},
			{"DELETE", "/admin/full-reductions/999999999", ""},
			{"GET", "/admin/group-buys", ""},
			{"POST", "/admin/group-buys", `{"name":"TF-WIRE-GB","spuId":"1","groupSize":3,"startTime":"2026-01-01 00:00:00","endTime":"2026-12-31 00:00:00"}`},
			{"PUT", "/admin/group-buys/999999999", `{"status":0}`},
			{"PUT", "/admin/group-buys/999999999/items", `{"items":[]}`},
			{"DELETE", "/admin/group-buys/999999999", ""},
			{"GET", "/admin/flash-sales", ""},
			{"POST", "/admin/flash-sales", `{"name":"TF-WIRE-FS","startTime":"2026-01-01 00:00:00","endTime":"2026-12-31 00:00:00"}`},
			{"PUT", "/admin/flash-sales/999999999", `{"status":0}`},
			{"PUT", "/admin/flash-sales/999999999/items", `{"items":[]}`},
			{"DELETE", "/admin/flash-sales/999999999", ""},
			{"GET", "/admin/bargains", ""},
			{"POST", "/admin/bargains", `{"name":"TF-WIRE-BA","spuId":"1","startTime":"2026-01-01 00:00:00","endTime":"2026-12-31 00:00:00"}`},
			{"PUT", "/admin/bargains/999999999", `{"status":0}`},
			{"PUT", "/admin/bargains/999999999/items", `{"items":[]}`},
			{"DELETE", "/admin/bargains/999999999", ""},
			{"GET", "/admin/assists", ""},
			{"POST", "/admin/assists", `{"name":"TF-WIRE-AS","rewardType":2,"pointAmount":100,"requiredCount":3,"startTime":"2026-01-01 00:00:00","endTime":"2026-12-31 00:00:00"}`},
			{"PUT", "/admin/assists/999999999", `{"status":0}`},
			{"DELETE", "/admin/assists/999999999", ""},
		}
		for _, c := range endpoints {
			body := doReq(ctx, t, base, token, c.method, c.path, c.body)
			t.Assert(!strings.Contains(body, "\"code\":10003"), true) // 未被鉴权拦下
			t.Assert(!strings.Contains(body, "\"code\":10005"), true) // 权限点已挂且超管直通
		}

		// 对照 1: 未登录 → 必须 10003（证明鉴权在位）
		body := doReq(ctx, t, base, "", "GET", "/admin/coupons", "")
		t.Assert(strings.Contains(body, "\"code\":10003"), true)

		// 清理探活可能创建成功的 TF-WIRE 数据（按精确名, 子表先删）
		_, _ = g.DB().Exec(ctx, "DELETE FROM promotion_activity_scope WHERE activity_id IN (SELECT id FROM promotion_activity WHERE name LIKE 'TF-WIRE%')")
		_, _ = g.DB().Exec(ctx, "DELETE FROM promotion_activity_ladder WHERE activity_id IN (SELECT id FROM promotion_activity WHERE name LIKE 'TF-WIRE%')")
		_, _ = g.DB().Exec(ctx, "DELETE FROM promotion_activity WHERE name LIKE 'TF-WIRE%'")
		_, _ = g.DB().Exec(ctx, "DELETE FROM coupon WHERE name LIKE 'TF-WIRE%'")
		_, _ = g.DB().Exec(ctx, "DELETE FROM group_buy_activity WHERE name LIKE 'TF-WIRE%'")
		_, _ = g.DB().Exec(ctx, "DELETE FROM flash_sale_activity WHERE name LIKE 'TF-WIRE%'")
		_, _ = g.DB().Exec(ctx, "DELETE FROM bargain_activity WHERE name LIKE 'TF-WIRE%'")
		_, _ = g.DB().Exec(ctx, "DELETE FROM assist_activity WHERE name LIKE 'TF-WIRE%'")
	})
}
