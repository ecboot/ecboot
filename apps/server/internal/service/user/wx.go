// wx.go 微信登录（mock 归并）——US5 / FR-011。
// 归并规则（spec 定档）: openid 命中直接登录；未命中且能取手机号：
// 手机号已有账号→绑定微信身份并登录（该账号已绑其他身份则 20004 拒绝）；
// 未命中→新建账号。休眠分级: 微信静默路径 ≥ 阈值 → 拒绝并引导短信通道。
package user

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/library/sms"
)

// WxIdentity 微信侧身份（真实渠道由 WxClient 实现 code2session 换取）。
type WxIdentity struct {
	Openid  string
	Unionid string
	Phone   string // getPhoneNumber 一键取号结果（mock 模式由命令显式传入）
}

// WxClient 微信能力抽象（真实接入另立特性, 只补实现）。
type WxClient interface {
	Code2Session(ctx context.Context, wxCode string) (*WxIdentity, error)
}

// mockWxClient 开发态实现: openid = mock-{wxCode}, 手机号由调用方传入。
type mockWxClient struct{}

func (m *mockWxClient) Code2Session(ctx context.Context, wxCode string) (*WxIdentity, error) {
	return &WxIdentity{Openid: "mock-" + wxCode}, nil
}

// ErrWxBindConflict 微信身份与既有账号绑定冲突（人工渠道）。
var ErrWxBindConflict = errcode.New(errcode.CodeWxBindConflict, "该手机号已绑定其他微信身份,请联系人工客服处理")

// WxLogin 微信登录（含归并）。phone/smsCode 为 mock 模式开发字段或归并核验字段。
func WxLogin(ctx context.Context, wxCode, phone, smsCode string, channel int) (*LoginOutcome, error) {
	var ident *WxIdentity
	var err error
	if mockEnabled(ctx) {
		ident, err = (&mockWxClient{}).Code2Session(ctx, wxCode)
	} else {
		return nil, errcode.New(errcode.CodeSystemError, "微信真实渠道未接入")
	}
	if err != nil {
		return nil, err
	}

	// ① openid 命中 → 直接登录（含休眠分级校验）
	record, err := dao.User.Ctx(ctx).Where(dao.User.Columns().WxOpenid, ident.Openid).One()
	if err != nil {
		return nil, err
	}
	if !record.IsEmpty() {
		if err = ensureNotDormant(record); err != nil {
			return nil, err
		}
		userId := record["id"].Int64()
		token, refresh, sErr := sessionManager(ctx).Create(ctx, userId)
		if sErr != nil {
			return nil, sErr
		}
		touchLogin(ctx, userId, channel, true, "")
		return &LoginOutcome{Token: token, RefreshToken: refresh, UserId: userId}, nil
	}

	// ② 未绑定 → 手机号优先归并（必须能取得手机号 + 短信码核验, mock 亦验——评审 C5 fail-closed）
	if phone == "" {
		return nil, errcode.New(errcode.CodeInvalidParam, "该微信身份未绑定,需提供手机号完成归并")
	}
	if !sms.ValidPhone(phone) {
		return nil, errcode.New(errcode.CodeInvalidParam, "手机号格式不正确")
	}
	if !sms.ConsumeSmsCode(ctx, phone, smsCode) {
		sms.CountFail(ctx, phone)
		return nil, errcode.New(errcode.CodeCaptchaError, "验证码错误或已过期")
	}
	pc := phoneCipher()
	phoneHash := pc.Hash(phone)
	target, err := dao.User.Ctx(ctx).Where(dao.User.Columns().PhoneHash, phoneHash).One()
	if err != nil {
		return nil, err
	}

	if !target.IsEmpty() {
		// 手机号已有账号 → 绑定微信身份（已绑其他身份则拒绝）
		if existing := target["wx_openid"].String(); existing != "" && existing != ident.Openid {
			return nil, ErrWxBindConflict
		}
		if err = ensureNotDormant(target); err != nil {
			return nil, err
		}
		userId := target["id"].Int64()
		if _, err = dao.User.Ctx(ctx).Where(dao.User.Columns().Id, userId).
			Data(g.Map{dao.User.Columns().WxOpenid: ident.Openid, dao.User.Columns().WxUnionid: ident.Unionid}).Update(); err != nil {
			return nil, err
		}
		token, refresh, sErr := sessionManager(ctx).Create(ctx, userId)
		if sErr != nil {
			return nil, sErr
		}
		touchLogin(ctx, userId, channel, true, "")
		return &LoginOutcome{Token: token, RefreshToken: refresh, UserId: userId}, nil
	}

	// ③ 全新用户 → 新建（含微信身份）
	if channel <= 0 {
		channel = 1
	}
	cipherText, encErr := pc.Encrypt(phone)
	if encErr != nil {
		return nil, encErr
	}
	now := time.Now()
	res, insErr := dao.User.Ctx(ctx).Data(g.Map{
		dao.User.Columns().Phone:           cipherText,
		dao.User.Columns().PhoneHash:       phoneHash,
		dao.User.Columns().WxOpenid:        ident.Openid,
		dao.User.Columns().WxUnionid:       ident.Unionid,
		dao.User.Columns().Nickname:        "微信用户" + wxCode,
		dao.User.Columns().RegisterChannel: channel,
		dao.User.Columns().ShareCode:       randomShareCode(),
		dao.User.Columns().LastActiveAt:    now,
		dao.User.Columns().LastLoginAt:     now,
	}).InsertAndGetId()
	if insErr != nil {
		return nil, insErr
	}
	token, refresh, sErr := sessionManager(ctx).Create(ctx, res)
	if sErr != nil {
		return nil, sErr
	}
	touchLogin(ctx, res, channel, true, "")
	return &LoginOutcome{Token: token, RefreshToken: refresh, UserId: res, IsNew: true}, nil
}

// ensureNotDormant 休眠分级: 微信静默路径 ≥ 一级阈值 → 拒绝并引导短信通道（FR-013）。
func ensureNotDormant(record gdb.Record) error {
	dormantDays := cfgInt(context.Background(), "dormant.tier1.days", defDormantDays)
	if last := record["last_active_at"].Time(); !last.IsZero() &&
		time.Since(last) > time.Duration(dormantDays)*24*time.Hour {
		return errcode.New(errcode.CodeInvalidParam, "账号长期未登录,请使用手机号验证码方式重新核验登录")
	}
	return nil
}
