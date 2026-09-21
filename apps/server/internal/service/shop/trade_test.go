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
	os.Setenv("ECBOOT_MOCK", "true")
	gdb.SetConfig(gdb.Config{
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
	UserId     int64
	OrderNo    string
}

// setupTradeFixture 建测试商品（上架+库存 100）并登记清理。
func setupTradeFixture(t *gtest.T) *TradeFixture {
	ctx := context.Background()
	f := &TradeFixture{Ctx: ctx}
	cleanup := func() {
		_, _ = g.DB().Exec(ctx, `
DELETE ps FROM product_sku ps JOIN product_spu s ON ps.spu_id=s.id WHERE s.name LIKE 'TF-%'`)
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE name LIKE 'TF-%'")
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_brand WHERE name LIKE 'TF-%'")
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_category WHERE name LIKE 'TF-%'")
		_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id IN (SELECT id FROM product_sku WHERE sku_no LIKE 'TF-SK%')")
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
		"name": "TF-商品", "category_id": f.CategoryId, "brand_id": f.BrandId,
		"images": `["http://img/tf.png"]`, "spec_definitions": `[{"name":"颜色","values":["黑"]}]`,
		"status": 1,
	}).InsertAndGetId()
	t.AssertNil(err)
	f.SpuId = spuId

	skuId, err := g.DB().Model("product_sku").Ctx(ctx).Data(g.Map{
		"sku_no": "TF-SKU-001", "spu_id": f.SpuId,
		"specs": `{"颜色":"黑"}`, "price": "10.00", "status": 1,
	}).InsertAndGetId()
	t.AssertNil(err)
	f.SkuId = skuId
	_, _ = g.DB().Exec(ctx, "INSERT INTO inventory (sku_id, total, locked) VALUES (?, 100, 0)", f.SkuId)

	// 上架
	_, _ = g.DB().Exec(ctx, "UPDATE product_spu SET status=1 WHERE id=?", f.SpuId)
	f.UserId = 999
	return f
}

// cleanupFixture 测试后清理（顺序: sku/库存 → spu → 品牌/分类）。
func cleanupFixture(ctx context.Context, f *TradeFixture) {
	_, _ = g.DB().Exec(ctx, "DELETE ps FROM product_sku ps JOIN product_spu s ON ps.spu_id=s.id WHERE s.name LIKE 'TF-%'")
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE name LIKE 'TF-%'")
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_brand WHERE name LIKE 'TF-%'")
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_category WHERE name LIKE 'TF-%'")
	_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id IN (SELECT id FROM product_sku WHERE sku_no LIKE 'TF-SK%')")
}
