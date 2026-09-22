// distribution_wiring_test.go 分销资金 23 端点可达性与权限挂载（017 批次 11）。
// I4 教训内化（批次 10）: 非超管对照是权限挂载的真实守卫（超管直通使漏挂结构性不可测）;
// doReq 支持 PUT/DELETE。
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

// TestDistributionEndpointsReachable 23 端点: 超管可达 + 非超管 10005 对照 + user 会员可达/未登录 10003。
func TestDistributionEndpointsReachable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		s := startTestServer(t, 38879)
		defer func() { _ = s.Shutdown() }()
		base := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())

		adminId := seedRouteAdmin(ctx, t, "ROUTE-DIST-ADM")
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM admin_user WHERE id=?", adminId) }()
		token, _, err := security.NewSessionManager("admin", 7).Create(ctx, adminId)
		t.AssertNil(err)

		// admin 12 端点: 超管 → 不回 10003/10005
		adminEndpoints := []struct{ method, path, body string }{
			{"GET", "/admin/distributors", ""},
			{"POST", "/admin/distributors/999999999/audit", `{"pass":true}`},
			{"POST", "/admin/distributors/999999999/freeze", `{"freeze":true}`},
			{"GET", "/admin/commission-rules", ""},
			{"POST", "/admin/commission-rules", `{"scopeType":1,"scopeId":"1","level1Rate":"5.00","level2Rate":"2.00"}`},
			{"PUT", "/admin/commission-rules/999999999", `{"level1Rate":"6.00"}`},
			{"DELETE", "/admin/commission-rules/999999999", ""},
			{"GET", "/admin/commission-records", ""},
			{"GET", "/admin/withdraws", ""},
			{"POST", "/admin/withdraws/WD-X/audit", `{"pass":false}`},
			{"POST", "/admin/withdraws/WD-X/pay", `{"success":false}`},
			{"GET", "/admin/invite-records", ""},
		}
		for _, c := range adminEndpoints {
			body := doReq(ctx, t, base, token, c.method, c.path, c.body)
			t.Assert(!strings.Contains(body, "\"code\":10003"), true)
			t.Assert(!strings.Contains(body, "\"code\":10005"), true)
		}

		// 非超管对照（I4 守卫）: 每类权限点至少一端点必须 10005
		noPermId := seedRouteAdmin2(ctx, t, "ROUTE-DIST-NOPERM")
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM admin_user WHERE id=?", noPermId) }()
		noPermToken, _, err := security.NewSessionManager("admin", 7).Create(ctx, noPermId)
		t.AssertNil(err)
		probe := []struct{ method, path, body string }{
			{"GET", "/admin/distributors", ""},
			{"POST", "/admin/distributors/999999999/audit", `{"pass":true}`},
			{"GET", "/admin/commission-rules", ""},
			{"POST", "/admin/commission-rules", `{"scopeType":1,"scopeId":"1","level1Rate":"5.00"}`},
			{"GET", "/admin/withdraws", ""},
			{"POST", "/admin/withdraws/WD-X/audit", `{"pass":true}`},
			{"POST", "/admin/withdraws/WD-X/pay", `{"success":true,"channelOrderNo":"X"}`},
			{"GET", "/admin/invite-records", ""},
		}
		for _, c := range probe {
			body := doReq(ctx, t, base, noPermToken, c.method, c.path, c.body)
			t.Assert(strings.Contains(body, "\"code\":10005"), true) // 权限点真实在位
		}

		// user 渠道: 会员凭证 → 不回 10003; 未登录 → 10003
		uid := seedRouteUser(ctx, t, "ROUTE-DIST-U1")
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE id=?", uid) }()
		userToken, _, err := security.NewSessionManager("user", 7).Create(ctx, uid)
		t.AssertNil(err)
		userEndpoints := []struct{ method, path, body string }{
			{"GET", "/user/distribution/status", ""},
			{"GET", "/user/distribution/relations", ""},
			{"GET", "/user/distribution/commission-rules?spuId=1", ""},
			{"GET", "/user/distribution/records", ""},
			{"GET", "/user/distribution/account", ""},
			{"GET", "/user/distribution/account/logs", ""},
			{"GET", "/user/distribution/withdraws", ""},
			{"GET", "/user/distribution/invite-records", ""},
			{"GET", "/user/share-code", ""},
			{"POST", "/user/distribution/withdraws", `{"amount":"1.00"}`},
		}
		for _, c := range userEndpoints {
			body := doReq(ctx, t, base, userToken, c.method, c.path, c.body)
			t.Assert(!strings.Contains(body, "\"code\":10003"), true)
		}
		body := doReq(ctx, t, base, "", "GET", "/user/distribution/account", "")
		t.Assert(strings.Contains(body, "\"code\":10003"), true)

		// 清理探活可能创建的数据（精确名/精确键）
		_, _ = g.DB().Exec(ctx, "DELETE FROM commission_rule WHERE scope_id=1 AND scope_type=1")
	})
}
