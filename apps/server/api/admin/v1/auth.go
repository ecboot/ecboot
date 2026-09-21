package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 后台登录（含登录审计）
	AdminLoginReq struct {
		g.Meta      `path:"/login" tags:"Admin" method:"POST" summary:"后台登录"`
		Username    string `json:"username" v:"required" dc:"登录名"`
		Password    string `json:"password" v:"required" dc:"密码"`
		CaptchaKey  string `json:"captchaKey" dc:"图形验证码key"`
		CaptchaCode string `json:"captchaCode" dc:"图形验证码"`
	}
	AdminLoginRes struct {
		Token        string `json:"token" dc:"访问凭证"`
		RefreshToken string `json:"refreshToken" dc:"刷新凭证"`
		RealName     string `json:"realName" dc:"姓名"`
		IsSuper      bool   `json:"isSuper" dc:"是否超管(跳过权限校验)"`
	}

	AdminTokenRefreshReq struct {
		g.Meta       `path:"/token/refresh" tags:"Admin" method:"POST" summary:"刷新后台凭证"`
		RefreshToken string `json:"refreshToken" v:"required" dc:"刷新凭证"`
	}
	AdminTokenRefreshRes struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refreshToken"`
	}

	AdminLogoutReq struct {
		g.Meta       `path:"/logout" tags:"Admin" method:"POST" summary:"后台登出"`
		RefreshToken string `json:"refreshToken" dc:"刷新凭证(双凭证同毁, 对齐user渠道)"`
	}
	AdminLogoutRes struct {
		Success bool `json:"success"`
	}

	AdminProfileReq struct {
		g.Meta `path:"/profile" tags:"Admin" method:"GET" summary:"个人信息"`
	}
	AdminProfileRes struct {
		Username string   `json:"username"`
		RealName string   `json:"realName"`
		Roles    []string `json:"roles" dc:"角色编码列表"`
	}

	AdminChangePasswordReq struct {
		g.Meta      `path:"/profile/password" tags:"Admin" method:"PUT" summary:"修改密码"`
		OldPassword string `json:"oldPassword" v:"required" dc:"原密码"`
		NewPassword string `json:"newPassword" v:"required" dc:"新密码"`
	}
	AdminChangePasswordRes struct {
		Success bool `json:"success"`
	}
)
