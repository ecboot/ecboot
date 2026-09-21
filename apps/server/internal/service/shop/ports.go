// ports.go 跨域接口锁定（依赖倒置）——shop 域定义、user 域实现并注册。
// 契约权威: specs/006-trade-loop/data-model.md §三。
// 事务感知: 全部方法第一个参数为 tx（*gdb.TX），由订单事务统一提交/回滚。
package shop

import "context"

// txContext 事务句柄（由 order 编排传入; 具体 *gdb.TX 在实现内断言）。
type TX any

// ICouponTrade 券核销/退回（下单事务内调用）。
type ICouponTrade interface {
	// ConsumeForOrder 核销用户券（校验归属+未用+未过期; 置已使用+绑单号）。
	ConsumeForOrder(tx TX, ctx context.Context, userId, userCouponId int64, orderNo string) error
	// ReturnForOrder 订单取消退回券（置可再用, 过期时间不变）。
	ReturnForOrder(tx TX, ctx context.Context, userCouponId int64) error
}

// IPointTrade 积分消耗/回退。
type IPointTrade interface {
	// ConsumeForOrder 下单消耗积分（校验余额, 扣减+流水; 可致负=欠款）。
	ConsumeForOrder(tx TX, ctx context.Context, userId int64, points int, orderNo string) error
	// RefundForOrder 退款回退积分（可致负）。
	RefundForOrder(tx TX, ctx context.Context, userId int64, points int, orderNo string) error
}

// IAccountTrade 佣金余额（资金语义: 冻结→核销→解冻→退回, bizType 6/7/8）。
type IAccountTrade interface {
	// FreezeForOrder 下单冻结（可用余额校验, balance→frozen）。
	FreezeForOrder(tx TX, ctx context.Context, userId int64, amountFen int64, orderNo string) error
	// SettleFrozen 支付成功核销冻结（frozen 扣减）。
	SettleFrozen(tx TX, ctx context.Context, userId int64, amountFen int64, orderNo string) error
	// Unfreeze 取消/失败解冻（frozen→balance）。
	Unfreeze(tx TX, ctx context.Context, userId int64, amountFen int64, orderNo string) error
	// RefundFrozen 售后退回余额（按分摊, 直接入 balance）。
	RefundFrozen(tx TX, ctx context.Context, userId int64, amountFen int64, orderNo string) error
}

// INotifyEnqueue 通知事件位（订单状态迁移触发; 投递属消息域）。
type INotifyEnqueue interface {
	OrderPaid(ctx context.Context, userId int64, orderNo string)
	OrderShipped(ctx context.Context, userId int64, orderNo string)
	OrderCancelled(ctx context.Context, userId int64, orderNo string)
	OrderCompleted(ctx context.Context, userId int64, orderNo string)
}

// 注册变量（user 域 bootstrap 装配时注入; 未注入则相关能力降级跳过并告警）。
var (
	CouponTrade  ICouponTrade
	PointTrade   IPointTrade
	AccountTrade IAccountTrade
	NotifyEnq    INotifyEnqueue
)
