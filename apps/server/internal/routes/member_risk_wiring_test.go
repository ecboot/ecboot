// member_risk_wiring_test.go 会员管理/风控/审计 13 端点可达性与权限挂载（018 批次 12）。
// 非超管 10005 对照（批次 10 I4 教训内化）: 漏挂/错挂 RequirePerm 即红。
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

// TestMemberRiskEndpointsReachable 13 端点可达性。
func TestMemberRiskEndpointsReachable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		s := startTestServer(t, 38881)
		defer func() { _ = s.Shutdown() }()
		base := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())

		adminId := seedRouteAdmin(ctx, t, "ROUTE-MR-ADM")
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM admin_user WHERE id=?", adminId) }()
		token, _, err := security.NewSessionManager("admin", 7).Create(ctx, adminId)
		t.AssertNil(err)

		endpoints := []struct{ method, path, body string }{
			{"GET", "/admin/members", ""},
			{"GET", "/admin/members/999999999", ""},
			{"POST", "/admin/members/999999999/disable", `{"disable":true}`},
			{"POST", "/admin/members/999999999/rebind-phone", `{"newPhone":"13900009999"}`},
			{"GET", "/admin/risk-rules", ""},
			{"POST", "/admin/risk-rules", `{"name":"TF-WIRE-RISK","ruleType":1,"conditionExpr":"user:1","action":1}`},
			{"PUT", "/admin/risk-rules/999999999", `{"name":"X"}`},
			{"GET", "/admin/risk-rules/999999999", ""},
			{"DELETE", "/admin/risk-rules/999999999", ""},
			{"GET", "/admin/risk-records", ""},
			{"POST", "/admin/risk-records/999999999/appeal", `{"pass":true}`},
			{"GET", "/admin/admin-login-logs", ""},
			{"GET", "/admin/operation-logs", ""},
		}
		for _, c := range endpoints {
			body := doReq(ctx, t, base, token, c.method, c.path, c.body)
			t.Assert(!strings.Contains(body, "\"code\":10003"), true)
			t.Assert(!strings.Contains(body, "\"code\":10005"), true)
		}

		// 非超管对照: 每类权限点必须 10005
		noPermId := seedRouteAdmin2(ctx, t, "ROUTE-MR-NOPERM")
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM admin_user WHERE id=?", noPermId) }()
		noPermToken, _, err := security.NewSessionManager("admin", 7).Create(ctx, noPermId)
		t.AssertNil(err)
		probe := []struct{ method, path, body string }{
			{"GET", "/admin/members", ""},
			{"POST", "/admin/members/999999999/disable", `{"disable":true}`},
			{"GET", "/admin/risk-rules", ""},
			{"POST", "/admin/risk-rules", `{"name":"X","ruleType":1,"conditionExpr":"c","action":1}`},
			{"GET", "/admin/risk-records", ""},
			{"POST", "/admin/risk-records/999999999/appeal", `{"pass":true}`},
			{"GET", "/admin/admin-login-logs", ""},
			{"GET", "/admin/operation-logs", ""},
		}
		for _, c := range probe {
			body := doReq(ctx, t, base, noPermToken, c.method, c.path, c.body)
			t.Assert(strings.Contains(body, "\"code\":10005"), true)
		}

		// 清理探活可能创建的规则（精确名）
		_, _ = g.DB().Exec(ctx, "DELETE FROM risk_rule WHERE name='TF-WIRE-RISK'")
	})
}

// TestDashboardEndpointsReachable 批次 13: 3 看板端点可达 + dashboard:read 挂载（非超管 10005 对照）。
func TestDashboardEndpointsReachable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		s := startTestServer(t, 38883)
		defer func() { _ = s.Shutdown() }()
		base := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())

		adminId := seedRouteAdmin(ctx, t, "ROUTE-DB-ADM")
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM admin_user WHERE id=?", adminId) }()
		token, _, err := security.NewSessionManager("admin", 7).Create(ctx, adminId)
		t.AssertNil(err)

		paths := []string{"/admin/dashboard/trade", "/admin/dashboard/member", "/admin/dashboard/product"}
		for _, p := range paths {
			body := doReq(ctx, t, base, token, "GET", p, "")
			t.Assert(!strings.Contains(body, "\"code\":10003"), true)
			t.Assert(!strings.Contains(body, "\"code\":10005"), true)
		}

		// 非超管对照
		noPermId := seedRouteAdmin2(ctx, t, "ROUTE-DB-NOPERM")
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM admin_user WHERE id=?", noPermId) }()
		noPermToken, _, err := security.NewSessionManager("admin", 7).Create(ctx, noPermId)
		t.AssertNil(err)
		for _, p := range paths {
			body := doReq(ctx, t, base, noPermToken, "GET", p, "")
			t.Assert(strings.Contains(body, "\"code\":10005"), true)
		}
	})
}
