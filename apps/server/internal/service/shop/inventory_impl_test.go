package shop

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// ---- 库存测试数据（自建+清理; 复用 009 的 SPU fixture） ----

// seedInvSku 建 SPU+SKU+库存行, 返回 skuId（用后 cleanupInvSku）。
func seedInvSku(ctx context.Context, t *gtest.T, suffix string, total, locked, warn int) int64 {
	spuId := seedFloorSpu(ctx, t, suffix, `["http://img/inv.png"]`, "10.00", 1)
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_sku WHERE sku_no=?", "TF-SKU-"+suffix)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO product_sku(sku_no,spu_id,specs,price,status) VALUES(?,?,?,?,1)",
		"TF-SKU-"+suffix, spuId, `{"规格":"A"}`, "10.00")
	t.AssertNil(err)
	skuId, _ := res.LastInsertId()

	_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id=?", skuId)
	_, err = g.DB().Exec(ctx,
		"INSERT INTO inventory(sku_id,total,locked,warn_count) VALUES(?,?,?,?)",
		skuId, total, locked, warn)
	t.AssertNil(err)
	return skuId
}

// cleanupInvSku 清理库存/流水/SKU/SPU 及引用行。
func cleanupInvSku(ctx context.Context, t *gtest.T, suffix string) {
	skuNo := "TF-SKU-" + suffix
	_, _ = g.DB().Exec(ctx, "DELETE FROM inventory_log WHERE sku_id IN (SELECT id FROM product_sku WHERE sku_no=?)", skuNo)
	_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id IN (SELECT id FROM product_sku WHERE sku_no=?)", skuNo)
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_sku WHERE sku_no=?", skuNo)
	cleanupFloorSpu(ctx, t, suffix)
}

// findInvItem 从列表中取指定 SKU 项。
func findInvItem(list []model.InventoryItem, skuId int64) *model.InventoryItem {
	for i := range list {
		if list[i].SkuId == skuId {
			return &list[i]
		}
	}
	return nil
}

// TestInventoryList 库存列表: 可售推导 = total - locked（FR-006）。
func TestInventoryList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const sfx = "inv_lst"
		skuId := seedInvSku(ctx, t, sfx, 10, 2, 5)
		defer cleanupInvSku(ctx, t, sfx)

		res, err := InventoryList(ctx, skuId, "", model.PageReq{Page: 1, PageSize: 20})
		t.AssertNil(err)
		it := findInvItem(res.List, skuId)
		t.Assert(it != nil, true)
		t.Assert(it.Total, 10)
		t.Assert(it.Locked, 2)
		t.Assert(it.Available, 8)
		t.Assert(it.WarnCount, 5)
		t.Assert(it.SkuNo, "TF-SKU-"+sfx)

		// 关键词筛选（sku_no 前缀）命中
		res, err = InventoryList(ctx, 0, "TF-SKU-"+sfx, model.PageReq{Page: 1, PageSize: 20})
		t.AssertNil(err)
		t.Assert(findInvItem(res.List, skuId) != nil, true)

		// 不存在的 skuId → 空列表
		res, err = InventoryList(ctx, 999999999, "", model.PageReq{Page: 1, PageSize: 20})
		t.AssertNil(err)
		t.Assert(findInvItem(res.List, 999999999), (*model.InventoryItem)(nil))
	})
}

// TestInventoryWarnings 预警: available <= warn_count 命中（含等于）（FR-007）。
func TestInventoryWarnings(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			sfxOk  = "inv_w_ok"  // avail 8 > warn 5 → 不命中
			sfxEq  = "inv_w_eq"  // avail 5 == warn 5 → 命中（含等于）
			sfxLow = "inv_w_low" // avail 3 <= warn 5 → 命中
		)
		okId := seedInvSku(ctx, t, sfxOk, 10, 2, 5)
		eqId := seedInvSku(ctx, t, sfxEq, 5, 0, 5)
		lowId := seedInvSku(ctx, t, sfxLow, 3, 0, 5)
		defer cleanupInvSku(ctx, t, sfxOk)
		defer cleanupInvSku(ctx, t, sfxEq)
		defer cleanupInvSku(ctx, t, sfxLow)

		res, err := InventoryWarnings(ctx, model.PageReq{Page: 1, PageSize: 100})
		t.AssertNil(err)
		t.Assert(findInvItem(res.List, eqId) != nil, true)                 // 等于阈值命中
		t.Assert(findInvItem(res.List, lowId) != nil, true)                // 低于阈值命中
		t.Assert(findInvItem(res.List, okId), (*model.InventoryItem)(nil)) // 高于阈值不命中
	})
}

// TestInventoryAdjust 调整: 正负生效 + 流水留痕 + 不可为负（FR-008）。
func TestInventoryAdjust(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const sfx = "inv_adj"
		skuId := seedInvSku(ctx, t, sfx, 10, 2, 5)
		defer cleanupInvSku(ctx, t, sfx)
		const op = "admin:1"

		// 正向调整 +5 → total 15, 流水一条（change_type=5）
		after, err := InventoryAdjust(ctx, skuId, 5, op, "补货")
		t.AssertNil(err)
		t.Assert(after, 15)
		rec, err := g.DB().GetOne(ctx,
			"SELECT change_type,quantity,total_after,locked_after,operator,remark FROM inventory_log WHERE sku_id=? ORDER BY id DESC LIMIT 1", skuId)
		t.AssertNil(err)
		t.Assert(rec["change_type"].Int(), 5)
		t.Assert(rec["quantity"].Int(), 5)
		t.Assert(rec["total_after"].Int(), 15)
		t.Assert(rec["locked_after"].Int(), 2)
		t.Assert(rec["operator"].String(), op)
		t.Assert(rec["remark"].String(), "补货")

		// 负向调整 -3 → total 12
		after, err = InventoryAdjust(ctx, skuId, -3, op, "盘点亏损")
		t.AssertNil(err)
		t.Assert(after, 12)
		rec, err = g.DB().GetOne(ctx, "SELECT quantity,total_after FROM inventory_log WHERE sku_id=? ORDER BY id DESC LIMIT 1", skuId)
		t.AssertNil(err)
		t.Assert(rec["quantity"].Int(), 3) // 绝对值为 3
		t.Assert(rec["total_after"].Int(), 12)

		// 致负 → 30007（不可为负）
		_, err = InventoryAdjust(ctx, skuId, -999, op, "超额扣减")
		t.Assert(errCode(err), errcode.CodeInventoryAdjust)
		// 值未被改动
		inv, err := g.DB().GetOne(ctx, "SELECT total FROM inventory WHERE sku_id=?", skuId)
		t.AssertNil(err)
		t.Assert(inv["total"].Int(), 12)

		// 不存在 → 30007
		_, err = InventoryAdjust(ctx, 999999999, 1, op, "")
		t.Assert(errCode(err), errcode.CodeInventoryAdjust)
	})
}
