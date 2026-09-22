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
	// InviteRecords 邀请激励记录（011-member-center 微扩: api 有该端点而接口缺定义, 同 D6/D7 模式）。
	InviteRecords(ctx context.Context, userId int64, page model.PageReq) (*model.PageResult[model.InviteRecordItem], error)

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

	// ---- 查询面（017 契约微扩: 端点已有而接口缺方法） ----
	// Records 我的佣金记录分页（status 0=全部）。
	Records(ctx context.Context, userId int64, status int, page model.PageReq) (*model.PageResult[DistRecordItem], error)
	// RuleQuery 佣金比例查询（按商品: 商品覆盖 > 分类默认; 未命中返回空列表）。
	RuleQuery(ctx context.Context, spuId int64) ([]DistRuleHit, error)
	// ShareCode 我的推广码（稳定: user.share_code; 为空则首访生成）。
	ShareCode(ctx context.Context, userId int64) (string, string, error)

	// ---- 账户与提现 ----
	Account(ctx context.Context, userId int64) (*model.DistAccount, error)
	AccountLogs(ctx context.Context, userId int64, bizType int, page model.PageReq) (*model.PageResult[model.AccountLogItem], error)
	// WithdrawApply 提现申请（余额校验 → 冻结 → 单据; 幂等由单号+状态机保证）。
	WithdrawApply(ctx context.Context, userId int64, amount string) (string, error)
	WithdrawList(ctx context.Context, userId int64, status int, page model.PageReq) (*model.PageResult[model.WithdrawItem], error)
}

// ---------- admin 侧（017 契约微扩 D3: 端点已有而接口缺管理方法, 批次 01/10 D3 同例） ----------

// AdminDistributorItem 后台推广员列表项。
// 注意: distribution_user 无 level 列（V29 预留未落地）→ api 的 Level 恒 0 占位（记账降级）。
type AdminDistributorItem struct {
	Id        int64  `json:"id"`
	UserId    int64  `json:"userId"`
	Nickname  string `json:"nickname" dc:"脱敏"`
	Status    int    `json:"status"`
	ApplyTime string `json:"applyTime"`
	AuditTime string `json:"auditTime"`
}

// AdminDistRuleItem 后台佣金规则列表项。
type AdminDistRuleItem struct {
	Id         int64  `json:"id"`
	ScopeType  int    `json:"scopeType" dc:"1分类 2商品"`
	ScopeId    int64  `json:"scopeId"`
	ScopeDesc  string `json:"scopeDesc" dc:"分类名/商品名"`
	Level1Rate string `json:"level1Rate"`
	Level2Rate string `json:"level2Rate"`
	Status     int    `json:"status"`
}

// AdminDistRecordItem 后台佣金记录项。
type AdminDistRecordItem struct {
	OrderNo     string `json:"orderNo"`
	Beneficiary string `json:"beneficiary" dc:"脱敏昵称"`
	Level       int    `json:"level"`
	BaseAmount  string `json:"baseAmount"`
	Amount      string `json:"amount"`
	Status      int    `json:"status"`
	SettleTime  string `json:"settleTime"`
}

// AdminInviteRecordItem 后台邀请激励记录项。
type AdminInviteRecordItem struct {
	Inviter    string `json:"inviter" dc:"脱敏昵称"`
	NewUser    string `json:"newUser" dc:"脱敏昵称"`
	RewardDesc string `json:"rewardDesc"`
	Trigger    int    `json:"trigger" dc:"1注册 2首单"`
	CreatedAt  string `json:"createdAt"`
}

// DistRecordItem C 端佣金记录项（api DistRecordItem 同形）。
type DistRecordItem struct {
	OrderNo    string `json:"orderNo"`
	Level      int    `json:"level"`
	Amount     string `json:"amount"`
	Status     int    `json:"status"`
	SettleTime string `json:"settleTime"`
	CreatedAt  string `json:"createdAt"`
}

// DistRuleHit C 端佣金比例查询命中项。
type DistRuleHit struct {
	ScopeDesc  string `json:"scopeDesc"`
	Level1Rate string `json:"level1Rate"`
	Level2Rate string `json:"level2Rate"`
}

// IDistributionAdminLogic 分销后台管理面（017 微扩; 与 C 端同域同实现）。
type IDistributionAdminLogic interface {
	AdminDistributorList(ctx context.Context, status int, keyword string, page model.PageReq) (*model.PageResult[AdminDistributorItem], error)
	AdminDistributorAudit(ctx context.Context, id int64, pass bool) error
	AdminDistributorFreeze(ctx context.Context, id int64, freeze bool) error
	AdminRuleList(ctx context.Context, page model.PageReq) (*model.PageResult[AdminDistRuleItem], error)
	AdminRuleCreate(ctx context.Context, scopeType int, scopeId int64, l1, l2 string) (int64, error)
	AdminRuleUpdate(ctx context.Context, id int64, l1, l2 string, status *int) error
	AdminRuleDelete(ctx context.Context, id int64) error
	AdminRecordList(ctx context.Context, status int, orderNo string, page model.PageReq) (*model.PageResult[AdminDistRecordItem], error)
	AdminWithdrawList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.WithdrawItem], error)
	AdminWithdrawAudit(ctx context.Context, withdrawNo string, pass bool, reason string) error
	AdminWithdrawPay(ctx context.Context, withdrawNo string, success bool, channelOrderNo, failReason string) error
	AdminInviteRecords(ctx context.Context, page model.PageReq) (*model.PageResult[AdminInviteRecordItem], error)
}
