// Package sms 短信发送组件：接口抽象 + 开发态 Mock。
// 真实渠道（阿里云/腾讯云）接入时新增实现即可，调用方零改动（spec FR-008）。
package sms

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/gogf/gf/v2/frame/g"
)

// Sender 短信发送抽象。
type Sender interface {
	// Send 发送短信验证码到指定手机号。
	Send(ctx context.Context, phone, code string) error
}

// CodeLength 短信验证码位数。
const CodeLength = 6

// GenerateCode 生成 6 位数字验证码。
func GenerateCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// MockSender 开发态实现：验证码写入 Redis `mock:sms:{phone}`（供联调取码端点）
// 并输出日志；生产环境由配置关闭 mock 后不会装配本实现。
type MockSender struct{}

func NewMockSender() *MockSender { return &MockSender{} }

func (s *MockSender) Send(ctx context.Context, phone, code string) error {
	g.Log().Infof(ctx, "[mock-sms] phone=%s code=%s", phone, code)
	_, err := g.Redis().Do(ctx, "SET", "mock:sms:"+phone, code, "EX", 300)
	return err
}
