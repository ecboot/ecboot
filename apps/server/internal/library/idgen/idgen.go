// Package idgen 分布式 ID 生成（AGENTS.md 基础库约定: 业务编号场景统一 sonyflake，
// 替代自增 ID 的订单号/SKU 编码/门店编码等）。
package idgen

import (
	"hash/fnv"
	"net"
	"os"
	"sync"

	"github.com/sony/sonyflake/v2"
)

// 惰性单例（sync.Once 防并发竞态, 与 service/user phoneCipher 同式）。
var (
	once    sync.Once
	flake   *sonyflake.Sonyflake
	initErr error
)

// machineID 机器标识（0~65535）: 优先私网 IPv4 低 16 位（sonyflake 默认口径）;
// 无可用私网地址（容器/CI/单机多网卡）时回退 hostname FNV 哈希——避免"一次失败被
// sync.Once 永久缓存、运行期 ID 生成全挂"（评审 I6）。回退值稳定（同机同值）。
func machineID() (int, error) {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}
			if ip4 := ipnet.IP.To4(); ip4 != nil && ip4.IsPrivate() {
				return int(ip4[2])<<8 | int(ip4[3]), nil
			}
		}
	}
	name, err := os.Hostname()
	if err != nil {
		return 0, err
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(name))
	return int(h.Sum32() & 0xFFFF), nil
}

// NextID 生成下一个分布式 ID（全局唯一、趋势递增）。
func NextID() (int64, error) {
	once.Do(func() {
		flake, initErr = sonyflake.New(sonyflake.Settings{MachineID: machineID})
	})
	if initErr != nil {
		return 0, initErr
	}
	return flake.NextID()
}
