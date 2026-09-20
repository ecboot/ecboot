package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 短信验证码登录（注册即登录; 003 设计, 路径迁移至 /user 分组）
	SmsLoginReq struct {
		g.Meta `path:"/login/sms" tags:"User" method:"POST" summary:"短信验证码登录接口"`

		PhoneNumber string `json:"phoneNumber" v:"required" dc:"手机号" d:"18888888888"`
		SmsCode     string `json:"smsCode" v:"required" dc:"短信验证码" d:"123456"`
		Channel     int    `json:"channel" dc:"注册渠道:1小程序 2H5" d:"1"`
	}
	SmsLoginRes struct {
		Token        string `json:"token" v:"required" dc:"访问凭证"`
		RefreshToken string `json:"refreshToken" v:"required" dc:"刷新凭证"`
		UserId       string `json:"userId" dc:"用户ID"`
		IsNew        bool   `json:"isNew" dc:"是否本次新注册"`
	}

	// 微信登录（手机号优先归并; 开发态 mock）
	WxLoginReq struct {
		g.Meta  `path:"/login/wx" tags:"User" method:"POST" summary:"微信登录接口"`
		WxCode  string `json:"wxCode" v:"required" dc:"微信授权码"`
		Phone   string `json:"phone" dc:"手机号(mock模式显式传入用于归并)"`
		SmsCode string `json:"smsCode" dc:"短信验证码(归并核验)"`
		Channel int    `json:"channel" dc:"注册渠道" d:"1"`
	}
	WxLoginRes struct {
		Token        string `json:"token" v:"required" dc:"访问凭证"`
		RefreshToken string `json:"refreshToken" v:"required" dc:"刷新凭证"`
		UserId       string `json:"userId" dc:"用户ID"`
		IsNew        bool   `json:"isNew" dc:"是否本次新注册"`
	}

	// 刷新访问凭证（双凭证会话）
	TokenRefreshReq struct {
		g.Meta       `path:"/token/refresh" tags:"User" method:"POST" summary:"刷新访问凭证"`
		RefreshToken string `json:"refreshToken" v:"required" dc:"刷新凭证"`
	}
	TokenRefreshRes struct {
		Token        string `json:"token" v:"required" dc:"新访问凭证"`
		RefreshToken string `json:"refreshToken" v:"required" dc:"新刷新凭证"`
	}

	// 登出（双凭证同失效）
	LogoutReq struct {
		g.Meta `path:"/logout" tags:"User" method:"POST" summary:"登出"`
	}
	LogoutRes struct {
		Success bool `json:"success" dc:"固定true"`
	}

	// 当前用户信息
	MeReq struct {
		g.Meta `path:"/me" tags:"User" method:"GET" summary:"当前用户信息"`
	}
	MeRes struct {
		UserId          string `json:"userId" dc:"用户ID"`
		Nickname        string `json:"nickname" dc:"昵称"`
		Avatar          string `json:"avatar" dc:"头像"`
		Phone           string `json:"phone" dc:"手机号(脱敏)"`
		Level           int    `json:"level" dc:"会员等级"`
		RegisterChannel int    `json:"registerChannel" dc:"注册渠道"`
	}
)
