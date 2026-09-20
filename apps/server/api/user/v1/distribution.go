package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 申请成为推广员
	DistApplyReq struct {
		g.Meta `path:"/distribution/apply" method:"POST" summary:"申请推广员"`
	}
	DistApplyRes struct {
		Status int `json:"status" dc:"1待审核"`
	}

	// 推广员状态与等级
	DistStatusReq struct {
		g.Meta `path:"/distribution/status" method:"GET" summary:"我的推广员状态"`
	}
	DistStatusRes struct {
		Status int `json:"status" dc:"0未申请 1待审核 2通过 3冻结"`
		Level  int `json:"level" dc:"推广员等级"`
	}

	// 我的邀请关系（直接上级 + 下级列表）
	DistRelationInvitee struct {
		UserId   string `json:"userId"`
		Nickname string `json:"nickname" dc:"昵称(脱敏)"`
		BindTime string `json:"bindTime"`
	}
	DistRelationReq struct {
		g.Meta `path:"/distribution/relations" method:"GET" summary:"我的邀请关系"`
		PageReq
	}
	DistRelationRes struct {
		Inviter  map[string]string     `json:"inviter" dc:"直接上级(无则为空)"`
		Invitees []DistRelationInvitee `json:"invitees" dc:"我邀请的下级"`
		Total    int64                 `json:"total" dc:"下级总数"`
	}

	// 可见佣金比例（按商品查询: 商品覆盖 > 分类默认）
	DistRuleItem struct {
		ScopeDesc  string `json:"scopeDesc" dc:"作用域描述(商品名/分类名)"`
		Level1Rate string `json:"level1Rate" dc:"一级比例%"`
		Level2Rate string `json:"level2Rate" dc:"二级比例%"`
	}
	DistRuleQueryReq struct {
		g.Meta `path:"/distribution/commission-rules" method:"GET" summary:"佣金比例查询"`
		SpuId  string `json:"spuId" dc:"商品ID"`
	}
	DistRuleQueryRes struct {
		List []DistRuleItem `json:"list"`
	}

	// 佣金记录
	DistRecordItem struct {
		OrderNo    string `json:"orderNo"`
		Level      int    `json:"level" dc:"1一级 2二级"`
		Amount     string `json:"amount" dc:"佣金(元)"`
		Status     int    `json:"status" dc:"1待结算 2已结算 3已失效 4欠款冲销"`
		SettleTime string `json:"settleTime" dc:"结算时间"`
		CreatedAt  string `json:"createdAt"`
	}
	DistRecordListReq struct {
		g.Meta `path:"/distribution/records" method:"GET" summary:"佣金记录"`
		Status int `json:"status" dc:"状态筛选"`
		PageReq
	}
	DistRecordListRes struct {
		PageRes
		List []DistRecordItem `json:"list"`
	}

	// 佣金账户
	DistAccountReq struct {
		g.Meta `path:"/distribution/account" method:"GET" summary:"佣金账户"`
	}
	DistAccountRes struct {
		Balance string `json:"balance" dc:"可用余额(可负)"`
		Frozen  string `json:"frozen" dc:"冻结金额"`
	}

	// 账户流水
	DistAccountLogItem struct {
		BizType      int    `json:"bizType" dc:"1入账 2提现冻结 3提现完成 4提现回退 5冲销 6消费冻结 7消费完成 8消费退回"`
		Amount       string `json:"amount" dc:"变动(有符号)"`
		BalanceAfter string `json:"balanceAfter"`
		BizNo        string `json:"bizNo" dc:"关联单据"`
		CreatedAt    string `json:"createdAt"`
	}
	DistAccountLogListReq struct {
		g.Meta  `path:"/distribution/account/logs" method:"GET" summary:"账户流水"`
		BizType int `json:"bizType" dc:"类型筛选"`
		PageReq
	}
	DistAccountLogListRes struct {
		PageRes
		List []DistAccountLogItem `json:"list"`
	}

	// 提现申请
	WithdrawCreateReq struct {
		g.Meta `path:"/distribution/withdraws" method:"POST" summary:"提现申请"`
		Amount string `json:"amount" v:"required" dc:"提现金额(元)"`
	}
	WithdrawCreateRes struct {
		WithdrawNo string `json:"withdrawNo" dc:"提现单号"`
	}

	WithdrawItem struct {
		WithdrawNo string `json:"withdrawNo"`
		Amount     string `json:"amount"`
		Status     int    `json:"status" dc:"10待审 20过审 30打款中 40成功 50拒绝 60失败回退"`
		CreatedAt  string `json:"createdAt"`
	}
	WithdrawListReq struct {
		g.Meta `path:"/distribution/withdraws" method:"GET" summary:"提现列表"`
		Status int `json:"status" dc:"状态筛选"`
		PageReq
	}
	WithdrawListRes struct {
		PageRes
		List []WithdrawItem `json:"list"`
	}

	// 邀请激励记录
	InviteRecordItem struct {
		NewUser    string `json:"newUser" dc:"新用户(脱敏昵称)"`
		RewardDesc string `json:"rewardDesc" dc:"奖励说明"`
		Status     int    `json:"status" dc:"1已发放"`
		CreatedAt  string `json:"createdAt"`
	}
	InviteRecordListReq struct {
		g.Meta `path:"/distribution/invite-records" method:"GET" summary:"邀请激励记录"`
		PageReq
	}
	InviteRecordListRes struct {
		PageRes
		List []InviteRecordItem `json:"list"`
	}

	// 我的推广码
	ShareCodeReq struct {
		g.Meta `path:"/share-code" method:"GET" summary:"我的推广码"`
	}
	ShareCodeRes struct {
		ShareCode string `json:"shareCode" dc:"推广码"`
		ShareLink string `json:"shareLink" dc:"分享链接"`
	}
)
