// profile.go 个人资料域——表: user（V1/V11/V19/V27/V28 列）。
// 规则: 手机号仅精确检索+脱敏展示（FR-004）；成长值只增不减；等级由成长值门槛判定。
package user

import "context"

// IProfileLogic 个人资料与等级。
type IProfileLogic interface {
	// Detail 查询资料（脱敏手机号、等级名称、成长值）。
	Detail(ctx context.Context, userId int64) (*ProfileDetail, error)
	// Update 修改昵称/头像/性别。
	Update(ctx context.Context, userId int64, in ProfileUpdateInput) error
	// GrowthAdd 成长值累加（只增不减; 签到/消费/分享等事件触发）。
	GrowthAdd(ctx context.Context, userId int64, value int) error
	// LevelRecalc 依据 user_level_rule 门槛重算等级（成长值变更后调用）。
	LevelRecalc(ctx context.Context, userId int64) error
	// LoginLogs 近 30 天登录记录（FR-012/FR-002）。
	LoginLogs(ctx context.Context, userId int64, page PageQuery) (*PageResult[LoginLogItem], error)
}

type ProfileDetail struct {
	UserId      string `json:"userId"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Gender      int    `json:"gender"`
	PhoneMasked string `json:"phone" dc:"138****1234"`
	Level       int64  `json:"level" dc:"等级ID"`
	LevelName   string `json:"levelName"`
	GrowthValue int    `json:"growthValue"`
}

type ProfileUpdateInput struct {
	Nickname string
	Avatar   string
	Gender   int
}

type LoginLogItem struct {
	Channel   int    `json:"channel"`
	Status    int    `json:"status"`
	Ip        string `json:"ip" dc:"脱敏"`
	UserAgent string `json:"userAgent"`
	CreatedAt string `json:"createdAt"`
}
