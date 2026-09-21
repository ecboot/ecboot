// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PointLog is the golang structure for table point_log.
type PointLog struct {
	Id           uint64      `json:"id"           orm:"id"            ` // 流水ID
	UserId       uint64      `json:"userId"       orm:"user_id"       ` // 用户ID
	BizType      int         `json:"bizType"      orm:"biz_type"      ` // 业务类型:1签到 2消费获得 3下单消耗 4退款回退 5分享获得 6评价获得 7注册赠送 8邀请奖励 9过期扣减
	Points       int         `json:"points"       orm:"points"        ` // 积分变动(有符号:获得正,消耗/回退负)
	BalanceAfter int         `json:"balanceAfter" orm:"balance_after" ` // 变动后余额快照(对账用)
	OrderNo      string      `json:"orderNo"      orm:"order_no"      ` // 关联订单号
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` // 创建时间
}
