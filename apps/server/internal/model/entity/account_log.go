// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AccountLog is the golang structure for table account_log.
type AccountLog struct {
	Id           uint64      `json:"id"           orm:"id"            ` // 流水ID
	UserId       uint64      `json:"userId"       orm:"user_id"       ` // 用户ID
	BizType      int         `json:"bizType"      orm:"biz_type"      ` // 业务类型:1佣金入账 2提现冻结 3提现完成 4提现失败回退 5冲销扣回 6余额消费冻结 7余额消费完成 8余额消费退回
	Amount       float64     `json:"amount"       orm:"amount"        ` // 变动金额(有符号:入账正,冻结/扣回负)
	BalanceAfter float64     `json:"balanceAfter" orm:"balance_after" ` // 变动后余额快照(对账用)
	FrozenAfter  float64     `json:"frozenAfter"  orm:"frozen_after"  ` // 变动后冻结快照(对账用)
	BizNo        string      `json:"bizNo"        orm:"biz_no"        ` // 关联业务单号(commission_record.id或withdraw_order.withdraw_no)
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` // 创建时间
}
