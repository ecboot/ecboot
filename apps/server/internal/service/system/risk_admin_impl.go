// risk_admin_impl.go IRiskAdminLogic 实现 + 风控评估器（018 批次 12）。
// 防线: 申诉条件更新（已结论拒绝重放）; 规则状态/软删条件更新; 评估失败不阻断主流程。
package system

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// RiskAdminLogicImpl IRiskAdminLogic 实现。
type RiskAdminLogicImpl struct{}

func NewRiskAdminLogic() *RiskAdminLogicImpl { return &RiskAdminLogicImpl{} }

var _ IRiskAdminLogic = (*RiskAdminLogicImpl)(nil)

func riskRuleOf(r gdb.Record) AdminRiskRuleItem {
	return AdminRiskRuleItem{
		Id: r["id"].Int64(), Name: r["name"].String(),
		RuleType: r["rule_type"].Int(), ConditionExpr: r["condition_expr"].String(),
		Action: r["action"].Int(), Status: r["status"].Int(),
	}
}

// AdminRuleList 规则列表。
func (i *RiskAdminLogicImpl) AdminRuleList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[AdminRiskRuleItem], error) {
	page = page.Normalized()
	m := g.DB().Model("risk_rule").Ctx(ctx).Where("deleted", 0)
	if status > 0 {
		m = m.Where("status", status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计规则失败")
	}
	recs, err := m.OrderDesc("id").Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询规则失败")
	}
	list := make([]AdminRiskRuleItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, riskRuleOf(r))
	}
	return &model.PageResult[AdminRiskRuleItem]{List: list, Total: int64(total)}, nil
}

func riskRuleCheck(ruleType, action int, conditionExpr string) error {
	if ruleType < 1 || ruleType > 5 {
		return errcode.New(errcode.CodeInvalidParam, "规则类型非法")
	}
	if action != 1 && action != 2 {
		return errcode.New(errcode.CodeInvalidParam, "处置动作非法")
	}
	if strings.TrimSpace(conditionExpr) == "" {
		return errcode.New(errcode.CodeInvalidParam, "规则条件不能为空")
	}
	return nil
}

// AdminRuleCreate 新增规则。
func (i *RiskAdminLogicImpl) AdminRuleCreate(ctx context.Context, name string, ruleType int, conditionExpr string, action int) (int64, error) {
	if err := riskRuleCheck(ruleType, action, conditionExpr); err != nil {
		return 0, err
	}
	res, err := g.DB().Model("risk_rule").Ctx(ctx).Data(g.Map{
		"name": name, "rule_type": ruleType,
		"condition_expr": conditionExpr, "action": action, "status": 1,
	}).Insert()
	if err != nil {
		return 0, gerror.Wrap(err, "创建规则失败")
	}
	id, err := res.LastInsertId()
	return id, gerror.Wrap(err, "读取规则ID失败")
}

// AdminRuleUpdate 修改规则（status 三态; 同值幂等）。
func (i *RiskAdminLogicImpl) AdminRuleUpdate(ctx context.Context, id int64, name, conditionExpr string, action int, status *int) error {
	n, err := g.DB().Model("risk_rule").Ctx(ctx).Where("id", id).Where("deleted", 0).Count()
	if err != nil {
		return gerror.Wrap(err, "查询规则失败")
	}
	if n == 0 {
		return errcode.New(errcode.CodeNotFound, "规则不存在")
	}
	data := g.Map{}
	if name != "" {
		data["name"] = name
	}
	if conditionExpr != "" {
		data["condition_expr"] = conditionExpr
	}
	if action != 0 {
		if action != 1 && action != 2 {
			return errcode.New(errcode.CodeInvalidParam, "处置动作非法")
		}
		data["action"] = action
	}
	if status != nil {
		data["status"] = *status
	}
	if len(data) == 0 {
		return nil
	}
	_, err = g.DB().Model("risk_rule").Ctx(ctx).Where("id", id).Where("deleted", 0).Data(data).Update()
	return gerror.Wrap(err, "修改规则失败")
}

// AdminRuleDelete 软删。
func (i *RiskAdminLogicImpl) AdminRuleDelete(ctx context.Context, id int64) error {
	res, err := g.DB().Model("risk_rule").Ctx(ctx).
		Where("id", id).Where("deleted", 0).Data(g.Map{"deleted": 1, "status": 0}).Update()
	if err != nil {
		return gerror.Wrap(err, "删除规则失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeNotFound, "规则不存在")
	}
	return nil
}

// AdminRuleDetail 规则详情。
func (i *RiskAdminLogicImpl) AdminRuleDetail(ctx context.Context, id int64) (*AdminRiskRuleItem, error) {
	r, err := g.DB().Model("risk_rule").Ctx(ctx).Where("id", id).Where("deleted", 0).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询规则失败")
	}
	if r.IsEmpty() {
		return nil, errcode.New(errcode.CodeNotFound, "规则不存在")
	}
	out := riskRuleOf(r)
	return &out, nil
}

// AdminRecordList 风控事件列表（用户/申诉状态筛选）。
func (i *RiskAdminLogicImpl) AdminRecordList(ctx context.Context, userId int64, appealStatus int, page model.PageReq) (*model.PageResult[AdminRiskRecordItem], error) {
	page = page.Normalized()
	m := g.DB().Model("risk_record").Ctx(ctx)
	if userId > 0 {
		m = m.Where("user_id", userId)
	}
	if appealStatus > 0 {
		m = m.Where("appeal_status", appealStatus)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计风控事件失败")
	}
	recs, err := m.OrderDesc("id").Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询风控事件失败")
	}
	out := make([]AdminRiskRecordItem, 0, len(recs))
	for _, r := range recs {
		ruleName := ""
		if rv, e := g.DB().Model("risk_rule").Ctx(ctx).Fields("name").
			Where("id", r["rule_id"].Int64()).One(); e == nil && !rv.IsEmpty() {
			ruleName = rv["name"].String()
		}
		out = append(out, AdminRiskRecordItem{
			Id: r["id"].Int64(), UserId: r["user_id"].Int64(), RuleName: ruleName,
			ObjectType: r["object_type"].Int(), ObjectNo: r["object_no"].String(),
			Action: r["action"].Int(), AppealStatus: r["appeal_status"].Int(),
			Remark: r["remark"].String(), CreatedAt: r["created_at"].String(),
		})
	}
	return &model.PageResult[AdminRiskRecordItem]{List: out, Total: int64(total)}, nil
}

// AdminAppeal 申诉处理（0→2通过 / 0→3驳回; 已有结论拒绝——条件更新）。
func (i *RiskAdminLogicImpl) AdminAppeal(ctx context.Context, id int64, pass bool, remark string) error {
	to := 3 // 驳回
	if pass {
		to = 2 // 通过
	}
	conclusion := "申诉驳回: " + remark
	if pass {
		conclusion = "申诉通过: " + remark
	}
	// 存在性前置（批次 10 I3 教训: affected=0 区分"不存在"与"已结论"——两码不同）
	n, err := g.DB().Model("risk_record").Ctx(ctx).Where("id", id).Count()
	if err != nil {
		return gerror.Wrap(err, "查询风控事件失败")
	}
	if n == 0 {
		return errcode.New(errcode.CodeNotFound, "风控事件不存在")
	}
	res, err := g.DB().Model("risk_record").Ctx(ctx).
		Where("id", id).Where("appeal_status", 0). // 无结论才可处理
		Data(g.Map{"appeal_status": to, "remark": conclusion}).Update()
	if err != nil {
		return gerror.Wrap(err, "处理申诉失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeStatusNotAllowed, "该事件已处理, 不可重放")
	}
	return nil
}

// ---------- 风控评估器（shop.IRiskHit 端口实现, bootstrap 装配） ----------

// RiskHitImpl 风控命中评估器（批次 09/11 预留端口的实现）。
// 评估失败**不阻断主流程**（降级放行 + 告警——批次 09 语义延续, 但装配后为真实判定）。
type RiskHitImpl struct{}

func NewRiskHitImpl() *RiskHitImpl { return &RiskHitImpl{} }

var _ RiskEvaluator = RiskHitImpl{}

// RiskEvaluator 风控评估接口（内部形态, 供 bootstrap 适配为 shop.IRiskHit）。
type RiskEvaluator interface {
	Hit(ctx context.Context, userId int64, ruleType int, payload string) (blocked bool, err error)
}

// Hit 评估: 启用的 rule_type 规则逐条判定。
//   - 黑名单(1): condition_expr 含用户标识（"user:<id>"）→ 命中;
//   - 其余类型（2~5 实时统计类）: V1 记录面完整、实时统计判定挂账增强 → 不命中。
//
// 命中: 按规则 action 落 risk_record（1拦截/2标记）并返回 blocked; 重复命中幂等（同用户同规则
// 已有未申诉拦截记录 → 不重复落, 但仍返回拦截）。
func (impl RiskHitImpl) Hit(ctx context.Context, userId int64, ruleType int, payload string) (bool, error) {
	rules, err := g.DB().Model("risk_rule").Ctx(ctx).
		Where("rule_type", ruleType).Where("status", 1).Where("deleted", 0).All()
	if err != nil {
		g.Log().Warningf(ctx, "[风控] 规则查询失败, 降级放行: user_id=%d err=%v", userId, err)
		return false, nil // 评估失败不阻断主流程（D2）
	}
	tag := fmt.Sprintf("user:%d", userId)
	for _, r := range rules {
		if r["rule_type"].Int() == 1 && // 黑名单: 精确标识命中
			strings.Contains(r["condition_expr"].String(), tag) {
			return impl.record(ctx, userId, r, payload)
		}
		// 类型 2~5: 实时统计判定挂账增强（V1 不命中——记录面完整, 判定后续增强）
	}
	return false, nil
}

// record 落风控事件（幂等: 同用户同规则已有未申诉拦截记录不重复落）。
func (impl RiskHitImpl) record(ctx context.Context, userId int64, rule gdb.Record, payload string) (bool, error) {
	dup, e := g.DB().Model("risk_record").Ctx(ctx).
		Where("user_id", userId).Where("rule_id", rule["id"].Int64()).
		Where("appeal_status", 0).Count()
	if e != nil {
		g.Log().Warningf(ctx, "[风控] 事件查重失败, 降级放行: user_id=%d err=%v", userId, e)
		return false, nil
	}
	blocked := rule["action"].Int() == 1
	if dup > 0 {
		return blocked, nil // 已有未处理事件, 不重复落
	}
	_, e = g.DB().Model("risk_record").Ctx(ctx).Data(g.Map{
		"user_id": userId, "rule_id": rule["id"].Int64(),
		"object_type": 1, // V1: 玩法/入口级命中, 关联对象类型默认订单域
		"object_no":   payload,
		"action":      rule["action"].Int(),
		"remark":      rule["condition_expr"].String(),
	}).Insert()
	if e != nil {
		// 落库失败仅告警, 不阻断主流程（plan D2）
		g.Log().Errorf(ctx, "[风控] 事件落库失败(不阻断): user_id=%d rule_id=%d err=%v", userId, rule["id"].Int64(), e)
	}
	return blocked, nil
}
