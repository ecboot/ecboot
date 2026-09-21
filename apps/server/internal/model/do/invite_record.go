// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// InviteRecord is the golang structure of table invite_record for DAO operations like Where/Data.
type InviteRecord struct {
	g.Meta        `orm:"table:invite_record, do:true"`
	Id            any         // 激励记录ID
	NewUserId     any         // 新用户ID(仅可被激励一次)
	InviterId     any         // 邀请人用户ID
	RewardType    any         // 奖励类型:1优惠券
	RewardTrigger any         // 奖励触发时机:1注册即发 2首单后发(防刷,需订单回调触发)
	RewardRef     any         // 奖励载体ID(如user_coupon.id)
	Status        any         // 状态:1已发放(发放与订单同事务,幂等)
	CreatedAt     *gtime.Time // 创建时间
}
