// point.go 积分与成长值账本——表: point_account / point_log / user_level_rule（V19/V26）。
// 规则(clarify 定档): 积分可消耗有流水（余额可负=欠款）; 成长值只增不减（双账本分离）;
// 积分滚动过期=自 last_earned_at 起 12 个月（过期任务调 ExpireDormant）; 批次制为升级路径（触发条件见 schema-design）。
package user

import (
	"context"

	"ecboot/internal/model"
)

// IPointLogic 积分账本（全部变动同事务: 账户余额 + 流水 balance_after 快照）。
type IPointLogic interface {
	Account(ctx context.Context, userId int64) (*model.PointAccountView, error)
	Logs(ctx context.Context, userId int64, bizType int, page model.PageReq) (*model.PageResult[model.PointLogItem], error)
	// Earn 获得（签到/消费/分享/评价/注册/邀请——bizType 1/2/5/6/7/8）; 幂等由调用方 bizNo 保证。
	Earn(ctx context.Context, userId int64, bizType int, points int, orderNo string) error
	// Consume 下单消耗（扣减, 可致负余额——欠款语义）。
	Consume(ctx context.Context, userId int64, points int, orderNo string) error
	// Refund 退款回退（bizType=4, 可致负）。
	Refund(ctx context.Context, userId int64, points int, orderNo string) error
	// ExpireDormant 滚动过期清零（定时任务: last_earned_at 超 12 个月; bizType=9）。
	ExpireDormant(ctx context.Context) (int64, error)
}
