// marketing_wiring_test.go 营销 C 端 12 端点的**可达性**（015 批次 09）。
//
// 为什么必须走 HTTP: 本批动作端点 `POST /shop/bargains/{id}/cut`、`POST /shop/assists/{id}/helpers`
// 与**公开的进度查询**共享前缀——前缀级白名单会把会员动作一并当公开路径直通（Auth 不注入 ctx userId）
// → 控制器 requireMember 恒判未登录（批次 08 的 C1 同型）。这是 service 层测试永远看不到的一跳。
// 守卫规则: 公开端点不得回 10003（游客可访问）; 会员端点**带合法凭证**不得回 10003（身份已注入）;
// 会员端点**不带凭证**必须回 10003（证明确实走了鉴权, 而非被白名单直通）。
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

// doReq 发一次请求并返回响应体（token 为空则不带凭证）。
func doReq(ctx context.Context, t *gtest.T, base, token, method, path, body string) string {
	cli := g.Client()
	if token != "" {
		cli = cli.SetHeader("Authorization", "Bearer "+token)
	}
	if method == "POST" {
		res, err := cli.Post(ctx, base+path, body)
		t.AssertNil(err)
		if res == nil {
			return ""
		}
		defer func() { _ = res.Close() }()
		return res.ReadAllString()
	}
	res, err := cli.Get(ctx, base+path)
	t.AssertNil(err)
	if res == nil {
		return ""
	}
	defer func() { _ = res.Close() }()
	return res.ReadAllString()
}

// TestMarketingEndpointsReachable 12 端点可达性。
func TestMarketingEndpointsReachable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		s := startTestServer(t, 38875)
		defer func() { _ = s.Shutdown() }()
		base := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())

		uid := seedRouteUser(ctx, t, "ROUTE-MKT-1")
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE id=?", uid) }()
		token, _, err := security.NewSessionManager("user", 7).Create(ctx, uid)
		t.AssertNil(err)

		// 公开端点（游客, 不带凭证）: 不得回"未登录"（业务错如"活动不存在"是正常的）
		public := []struct{ method, path, body string }{
			{"GET", "/shop/index", ""},
			{"GET", "/shop/activities/flash-sales", ""},
			{"GET", "/shop/activities/group-buys", ""},
			{"GET", "/shop/activities/bargains", ""},
			{"GET", "/shop/activities/assists", ""},
			{"GET", "/shop/full-reductions", ""},
			{"GET", "/shop/bargains/999999999", ""}, // 砍价进度
			{"GET", "/shop/assists/999999999", ""},  // 助力进度
		}
		for _, c := range public {
			body := doReq(ctx, t, base, "", c.method, c.path, c.body)
			t.Assert(strings.Contains(body, "\"code\":10003"), false)
		}

		// 会员动作端点（带合法凭证）: 身份必须已注入 ctx（否则 10003）
		member := []struct{ method, path, body string }{
			{"POST", "/shop/bargains", `{"bargainItemId":"999999999"}`},
			{"POST", "/shop/bargains/999999999/cut", ""},
			{"POST", "/shop/assists", `{"activityId":"999999999"}`},
			{"POST", "/shop/assists/999999999/helpers", ""},
		}
		for _, c := range member {
			body := doReq(ctx, t, base, token, c.method, c.path, c.body)
			t.Assert(strings.Contains(body, "\"code\":10003"), false)
		}

		// 对照: 会员动作**不带凭证** → 必须 10003（证明这条路径确实走了鉴权, 而非被白名单直通）。
		// M12（015 评审）: 原先只对照 4 个会员端点中的 1 个, 现全部覆盖。
		for _, c := range member {
			body := doReq(ctx, t, base, "", c.method, c.path, c.body)
			t.Assert(strings.Contains(body, "\"code\":10003"), true)
		}
	})
}
