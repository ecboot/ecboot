package shop

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// ---- 测试数据（自建+清理; init 基座见 trade_test.go） ----

// seedLogistics 建测试物流公司（用后 cleanupLogistics; 按 code 双路清理）。
func seedLogistics(ctx context.Context, t *gtest.T, code string, status int) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM `logistics_company` WHERE code=?", code)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO `logistics_company`(code,name,tracking_rule,sort,status) VALUES(?,?,?,?,?)",
		code, code+"-name", "", 0, status)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

func cleanupLogistics(ctx context.Context, t *gtest.T, codes ...string) {
	for _, c := range codes {
		_, _ = g.DB().Exec(ctx, "DELETE FROM `logistics_company` WHERE code=?", c)
	}
}

// ---- US1 物流公司字典 ----

// TestLogisticsCrud 物流公司 CRUD: 编码唯一 + 编码不可改 + 停用保留（FR-001~006）。
func TestLogisticsCrud(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const code = "T_SF"
		defer cleanupLogistics(ctx, t, code, "T_NEWCODE")

		// 创建成功
		id, err := LogisticsCreate(ctx, model.LogisticsCompanyInput{
			Code: code, Name: "顺丰速运", TrackingRule: "SF+12位数字", Status: 1,
		})
		t.AssertNil(err)
		t.AssertGT(id, 0)

		// 编码重复 → 40012
		_, err = LogisticsCreate(ctx, model.LogisticsCompanyInput{Code: code, Name: "另一家", Status: 1})
		t.Assert(errCode(err), errcode.CodeLogisticsCodeTaken)

		// 详情
		d, err := LogisticsDetail(ctx, id)
		t.AssertNil(err)
		t.Assert(d.Code, code)
		t.Assert(d.Name, "顺丰速运")
		t.Assert(d.TrackingRule, "SF+12位数字")
		t.Assert(d.Status, 1)

		// 修改: 名称/规则/状态生效; 编码不可改（入参携带亦被忽略, FR-003）
		t.AssertNil(LogisticsUpdate(ctx, id, model.LogisticsCompanyInput{
			Code: "T_NEWCODE", Name: "顺丰", TrackingRule: "SF+13位", Status: 0,
		}))
		d, err = LogisticsDetail(ctx, id)
		t.AssertNil(err)
		t.Assert(d.Name, "顺丰")
		t.Assert(d.TrackingRule, "SF+13位")
		t.Assert(d.Status, 0)
		t.Assert(d.Code, code) // 编码未变

		// 停用保留: 不筛选列表仍含（FR-006）
		res, err := LogisticsList(ctx, 0, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		found := false
		for _, it := range res.List {
			if it.Id == id {
				found = true
				t.Assert(it.Status, 0)
			}
		}
		t.Assert(found, true)

		// status=1 筛选不含停用项
		res, err = LogisticsList(ctx, 1, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		for _, it := range res.List {
			t.AssertNE(it.Id, id)
		}

		// 详情不存在 → 10006
		_, err = LogisticsDetail(ctx, 999999999)
		t.Assert(errCode(err), errcode.CodeNotFound)

		// 修改不存在 → 10006
		err = LogisticsUpdate(ctx, 999999999, model.LogisticsCompanyInput{Name: "x", Status: 1})
		t.Assert(errCode(err), errcode.CodeNotFound)

		// 软删 → 列表与详情均不可见
		t.AssertNil(LogisticsDelete(ctx, id))
		_, err = LogisticsDetail(ctx, id)
		t.Assert(errCode(err), errcode.CodeNotFound)
		res, err = LogisticsList(ctx, 0, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		for _, it := range res.List {
			t.AssertNE(it.Id, id)
		}

		// 重复删除 → 10006
		err = LogisticsDelete(ctx, id)
		t.Assert(errCode(err), errcode.CodeNotFound)
	})
}

// TestLogisticsListPaging 列表筛选与分页（FR-001）。
func TestLogisticsListPaging(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			c1 = "T_lst1"
			c2 = "T_lst2"
		)
		defer cleanupLogistics(ctx, t, c1, c2)
		seedLogistics(ctx, t, c1, 1)
		seedLogistics(ctx, t, c2, 1)

		// 分页: pageSize=1 → 单行但总数>=2
		res, err := LogisticsList(ctx, 1, model.PageReq{Page: 1, PageSize: 1})
		t.AssertNil(err)
		t.Assert(len(res.List), 1)
		t.AssertGE(res.Total, 2)
	})
}
