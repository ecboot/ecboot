// distribution.go 分销域——表: user_relation(V16/V23 锁定) / distribution_user(V29 等级) /
// commission_rule(V16 两级作用域) / commission_record(V22 冲销) / user_account(V25 可负) /
// account_log / withdraw_order(V16 幂等) / invite_record(V28 双时机) / share_record(V28) / user.share_code。
// 合规红线(AGENTS): 关系链两级封顶(ADR-0003 结构强制); 自邀/互环应用层拒绝(契约 R6);
// 佣金基数=订单项实付(不含运费/积分抵扣); 归因(分享>关系链>自然流量, 窗口配置化)。
package user

import "context"

// IDistributionLogic 分销。
type IDistributionLogic interface {
	// ---- 关系与资质 ----
	// Apply 申请推广员（已申请/已是 → 60001）。
	Apply(ctx context.Context, userId int64) error
	// Status 我的推广员状态与等级。
	Status(ctx context.Context, userId int64) (*DistStatus, error)
	// Relations 我的邀请关系（上级+下级分页）。
	Relations(ctx context.Context, userId int64, page PageQuery) (*DistRelationResult, error)
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
	Account(ctx context.Context, userId int64) (*DistAccount, error)
	AccountLogs(ctx context.Context, userId int64, bizType int, page PageQuery) (*PageResult[AccountLogItem], error)
	// WithdrawApply 提现申请（余额校验 → 冻结 → 单据; 幂等由单号+状态机保证）。
	WithdrawApply(ctx context.Context, userId int64, amount string) (string, error)
	WithdrawList(ctx context.Context, userId int64, status int, page PageQuery) (*PageResult[WithdrawItem], error)
}

type DistStatus struct {
	Status int `json:"status" dc:"0未申请 1待审核 2通过 3冻结"`
	Level  int `json:"level" dc:"推广员等级(V29 预留)"`
}

type DistRelationResult struct {
	Inviter  map[string]string `json:"inviter" dc:"直接上级(脱敏)"`
	Invitees []DistInvitee     `json:"invitees"`
	Total    int64             `json:"total"`
}

type DistInvitee struct {
	UserId   int64  `json:"userId"`
	Nickname string `json:"nickname" dc:"脱敏"`
	BindTime string `json:"bindTime"`
}

type DistAccount struct {
	Balance string `json:"balance" dc:"可负"`
	Frozen  string `json:"frozen"`
}

type AccountLogItem struct {
	BizType      int    `json:"bizType" dc:"1入账 2提现冻结 3完成 4回退 5冲销 6消费冻结 7消费完成 8消费退回"`
	Amount       string `json:"amount" dc:"有符号"`
	BalanceAfter string `json:"balanceAfter"`
	BizNo        string `json:"bizNo"`
	CreatedAt    string `json:"createdAt"`
}

type WithdrawItem struct {
	WithdrawNo string `json:"withdrawNo"`
	Amount     string `json:"amount"`
	Status     int    `json:"status" dc:"10待审 20过审 30打款中 40成功 50拒绝 60失败回退"`
	CreatedAt  string `json:"createdAt"`
}
