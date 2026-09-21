package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	GetSmsCodeReq struct {
		g.Meta `path:"/sms/code" tags:"Common" method:"POST" summary:"获取短信验证码接口"`

		PhoneNumber string `json:"phoneNumber" v:"required" dc:"手机号"`
		Template    string `json:"template" dc:"短信模板" d:"login"`
		CaptchaCode string `json:"captchaCode" v:"required" dc:"图形验证码"`
		CaptchaKey  string `json:"captchaKey" v:"required" dc:"图形验证码key"`
	}

	GetSmsCodeRes struct {
		ExpiresIn int `json:"expiresIn" dc:"验证码有效期(秒)"`
	}
)
