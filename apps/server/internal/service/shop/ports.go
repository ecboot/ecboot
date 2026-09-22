// ports.go 跨域接口锁定（依赖倒置）——shop 域定义、user 域实现并注册。
// 契约权威: specs/006-trade-loop/data-model.md §三。
// 事务感知: 全部方法第一个参数为 tx（*gdb.TX），由订单事务统一提交/回滚。
package shop

import (
	"context"

	"ecboot/internal/model"
)

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

// ICouponQuery 券**查询**口（012 评审 I5）: 结算试算要回"当前可用券列表", 属 shop 域端点的出参义务。
// 与 ICouponTrade（下单事务内的核销/退回口）刻意分开: 查询口无事务语义, 本批即可接;
// 核销口随"下单券/积分/余额三段联动"一并延后（PROGRESS §五 记账）。
type ICouponQuery interface {
	// UsableForOrder 返回该用户对给定商品金额（元, 字符串）可用的券, 抵扣降序。
	UsableForOrder(ctx context.Context, userId int64, goodsAmountYuan string) ([]model.UsableCouponBrief, error)
}

// ICommissionReverse 佣金冲销事件出口（013 售后域）: 售后完成时投递, 避免"退了款而佣金仍计提"。
// 归属: shop 域定义, 分销/资金域（批次 11）实现并装配; 未注入则降级跳过并告警（同下方注册变量约定）。
// 本批只**投递事件**, 不落账（结算语义属批次 11, 与 Confirm 的"计提由分销结算推进"同源）。
type ICommissionReverse interface {
	// ReverseForAfterSale 订单售后完成时投递冲销意图（refundFen 为本次退款金额, 用于按额冲销）。
	ReverseForAfterSale(ctx context.Context, orderNo, afterSaleNo string, refundFen int64)
}

// IRiskHit 风控命中判定（015 营销 C 端）: 帮砍/助力入口调用。
// 归属: shop 域定义; 实现属**批次 12**（风控规则引擎）→ 未装配时**放行 + 告警**（不阻塞玩法可用）。
// 依据: 迁移 000029 两处注释明写"砍价/助力是被刷重灾区——帮砍/助力入口须挂风控"（schema 级要求）。
type IRiskHit interface {
	// Hit 判定是否拦截（返回 true 表示拦截; ruleType 沿用 risk_rule.rule_type 口径）。
	Hit(ctx context.Context, userId int64, ruleType int, payload string) (blocked bool, err error)
}

// IAssistReward 助力发奖意图（015 营销 C 端）: 达标时投递。
// 归属: shop 域定义; 实际发放（发券/发积分）跨 user 域, 属后续批次 → 未装配时告警降级。
// 与 ICommissionReverse 同型: 只**投递意图**, 不落账（结算/发放后置）。
type IAssistReward interface {
	// GrantForAssist 投递发奖意图（rewardType 1 券 / 2 积分; rewardRef 为券 id 或积分点数载体）。
	GrantForAssist(ctx context.Context, userId, activityId, recordId int64, rewardType int, rewardRef int64)
}

// 注册变量（user 域 bootstrap 装配时注入; 未注入则相关能力降级跳过并告警）。
var (
	CouponTrade       ICouponTrade
	PointTrade        IPointTrade
	AccountTrade      IAccountTrade
	NotifyEnq         INotifyEnqueue
	CouponQuery       ICouponQuery
	CommissionReverse ICommissionReverse
	RiskHit           IRiskHit
	AssistReward      IAssistReward
)
