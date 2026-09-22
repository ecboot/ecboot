// distribution_impl.go IDistributionLogic / IDistributionAdminLogic 实现（017 批次 11）。
// 红线（宪法 #1）: 关系链两级封顶——user_relation 仅存直接上级（无祖父列, ADR-0003 结构强制）,
// 二级解析=两次单列查询到顶即止, 任何路径不得第三级。
// 资金铁律（D1）: 账户变更=单行条件 UPDATE 判行数 + 同事务 account_log 双快照流水;
// 提现状态机不可逆迁移全条件更新; 打款幂等 uk_channel_order 1062 → 业务码。
package user

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/library/idgen"
	"ecboot/internal/library/money"
	"ecboot/internal/model"
)

// DistributionLogicImpl IDistributionLogic 实现。
type DistributionLogicImpl struct{}

func NewDistributionLogic() *DistributionLogicImpl { return &DistributionLogicImpl{} }

// DistributionAdminLogicImpl IDistributionAdminLogic 实现。
type DistributionAdminLogicImpl struct{}

func NewDistributionAdminLogic() *DistributionAdminLogicImpl { return &DistributionAdminLogicImpl{} }

// 编译期接口满足性断言（批次 10 M7 教训: 不靠文字承诺）。
var (
	_ IDistributionLogic      = (*DistributionLogicImpl)(nil)
	_ IDistributionAdminLogic = (*DistributionAdminLogicImpl)(nil)
)

// maskNick 昵称脱敏（分销域内自带, 与 shop 域同口径: 首尾保留+中间星号; 空名占位）。
func maskNick(nick string) string {
	r := []rune(nick)
	switch {
	case len(r) == 0:
		return "用户"
	case len(r) <= 2:
		return string(r[:1]) + "*"
	default:
		return string(r[:1]) + "*" + string(r[len(r)-1:])
	}
}

// ---------- US1 关系与资质 ----------

// Apply 申请推广员（已申请/已是 → 60001; uk_user 唯一键兜底并发）。
func (i *DistributionLogicImpl) Apply(ctx context.Context, userId int64) error {
	n, err := dao.DistributionUser.Ctx(ctx).Where(dao.DistributionUser.Columns().UserId, userId).Count()
	if err != nil {
		return gerror.Wrap(err, "查询推广员资质失败")
	}
	if n > 0 {
		return errcode.New(errcode.CodeDistAlreadyApplied, "已申请或已是推广员")
	}
	_, err = dao.DistributionUser.Ctx(ctx).Data(g.Map{
		"user_id": userId, "status": 1,
	}).Insert()
	if err != nil {
		if isDupKeyUser(err) {
			return errcode.New(errcode.CodeDistAlreadyApplied, "已申请或已是推广员")
		}
		return gerror.Wrap(err, "创建推广员申请失败")
	}
	return nil
}

// isDupKeyUser user 域内 1062 判断（与 shop 域 isDupKey 同形态, 域间不互引）。
func isDupKeyUser(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "1062") || strings.Contains(s, "Duplicate entry")
}

// Status 我的推广员状态（0 未申请）。
func (i *DistributionLogicImpl) Status(ctx context.Context, userId int64) (*model.DistStatus, error) {
	r, err := dao.DistributionUser.Ctx(ctx).
		Where(dao.DistributionUser.Columns().UserId, userId).
		Where(dao.DistributionUser.Columns().Deleted, 0).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询推广员资质失败")
	}
	if r.IsEmpty() {
		return &model.DistStatus{Status: 0}, nil // 未申请
	}
	return &model.DistStatus{Status: r["status"].Int(), Level: 0}, nil // Level: V29 预留无表结构 → 恒 0（D3 记账降级）
}

// Relations 我的邀请关系（上级 + 下级分页; 昵称脱敏）。
func (i *DistributionLogicImpl) Relations(ctx context.Context, userId int64, page model.PageReq) (*model.DistRelationResult, error) {
	page = page.Normalized()
	out := &model.DistRelationResult{Invitees: []model.DistInvitee{}}

	// 上级（直接一级; 红线: 不解析更上层）
	my, err := dao.UserRelation.Ctx(ctx).Where(dao.UserRelation.Columns().UserId, userId).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询关系链失败")
	}
	if !my.IsEmpty() {
		inviterId := my["inviter_id"].Int64()
		u, e := dao.User.Ctx(ctx).Fields(dao.User.Columns().Nickname).
			Where(dao.User.Columns().Id, inviterId).One()
		if e != nil {
			return nil, gerror.Wrap(e, "查询上级失败")
		}
		if !u.IsEmpty() {
			out.Inviter = map[string]string{"nickname": maskNick(u["nickname"].String())}
		}
	}

	// 下级分页
	total, err := dao.UserRelation.Ctx(ctx).Where(dao.UserRelation.Columns().InviterId, userId).Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计下级失败")
	}
	out.Total = int64(total)
	kids, err := dao.UserRelation.Ctx(ctx).
		Where(dao.UserRelation.Columns().InviterId, userId).
		OrderDesc(dao.UserRelation.Columns().Id).
		Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询下级失败")
	}
	for _, r := range kids {
		uid := r["user_id"].Int64()
		nick := ""
		u, e := dao.User.Ctx(ctx).Fields(dao.User.Columns().Nickname).
			Where(dao.User.Columns().Id, uid).One()
		if e == nil && !u.IsEmpty() {
			nick = maskNick(u["nickname"].String())
		}
		out.Invitees = append(out.Invitees, model.DistInvitee{
			UserId: uid, Nickname: nick, BindTime: r["bind_time"].String(),
		})
	}
	return out, nil
}

// BindRelation 绑定直接上级（注册/首次归因; 一人一链: uk_user 兜底 + 自邀拒绝）。
// 红线: 只落一级关系, 不存在"祖父"写入路径。
func (i *DistributionLogicImpl) BindRelation(ctx context.Context, userId, inviterId int64, channel int) error {
	if userId == inviterId {
		return errcode.New(errcode.CodeInvalidParam, "不能绑定自己为上级")
	}
	if channel < 1 || channel > 2 {
		channel = 1
	}
	n, err := dao.UserRelation.Ctx(ctx).Where(dao.UserRelation.Columns().UserId, userId).Count()
	if err != nil {
		return gerror.Wrap(err, "查询关系链失败")
	}
	if n > 0 {
		return errcode.New(errcode.CodeInvalidParam, "已绑定上级, 一人一链不可更改")
	}
	if _, err = dao.UserRelation.Ctx(ctx).Data(g.Map{
		"user_id": userId, "inviter_id": inviterId, "bind_channel": channel,
	}).Insert(); err != nil {
		if isDupKeyUser(err) { // 并发双绑: uk_user 兜底
			return errcode.New(errcode.CodeInvalidParam, "已绑定上级, 一人一链不可更改")
		}
		return gerror.Wrap(err, "绑定关系失败")
	}
	return nil
}

// InviteRecords 邀请激励记录（我邀请的新用户的激励流水）。
func (i *DistributionLogicImpl) InviteRecords(ctx context.Context, userId int64, page model.PageReq) (*model.PageResult[model.InviteRecordItem], error) {
	page = page.Normalized()
	cols := dao.InviteRecord.Columns()
	m := dao.InviteRecord.Ctx(ctx).Where(cols.InviterId, userId)
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计邀请激励失败")
	}
	recs, err := m.OrderDesc(cols.Id).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询邀请激励失败")
	}
	list := make([]model.InviteRecordItem, 0, len(recs))
	for _, r := range recs {
		nick := ""
		u, e := dao.User.Ctx(ctx).Fields(dao.User.Columns().Nickname).
			Where(dao.User.Columns().Id, r[cols.NewUserId].Int64()).One()
		if e == nil && !u.IsEmpty() {
			nick = maskNick(u["nickname"].String())
		}
		list = append(list, model.InviteRecordItem{
			NewUser: nick, RewardDesc: "邀请注册奖励", Status: r[cols.Status].Int(),
			CreatedAt: r[cols.CreatedAt].String(),
		})
	}
	return &model.PageResult[model.InviteRecordItem]{List: list, Total: int64(total)}, nil
}

// ---------- 查询面 ----------

// Records 我的佣金记录分页（status 0=全部）。
func (i *DistributionLogicImpl) Records(ctx context.Context, userId int64, status int, page model.PageReq) (*model.PageResult[DistRecordItem], error) {
	page = page.Normalized()
	m := g.DB().Model("commission_record").Ctx(ctx).
		Where("beneficiary_user_id", userId)
	if status > 0 {
		m = m.Where("status", status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计佣金记录失败")
	}
	recs, err := m.OrderDesc("id").Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询佣金记录失败")
	}
	list := make([]DistRecordItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, DistRecordItem{
			OrderNo: r["order_no"].String(), Level: r["level"].Int(),
			Amount: r["amount"].String(), Status: r["status"].Int(),
			SettleTime: r["settle_time"].String(), CreatedAt: r["created_at"].String(),
		})
	}
	return &model.PageResult[DistRecordItem]{List: list, Total: int64(total)}, nil
}

// RuleQuery 佣金比例查询（按商品: 商品覆盖 > 分类默认; 未命中空列表）。
func (i *DistributionLogicImpl) RuleQuery(ctx context.Context, spuId int64) ([]DistRuleHit, error) {
	spu, err := dao.ProductSpu.Ctx(ctx).
		Fields(dao.ProductSpu.Columns().Name, dao.ProductSpu.Columns().CategoryId).
		Where(dao.ProductSpu.Columns().Id, spuId).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询商品失败")
	}
	// 商品覆盖优先
	if !spu.IsEmpty() {
		if hit := distRuleHit(ctx, 2, spuId, spu["name"].String()); hit != nil {
			return []DistRuleHit{*hit}, nil
		}
		if hit := distRuleHit(ctx, 1, spu["category_id"].Int64(), "分类默认"); hit != nil {
			return []DistRuleHit{*hit}, nil
		}
	}
	return []DistRuleHit{}, nil
}

func distRuleHit(ctx context.Context, scopeType int, scopeId int64, desc string) *DistRuleHit {
	r, err := g.DB().Model("commission_rule").Ctx(ctx).
		Where("scope_type", scopeType).Where("scope_id", scopeId).
		Where("status", 1).Where("deleted", 0).One()
	if err != nil || r.IsEmpty() {
		return nil
	}
	return &DistRuleHit{
		ScopeDesc: desc, Level1Rate: r["level1_rate"].String(), Level2Rate: r["level2_rate"].String(),
	}
}

// ShareCode 我的推广码（稳定: user.share_code 为空则首访生成）。
func (i *DistributionLogicImpl) ShareCode(ctx context.Context, userId int64) (string, string, error) {
	cols := dao.User.Columns()
	v, err := dao.User.Ctx(ctx).Where(cols.Id, userId).Value(cols.ShareCode)
	if err != nil {
		return "", "", gerror.Wrap(err, "查询推广码失败")
	}
	code := v.String()
	if code == "" {
		next, e := idgen.NextID()
		if e != nil {
			return "", "", gerror.Wrap(e, "生成推广码失败")
		}
		// 列宽 varchar(16) UNIQUE: base36(sonyflake) 全局唯一且 ≤13 字符 + 前缀 = 14
		code = "S" + strconv.FormatInt(next, 36)
		if _, e = dao.User.Ctx(ctx).Where(cols.Id, userId).Data(g.Map{cols.ShareCode: code}).Update(); e != nil {
			return "", "", gerror.Wrap(e, "写入推广码失败")
		}
	}
	return code, "/register?code=" + code, nil
}

// ---------- US2 佣金（计提/结算/冲销——事件入口, 由 shop 域端口投递） ----------

// SettleOrder 订单佣金计提（确认收货事件; 订单项粒度幂等）。
// 归因: 分享窗口内归因人 > 关系链两级 > 不计佣。冻结/待审/软删推广员不计提。
func (i *DistributionLogicImpl) SettleOrder(ctx context.Context, orderNo string) error {
	items, err := g.DB().Model("trade_order_item").Ctx(ctx).
		Where("order_no", orderNo).All()
	if err != nil {
		return gerror.Wrap(err, "查询订单项失败")
	}
	order, err := g.DB().Model("trade_order").Ctx(ctx).Where("order_no", orderNo).One()
	if err != nil || order.IsEmpty() {
		return gerror.Wrap(err, "查询订单失败")
	}
	buyer := order["user_id"].Int64()
	for _, it := range items {
		base, e := money.FromYuanString(it["pay_amount"].String()) // 基数=行实付（不含运费/积分, 表口径）
		if e != nil {
			continue
		}
		if base <= 0 {
			continue
		}
		spuId := it["spu_id"].Int64()
		bene, level, rate := i.attribute(ctx, buyer, spuId) // (受益人, 层级, 比例%)
		if bene == 0 || rate <= 0 {
			continue
		}
		amountFen := base * rate / 100 // 四舍五入见下方修正
		if rem := (base * rate) % 100; rem >= 50 {
			amountFen++ // base*rate/100 的小数部分 ≥0.5 分则进位（整数域四舍五入到分）
		}
		if amountFen <= 0 {
			continue
		}
		// 计提幂等: (order_item_id, beneficiary, level) 业务键同事务查重（事件重放不重复落账）
		dup, e := g.DB().Model("commission_record").Ctx(ctx).
			Where("order_item_id", it["id"].Int64()).
			Where("beneficiary_user_id", bene).
			Where("level", level).Count()
		if e != nil {
			return gerror.Wrap(e, "查询佣金记录失败")
		}
		if dup > 0 {
			continue
		}
		if _, e = g.DB().Model("commission_record").Ctx(ctx).Data(g.Map{
			"order_no":            orderNo,
			"order_item_id":       it["id"].Int64(),
			"beneficiary_user_id": bene,
			"level":               level,
			"base_amount":         it["pay_amount"].String(),
			"rate":                rate,
			"amount":              money.ToYuanString(amountFen),
			"status":              1, // 待结算（保护期满由 ConfirmSettle 入账）
		}).Insert(); e != nil {
			return gerror.Wrap(e, "写入佣金记录失败")
		}
	}
	return nil
}

// attribute 归因: 分享窗口内最近分享人 > 关系链一级 > 二级（0/0/0=不计佣）。
// 冻结(3)/待审(1)/软删推广员不计提; 红线: 二级=上级的上级, 两次单列查询到顶即止。
func (i *DistributionLogicImpl) attribute(ctx context.Context, buyer, spuId int64) (int64, int, int64) {
	// ① 分享窗口归因（V1: share_record 中该买家对该 SPU 窗口内最近一次分享人）
	share, err := g.DB().Model("share_record").Ctx(ctx).
		Where("spu_id", spuId).Where("buyer_user_id", buyer).
		OrderDesc("id").Limit(1).One()
	if err == nil && !share.IsEmpty() {
		if uid := i.activeDistributor(ctx, share["user_id"].Int64()); uid > 0 {
			return uid, 1, i.rateOf(ctx, spuId, 1)
		}
	}
	// ② 关系链: 一级=直接上级; 二级=上级的上级（两次单列查询, 到顶即止——红线）
	my, err := dao.UserRelation.Ctx(ctx).Where(dao.UserRelation.Columns().UserId, buyer).One()
	if err != nil || my.IsEmpty() {
		return 0, 0, 0
	}
	l1 := my["inviter_id"].Int64()
	if uid := i.activeDistributor(ctx, l1); uid > 0 {
		return l1, 1, i.rateOf(ctx, spuId, 1)
	}
	grand, err := dao.UserRelation.Ctx(ctx).Where(dao.UserRelation.Columns().UserId, l1).One()
	if err != nil || grand.IsEmpty() {
		return 0, 0, 0
	}
	l2 := grand["inviter_id"].Int64()
	if uid := i.activeDistributor(ctx, l2); uid > 0 {
		return l2, 2, i.rateOf(ctx, spuId, 2)
	}
	return 0, 0, 0
}

// activeDistributor 用户是否"通过且未冻结未删"的推广员。
func (i *DistributionLogicImpl) activeDistributor(ctx context.Context, userId int64) int64 {
	if userId <= 0 {
		return 0
	}
	n, err := dao.DistributionUser.Ctx(ctx).
		Where(dao.DistributionUser.Columns().UserId, userId).
		Where(dao.DistributionUser.Columns().Status, 2).
		Where(dao.DistributionUser.Columns().Deleted, 0).Count()
	if err != nil || n == 0 {
		return 0
	}
	return userId
}

// rateOf 规则命中（商品覆盖 > 分类默认; 未命中 0）。
func (i *DistributionLogicImpl) rateOf(ctx context.Context, spuId int64, level int) int64 {
	spu, err := dao.ProductSpu.Ctx(ctx).
		Fields(dao.ProductSpu.Columns().CategoryId).
		Where(dao.ProductSpu.Columns().Id, spuId).One()
	if err != nil || spu.IsEmpty() {
		return 0
	}
	col := "level1_rate"
	if level == 2 {
		col = "level2_rate"
	}
	// 商品覆盖
	v, err := g.DB().Model("commission_rule").Ctx(ctx).
		Where("scope_type", 2).Where("scope_id", spuId).
		Where("status", 1).Where("deleted", 0).Value(col)
	if err == nil && v != nil && !v.IsNil() {
		return v.Int64()
	}
	// 分类默认
	v, err = g.DB().Model("commission_rule").Ctx(ctx).
		Where("scope_type", 1).Where("scope_id", spu["category_id"].Int64()).
		Where("status", 1).Where("deleted", 0).Value(col)
	if err == nil && v != nil && !v.IsNil() {
		return v.Int64()
	}
	return 0
}

// ConfirmSettle 结算保护期保护期满批量入账（D7; V1 常量保护期 7 天）。返回迁移条数。
func (i *DistributionLogicImpl) ConfirmSettle(ctx context.Context) (int64, error) {
	due, err := g.DB().Model("commission_record").Ctx(ctx).
		Where("status", 1).
		Where("created_at <= DATE_SUB(NOW(), INTERVAL 7 DAY)").All()
	if err != nil {
		return 0, gerror.Wrap(err, "查询到期佣金失败")
	}
	var settled int64
	for _, r := range due {
		e := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// 条件迁移（判行数——并发/重复任务只入账一次）
			res, e := tx.Model("commission_record").Ctx(ctx).
				Where("id", r["id"].Int64()).Where("status", 1).
				Data(g.Map{"status": 2, "settle_time": gtime.Now()}).Update()
			if e != nil {
				return gerror.Wrap(e, "推进佣金状态失败")
			}
			if n, _ := res.RowsAffected(); n == 0 {
				return nil // 已被并发任务处理
			}
			return creditAccount(ctx, tx, r["beneficiary_user_id"].Int64(),
				r["amount"].String(), 1, fmt.Sprintf("%d", r["id"].Int64()))
		})
		if e != nil {
			return settled, e
		}
		settled++
	}
	return settled, nil
}

// creditAccount 账户入账（幂等建户）+ 双快照流水（同事务; D1 铁律）。
func creditAccount(ctx context.Context, tx gdb.TX, userId int64, amount string, bizType int, bizNo string) error {
	if _, err := tx.Model("user_account").Ctx(ctx).
		Data(g.Map{"user_id": userId}).Insert(); err != nil {
		if !isDupKeyUser(err) {
			return gerror.Wrap(err, "建户失败")
		}
	}
	// 条件加余额（佣金入账恒正, 无下限需求; 行锁由 UPDATE 语义保证）
	res, err := tx.Model("user_account").Ctx(ctx).
		Where("user_id", userId).
		Data(g.Map{"balance": gdb.Raw("balance + " + amount)}).Update()
	if err != nil {
		return gerror.Wrap(err, "入账失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeInventoryAdjust, "账户不存在, 入账失败")
	}
	return writeAccountLog(ctx, tx, userId, amount, bizType, bizNo)
}

// writeAccountLog 双快照流水（读更新后的行; 同事务）。
func writeAccountLog(ctx context.Context, tx gdb.TX, userId int64, amount string, bizType int, bizNo string) error {
	acc, err := tx.Model("user_account").Ctx(ctx).
		Fields("balance", "frozen").Where("user_id", userId).One()
	if err != nil {
		return gerror.Wrap(err, "读取账户快照失败")
	}
	_, err = tx.Model("account_log").Ctx(ctx).Data(g.Map{
		"user_id": userId, "biz_type": bizType, "amount": amount,
		"balance_after": acc["balance"].String(), "frozen_after": acc["frozen"].String(),
		"biz_no": bizNo,
	}).Insert()
	return gerror.Wrap(err, "写入流水失败")
}

// ReverseOnRefund 售后退款冲销（FR-019 / ICommissionReverse 消费侧）:
// 未结算(1)置 3 失效; 已结算(2)生成负额冲销记录(4 欠款冲销中, reversal_of_id 回指) + 余额扣回(可负)。
func (i *DistributionLogicImpl) ReverseOnRefund(ctx context.Context, orderItemId int64) error {
	recs, err := g.DB().Model("commission_record").Ctx(ctx).
		Where("order_item_id", orderItemId).
		WhereIn("status", []int{1, 2}).All()
	if err != nil {
		return gerror.Wrap(err, "查询佣金记录失败")
	}
	for _, r := range recs {
		if r["status"].Int() == 1 {
			// 未结算: 条件置失效（判行数防重放）
			res, e := g.DB().Model("commission_record").Ctx(ctx).
				Where("id", r["id"].Int64()).Where("status", 1).
				Data(g.Map{"status": 3}).Update()
			if e != nil {
				return gerror.Wrap(e, "失效佣金失败")
			}
			if n, _ := res.RowsAffected(); n > 0 {
				continue
			}
		}
		// 已结算（或置失效输给了并发）: 负额冲销记录 + 扣回（幂等: 原记录已有回指则跳过）
		dup, e := g.DB().Model("commission_record").Ctx(ctx).
			Where("reversal_of_id", r["id"].Int64()).Count()
		if e != nil {
			return gerror.Wrap(e, "查询冲销记录失败")
		}
		if dup > 0 {
			continue
		}
		amount := r["amount"].String()
		e = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			res, e := tx.Model("commission_record").Ctx(ctx).Data(g.Map{
				"order_no":            r["order_no"].String(),
				"order_item_id":       r["order_item_id"].Int64(),
				"beneficiary_user_id": r["beneficiary_user_id"].Int64(),
				"level":               r["level"].Int(),
				"base_amount":         r["base_amount"].String(),
				"rate":                r["rate"].Int64(),
				"amount":              "-" + amount, // 负额冲销
				"status":              4,
				"reversal_of_id":      r["id"].Int64(),
			}).Insert()
			if e != nil {
				return gerror.Wrap(e, "写入冲销记录失败")
			}
			revId, e := res.LastInsertId()
			if e != nil {
				return gerror.Wrap(e, "读取冲销ID失败")
			}
			// 原记录回指 + 置冲销中（条件更新防重放）
			if _, e = tx.Model("commission_record").Ctx(ctx).
				Where("id", r["id"].Int64()).Where("status", 2).Where("reversal_record_id", nil).
				Data(g.Map{"status": 4, "reversal_record_id": revId}).Update(); e != nil {
				return gerror.Wrap(e, "回指原记录失败")
			}
			// 余额扣回（可负——user_account.balance 有符号, 列注释"欠款为负"）
			if _, e = tx.Model("user_account").Ctx(ctx).
				Where("user_id", r["beneficiary_user_id"].Int64()).
				Data(g.Map{"balance": gdb.Raw("balance - " + amount)}).Update(); e != nil {
				return gerror.Wrap(e, "扣回佣金失败")
			}
			return writeAccountLog(ctx, tx, r["beneficiary_user_id"].Int64(),
				"-"+amount, 5, fmt.Sprintf("%d", revId))
		})
		if e != nil {
			return e
		}
	}
	return nil
}

// ---------- 账户与提现 ----------

// Account 佣金账户（无账户=0/0）。
func (i *DistributionLogicImpl) Account(ctx context.Context, userId int64) (*model.DistAccount, error) {
	acc, err := g.DB().Model("user_account").Ctx(ctx).
		Where("user_id", userId).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询账户失败")
	}
	if acc.IsEmpty() {
		return &model.DistAccount{Balance: "0.00", Frozen: "0.00"}, nil
	}
	return &model.DistAccount{Balance: acc["balance"].String(), Frozen: acc["frozen"].String()}, nil
}

// AccountLogs 账户流水分页（bizType 0=全部）。
func (i *DistributionLogicImpl) AccountLogs(ctx context.Context, userId int64, bizType int, page model.PageReq) (*model.PageResult[model.AccountLogItem], error) {
	page = page.Normalized()
	m := g.DB().Model("account_log").Ctx(ctx).Where("user_id", userId)
	if bizType > 0 {
		m = m.Where("biz_type", bizType)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计流水失败")
	}
	recs, err := m.OrderDesc("id").Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询流水失败")
	}
	list := make([]model.AccountLogItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, model.AccountLogItem{
			BizType: r["biz_type"].Int(), Amount: r["amount"].String(),
			BalanceAfter: r["balance_after"].String(), BizNo: r["biz_no"].String(),
			CreatedAt: r["created_at"].String(),
		})
	}
	return &model.PageResult[model.AccountLogItem]{List: list, Total: int64(total)}, nil
}

// nextWithdrawNo 提现单号（sonyflake, 与售后单号同风格）。
func nextWithdrawNo() (string, error) {
	next, err := idgen.NextID()
	if err != nil {
		return "", gerror.Wrap(err, "生成提现单号失败")
	}
	return "WD" + fmt.Sprint(next), nil
}

// WithdrawApply 提现申请（D1: 余额条件冻结判行数——并发双申请只冻结得起的通过）。
func (i *DistributionLogicImpl) WithdrawApply(ctx context.Context, userId int64, amount string) (string, error) {
	fen, err := money.FromYuanString(amount)
	if err != nil || fen <= 0 {
		return "", errcode.New(errcode.CodeInvalidParam, "提现金额非法")
	}
	if _, err = g.DB().Model("user_account").Ctx(ctx).
		Data(g.Map{"user_id": userId}).Insert(); err != nil {
		if !isDupKeyUser(err) {
			return "", gerror.Wrap(err, "建户失败")
		}
	}
	no, err := nextWithdrawNo()
	if err != nil {
		return "", err
	}
	// 提现单先行（10 待审核; 冻结失败则条件作废——简化: 先冻结后建单, 失败无单可留）
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 条件冻结: balance >= amount 才通过（并发双申请只成一笔——SC-3）
		res, e := tx.Model("user_account").Ctx(ctx).
			Where("user_id", userId).
			Where("balance >= ?", amount).
			Data(g.Map{
				"balance": gdb.Raw("balance - " + amount),
				"frozen":  gdb.Raw("frozen + " + amount),
			}).Update()
		if e != nil {
			return gerror.Wrap(e, "冻结余额失败")
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return errcode.New(errcode.CodeInvalidParam, "可用余额不足")
		}
		if _, e = tx.Model("withdraw_order").Ctx(ctx).Data(g.Map{
			"withdraw_no": no, "user_id": userId, "amount": amount, "status": 10,
		}).Insert(); e != nil {
			return gerror.Wrap(e, "创建提现单失败")
		}
		return writeAccountLog(ctx, tx, userId, "-"+amount, 2, no)
	})
	if err != nil {
		return "", err
	}
	return no, nil
}

// WithdrawList 我的提现列表。
func (i *DistributionLogicImpl) WithdrawList(ctx context.Context, userId int64, status int, page model.PageReq) (*model.PageResult[model.WithdrawItem], error) {
	page = page.Normalized()
	m := g.DB().Model("withdraw_order").Ctx(ctx).Where("user_id", userId)
	if status > 0 {
		m = m.Where("status", status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计提现单失败")
	}
	recs, err := m.OrderDesc("id").Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询提现单失败")
	}
	list := make([]model.WithdrawItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, model.WithdrawItem{
			WithdrawNo: r["withdraw_no"].String(), Amount: r["amount"].String(),
			Status: r["status"].Int(), CreatedAt: r["created_at"].String(),
		})
	}
	return &model.PageResult[model.WithdrawItem]{List: list, Total: int64(total)}, nil
}


// ShareReport 分享行为上报（归因窗口起点; 游客可报 user_id 可空; 只追加 share_record）。
func (i *DistributionLogicImpl) ShareReport(ctx context.Context, userId int64, spuId int64, channel int, scene string) error {
	data := g.Map{"spu_id": spuId, "channel": channel, "scene": scene}
	if userId > 0 {
		data["user_id"] = userId
	}
	if _, err := g.DB().Model("share_record").Ctx(ctx).Data(data).Insert(); err != nil {
		return gerror.Wrap(err, "记录分享失败")
	}
	return nil
}

// ---------- admin 侧实现 ----------

// AdminDistributorList 推广员列表（状态/关键词筛选; 昵称脱敏; Level 恒 0 占位）。
func (i *DistributionAdminLogicImpl) AdminDistributorList(ctx context.Context, status int, keyword string, page model.PageReq) (*model.PageResult[AdminDistributorItem], error) {
	page = page.Normalized()
	m := g.DB().Model("distribution_user").Ctx(ctx).Where("deleted", 0)
	if status > 0 {
		m = m.Where("status", status)
	}
	if keyword != "" {
		if idv, e := g.DB().Model("user").Ctx(ctx).Where("id", keyword).Value("id"); e == nil && idv != nil {
			m = m.Where("user_id", idv.Int64())
		} else {
			ids, e := g.DB().Model("user").Ctx(ctx).
				WhereLike("nickname", "%"+keyword+"%").Fields("id").Limit(200).All()
			if e != nil {
				return nil, gerror.Wrap(e, "查询用户失败")
			}
			idsList := []int64{}
			for _, r := range ids {
				idsList = append(idsList, r["id"].Int64())
			}
			if len(idsList) == 0 {
				return &model.PageResult[AdminDistributorItem]{List: []AdminDistributorItem{}}, nil
			}
			m = m.WhereIn("user_id", idsList)
		}
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计推广员失败")
	}
	recs, err := m.OrderDesc("id").Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询推广员失败")
	}
	list := make([]AdminDistributorItem, 0, len(recs))
	for _, r := range recs {
		nick := ""
		if u, e := g.DB().Model("user").Ctx(ctx).Fields("nickname").
			Where("id", r["user_id"].Int64()).One(); e == nil && !u.IsEmpty() {
			nick = maskNick(u["nickname"].String())
		}
		at := ""
		if !r["audit_time"].IsNil() {
			at = r["audit_time"].String()
		}
		list = append(list, AdminDistributorItem{
			Id: r["id"].Int64(), UserId: r["user_id"].Int64(), Nickname: nick,
			Status: r["status"].Int(), ApplyTime: r["apply_time"].String(), AuditTime: at,
		})
	}
	return &model.PageResult[AdminDistributorItem]{List: list, Total: int64(total)}, nil
}

// AdminDistributorAudit 推广员审核（1→2 通过; 拒绝置软删（终态, 语义=不通过）; 条件更新判行数, 重放拒绝）。
func (i *DistributionAdminLogicImpl) AdminDistributorAudit(ctx context.Context, id int64, pass bool) error {
	if pass {
		res, err := g.DB().Model("distribution_user").Ctx(ctx).
			Where("id", id).Where("status", 1).
			Data(g.Map{"status": 2, "audit_time": gtime.Now()}).Update()
		if err != nil {
			return gerror.Wrap(err, "审核通过失败")
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可审核")
		}
		return nil
	}
	res, err := g.DB().Model("distribution_user").Ctx(ctx).
		Where("id", id).Where("status", 1).
		Data(g.Map{"deleted": 1}).Update()
	if err != nil {
		return gerror.Wrap(err, "审核拒绝失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可审核")
	}
	return nil
}

// AdminDistributorFreeze 冻结(2→3)/解冻(3→2)——条件更新, 非法迁移拒绝。
func (i *DistributionAdminLogicImpl) AdminDistributorFreeze(ctx context.Context, id int64, freeze bool) error {
	from, to := 2, 3
	if !freeze {
		from, to = 3, 2
	}
	res, err := g.DB().Model("distribution_user").Ctx(ctx).
		Where("id", id).Where("status", from).Where("deleted", 0).
		Data(g.Map{"status": to}).Update()
	if err != nil {
		return gerror.Wrap(err, "变更推广员状态失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可冻结/解冻")
	}
	return nil
}

// AdminRuleList 佣金规则列表（作用域名称装配: 商品名/分类名）。
func (i *DistributionAdminLogicImpl) AdminRuleList(ctx context.Context, page model.PageReq) (*model.PageResult[AdminDistRuleItem], error) {
	page = page.Normalized()
	m := g.DB().Model("commission_rule").Ctx(ctx).Where("deleted", 0)
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计规则失败")
	}
	recs, err := m.OrderDesc("id").Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询规则失败")
	}
	out := make([]AdminDistRuleItem, 0, len(recs))
	for _, r := range recs {
		desc := ""
		if r["scope_type"].Int() == 2 {
			v, e := g.DB().Model("product_spu").Ctx(ctx).Fields("name").Where("id", r["scope_id"].Int64()).One()
			if e == nil && !v.IsEmpty() {
				desc = v["name"].String()
			}
		} else {
			v, e := g.DB().Model("product_category").Ctx(ctx).Fields("name").Where("id", r["scope_id"].Int64()).One()
			if e == nil && !v.IsEmpty() {
				desc = v["name"].String()
			}
		}
		out = append(out, AdminDistRuleItem{
			Id: r["id"].Int64(), ScopeType: r["scope_type"].Int(), ScopeId: r["scope_id"].Int64(),
			ScopeDesc: desc, Level1Rate: r["level1_rate"].String(), Level2Rate: r["level2_rate"].String(),
			Status: r["status"].Int(),
		})
	}
	return &model.PageResult[AdminDistRuleItem]{List: out, Total: int64(total)}, nil
}

// distRateCheck 比例校验（0-100, ≤2 位小数, 复用两位小数解析量纲: 100.00% = 10000）。
func distRateCheck(v string) (string, error) {
	if v == "" {
		return "0.00", nil
	}
	f, err := money.FromYuanString(v)
	if err != nil {
		return "", errcode.New(errcode.CodeInvalidParam, "比例格式非法")
	}
	if f < 0 || f > 10000 {
		return "", errcode.New(errcode.CodeInvalidParam, "比例须在 0~100 之间")
	}
	return v, nil
}

// AdminRuleCreate 创建规则（uk_scope 1062 → 业务码）。
func (i *DistributionAdminLogicImpl) AdminRuleCreate(ctx context.Context, scopeType int, scopeId int64, l1, l2 string) (int64, error) {
	if _, err := distRateCheck(l1); err != nil {
		return 0, err
	}
	if _, err := distRateCheck(l2); err != nil {
		return 0, err
	}
	res, err := g.DB().Model("commission_rule").Ctx(ctx).Data(g.Map{
		"scope_type": scopeType, "scope_id": scopeId,
		"level1_rate": l1, "level2_rate": l2, "status": 1,
	}).Insert()
	if err != nil {
		if isDupKeyUser(err) {
			return 0, errcode.New(errcode.CodeInvalidParam, "该作用域已有规则")
		}
		return 0, gerror.Wrap(err, "创建规则失败")
	}
	id, err := res.LastInsertId()
	return id, gerror.Wrap(err, "读取规则ID失败")
}

// AdminRuleUpdate 修改规则（status *int 三态; 同值幂等——批次 10 I3/I6 教训内化）。
func (i *DistributionAdminLogicImpl) AdminRuleUpdate(ctx context.Context, id int64, l1, l2 string, status *int) error {
	n, err := g.DB().Model("commission_rule").Ctx(ctx).Where("id", id).Where("deleted", 0).Count()
	if err != nil {
		return gerror.Wrap(err, "查询规则失败")
	}
	if n == 0 {
		return errcode.New(errcode.CodeNotFound, "规则不存在")
	}
	data := g.Map{}
	if l1 != "" {
		if _, err = distRateCheck(l1); err != nil {
			return err
		}
		data["level1_rate"] = l1
	}
	if l2 != "" {
		if _, err = distRateCheck(l2); err != nil {
			return err
		}
		data["level2_rate"] = l2
	}
	if status != nil {
		data["status"] = *status
	}
	if len(data) == 0 {
		return nil
	}
	_, err = g.DB().Model("commission_rule").Ctx(ctx).Where("id", id).Data(data).Update()
	return gerror.Wrap(err, "修改规则失败")
}

// AdminRuleDelete 软删。
func (i *DistributionAdminLogicImpl) AdminRuleDelete(ctx context.Context, id int64) error {
	res, err := g.DB().Model("commission_rule").Ctx(ctx).
		Where("id", id).Where("deleted", 0).Data(g.Map{"deleted": 1}).Update()
	if err != nil {
		return gerror.Wrap(err, "删除规则失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeNotFound, "规则不存在")
	}
	return nil
}

// AdminRecordList 全局佣金记录（状态/订单号筛选; 受益人脱敏）。
func (i *DistributionAdminLogicImpl) AdminRecordList(ctx context.Context, status int, orderNo string, page model.PageReq) (*model.PageResult[AdminDistRecordItem], error) {
	page = page.Normalized()
	m := g.DB().Model("commission_record").Ctx(ctx)
	if status > 0 {
		m = m.Where("status", status)
	}
	if orderNo != "" {
		m = m.Where("order_no", orderNo)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计佣金记录失败")
	}
	recs, err := m.OrderDesc("id").Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询佣金记录失败")
	}
	out := make([]AdminDistRecordItem, 0, len(recs))
	for _, r := range recs {
		nick := ""
		if u, e := g.DB().Model("user").Ctx(ctx).Fields("nickname").
			Where("id", r["beneficiary_user_id"].Int64()).One(); e == nil && !u.IsEmpty() {
			nick = maskNick(u["nickname"].String())
		}
		st := ""
		if !r["settle_time"].IsNil() {
			st = r["settle_time"].String()
		}
		out = append(out, AdminDistRecordItem{
			OrderNo: r["order_no"].String(), Beneficiary: nick, Level: r["level"].Int(),
			BaseAmount: r["base_amount"].String(), Amount: r["amount"].String(),
			Status: r["status"].Int(), SettleTime: st,
		})
	}
	return &model.PageResult[AdminDistRecordItem]{List: out, Total: int64(total)}, nil
}

// AdminWithdrawList 提现列表。
func (i *DistributionAdminLogicImpl) AdminWithdrawList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.WithdrawItem], error) {
	page = page.Normalized()
	m := g.DB().Model("withdraw_order").Ctx(ctx)
	if status > 0 {
		m = m.Where("status", status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计提现单失败")
	}
	recs, err := m.OrderDesc("id").Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询提现单失败")
	}
	list := make([]model.WithdrawItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, model.WithdrawItem{
			WithdrawNo: r["withdraw_no"].String(), Amount: r["amount"].String(),
			Status: r["status"].Int(), CreatedAt: r["created_at"].String(),
		})
	}
	return &model.PageResult[model.WithdrawItem]{List: list, Total: int64(total)}, nil
}

// withdrawByNo 按单号取提现单（不存在 → 40005）。
func withdrawByNo(ctx context.Context, withdrawNo string) (gdb.Record, error) {
	r, err := g.DB().Model("withdraw_order").Ctx(ctx).Where("withdraw_no", withdrawNo).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询提现单失败")
	}
	if r.IsEmpty() {
		return nil, errcode.New(errcode.CodeOrderNotFound, "提现单不存在")
	}
	return r, nil
}

// unfreeze 冻结回退余额（拒绝/打款失败共用）+ 流水 biz_type=4。
func unfreeze(ctx context.Context, tx gdb.TX, userId int64, amount, withdrawNo string) error {
	res, err := tx.Model("user_account").Ctx(ctx).
		Where("user_id", userId).Where("frozen >= ?", amount).
		Data(g.Map{
			"frozen":  gdb.Raw("frozen - " + amount),
			"balance": gdb.Raw("balance + " + amount),
		}).Update()
	if err != nil {
		return gerror.Wrap(err, "回退冻结失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeInventoryAdjust, "冻结不足, 回退失败")
	}
	return writeAccountLog(ctx, tx, userId, "+"+amount, 4, withdrawNo)
}

// AdminWithdrawAudit 提现审核（10→20 通过 / 10→50 拒绝回退; 条件状态机不可逆）。
func (i *DistributionAdminLogicImpl) AdminWithdrawAudit(ctx context.Context, withdrawNo string, pass bool, reason string) error {
	r, err := withdrawByNo(ctx, withdrawNo)
	if err != nil {
		return err
	}
	if !pass {
		return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			res, e := tx.Model("withdraw_order").Ctx(ctx).
				Where("withdraw_no", withdrawNo).Where("status", 10).
				Data(g.Map{"status": 50, "fail_reason": reason, "audit_time": gtime.Now()}).Update()
			if e != nil {
				return gerror.Wrap(e, "审核拒绝失败")
			}
			if n, _ := res.RowsAffected(); n == 0 {
				return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可审核")
			}
			return unfreeze(ctx, tx, r["user_id"].Int64(), r["amount"].String(), withdrawNo)
		})
	}
	res, err := g.DB().Model("withdraw_order").Ctx(ctx).
		Where("withdraw_no", withdrawNo).Where("status", 10).
		Data(g.Map{"status": 20, "audit_time": gtime.Now()}).Update()
	if err != nil {
		return gerror.Wrap(err, "审核通过失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可审核")
	}
	return nil
}

// AdminWithdrawPay 打款登记（20→40 成功核销 / 20→60 失败回退; uk_channel_order 幂等防重复打款）。
func (i *DistributionAdminLogicImpl) AdminWithdrawPay(ctx context.Context, withdrawNo string, success bool, channelOrderNo, failReason string) error {
	r, err := withdrawByNo(ctx, withdrawNo)
	if err != nil {
		return err
	}
	if success {
		if channelOrderNo == "" {
			return errcode.New(errcode.CodeInvalidParam, "打款成功必须登记渠道单号")
		}
		return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			res, e := tx.Model("withdraw_order").Ctx(ctx).
				Where("withdraw_no", withdrawNo).Where("status", 20).
				Data(g.Map{"status": 40, "channel_order_no": channelOrderNo, "pay_time": gtime.Now()}).Update()
			if e != nil {
				if isDupKeyUser(e) { // 渠道单号唯一 → 疑似重复打款（红线幂等）
					return errcode.New(errcode.CodeStatusNotAllowed, "渠道单号已存在, 疑似重复打款")
				}
				return gerror.Wrap(e, "打款核销失败")
			}
			if n, _ := res.RowsAffected(); n == 0 {
				return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可打款")
			}
			if _, e = tx.Model("user_account").Ctx(ctx).
				Where("user_id", r["user_id"].Int64()).Where("frozen >= ?", r["amount"].String()).
				Data(g.Map{"frozen": gdb.Raw("frozen - " + r["amount"].String())}).Update(); e != nil {
				return gerror.Wrap(e, "核销冻结失败")
			}
			return writeAccountLog(ctx, tx, r["user_id"].Int64(), "-"+r["amount"].String(), 3, withdrawNo)
		})
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, e := tx.Model("withdraw_order").Ctx(ctx).
			Where("withdraw_no", withdrawNo).Where("status", 20).
			Data(g.Map{"status": 60, "fail_reason": failReason}).Update()
		if e != nil {
			return gerror.Wrap(e, "登记打款失败")
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可打款")
		}
		return unfreeze(ctx, tx, r["user_id"].Int64(), r["amount"].String(), withdrawNo)
	})
}

// AdminInviteRecords 全局邀请激励记录（昵称脱敏; trigger 恒 1=注册时机, 首单时机 V1 未启用记账）。
func (i *DistributionAdminLogicImpl) AdminInviteRecords(ctx context.Context, page model.PageReq) (*model.PageResult[AdminInviteRecordItem], error) {
	page = page.Normalized()
	m := g.DB().Model("invite_record").Ctx(ctx)
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计邀请激励失败")
	}
	recs, err := m.OrderDesc("id").Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询邀请激励失败")
	}
	out := make([]AdminInviteRecordItem, 0, len(recs))
	for _, r := range recs {
		nick := func(uid int64) string {
			u, e := g.DB().Model("user").Ctx(ctx).Fields("nickname").Where("id", uid).One()
			if e != nil || u.IsEmpty() {
				return ""
			}
			return maskNick(u["nickname"].String())
		}
		out = append(out, AdminInviteRecordItem{
			Inviter: nick(r["inviter_id"].Int64()), NewUser: nick(r["new_user_id"].Int64()),
			RewardDesc: "邀请注册奖励", Trigger: 1, CreatedAt: r["created_at"].String(),
		})
	}
	return &model.PageResult[AdminInviteRecordItem]{List: out, Total: int64(total)}, nil
}
