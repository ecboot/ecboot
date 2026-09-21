// Package paychannel 支付渠道适配层：渠道差异隔离在实现内，
// 上层（service/shop.pay_impl）只面向 Channel 接口。真实微信支付另立特性实现。
package paychannel

import (
	"encoding/json"
	"errors"
)

// Channel 支付渠道抽象。
type Channel interface {
	// Name 渠道标识。
	Name() string
	// CreateParams 生成渠道唤起参数（payNo/金额分/商品描述）。
	CreateParams(payNo string, amountFen int64, subject string) (map[string]any, error)
	// ParseNotify 解析并验签回调报文（验签失败返回错误）。
	ParseNotify(rawBody []byte) (*NotifyPayload, error)
	// Refund 发起退款（幂等: outRefundNo）——mock 直接成功。
	Refund(outRefundNo string, amountFen int64) error
}

// NotifyPayload 解析后的回调载荷。
type NotifyPayload struct {
	PayNo          string `json:"payNo"`
	ChannelTradeNo string `json:"channelTradeNo"`
	AmountFen      int64  `json:"amountFen"`
	Success        bool   `json:"success"`
}

// MockChannel 开发态渠道：不做真实验签（验签为占位通过）。
type MockChannel struct{}

func NewMock() *MockChannel { return &MockChannel{} }

func (m *MockChannel) Name() string { return "mock" }

func (m *MockChannel) CreateParams(payNo string, amountFen int64, subject string) (map[string]any, error) {
	return map[string]any{
		"payNo":     payNo,
		"amountFen": amountFen,
		"subject":   subject,
		"mock":      true,
	}, nil
}

// ParseNotify mock 回调报文: {"payNo":"...","channelTradeNo":"...","amountFen":0,"success":true}
func (m *MockChannel) ParseNotify(rawBody []byte) (*NotifyPayload, error) {
	if len(rawBody) == 0 {
		return nil, errors.New("回调报文为空")
	}
	var p NotifyPayload
	if err := json.Unmarshal(rawBody, &p); err != nil {
		return nil, err
	}
	if p.PayNo == "" {
		return nil, errors.New("回调缺少 payNo")
	}
	// mock 渠道验签占位: 真实渠道在此校验签名
	return &p, nil
}

func (m *MockChannel) Refund(outRefundNo string, amountFen int64) error {
	return nil // mock 退款直接成功
}
