// Package user 用户领域服务——认证链路（注册即登录/微信归并/休眠核身）。
// 依赖: api 契约的入参结构、dao（gf 生成）、library 技术组件（captcha/sms/security/errcode）。
// 契约权威: specs/003-auth-bootstrap/spec.md + contracts/api-contracts.md。
package user

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/library/captcha"
	"ecboot/internal/library/security"
	"ecboot/internal/library/sms"
)

// ---------- Redis 键与配置常量（data-model §一） ----------

const (
	smsCodeKeyPrefix  = "captcha:sms:"
	smsIntervalPrefix = "captcha:sms:interval:"
	failCountPrefix   = "captcha:fail:"
	mockSmsPrefix     = "mock:sms:"
)

// 配置兜底默认（system_config 是覆盖层, 缺失/停用回退此处——宪法 V）。
const (
	defSmsTTLSeconds  = 300
	defResendSeconds  = 60
	defFailMax        = 5
	defLockSeconds    = 1800
	defSessionTTLDays = 7
	defDormantDays    = 90
)

// cfgInt 读 system_config 整型配置（停用/缺失回退默认值）。
func cfgInt(ctx context.Context, code string, fallback int) int {
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

// phoneCipher 惰性单例：密钥经环境变量 PHONE_KEY 注入（Base64, 32 字节），
// 未配置时使用开发态默认密钥（仅限本地开发, 生产 MUST 配置环境变量）。
var phoneCipherInstance *security.PhoneCipher

func phoneCipher() *security.PhoneCipher {
	if phoneCipherInstance == nil {
		key := g.Cfg().MustGet(gctxNew(), "security.phoneKey", "MDEyMzQ1Njc4OWFiY2RlZjAxMjM0NTY3ODlhYmNkZWY=").String()
		pc, err := security.NewPhoneCipher(key)
		if err != nil {
			panic(fmt.Sprintf("phone key 初始化失败: %v", err))
		}
		phoneCipherInstance = pc
	}
	return phoneCipherInstance
}

func gctxNew() context.Context { return context.Background() }

// sessionManager 会话管理器单例（TTL 走配置）。
func sessionManager(ctx context.Context) *security.SessionManager {
	return security.NewSessionManager(cfgInt(ctx, "session.ttl_days", defSessionTTLDays))
}

// mockEnabled 是否开发态 mock（决定发码走 Mock、微信走 MockWxClient）。
func mockEnabled(ctx context.Context) bool {
	return g.Cfg().MustGet(ctx, "wechat.mockEnabled", true).Bool()
}

// ---------- 短信验证码发送（US3 / FR-007~009, FR-014） ----------

// SendSmsCode 发送短信验证码：格式校验→图形码校验→频控→生成落 Redis→Mock 发送。
func SendSmsCode(ctx context.Context, phone, captchaKey, captchaCode string) (expiresIn int, err error) {
	if !validPhone(phone) {
		return 0, errcode.New(errcode.CodeInvalidParam, "手机号格式不正确")
	}
	ttl := cfgInt(ctx, "captcha.sms.ttl_seconds", defSmsTTLSeconds)

	// 锁定检查
	locked, _ := g.Redis().Do(ctx, "EXISTS", failCountPrefix+phone)
	if locked.Bool() {
		return 0, errcode.New(errcode.CodeLocked, "尝试次数超限,请稍后再试")
	}

	// 图形码一次性校验
	ok, verr := captcha.Verify(ctx, captchaKey, captchaCode)
	if verr != nil {
		return 0, verr
	}
	if !ok {
		return 0, errcode.New(errcode.CodeCaptchaError, "图形验证码错误或已过期")
	}

	// 重发间隔锁
	resend := cfgInt(ctx, "captcha.sms.resend_seconds", defResendSeconds)
	v, err := g.Redis().Do(ctx, "SET", smsIntervalPrefix+phone, 1, "NX", "EX", resend)
	if err != nil {
		return 0, err
	}
	if v.String() != "OK" {
		return 0, errcode.New(errcode.CodeTooFrequent, "发送过于频繁,请稍后再试")
	}

	// 生成短信码并落 Redis
	code, err := sms.GenerateCode()
	if err != nil {
		return 0, err
	}
	if _, err = g.Redis().Do(ctx, "SET", smsCodeKeyPrefix+phone, code, "EX", ttl); err != nil {
		return 0, err
	}

	// 发送（开发态 mock: 写 mock:sms:{phone} 供取码端点）
	sender := sms.NewMockSender()
	if err = sender.Send(ctx, phone, code); err != nil {
		return 0, err
	}
	return ttl, nil
}

// ---------- 短信验证码登录（注册即登录, US3 / FR-010/012/013） ----------

type LoginOutcome struct {
	Token        string
	RefreshToken string
	UserId       int64
	IsNew        bool
}

// SmsLogin 短信验证码登录：校验(一次性)→查/建账号→休眠核身→时间戳+日志→发会话。
// 码校验失败计 fail 计数（达上限锁定, FR-009）。
func SmsLogin(ctx context.Context, phone, smsCode string, channel int) (*LoginOutcome, error) {
	if !validPhone(phone) {
		return nil, errcode.New(errcode.CodeInvalidParam, "手机号格式不正确")
	}

	// 锁定检查
	locked, _ := g.Redis().Do(ctx, "EXISTS", failCountPrefix+phone)
	if locked.Bool() {
		return nil, errcode.New(errcode.CodeLocked, "尝试次数超限,请稍后再试")
	}

	// 短信码一次性校验（GETDEL）; 失败计数
	ok, verr := g.Redis().Do(ctx, "GETDEL", smsCodeKeyPrefix+phone)
	if verr != nil {
		return nil, verr
	}
	if ok.String() == "" || ok.String() != smsCode {
		countFail(ctx, phone)
		return nil, errcode.New(errcode.CodeCaptchaError, "验证码错误或已过期")
	}

	// 查/建账号
	pc := phoneCipher()
	phoneHash := pc.Hash(phone)
	record, err := dao.User.Ctx(ctx).Where(dao.User.Columns().PhoneHash, phoneHash).One()
	if err != nil {
		return nil, err
	}

	isNew := false
	var userId int64
	if record.IsEmpty() {
		cipherText, encErr := pc.Encrypt(phone)
		if encErr != nil {
			return nil, encErr
		}
		if channel <= 0 {
			channel = 1
		}
		data := g.Map{
			dao.User.Columns().Phone:           cipherText,
			dao.User.Columns().PhoneHash:       phoneHash,
			dao.User.Columns().Nickname:        "用户" + phone[len(phone)-4:],
			dao.User.Columns().RegisterChannel: channel,
			dao.User.Columns().ShareCode:       randomShareCode(),
			dao.User.Columns().LastActiveAt:    time.Now(),
			dao.User.Columns().LastLoginAt:     time.Now(),
		}
		res, insErr := dao.User.Ctx(ctx).Data(data).InsertAndGetId()
		if insErr != nil {
			// 并发注册唯一索引兜底: 冲突则按已注册继续
			rec2, qErr := dao.User.Ctx(ctx).Where(dao.User.Columns().PhoneHash, phoneHash).One()
			if qErr != nil || rec2.IsEmpty() {
				return nil, insErr
			}
			record = rec2
		} else {
			userId = res
			isNew = true
			record, _ = dao.User.Ctx(ctx).Where(dao.User.Columns().Id, userId).One()
		}
	}

	if record.IsEmpty() {
		return nil, errcode.New(errcode.CodeOrderNotFound, "订单不存在") // 不可达保护
	}
	userId = record["id"].Int64()

	// 禁用校验
	if record["status"].Int() == 2 {
		return nil, errcode.New(errcode.CodeUserDisabled, "账号已被禁用,请联系客服")
	}

	// 休眠分级核身: 短信码登录=完整核验, 天然满足一级要求（微信静默路径见 WxLogin）
	dormantDays := cfgInt(ctx, "dormant.tier1.days", defDormantDays)
	if last := record["last_active_at"].Time(); !last.IsZero() &&
		time.Since(last) > time.Duration(dormantDays)*24*time.Hour {
		// 短信通道即强核身——放行; 微信静默通道在 WxLogin 中拒绝
		_ = dormantDays
	}

	// 会话发放
	sm := sessionManager(ctx)
	token, refreshToken, sErr := sm.Create(ctx, userId)
	if sErr != nil {
		return nil, sErr
	}

	// 时间戳与登录日志（成功）
	touchLogin(ctx, userId, channel, true, "")

	return &LoginOutcome{Token: token, RefreshToken: refreshToken, UserId: userId, IsNew: isNew}, nil
}

// countFail 登录/校验失败计数（达上限写锁定键, TTL=锁定时长）。
func countFail(ctx context.Context, phone string) {
	maxFails := cfgInt(ctx, "captcha.fail.max", defFailMax)
	lockSeconds := cfgInt(ctx, "captcha.fail.lock_seconds", defLockSeconds)
	key := failCountPrefix + phone
	v, err := g.Redis().Do(ctx, "INCR", key)
	if err != nil {
		return
	}
	if v.Int() == 1 {
		_, _ = g.Redis().Do(ctx, "EXPIRE", key, lockSeconds)
		return
	}
	if v.Int() >= maxFails {
		_, _ = g.Redis().Do(ctx, "EXPIRE", key, lockSeconds)
	}
}

// touchLogin 更新最后登录/活跃时间并写登录日志。
func touchLogin(ctx context.Context, userId int64, channel int, success bool, ip string) {
	now := time.Now()
	if success {
		_, _ = dao.User.Ctx(ctx).Where(dao.User.Columns().Id, userId).
			Data(g.Map{
				dao.User.Columns().LastLoginAt:  now,
				dao.User.Columns().LastActiveAt: now,
			}).Update()
	}
	status := 1
	if !success {
		status = 2
	}
	uid := any(userId)
	if userId == 0 {
		uid = nil
	}
	_, _ = dao.UserLoginLog.Ctx(ctx).Data(g.Map{
		dao.UserLoginLog.Columns().UserId:    uid,
		dao.UserLoginLog.Columns().LoginChannel: channel,
		dao.UserLoginLog.Columns().LoginStatus:  status,
		dao.UserLoginLog.Columns().Ip:           ip,
	}).Insert()
}

// validPhone 手机号格式（FR-014）。
func validPhone(phone string) bool {
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

// randomShareCode 个人推广码（8 位大写字母数字; 唯一冲突重试由调用方处理）。
func randomShareCode() string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 8)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

var _ = grand.Perm // 保留 gf 随机引用

// PhoneCipherFor 暴露手机号加密器（受控场景: 展示脱敏前的解密、后台改绑）。
func PhoneCipherFor(ctx context.Context) *security.PhoneCipher {
	return phoneCipher()
}
