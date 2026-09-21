package system

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
)

// findConfig 取某配置当前行。
func findConfig(ctx context.Context, t *gtest.T, code string) gdb.Record {
	rec, err := g.DB().GetOne(ctx, "SELECT * FROM system_config WHERE code=?", code)
	t.AssertNil(err)
	return rec
}

// TestConfigList 配置列表（FR-019）: 返回种子全项, 字段完整。
func TestConfigList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		items, err := ConfigList(ctx)
		t.AssertNil(err)
		t.AssertGE(len(items), 6) // 000030 种子 6 项
		found := false
		for _, it := range items {
			if it.Code == "order.timeout.minutes" {
				found = true
				t.Assert(it.ValueType, 1)
				t.Assert(it.Description != "", true)
			}
		}
		t.Assert(found, true)
	})
}

// TestConfigUpdate 改值与启停（FR-020/021）: 类型校验 + 停用语义 + 测试后还原。
func TestConfigUpdate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		orig := findConfig(ctx, t, "order.timeout.minutes")
		defer func() {
			_, _ = g.DB().Exec(ctx,
				"UPDATE system_config SET value=?, status=? WHERE code='order.timeout.minutes'",
				orig["value"], orig["status"])
		}()

		// 整数型传非数字 → 80007
		err := ConfigUpdate(ctx, "order.timeout.minutes", "abc", 1)
		t.Assert(errCode(err), errcode.CodeConfigValueInvalid)

		// 合法修改 → List 反映
		t.AssertNil(ConfigUpdate(ctx, "order.timeout.minutes", "45", 1))
		items, err := ConfigList(ctx)
		t.AssertNil(err)
		for _, it := range items {
			if it.Code == "order.timeout.minutes" {
				t.Assert(it.Value, "45")
				t.Assert(it.Status, 1)
			}
		}

		// 停用（覆盖层语义: status=0, 读取方回退代码默认）
		t.AssertNil(ConfigUpdate(ctx, "order.timeout.minutes", "45", 0))
		rec := findConfig(ctx, t, "order.timeout.minutes")
		t.Assert(rec["status"].Int(), 0)

		// 布尔型校验: 传非法串 → 80007
		err = ConfigUpdate(ctx, "account_payment.enabled", "not-bool", 1)
		t.Assert(errCode(err), errcode.CodeConfigValueInvalid)

		// 不存在的配置 → 10006
		err = ConfigUpdate(ctx, "no.such.config", "1", 1)
		t.Assert(errCode(err), errcode.CodeNotFound)
	})
}
