package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	GetSmsCodeReq struct {
		g.Meta `path:"/sms/code" tags:"Common" method:"POST" summary:"获取短信验证码接口"`

		PhoneNumber string `json:"phoneNumber" v:"required" dc:"手机号" d:"18888888888"`
		Template    string `json:"template" v:"required" dc:"短信模板" d:"login"`
		CaptchaCode string `json:"captchaCode" v:"required" dc:"验证码"`
		CaptchaKey  string `json:"captchaKey" v:"required" dc:"验证码key"`
	}

	GetSmsCodeRes struct {
		SmsCodeKey string `json:"smsCodeKey" v:"required" dc:"短信验证码key"`
	}
)
