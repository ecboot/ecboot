// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// WithdrawOrder is the golang structure of table withdraw_order for DAO operations like Where/Data.
type WithdrawOrder struct {
	g.Meta          `orm:"table:withdraw_order, do:true"`
	Id              any         // 提现单ID
	WithdrawNo      any         // 提现单号(全局唯一)
	UserId          any         // 用户ID
	Amount          any         // 提现金额
	WithdrawChannel any         // 渠道:1微信商家转账(V1单渠道)
	ChannelOrderNo  any         // 渠道打款单号(回填;与渠道联合唯一,防重复打款)
	Status          any         // 状态:10待审核 20审核通过 30打款中 40成功 50审核拒绝 60打款失败已回退
	AuditTime       *gtime.Time // 审核时间
	PayTime         *gtime.Time // 打款成功时间
	FailReason      any         // 拒绝/失败原因
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
}
