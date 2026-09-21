// Package idgen 分布式 ID 生成（AGENTS.md 基础库约定: 业务编号场景统一 sonyflake，
// 替代自增 ID 的订单号/SKU 编码/门店编码等）。
package idgen

import (
	"sync"

	"github.com/sony/sonyflake/v2"
)

// 惰性单例（sync.Once 防并发竞态, 与 service/user phoneCipher 同式）。
var (
	once    sync.Once
	flake   *sonyflake.Sonyflake
	initErr error
)

// NextID 生成下一个分布式 ID（全局唯一、趋势递增；机器位取本机私有 IP 低 16 位）。
func NextID() (int64, error) {
	once.Do(func() {
		flake, initErr = sonyflake.New(sonyflake.Settings{})
	})
	if initErr != nil {
		return 0, initErr
	}
	return flake.NextID()
}
