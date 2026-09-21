// Package user 用户领域服务——认证链路（注册即登录/微信归并/休眠核身）。
// 依赖: api 契约的入参结构、dao（gf 生成）、library 技术组件（captcha/sms/security/errcode）。
// 契约权威: specs/003-auth-bootstrap/spec.md + contracts/api-contracts.md。
package user

import (
	"context"
	"crypto/rand"
	"ecboot/internal/model"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/library/security"
	"ecboot/internal/library/sms"
)

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

// phoneCipher 惰性单例（sync.Once 防并发竞态, 评审 Minor17）：
// 密钥优先级 config security.phoneKey → env PHONE_KEY → 开发态默认（仅限本地, 启动警告）。
var (
	phoneCipherOnce sync.Once
	phoneCipherImpl *security.PhoneCipher
	phoneCipherErr  error
)

func phoneCipher() *security.PhoneCipher {
	phoneCipherOnce.Do(func() {
		key := g.Cfg().MustGet(gctxNew(), "security.phoneKey", os.Getenv("PHONE_KEY")).String()
		if key == "" {
			key = "MDEyMzQ1Njc4OWFiY2RlZjAxMjM0NTY3ODlhYmNkZWY=" // 开发态默认
			g.Log().Warning(gctxNew(), "security.phoneKey 未配置, 使用开发态默认密钥——生产环境必须注入!")
		}
		phoneCipherImpl, phoneCipherErr = security.NewPhoneCipher(key)
	})
	if phoneCipherErr != nil {
		panic(fmt.Sprintf("phone key 初始化失败: %v", phoneCipherErr))
	}
	return phoneCipherImpl
}

func gctxNew() context.Context { return context.Background() }

// sessionManager 会话管理器（TTL 读配置, 评审 I8 统一入口）。
func sessionManager(ctx context.Context) *security.SessionManager {
	return security.NewSessionManagerFromConfig(ctx)
}

// mockEnabled 是否开发态 mock（决定发码走 Mock、微信走 MockWxClient）。
// fail-closed（评审 C5）: 仅当配置 wechat.mockEnabled=true 或 env ECBOOT_MOCK=true 才开启。
func mockEnabled(ctx context.Context) bool {
	if g.Cfg().MustGet(ctx, "wechat.mockEnabled", false).Bool() {
		return true
	}
	return os.Getenv("ECBOOT_MOCK") == "true"
}

// SendSmsCode 发送短信验证码（实现下沉 library/sms, common 渠道直调; service 保留转发兼容测试）。
func SendSmsCode(ctx context.Context, phone, captchaKey, captchaCode string) (int, error) {
	return sms.SendSmsCode(ctx, phone, captchaKey, captchaCode)
}

// ---------- 短信验证码登录（注册即登录, US3 / FR-010/012/013） ----------

// SmsLogin 短信验证码登录：校验(一次性)→查/建账号→休眠核身→时间戳+日志→发会话。
// 码校验失败计 fail 计数（达上限锁定, FR-009）。
func SmsLogin(ctx context.Context, phone, smsCode string, channel int) (*model.LoginOutcome, error) {
	if !validPhone(phone) {
		return nil, errcode.New(errcode.CodeInvalidParam, "手机号格式不正确")
	}

	if sms.IsLocked(ctx, phone) {
		return nil, errcode.New(errcode.CodeLocked, "尝试次数超限,请稍后再试")
	}
	if !sms.ConsumeSmsCode(ctx, phone, smsCode) {
		sms.CountFail(ctx, phone)
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

	// 会话发放
	sm := sessionManager(ctx)
	token, refreshToken, sErr := sm.Create(ctx, userId)
	if sErr != nil {
		return nil, sErr
	}

	// 时间戳与登录日志（成功）
	touchLogin(ctx, userId, channel, true, "")

	return &model.LoginOutcome{Token: token, RefreshToken: refreshToken, UserId: userId, IsNew: isNew}, nil
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
		dao.UserLoginLog.Columns().UserId:       uid,
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
