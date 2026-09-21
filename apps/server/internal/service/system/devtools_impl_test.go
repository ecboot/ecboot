package system

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
)

// TestMockLatestSms 调试端点（FR-023, research D4）: mock 态返回最近模拟短信; 无记录报错。
func TestMockLatestSms(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900009999"
		_, _ = g.Redis().Do(ctx, "DEL", "mock:sms:"+phone)
		defer func() { _, _ = g.Redis().Do(ctx, "DEL", "mock:sms:"+phone) }()

		// 测试基座 init 已设 ECBOOT_MOCK=true（非生产态）→ 可读取
		_, _ = g.Redis().Do(ctx, "SET", "mock:sms:"+phone, "654321", "EX", 300)
		code, err := MockLatestSms(ctx, phone)
		t.AssertNil(err)
		t.Assert(code, "654321")

		// 无记录 → 未找到
		_, _ = g.Redis().Do(ctx, "DEL", "mock:sms:"+phone)
		_, err = MockLatestSms(ctx, phone)
		t.Assert(errCode(err), errcode.CodeNotFound)
	})
}
