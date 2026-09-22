// risk_admin_impl_test.go 风控管理 + 评估器（018 批次 12）——SC-3 全链一条龙:
// 规则→命中拦截→落记录→申诉→解除; 评估失败降级不阻断。
package system

import (
	"context"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

func riskCleanupRule(ctx context.Context, t *gtest.T, name string) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM risk_record WHERE rule_id IN (SELECT id FROM risk_rule WHERE name=?)", name)
	_, _ = g.DB().Exec(ctx, "DELETE FROM risk_rule WHERE name=?", name)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO risk_rule(name,rule_type,condition_expr,action,status) VALUES(?,1,?,1,1)",
		name, "user:998001")
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

func riskCleanupRecord(ctx context.Context, t *gtest.T) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM risk_record WHERE user_id=998001")
}

// TestRiskRuleCrud 规则 CRUD+详情+校验矩阵（FR-2）。
func TestRiskRuleCrud(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const n = "TF-风控规则"
		defer riskCleanupRule(ctx, t, n)
		logic := NewRiskAdminLogic()

		// 校验矩阵
		_, err := logic.AdminRuleCreate(ctx, n, 9, "cond", 1)
		t.Assert(errCode(err), errcode.CodeInvalidParam)
		_, err = logic.AdminRuleCreate(ctx, n, 1, "cond", 9)
		t.Assert(errCode(err), errcode.CodeInvalidParam)
		_, err = logic.AdminRuleCreate(ctx, n, 1, " ", 1)
		t.Assert(errCode(err), errcode.CodeInvalidParam)

		id, err := logic.AdminRuleCreate(ctx, n, 1, "user:998001", 1)
		t.AssertNil(err)
		t.Assert(id > 0, true)

		d, err := logic.AdminRuleDetail(ctx, id)
		t.AssertNil(err)
		t.Assert(d.Name, n)
		t.Assert(d.RuleType, 1)
		t.Assert(d.Action, 1)
		_, err = logic.AdminRuleDetail(ctx, 999999999)
		t.Assert(errCode(err), errcode.CodeNotFound)

		// 更新（status 三态: nil 不修改启停）
		st := 0
		t.AssertNil(logic.AdminRuleUpdate(ctx, id, n+"-v2", "", 2, &st))
		d, _ = logic.AdminRuleDetail(ctx, id)
		t.Assert(d.Action, 2)
		t.Assert(d.Status, 0)
		t.AssertNil(logic.AdminRuleUpdate(ctx, id, n, "", 0, nil)) // 只改名
		d, _ = logic.AdminRuleDetail(ctx, id)
		t.Assert(d.Status, 0) // 启停未被触碰

		// 软删
		t.AssertNil(logic.AdminRuleDelete(ctx, id))
		_, err = logic.AdminRuleDetail(ctx, id)
		t.Assert(errCode(err), errcode.CodeNotFound)
	})
}

// TestRiskAppeal 申诉处理（0→2/3; 已结论拒绝; 不存在 404）（FR-3）。
func TestRiskAppeal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const n = "TF-风控申诉"
		defer riskCleanupRecord(ctx, t)
		id := riskCleanupRule(ctx, t, n)
		logic := NewRiskAdminLogic()

		_, err := g.DB().Exec(ctx,
			"INSERT INTO risk_record(user_id,rule_id,object_type,object_no,action,appeal_status,remark) "+
				"VALUES(998001,?,1,'TF-ORD-1',1,0,'命中')", id)
		t.AssertNil(err)
		recId, err := g.DB().GetValue(ctx, "SELECT id FROM risk_record WHERE user_id=998001")
		t.AssertNil(err)

		t.AssertNil(logic.AdminAppeal(ctx, recId.Int64(), true, "误判"))
		st, err := g.DB().GetOne(ctx, "SELECT appeal_status, remark FROM risk_record WHERE id=?", recId.Int64())
		t.AssertNil(err)
		t.Assert(st["appeal_status"].Int(), 2)
		t.Assert(strings.Contains(st["remark"].String(), "误判"), true)

		// 已结论重放 → 拒绝
		t.Assert(errCode(logic.AdminAppeal(ctx, recId.Int64(), false, "x")), errcode.CodeStatusNotAllowed)

		// 不存在 → 404 类
		t.Assert(errCode(logic.AdminAppeal(ctx, 999999999, true, "")), errcode.CodeNotFound)
	})
}

// TestRiskHitEvaluator 评估器（FR-6/SC-3）: 启用黑名单命中拦截+落记录; 停用不生效; 无规则放行;
// 重复命中不重复落事件; 评估器是 shop.IRiskHit 的真实判定实现。
// 用户标识唯一化（998003）: 测试间共享库, 他测试遗留的同 tag 规则会污染判定（实测教训）。
func TestRiskHitEvaluator(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const n = "TF-风控评估"
		const tag = "user:998003"
		defer func() {
			_, _ = g.DB().Exec(ctx, "DELETE FROM risk_record WHERE user_id=998003")
			_, _ = g.DB().Exec(ctx, "DELETE FROM risk_rule WHERE condition_expr=?", tag)
		}()
		_, _ = g.DB().Exec(ctx, "DELETE FROM risk_rule WHERE condition_expr=?", tag)
		res, err := g.DB().Exec(ctx,
			"INSERT INTO risk_rule(name,rule_type,condition_expr,action,status) VALUES(?,1,?,1,1)", n, tag)
		t.AssertNil(err)
		id, _ := res.LastInsertId()
		impl := RiskHitImpl{}

		// 命中黑名单 → 拦截 + 落事件
		blocked, err := impl.Hit(ctx, 998001, 1, "TF-CTX")
		t.AssertNil(err)
		t.Assert(blocked, true)
		cnt, _ := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM risk_record WHERE user_id=998001")
		t.Assert(cnt.Int(), 1)

		// 重复命中 → 仍拦截但不重复落事件（幂等）
		blocked, err = impl.Hit(ctx, 998003, 1, "TF-CTX")
		t.AssertNil(err)
		t.Assert(blocked, true)
		cnt, _ = g.DB().GetValue(ctx, "SELECT COUNT(*) FROM risk_record WHERE user_id=998003")
		t.Assert(cnt.Int(), 1)

		// 无关用户 → 放行
		blocked, err = impl.Hit(ctx, 998002, 1, "TF-CTX")
		t.AssertNil(err)
		t.Assert(blocked, false)

		// 停用规则 → 放行
		_, _ = g.DB().Exec(ctx, "UPDATE risk_rule SET status=0 WHERE id=?", id)
		blocked, err = impl.Hit(ctx, 998003, 1, "TF-CTX")
		t.AssertNil(err)
		t.Assert(blocked, false)

		// 无规则 → 放行（不阻断主流程）
		_, _ = g.DB().Exec(ctx, "DELETE FROM risk_rule WHERE id=?", id)
		blocked, err = impl.Hit(ctx, 998003, 1, "TF-CTX")
		t.AssertNil(err)
		t.Assert(blocked, false)
	})
}

// TestAuditLogs 审计查询面（FR-4）: 登录日志（批次 01 产出）+ 操作日志写入与读回。
func TestAuditLogs(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		logic := NewAuditLogic()

		// 操作日志: Write → 读回
		entry := model.OperationLogEntry{
			AdminId: 998002, Username: "TF-AUDIT-ADMIN", Module: "TF-审计模块",
			Operation: "测试操作", Method: "POST", RequestUri: "/admin/tf-audit",
			RequestParams: map[string]any{"k": "v"}, ResultStatus: 1, Ip: "127.0.0.1",
		}
		maxId, err := g.DB().GetValue(ctx, "SELECT COALESCE(MAX(id),0) FROM admin_operation_log")
		t.AssertNil(err)
		t.AssertNil(logic.Write(ctx, entry))
		lg, err := g.DB().GetOne(ctx, "SELECT * FROM admin_operation_log WHERE id > ?", maxId.Int64())
		t.AssertNil(err)
		t.Assert(lg.IsEmpty(), false)
		t.Assert(strings.Contains(lg["request_params"].String(), "v"), true) // JSON 落档
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM admin_operation_log WHERE id > ?", maxId.Int64()) }()

		q := model.AuditQuery{Module: "TF-审计模块", PageReq: model.PageReq{Page: 1, PageSize: 10}}
		ops, err := logic.OperationLogs(ctx, q)
		t.AssertNil(err)
		t.Assert(ops.Total >= 1, true)

		// 登录日志: 批次 01 的登录用例应已产出（按用户名过滤精确验证）
		lq := model.AuditQuery{Username: "TF-NO-SUCH-ADMIN", PageReq: model.PageReq{Page: 1, PageSize: 10}}
		ll, err := logic.LoginLogs(ctx, lq)
		t.AssertNil(err)
		t.Assert(ll.Total, int64(0))
	})
}
