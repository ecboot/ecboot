// ports_wiring_test.go 跨域端口装配断言（017 批次 11 评审 I6）:
// init() 装配必须在编译/测试期可证——此前 CommissionSettle/CommissionReverse 装配后
// 全库零验证（SC-4"有装配测试"声明不实）。本文件 import bootstrap 触发 init()。
package bootstrap

import (
	"testing"

	"ecboot/internal/service/shop"
)

// TestDistributionPortsWired 分销两端口的装配存在性（I6 守卫）。
func TestDistributionPortsWired(t *testing.T) {
	if shop.CommissionSettle == nil {
		t.Fatal("shop.CommissionSettle 未装配（确认收货→佣金计提事件会静默降级）")
	}
	if shop.CommissionReverse == nil {
		t.Fatal("shop.CommissionReverse 未装配（售后完成→佣金冲销事件会静默降级）")
	}
}
