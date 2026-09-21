package idgen

import (
	"sync"
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

// TestNextIDConcurrent 并发唯一性（评审 M7）: 并发生成无重复（-race 下亦应通过）。
func TestNextIDConcurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		const n = 64
		seen := make(chan int64, n)
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				id, err := NextID()
				t.AssertNil(err)
				seen <- id
			}()
		}
		wg.Wait()
		close(seen)
		uniq := map[int64]bool{}
		for id := range seen {
			t.Assert(uniq[id], false)
			uniq[id] = true
		}
		t.Assert(len(uniq), n)
	})
}
