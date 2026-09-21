// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AccountLog is the golang structure of table account_log for DAO operations like Where/Data.
type AccountLog struct {
	g.Meta       `orm:"table:account_log, do:true"`
	Id           any         // 流水ID
	UserId       any         // 用户ID
	BizType      any         // 业务类型:1佣金入账 2提现冻结 3提现完成 4提现失败回退 5冲销扣回 6余额消费冻结 7余额消费完成 8余额消费退回
	Amount       any         // 变动金额(有符号:入账正,冻结/扣回负)
	BalanceAfter any         // 变动后余额快照(对账用)
	FrozenAfter  any         // 变动后冻结快照(对账用)
	BizNo        any         // 关联业务单号(commission_record.id或withdraw_order.withdraw_no)
	CreatedAt    *gtime.Time // 创建时间
}
