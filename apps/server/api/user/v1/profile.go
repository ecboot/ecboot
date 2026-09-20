package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	ProfileDetailReq struct {
		g.Meta `path:"/profile" method:"GET" summary:"个人资料"`
	}
	ProfileDetailRes struct {
		UserId      string `json:"userId"`
		Nickname    string `json:"nickname" dc:"昵称"`
		Avatar      string `json:"avatar" dc:"头像"`
		Gender      int    `json:"gender" dc:"性别:0未知 1男 2女"`
		Phone       string `json:"phone" dc:"手机号(脱敏)"`
		HasPassword bool   `json:"hasPassword" dc:"是否已设密码"`
		GrowthValue int    `json:"growthValue" dc:"成长值"`
		Level       int    `json:"level" dc:"会员等级"`
		LevelName   string `json:"levelName" dc:"等级名称"`
	}

	ProfileUpdateReq struct {
		g.Meta   `path:"/profile" method:"PUT" summary:"修改个人资料"`
		Nickname string `json:"nickname" dc:"昵称"`
		Avatar   string `json:"avatar" dc:"头像"`
		Gender   int    `json:"gender" dc:"性别"`
	}
	ProfileUpdateRes struct {
		Success bool `json:"success"`
	}

	PointAccountReq struct {
		g.Meta `path:"/point-account" method:"GET" summary:"积分账户"`
	}
	PointAccountRes struct {
		Balance       int    `json:"balance" dc:"积分余额(可为负)"`
		LastEarnedAt  string `json:"lastEarnedAt" dc:"最后获得时间(滚动有效期口径)"`
	}

	PointLogListReq struct {
		g.Meta  `path:"/point-logs" method:"GET" summary:"积分流水"`
		BizType int `json:"bizType" dc:"类型筛选:1签到 2消费获得 3下单消耗 4退款回退 5分享获得 6评价获得 7注册赠送 8邀请奖励 9过期扣减"`
		PageReq
	}
	PointLogItem struct {
		BizType      int    `json:"bizType"`
		Points       int    `json:"points" dc:"变动(有符号)"`
		BalanceAfter int    `json:"balanceAfter" dc:"变动后余额"`
		OrderNo      string `json:"orderNo" dc:"关联订单号"`
		CreatedAt    string `json:"createdAt"`
	}
	PointLogListRes struct {
		PageRes
		List []PointLogItem `json:"list"`
	}
)
