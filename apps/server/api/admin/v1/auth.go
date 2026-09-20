package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	LoginReq struct {
		g.Meta `path:"/login" tags:"Admin" method:"POST" summary:"登录接口"`

		Username     string  `json:"username" v:"required" dc:"用户名"`
		Password    string  `json:"password" v:"required" dc:"密码"`
		CaptchaCode string `json:"captchaCode" v:"required" dc:"验证码"`
		CaptchaKey  string `json:"captchaKey" v:"required" dc:"验证码key"`
	}

	LoginRes struct {
		Token        string `json:"token" v:"required" dc:"登录凭证"`
		RefreshToken string `json:"refreshToken" v:"required" dc:"刷新凭证，用于刷新登录凭证"`
	}
)
