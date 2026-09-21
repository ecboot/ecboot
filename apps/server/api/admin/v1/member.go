package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/model"
)

type (
	// 会员列表（手机号精确检索; 脱敏展示）
	AdminMemberListReq struct {
		g.Meta  `path:"/members" method:"GET" summary:"会员列表"`
		Phone   string `json:"phone" dc:"完整手机号精确匹配"`
		Status  int    `json:"status" dc:"1正常 2禁用"`
		Keyword string `json:"keyword" dc:"昵称/ID"`
		model.PageReq
	}
	AdminMemberItem struct {
		UserId    string `json:"userId"`
		Nickname  string `json:"nickname"`
		Phone     string `json:"phone" dc:"手机号(脱敏)"`
		Level     int    `json:"level" dc:"会员等级"`
		Status    int    `json:"status" dc:"1正常 2禁用"`
		CreatedAt string `json:"createdAt" dc:"注册时间"`
	}
	AdminMemberListRes struct {
		model.PageRes
		List []AdminMemberItem `json:"list"`
	}

	// 会员详情（含资产概要）
	AdminMemberDetailReq struct {
		g.Meta `path:"/members/{userId}" method:"GET" summary:"会员详情"`
		UserId string `json:"userId" v:"required" dc:"用户ID"`
	}
	AdminMemberDetailRes struct {
		UserId      string            `json:"userId"`
		Nickname    string            `json:"nickname"`
		Phone       string            `json:"phone" dc:"脱敏"`
		Gender      int               `json:"gender"`
		Level       int               `json:"level"`
		GrowthValue int               `json:"growthValue"`
		Status      int               `json:"status"`
		Assets      map[string]string `json:"assets" dc:"资产概要{points,balance,coupons}"`
		OrderStats  map[string]int64  `json:"orderStats" dc:"订单统计{total,finished,refunded}"`
		CreatedAt   string            `json:"createdAt"`
	}

	// 禁用/启用（审计留痕）
	// 权限: member:update
	AdminMemberDisableReq struct {
		g.Meta  `path:"/members/{userId}/disable" method:"POST" summary:"禁用/启用会员"`
		UserId  string `json:"userId" v:"required" dc:"用户ID"`
		Disable bool   `json:"disable" dc:"true禁用 false启用"`
		Reason  string `json:"reason" dc:"原因"`
	}
	AdminMemberDisableRes struct {
		Success bool `json:"success"`
	}

	// 改绑手机号（旧号解占; 审计留痕）
	// 权限: member:update
	AdminMemberRebindPhoneReq struct {
		g.Meta   `path:"/members/{userId}/rebind-phone" method:"POST" summary:"改绑手机号"`
		UserId   string `json:"userId" v:"required" dc:"用户ID"`
		NewPhone string `json:"newPhone" v:"required" dc:"新手机号"`
	}
	AdminMemberRebindPhoneRes struct {
		Success bool `json:"success"`
	}
)
