// member_admin.go 会员管理（admin, 018 批次 12）——接口定义。
// 规则: 禁用/启用=条件更新判行数（同值幂等）; 改绑手机号 uk_phone_hash 1062→业务码;
// 列表手机号脱敏（批次 05 I1 口径: 脱敏仅展示层）。
package user

import (
	"context"

	"ecboot/internal/model"
)

// AdminMemberItem 后台会员列表/详情项。
type AdminMemberItem struct {
	UserId      int64             `json:"userId"`
	Nickname    string            `json:"nickname"`
	Phone       string            `json:"phone" dc:"脱敏"`
	Gender      int               `json:"gender"`
	Level       int64             `json:"level" dc:"等级ID(user_level_rule.id)"`
	GrowthValue int               `json:"growthValue"`
	Status      int               `json:"status" dc:"1正常 2禁用"`
	Assets      map[string]string `json:"assets" dc:"资产概要{points,balance,coupons}"`
	OrderStats  map[string]int64  `json:"orderStats" dc:"订单统计{total,finished,refunded}"`
	CreatedAt   string            `json:"createdAt"`
}

// IMemberAdminLogic 会员管理（admin）。
type IMemberAdminLogic interface {
	// AdminList 会员列表（phone 精确 / status / keyword 昵称; 分页; 手机号脱敏）。
	AdminList(ctx context.Context, phone string, status int, keyword string, page model.PageReq) (*model.PageResult[AdminMemberItem], error)
	// AdminDetail 会员详情（含资产概要与订单统计）。
	AdminDetail(ctx context.Context, userId int64) (*AdminMemberItem, error)
	// Disable 禁用/启用（条件更新判行数; 同值幂等）。
	Disable(ctx context.Context, userId int64, disable bool, reason string) error
	// RebindPhone 改绑手机号（uk_phone_hash 1062→业务码）。
	RebindPhone(ctx context.Context, userId int64, newPhone string) error
}
