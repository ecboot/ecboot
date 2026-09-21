package shop

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// ---- 010 评审修复轮的回归断言（连线把 005 遗产"放活"后暴露的零值覆盖类缺陷） ----

func cleanupCatByName(ctx context.Context, t *gtest.T, name string) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_category WHERE name=?", name)
}

// TestCategoryCreateDefaultsEnabled C1: 新建分类默认启用（不会被零值覆盖成禁用）。
func TestCategoryCreateDefaultsEnabled(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const name = "T-CatDefault"
		defer cleanupCatByName(ctx, t, name)
		cleanupCatByName(ctx, t, name)

		// api 的 Create 契约无 status 字段 → 控制器/服务以零值调用（模拟真实链路）
		id, err := NewProductLogic().AdminCategoryCreate(ctx, model.CategoryInput{
			ParentId: 0, Name: name, Icon: "", Level: 3, Sort: 1,
		})
		t.AssertNil(err)
		rec, err := g.DB().GetOne(ctx, "SELECT status, level FROM product_category WHERE id=?", id)
		t.AssertNil(err)
		t.Assert(rec["status"].Int(), 1) // C1: 必须是启用（此前被写成 0 → C 端不可见）
		t.Assert(rec["level"].Int(), 3)
	})
}

// TestCategoryUpdateKeepsUnsetCols C2: 更新（api 契约无 parentId/level）不得清零这两列。
func TestCategoryUpdateKeepsUnsetCols(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const name = "T-CatKeep"
		defer cleanupCatByName(ctx, t, name)
		cleanupCatByName(ctx, t, name)

		id, err := NewProductLogic().AdminCategoryCreate(ctx, model.CategoryInput{
			ParentId: 0, Name: name, Level: 3, Sort: 1,
		})
		t.AssertNil(err)
		// 造非零 parent/level（模拟子分类）
		_, err = g.DB().Exec(ctx, "UPDATE product_category SET parent_id=1, level=2 WHERE id=?", id)
		t.AssertNil(err)

		// 仅改名（api Update 契约只有 Id/Name/Icon/Sort/Status）
		t.AssertNil(NewProductLogic().AdminCategoryUpdate(ctx, id, model.CategoryInput{
			Name: name + "_x", Sort: 9, Status: 1,
		}))
		rec, err := g.DB().GetOne(ctx, "SELECT name, sort, parent_id, level, status FROM product_category WHERE id=?", id)
		t.AssertNil(err)
		t.Assert(rec["name"].String(), name+"_x")
		t.Assert(rec["sort"].Int(), 9)
		t.Assert(rec["parent_id"].Int(), 1) // C2: 不得被清零
		t.Assert(rec["level"].Int(), 2)     // C2: 不得被清零
		t.Assert(rec["status"].Int(), 1)
	})
}

// TestBrandListExcludesDeleted I5: 后台品牌列表须过滤软删。
func TestBrandListExcludesDeleted(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const name = "T-BrandDel"
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_brand WHERE name=?", name)
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM product_brand WHERE name=?", name) }()

		id, err := NewProductLogic().AdminBrandCreate(ctx, model.BrandInput{Name: name, Status: 1})
		t.AssertNil(err)
		t.AssertNil(NewProductLogic().AdminBrandDelete(ctx, id))

		res, err := NewProductLogic().AdminBrandList(ctx, 0, model.PageReq{Page: 1, PageSize: 100})
		t.AssertNil(err)
		for _, it := range res.List {
			t.AssertNE(it.Id, id) // I5: 软删品牌不得出现在后台列表
		}
	})
}

// TestInventoryAdjustAvailableGuard I1: 扣减超过可售（total-locked）→ 30007（而非撞表级 CHECK 报系统错误）。
func TestInventoryAdjustAvailableGuard(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const sfx = "inv_guard"
		skuId := seedInvSku(ctx, t, sfx, 10, 8, 0) // total=10, locked=8 → 可售 2
		defer cleanupInvSku(ctx, t, sfx)

		// 扣减 3 > 可售 2 → 30007（此前会撞 CHECK(chk_locked_le_total) 返回无业务码错误）
		_, err := InventoryAdjust(ctx, skuId, -3, "admin:1", "over-available")
		t.Assert(errCode(err), errcode.CodeInventoryAdjust)
		// 值未被改动
		inv, err := g.DB().GetOne(ctx, "SELECT total FROM inventory WHERE sku_id=?", skuId)
		t.AssertNil(err)
		t.Assert(inv["total"].Int(), 10)

		// 恰好扣到可售为 0 → 成功（total=8, locked=8, available=0）
		after, err := InventoryAdjust(ctx, skuId, -2, "admin:1", "to-zero-available")
		t.AssertNil(err)
		t.Assert(after, 8)

		// delta=0 → 10001
		_, err = InventoryAdjust(ctx, skuId, 0, "admin:1", "noop")
		t.Assert(errCode(err), errcode.CodeInvalidParam)
	})
}
