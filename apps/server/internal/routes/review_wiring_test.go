// review_wiring_test.go 评价域 4 端点的**可达性**（014 评审 C1 的回归守卫）。
//
// 为什么必须在 HTTP 层测: C1 的缺陷是"鉴权中间件把该端点当公开路径直通、**不注入 ctx userId**"
// → 控制器 requireMember 恒判未登录（10003）。service 层测试直调 service（自带 userId 参数），
// 永远看不到这一跳; 只有走真实路由 + 真实中间件链 + 真实 controller + **合法会员凭证**才能发现。
// 教训（与批次 06 的 C4a/C4b 同源）: "桩清零"只证明没有 CodeNotImplemented, 不证明端点可达。
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

// seedRouteUser 造一个会员（会话只需 userId, 无外键依赖）。
func seedRouteUser(ctx context.Context, t *gtest.T, phoneHash string) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE phone_hash=?", phoneHash)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO `user`(nickname,phone,phone_hash,growth_value,status) VALUES('端点探活','x',?,0,1)", phoneHash)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

// TestReviewEndpointsReachable 评价域 4 端点对**合法会员**必须可达（不得回 10003）。
func TestReviewEndpointsReachable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		s := startTestServer(t, 38873)
		defer func() { _ = s.Shutdown() }()
		base := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())

		uid := seedRouteUser(ctx, t, "ROUTE-REVIEW-1")
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE id=?", uid) }()
		token, _, err := security.NewSessionManager("user", 7).Create(ctx, uid)
		t.AssertNil(err)

		cases := []struct {
			method, path, body string
		}{
			{"GET", "/shop/reviews/mine", ""},                          // C1 受害者①
			{"POST", "/shop/reviews/1/extra", `{"content":"端点探活"}`},    // C1 受害者②
			{"POST", "/shop/reviews", `{"orderItemId":"1","score":5}`}, // 提交评价（原精确例外）
			{"GET", "/shop/products/970000000/reviews", ""},            // 公开的商品评价列表
			{"GET", "/shop/orders", ""},                                // 对照: 已知需会员的端点
		}
		for _, c := range cases {
			cli := g.Client().SetHeader("Authorization", "Bearer "+token)
			var body string
			var e error
			if c.method == "POST" {
				res, re := cli.Post(ctx, base+c.path, c.body)
				e = re
				if res != nil {
					body = res.ReadAllString()
					_ = res.Close()
				}
			} else {
				res, re := cli.Get(ctx, base+c.path)
				e = re
				if res != nil {
					body = res.ReadAllString()
					_ = res.Close()
				}
			}
			t.AssertNil(e)
			// 任何端点回"未登录/凭证失效"都意味着中间件没把会员身份注入 ctx（C1 的形态）
			t.Assert(strings.Contains(body, "\"code\":10003"), false)
			// 且必须真的走到业务层（而非在白名单里空转后报别的框架错误）
			t.Assert(strings.Contains(body, "\"code\":10002"), false)
		}
	})
}
