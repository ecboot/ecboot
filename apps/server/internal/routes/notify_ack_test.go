// notify_ack_test.go 回调应答不被统一响应包装（012 评审 I10 回归）。
// 为什么必须走真实 HTTP: I10 的缺陷只在"控制器返回值 → Response 中间件包装"这条链上可见——
// 原实现 `return &v1.PayNotifyRes{Code:"FAIL"}, nil` 会被中间件包成三段式（恒 code:0）,
// 渠道看到的是"成功", "失败→渠道重试"的语义直接失效。service 层测试覆盖不到这一跳。
// 这是本仓第一个 HTTP 端点级测试（真实路由 + 真实中间件链）。
package routes

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/database/gdb"
	gredis "github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/testutil"
)

func init() {
	time.Local = time.UTC // 与 main.go 及各测试基座同一口径
	_ = os.Setenv("ECBOOT_MOCK", "true")
	_ = gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{{Link: testutil.DSN()}},
	})
	gredis.SetConfig(&gredis.Config{Address: testutil.RedisAddr, Db: 0})
}

// startTestServer 起一个只监听回环地址的实例（含完整中间件链与路由）。
func startTestServer(t *gtest.T, port int) *ghttp.Server {
	s := g.Server(fmt.Sprintf("ecboot-test-%d", port))
	s.SetAddr("127.0.0.1")
	s.SetPort(port)
	s.SetDumpRouterMap(false)
	s.SetAccessLogEnabled(false)
	if err := s.SetLogPath(os.TempDir()); err != nil { // 不往仓库里写日志
		t.Fatal(err)
	}
	RouterGroup(s)
	go func() { _ = s.Start() }()
	for i := 0; i < 100 && s.GetListenedPort() <= 0; i++ {
		time.Sleep(50 * time.Millisecond)
	}
	if s.GetListenedPort() <= 0 {
		t.Fatal("测试服务器未监听")
	}
	return s
}

// TestPayNotifyAckNotWrapped 回调应答必须是渠道自己的裸 JSON（I10）。
// 用非法报文触发 FAIL 分支即可（无需任何业务 fixture）: 应答体必须是 `{"code":"FAIL","message":"处理失败"}`,
// 而不是统一三段式（`{"code":0,...,"data":{"code":"FAIL"...}}`）; 同时顺带验证该路径在鉴权白名单内可达。
func TestPayNotifyAckNotWrapped(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := startTestServer(t, 38871)
		defer func() { _ = s.Shutdown() }()

		res, err := g.Client().Post(context.Background(),
			fmt.Sprintf("http://127.0.0.1:%d/shop/pay/notify", s.GetListenedPort()), "{not-json-payload")
		t.AssertNil(err)
		defer func() { _ = res.Close() }()
		body := res.ReadAllString()
		t.Assert(body, `{"code":"FAIL","message":"处理失败"}`)

		// 退款回调同款
		res2, err := g.Client().Post(context.Background(),
			fmt.Sprintf("http://127.0.0.1:%d/shop/refund/notify", s.GetListenedPort()), "{not-json-payload")
		t.AssertNil(err)
		defer func() { _ = res2.Close() }()
		t.Assert(res2.ReadAllString(), `{"code":"FAIL","message":"处理失败"}`)
	})
}
