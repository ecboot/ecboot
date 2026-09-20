package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	GetCaptchaReq struct {
		g.Meta `path:"/captcha" tags:"Common" method:"GET" summary:"获取验证码接口"`
	}

	GetCaptchaRes struct {
		CaptchaKey string `json:"captchaKey" v:"required" dc:"验证码key"`
		CaptchaImg string `json:"captchaImg" v:"required" dc:"验证码图片"`
	}
)

type (
	// 图形验证码独立校验（一次性; 随 003 实现补齐逻辑）
	VerifyCaptchaReq struct {
		g.Meta `path:"/captcha/verify" tags:"Common" method:"POST" summary:"图形验证码校验"`

		CaptchaKey  string `json:"captchaKey" v:"required" dc:"验证码key"`
		CaptchaCode string `json:"captchaCode" v:"required" dc:"验证码答案"`
	}
	VerifyCaptchaRes struct {
		Success bool `json:"success" dc:"校验通过"`
	}

	// 联调取码（仅 mock 模式注册路由; 生产不存在该端点）
	MockLatestSmsReq struct {
		g.Meta `path:"/captcha/sms/mock-latest" tags:"Common" method:"GET" summary:"最近短信验证码(mock)"`
		PhoneNumber string `json:"phoneNumber" v:"required" dc:"手机号"`
	}
	MockLatestSmsRes struct {
		SmsCode string `json:"smsCode" dc:"短信验证码"`
	}
)
