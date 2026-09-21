package shop

import (
	"context"
	"os"
	"time"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/database/gdb"
	gredis "github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func init() {
	// 时区口径统一（009 评审 C2）: 与库内 UTC 墙钟一致（见 main.go）
	time.Local = time.UTC
	_ = os.Setenv("ECBOOT_MOCK", "true")
	_ = gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			{
				Link: "mysql:myuser:secret@tcp(127.0.0.1:13306)/mydatabase",
			},
		},
	})
	gredis.SetConfig(&gredis.Config{
		Address: "127.0.0.1:6379",
		Db:      0,
	})
}

// ---- 测试数据 builder（自建+清理, research D6） ----

type TradeFixture struct {
	Ctx        context.Context
	CategoryId int64
	BrandId    int64
	SpuId      int64
	SkuId      int64 // 库存 100 可售
	OrderNo    string
}

// tradeTestSpuName 本 fixture 的商品名（清理一律按**精确名**匹配）。
const tradeTestSpuName = "TF-商品"

// setupTradeFixture 建测试商品（上架+库存 100）并登记清理。
// 012 修复轮: ① product_spu.spu_no 自 000002 起即非空且无默认值（本 fixture 漏写 → INSERT 必失败,
// 因长期无调用者而未暴露）; ② 清理顺序——inventory 的子查询依赖 product_sku, 原顺序先删 SKU 致库存行泄漏。
// ③ 清理判据由 `LIKE 'TF-%'` 收紧为精确名: 原模糊模式会**跨包误删**同名前缀的 fixture 数据
// （go test ./... 并行跑包, 实测把 controller 包的下单 fixture 删掉致偶发失败）。
func setupTradeFixture(t *gtest.T) *TradeFixture {
	ctx := context.Background()
	f := &TradeFixture{Ctx: ctx}
	cleanup := func() {
		_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id IN (SELECT id FROM product_sku WHERE sku_no LIKE 'TF-SKU%')")
		_, _ = g.DB().Exec(ctx, `
DELETE ps FROM product_sku ps JOIN product_spu s ON ps.spu_id=s.id WHERE s.name = ?`, tradeTestSpuName)
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE name = ?", tradeTestSpuName)
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_brand WHERE name = 'TF-品牌'")
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_category WHERE name = 'TF-分类'")
	}
	cleanup()

	if _, err := g.DB().Exec(ctx,
		"INSERT INTO product_brand (name, status) VALUES ('TF-品牌', 1)"); err != nil {
		t.Fatal(err)
	}
	rec, err := g.DB().GetOne(ctx, "SELECT id FROM product_brand WHERE name='TF-品牌'")
	t.AssertNil(err)
	f.BrandId = rec["id"].Int64()

	catId, err := g.DB().Model("product_category").Ctx(ctx).Data(g.Map{
		"parent_id": 0, "name": "TF-分类", "level": 3, "status": 1,
	}).InsertAndGetId()
	t.AssertNil(err)
	f.CategoryId = catId

	spuId, err := g.DB().Model("product_spu").Ctx(ctx).Data(g.Map{
		"spu_no": "TF-SPU-001",
		"name":   tradeTestSpuName, "category_id": f.CategoryId, "brand_id": f.BrandId,
		"images": `["http://img/tf.png"]`, "spec_definitions": `[{"name":"颜色","values":["黑"]}]`,
		"status": 1,
	}).InsertAndGetId()
	t.AssertNil(err)
	f.SpuId = spuId

	skuId, err := g.DB().Model("product_sku").Ctx(ctx).Data(g.Map{
		"sku_no": "TF-SKU-001", "spu_id": f.SpuId, "name": "TF-商品 黑",
		"specs": `{"颜色":"黑"}`, "price": "10.00", "status": 1,
	}).InsertAndGetId()
	t.AssertNil(err)
	f.SkuId = skuId
	_, _ = g.DB().Exec(ctx,
		"INSERT INTO inventory (sku_id, total, locked) VALUES (?, 100, 0) ON DUPLICATE KEY UPDATE total=100, locked=0", f.SkuId)

	// 上架
	_, _ = g.DB().Exec(ctx, "UPDATE product_spu SET status=1 WHERE id=?", f.SpuId)
	return f
}

// cleanupFixture 测试后清理（顺序: 库存 → sku → spu → 品牌/分类; 判据为精确名, 见 setupTradeFixture）。
func cleanupFixture(ctx context.Context, f *TradeFixture) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id IN (SELECT id FROM product_sku WHERE sku_no LIKE 'TF-SKU%')")
	_, _ = g.DB().Exec(ctx, `DELETE ps FROM product_sku ps JOIN product_spu s ON ps.spu_id=s.id WHERE s.name = ?`, tradeTestSpuName)
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE name = ?", tradeTestSpuName)
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_brand WHERE name = 'TF-品牌'")
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_category WHERE name = 'TF-分类'")
}
