package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/model"
)

type (
	// 推广员列表（审核/冻结入口）
	AdminDistributorListReq struct {
		g.Meta  `path:"/distributors" method:"GET" summary:"推广员列表"`
		Status  int    `json:"status" dc:"1待审 2通过 3冻结"`
		Keyword string `json:"keyword" dc:"用户昵称/ID"`
		model.PageReq
	}
	AdminDistributorItem struct {
		Id        string `json:"id"`
		UserId    string `json:"userId"`
		Nickname  string `json:"nickname" dc:"昵称(脱敏)"`
		Level     int    `json:"level" dc:"推广员等级"`
		Status    int    `json:"status"`
		ApplyTime string `json:"applyTime"`
		AuditTime string `json:"auditTime"`
	}
	AdminDistributorListRes struct {
		model.PageRes
		List []AdminDistributorItem `json:"list"`
	}

	// 推广员审核
	// 权限: distribution:audit
	AdminDistributorAuditReq struct {
		g.Meta `path:"/distributors/{id}/audit" method:"POST" summary:"推广员审核"`
		Id     string `json:"id" v:"required" dc:"推广员记录ID"`
		Pass   bool   `json:"pass" dc:"通过/拒绝"`
	}
	AdminDistributorAuditRes struct {
		Success bool `json:"success"`
	}

	// 冻结/解冻
	// 权限: distribution:audit
	AdminDistributorFreezeReq struct {
		g.Meta `path:"/distributors/{id}/freeze" method:"POST" summary:"冻结/解冻推广员"`
		Id     string `json:"id" v:"required" dc:"推广员记录ID"`
		Freeze bool   `json:"freeze" dc:"true冻结 false解冻"`
	}
	AdminDistributorFreezeRes struct {
		Success bool `json:"success"`
	}

	// 佣金规则列表
	AdminDistRuleListReq struct {
		g.Meta `path:"/commission-rules" method:"GET" summary:"佣金规则列表"`
		model.PageReq
	}
	AdminDistRuleItem struct {
		Id         string `json:"id"`
		ScopeType  int    `json:"scopeType" dc:"1分类 2商品"`
		ScopeDesc  string `json:"scopeDesc" dc:"作用域名称"`
		Level1Rate string `json:"level1Rate" dc:"一级比例%"`
		Level2Rate string `json:"level2Rate" dc:"二级比例%"`
		Status     int    `json:"status"`
	}
	AdminDistRuleListRes struct {
		model.PageRes
		List []AdminDistRuleItem `json:"list"`
	}

	// 创建佣金规则（作用域唯一）
	// 权限: distribution:rule:manage
	AdminDistRuleCreateReq struct {
		g.Meta     `path:"/commission-rules" method:"POST" summary:"创建佣金规则"`
		ScopeType  int    `json:"scopeType" v:"required|in:1,2" dc:"1分类 2商品"`
		ScopeId    string `json:"scopeId" v:"required" dc:"作用域目标ID"`
		Level1Rate string `json:"level1Rate" v:"required" dc:"一级比例%"`
		Level2Rate string `json:"level2Rate" dc:"二级比例%"`
	}
	AdminDistRuleCreateRes struct {
		Id string `json:"id"`
	}

	// 权限: distribution:rule:manage
	AdminDistRuleUpdateReq struct {
		g.Meta     `path:"/commission-rules/{id}" method:"PUT" summary:"修改佣金规则"`
		Id         string `json:"id" v:"required" dc:"规则ID"`
		Level1Rate string `json:"level1Rate" dc:"一级比例%"`
		Level2Rate string `json:"level2Rate" dc:"二级比例%"`
		Status     int    `json:"status" dc:"状态"`
	}
	AdminDistRuleUpdateRes struct {
		Success bool `json:"success"`
	}

	// 权限: distribution:rule:manage
	AdminDistRuleDeleteReq struct {
		g.Meta `path:"/commission-rules/{id}" method:"DELETE" summary:"删除佣金规则(软删)"`
		Id     string `json:"id" v:"required" dc:"规则ID"`
	}
	AdminDistRuleDeleteRes struct {
		Success bool `json:"success"`
	}

	// 佣金记录（全局）
	AdminDistRecordListReq struct {
		g.Meta  `path:"/commission-records" method:"GET" summary:"佣金记录"`
		Status  int    `json:"status" dc:"状态筛选"`
		OrderNo string `json:"orderNo" dc:"订单号"`
		model.PageReq
	}
	AdminDistRecordItem struct {
		OrderNo     string `json:"orderNo"`
		Beneficiary string `json:"beneficiary" dc:"受益人(脱敏昵称)"`
		Level       int    `json:"level" dc:"1一级 2二级"`
		BaseAmount  string `json:"baseAmount" dc:"计佣基数"`
		Amount      string `json:"amount"`
		Status      int    `json:"status"`
		SettleTime  string `json:"settleTime"`
	}
	AdminDistRecordListRes struct {
		model.PageRes
		List []AdminDistRecordItem `json:"list"`
	}

	// 提现列表
	AdminWithdrawListReq struct {
		g.Meta `path:"/withdraws" method:"GET" summary:"提现列表"`
		Status int `json:"status" dc:"状态筛选"`
		model.PageReq
	}
	AdminWithdrawItem struct {
		WithdrawNo string `json:"withdrawNo"`
		UserId     string `json:"userId"`
		Amount     string `json:"amount"`
		Status     int    `json:"status" dc:"10待审 20过审 30打款中 40成功 50拒绝 60失败回退"`
		CreatedAt  string `json:"createdAt"`
	}
	AdminWithdrawListRes struct {
		model.PageRes
		List []AdminWithdrawItem `json:"list"`
	}

	// 提现审核（通过→打款中冻结; 拒绝→回退）
	// 权限: distribution:withdraw:audit
	AdminWithdrawAuditReq struct {
		g.Meta     `path:"/withdraws/{withdrawNo}/audit" method:"POST" summary:"提现审核"`
		WithdrawNo string `json:"withdrawNo" v:"required" dc:"提现单号"`
		Pass       bool   `json:"pass" dc:"通过/拒绝"`
		Reason     string `json:"reason" dc:"拒绝原因"`
	}
	AdminWithdrawAuditRes struct {
		Success bool `json:"success"`
	}

	// 打款结果登记（成功核销 / 失败回退）
	// 权限: distribution:withdraw:pay
	AdminWithdrawPayReq struct {
		g.Meta         `path:"/withdraws/{withdrawNo}/pay" method:"POST" summary:"提现打款登记"`
		WithdrawNo     string `json:"withdrawNo" v:"required" dc:"提现单号"`
		Success        bool   `json:"success" dc:"打款结果"`
		ChannelOrderNo string `json:"channelOrderNo" dc:"渠道打款单号(成功必填,唯一幂等)"`
		FailReason     string `json:"failReason" dc:"失败原因"`
	}
	AdminWithdrawPayRes struct {
		Success bool `json:"success"`
	}

	// 邀请激励记录
	AdminInviteRecordListReq struct {
		g.Meta `path:"/invite-records" method:"GET" summary:"邀请激励记录"`
		model.PageReq
	}
	AdminInviteRecordItem struct {
		Inviter    string `json:"inviter" dc:"邀请人(脱敏昵称)"`
		NewUser    string `json:"newUser" dc:"新用户(脱敏昵称)"`
		RewardDesc string `json:"rewardDesc"`
		Trigger    int    `json:"trigger" dc:"1注册 2首单"`
		CreatedAt  string `json:"createdAt"`
	}
	AdminInviteRecordListRes struct {
		model.PageRes
		List []AdminInviteRecordItem `json:"list"`
	}
)
