// distribution.go 分销域——表: user_relation(V16/V23 锁定) / distribution_user(V29 等级) /
// commission_rule(V16 两级作用域) / commission_record(V22 冲销) / user_account(V25 可负) /
// account_log / withdraw_order(V16 幂等) / invite_record(V28 双时机) / share_record(V28) / user.share_code。
// 合规红线(AGENTS): 关系链两级封顶(ADR-0003 结构强制); 自邀/互环应用层拒绝(契约 R6);
// 佣金基数=订单项实付(不含运费/积分抵扣); 归因(分享>关系链>自然流量, 窗口配置化)。
package user

import (
	"context"

	"ecboot/internal/model"
)

// IDistributionLogic 分销。
type IDistributionLogic interface {
	// ---- 关系与资质 ----
	// Apply 申请推广员（已申请/已是 → 60001）。
	Apply(ctx context.Context, userId int64) error
	// Status 我的推广员状态与等级。
	Status(ctx context.Context, userId int64) (*model.DistStatus, error)
	// Relations 我的邀请关系（上级+下级分页）。
	Relations(ctx context.Context, userId int64, page model.PageReq) (*model.DistRelationResult, error)
	// BindRelation 绑定直接上级（注册/首次归因时; 保护期内可换绑, 期满锁定 V23）。
	BindRelation(ctx context.Context, userId, inviterId int64, channel int) error

	// ---- 归因与佣金 ----
	// ShareReport 分享行为上报（归因窗口起点; 游客可报, user_id 可空）。
	ShareReport(ctx context.Context, userId int64, spuId int64, channel int, scene string) error
	// SettleOrder 订单佣金计提（确认收货事件触发）:
	// 归因判定(分享窗口内归因人 > 关系链两级 > 自然流量不计) → 按规则生成待结算记录（自购返佣: 受益人=本人一级）。
	SettleOrder(ctx context.Context, orderNo string) error
	// ConfirmSettle 结算保护期(配置 commission.settle.protect_days)满自动入账（定时任务）。
	ConfirmSettle(ctx context.Context) (int64, error)
	// ReverseOnRefund 售后退款冲销（售后完成事件: 未结算置失效 / 已结算生成负额冲销记录+扣回余额, FR-019）。
	ReverseOnRefund(ctx context.Context, orderItemId int64) error

	// ---- 账户与提现 ----
	Account(ctx context.Context, userId int64) (*model.DistAccount, error)
	AccountLogs(ctx context.Context, userId int64, bizType int, page model.PageReq) (*model.PageResult[model.AccountLogItem], error)
	// WithdrawApply 提现申请（余额校验 → 冻结 → 单据; 幂等由单号+状态机保证）。
	WithdrawApply(ctx context.Context, userId int64, amount string) (string, error)
	WithdrawList(ctx context.Context, userId int64, status int, page model.PageReq) (*model.PageResult[model.WithdrawItem], error)
}
