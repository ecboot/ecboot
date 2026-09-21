// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PointLog is the golang structure of table point_log for DAO operations like Where/Data.
type PointLog struct {
	g.Meta       `orm:"table:point_log, do:true"`
	Id           any         // 流水ID
	UserId       any         // 用户ID
	BizType      any         // 业务类型:1签到 2消费获得 3下单消耗 4退款回退 5分享获得 6评价获得 7注册赠送 8邀请奖励 9过期扣减
	Points       any         // 积分变动(有符号:获得正,消耗/回退负)
	BalanceAfter any         // 变动后余额快照(对账用)
	OrderNo      any         // 关联订单号
	CreatedAt    *gtime.Time // 创建时间
}
