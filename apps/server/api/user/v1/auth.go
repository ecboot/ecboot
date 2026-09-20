package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	SmsLoginReq struct {
		g.Meta `path:"/login/sms" tags:"Admin" method:"POST" summary:"短信验证码登录接口"`

		PhoneNumber string `json:"phoneNumber" v:"required" dc:"手机号" d:"18888888888"`
		SmsCode     string `json:"smsCode" v:"required" dc:"短信验证码" d:"123456"`
	}

	SmsLoginRes struct {
		Token        string `json:"token" v:"required" dc:"登录凭证"`
		RefreshToken string `json:"refreshToken" v:"required" dc:"刷新凭证，用于刷新登录凭证"`
	}
)
