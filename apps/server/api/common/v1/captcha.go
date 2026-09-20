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
