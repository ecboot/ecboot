// service.go 短信验证码发送编排（技术编排, 无业务归属——common 渠道直接调用）。
// 流程: 格式校验 → 锁定检查 → 图形码一次性校验 → 频控 → 生成落 Redis → Mock 发送。
package sms

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/errcode"
	"ecboot/internal/library/captcha"
)

const (
	smsCodeKeyPrefix  = "captcha:sms:"
	smsIntervalPrefix = "captcha:sms:interval:"
	failCountPrefix   = "captcha:fail:"
	mockSmsPrefix     = "mock:sms:"
)

// 配置兜底默认（system_config 是覆盖层）。
const (
	DefSmsTTLSeconds = 300
	DefResendSeconds = 60
	DefFailMax       = 5
	DefLockSeconds   = 1800
)

// CfgInt 读 system_config 整型配置（停用/缺失回退默认值）——导出供会话/休眠复用。
func CfgInt(ctx context.Context, code string, fallback int) int {
	v, err := g.DB().GetOne(ctx,
		"SELECT value FROM system_config WHERE code=? AND status=1 AND deleted=0", code)
	if err != nil || v.IsEmpty() {
		return fallback
	}
	if n := v["value"].Int(); n > 0 {
		return n
	}
	return fallback
}

// ValidPhone 手机号格式（FR-014: 非法不入发码流程）。
func ValidPhone(phone string) bool {
	if len(phone) != 11 || phone[0] != '1' {
		return false
	}
	for i := 1; i < len(phone); i++ {
		if phone[i] < '0' || phone[i] > '9' {
			return false
		}
	}
	return phone[1] >= '3'
}

// SendSmsCode 发送短信验证码（返回有效期秒）。
// 防刷次序: 锁定检查 → 图形码(一次性) → 频控锁 → 生成/落 Redis → 发送。
func SendSmsCode(ctx context.Context, phone, captchaKey, captchaCode string) (expiresIn int, err error) {
	if !ValidPhone(phone) {
		return 0, errcode.New(errcode.CodeInvalidParam, "手机号格式不正确")
	}
	ttl := CfgInt(ctx, "captcha.sms.ttl_seconds", DefSmsTTLSeconds)

	if locked, _ := g.Redis().Do(ctx, "EXISTS", failCountPrefix+phone); locked.Bool() {
		return 0, errcode.New(errcode.CodeLocked, "尝试次数超限,请稍后再试")
	}

	ok, verr := captcha.Verify(ctx, captchaKey, captchaCode)
	if verr != nil {
		return 0, verr
	}
	if !ok {
		return 0, errcode.New(errcode.CodeCaptchaError, "图形验证码错误或已过期")
	}

	resend := CfgInt(ctx, "captcha.sms.resend_seconds", DefResendSeconds)
	v, err := g.Redis().Do(ctx, "SET", smsIntervalPrefix+phone, 1, "NX", "EX", resend)
	if err != nil {
		return 0, err
	}
	if v.String() != "OK" {
		return 0, errcode.New(errcode.CodeTooFrequent, "发送过于频繁,请稍后再试")
	}

	code, err := GenerateCode()
	if err != nil {
		return 0, err
	}
	if _, err = g.Redis().Do(ctx, "SET", smsCodeKeyPrefix+phone, code, "EX", ttl); err != nil {
		return 0, err
	}

	if err = NewMockSender().Send(ctx, phone, code); err != nil {
		return 0, err
	}
	return ttl, nil
}

// CountFail 登录/校验失败计数（Lua 原子 INCR+EXPIRE, 达上限即锁定窗口）。
func CountFail(ctx context.Context, phone string) {
	maxFails := CfgInt(ctx, "captcha.fail.max", DefFailMax)
	lockSeconds := CfgInt(ctx, "captcha.fail.lock_seconds", DefLockSeconds)
	script := `
local n = redis.call('INCR', KEYS[1])
if n == 1 or n >= tonumber(ARGV[2]) then
  redis.call('EXPIRE', KEYS[1], ARGV[1])
end
return n`
	_, _ = g.Redis().Do(ctx, "EVAL", script, 1, failCountPrefix+phone, lockSeconds, maxFails)
}

// IsLocked 手机号是否处于失败锁定窗口（计数 ≥ 上限, 非"键存在"——
// 计数键自首次失败即存在并带 TTL, 只有达到上限才算锁定, 评审根因修复）。
func IsLocked(ctx context.Context, phone string) bool {
	v, err := g.Redis().Do(ctx, "GET", failCountPrefix+phone)
	if err != nil || v == nil {
		return false
	}
	return v.Int() >= CfgInt(ctx, "captcha.fail.max", DefFailMax)
}

// ConsumeSmsCode 一次性消费短信码（GETDEL; 匹配返回 true, 失败由调用方 CountFail）。
func ConsumeSmsCode(ctx context.Context, phone, code string) bool {
	if phone == "" || code == "" {
		return false
	}
	v, err := g.Redis().Do(ctx, "GETDEL", smsCodeKeyPrefix+phone)
	if err != nil || v.String() == "" {
		return false
	}
	return v.String() == code
}
