package idgen

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// TestNextID 唯一性（基础库约定: 业务编号场景）: 连续生成的 ID 不相等且为正。
func TestNextID(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a, err := NextID()
		t.AssertNil(err)
		b, err := NextID()
		t.AssertNil(err)
		t.AssertNE(a, b)
		t.AssertGT(a, 0)
		t.AssertGT(b, 0)
	})
}
