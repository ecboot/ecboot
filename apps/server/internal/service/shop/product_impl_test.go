package shop

import (
	"context"
	"errors"
	"testing"

	"github.com/gogf/gf/v2/errors/gerror"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/database/gdb"
	gredis "github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/model"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"

	"ecboot/internal/testutil"
)

func init() {
	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			{
				Link: testutil.DSN(),
			},
		},
	})
	gredis.SetConfig(&gredis.Config{
		Address: "127.0.0.1:6379",
		Db:      0,
	})
	_ = g.Cfg().GetAdapter()
}

// ---- 测试数据工具（自建+清理, research D6） ----

const (
	testBrandName = "T-测试品牌"
	testCatName   = "T-测试三级分类"
	testSpuName   = "T-测试商品"
)

func catalogSuite(t *gtest.T) (ctx context.Context, brandId, catId, spuId int64) {
	ctx = context.Background()
	// 前置清理（按测试命名约定; SKU 经 SPU 关联删, 连带清残留编码行）
	g.DB().Exec(ctx, "DELETE ps FROM product_sku ps JOIN product_spu s ON ps.spu_id=s.id WHERE s.name LIKE 'T-测试%'")
	g.DB().Exec(ctx, "DELETE FROM product_spu WHERE name LIKE 'T-测试%'")
	g.DB().Exec(ctx, "DELETE FROM product_sku WHERE sku_no LIKE 'T-SKU%'")
	g.DB().Exec(ctx, "DELETE FROM product_brand WHERE name LIKE 'T-测试%'")
	g.DB().Exec(ctx, "DELETE FROM product_category WHERE name LIKE 'T-测试%'")

	// 品牌与三级分类（level=3, parent 用 0 简化树断言独立）
	res, err := g.DB().Model("product_brand").Ctx(ctx).Data(g.Map{"name": testBrandName, "status": 1}).InsertAndGetId()
	t.AssertNil(err)
	brandId = res
	res, err = g.DB().Model("product_category").Ctx(ctx).Data(g.Map{
		"parent_id": 0, "name": testCatName, "level": 3, "sort": 999, "status": 1,
	}).InsertAndGetId()
	t.AssertNil(err)
	catId = res
	return
}

func catalogTeardown(ctx context.Context, t *gtest.T, spuId int64) {
	g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id IN (SELECT id FROM (SELECT id FROM product_sku WHERE spu_id=?) x)", spuId)
	g.DB().Exec(ctx, "DELETE FROM product_sku WHERE spu_id=?", spuId)
	g.DB().Exec(ctx, "DELETE FROM product_spu WHERE id=?", spuId)
	g.DB().Exec(ctx, "DELETE FROM product_brand WHERE name=?", testBrandName)
	g.DB().Exec(ctx, "DELETE FROM product_category WHERE name=?", testCatName)
}

// svcIProduct 便捷引用（实现注册后由实现提供; 测试直接构造实现结构）。
func TestCategoryLifecycle(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		impl := NewProductLogic()
		// 创建三级分类
		catId, err := impl.AdminCategoryCreate(ctx, model.CategoryInput{ParentId: 0, Name: "T-分类A", Level: 3, Sort: 99})
		t.AssertNil(err)
		// 更新
		t.AssertNil(impl.AdminCategoryUpdate(ctx, catId, model.CategoryInput{Name: "T-分类A2", Sort: 98}))
		// 树中可见
		tree, err := impl.AdminCategoryTree(ctx)
		t.AssertNil(err)
		found := false
		var walk func(nodes []model.AdminCategoryNode)
		walk = func(nodes []model.AdminCategoryNode) {
			for _, n := range nodes {
				if n.Id == catId {
					found = true
				}
				walk(n.Children)
			}
		}
		walk(tree)
		t.Assert(found, true)
		// 删除（无引用）成功
		t.AssertNil(impl.AdminCategoryDelete(ctx, catId))
		// 引用禁删：重建分类并挂一个商品
		catId2, err := impl.AdminCategoryCreate(ctx, model.CategoryInput{ParentId: 0, Name: "T-分类B", Level: 3})
		t.AssertNil(err)
		spuId, _, err := impl.AdminProductCreate(ctx, model.SpuInput{Name: "T-引用商品", CategoryId: catId2, Images: []string{"x"}})
		t.AssertNil(err)
		err = impl.AdminCategoryDelete(ctx, catId2)
		t.Assert(errCode(err), errcode.CodeCategoryInUse)
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE id=?", spuId)
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_category WHERE id IN (?, ?)", catId, catId2)
	})
}

func TestBrandLifecycle(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_brand WHERE name LIKE 'T-品牌%'")

		impl := NewProductLogic()

		// ① 创建
		id, err := impl.AdminBrandCreate(ctx, model.BrandInput{Name: "T-品牌X", Status: 1})
		t.AssertNil(err)

		// ② 名称唯一: 重复创建被 uk_name 拒绝（FR-008 数据库保证）
		_, err = impl.AdminBrandCreate(ctx, model.BrandInput{Name: "T-品牌X", Status: 1})
		t.Assert(err != nil, true)

		// ③ 更新改名 + 启用
		t.AssertNil(impl.AdminBrandUpdate(ctx, id, model.BrandInput{Name: "T-品牌X2", Status: 1}))

		// ④ 软删后 name 仍占用（ADR-0002 严格唯一）: 软删行 name='T-品牌X2', 再建 'T-品牌X2' 被拒
		t.AssertNil(impl.AdminBrandDelete(ctx, id))
		_, err = impl.AdminBrandCreate(ctx, model.BrandInput{Name: "T-品牌X2", Status: 1})
		t.Assert(err != nil, true)

		// 清理
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_brand WHERE id=?", id)
	})
}

// TestSpuSkuLifecycle 核心生命周期（US2 主链 + SC-003 三防御）。
func TestSpuSkuLifecycle(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE name LIKE 'T-测试%'")
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_sku WHERE sku_no LIKE 'T-SKU%'")

		impl := NewProductLogic()
		// 前置: 分类
		catId, err := impl.AdminCategoryCreate(ctx, model.CategoryInput{ParentId: 0, Name: "T-目录分类", Level: 3})
		t.AssertNil(err)
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM product_category WHERE id=?", catId) }()

		// 创建 SPU（规格定义: 颜色）
		spuId, spuNo, err := impl.AdminProductCreate(ctx, model.SpuInput{
			Name:            "T-测试SPU",
			CategoryId:      catId,
			Images:          []string{"http://img/1.png"},
			SpecDefinitions: []map[string]any{{"name": "颜色", "values": []string{"黑", "白"}}},
		})
		t.AssertNil(err)
		t.Assert(spuNo != "", true)

		// SKU 创建（两个规格组合）
		sku1, err := impl.AdminSkuCreate(ctx, spuId, model.SkuInput{
			Specs: map[string]string{"颜色": "黑"}, Price: "99.00", CostPrice: "60.00",
		})
		t.AssertNil(err)
		sku2, err := impl.AdminSkuCreate(ctx, spuId, model.SkuInput{
			Specs: map[string]string{"颜色": "白"}, Price: "129.00",
		})
		t.AssertNil(err)

		// 防御①: 重复规格组合 → 30006
		_, err = impl.AdminSkuCreate(ctx, spuId, model.SkuInput{Specs: map[string]string{"颜色": "黑"}, Price: "1.00"})
		t.Assert(errCode(err), errcode.CodeSpecDup)

		// 价格冗余刷新断言: min=99 max=129
		rec, err := g.DB().GetOne(ctx, "SELECT price_min, price_max FROM product_spu WHERE id=?", spuId)
		t.AssertNil(err)
		t.Assert(rec["price_min"].String(), "99.00")
		t.Assert(rec["price_max"].String(), "129.00")

		// 防御②: 无启用 SKU 禁上架——先全禁用
		t.AssertNil(impl.AdminSkuStatus(ctx, sku1, 0))
		t.AssertNil(impl.AdminSkuStatus(ctx, sku2, 0))
		err = impl.AdminProductStatus(ctx, spuId, 1)
		t.Assert(errCode(err), errcode.CodeNoEnableSku)

		// 启用一个 → 上架成功
		t.AssertNil(impl.AdminSkuStatus(ctx, sku1, 1))
		t.AssertNil(impl.AdminProductStatus(ctx, spuId, 1))

		// C 端可见 + 详情可售口径
		cards, err := impl.Products(ctx, model.ProductQuery{CategoryId: catId})
		t.AssertNil(err)
		found := false
		for _, c := range cards.List {
			if c.SpuId == spuId {
				found = true
			}
		}
		t.Assert(found, true)
		detail, err := impl.AdminProductDetail(ctx, spuId)
		t.AssertNil(err)
		t.Assert(len(detail.Skus), 2)
		// 按 SKU 定位成本断言（管理视角返回全量 SKU, 顺序不保证）
		var sku1Detail *model.AdminSkuDetail
		for idx := range detail.Skus {
			if detail.Skus[idx].SkuId == sku1 {
				sku1Detail = &detail.Skus[idx]
			}
		}
		t.Assert(sku1Detail != nil, true)
		t.Assert(sku1Detail.CostPrice, "60.00") // 管理视角含成本

		// 限售设置（全量替换语义）
		t.AssertNil(impl.AdminProductRestrict(ctx, spuId, []string{"650000", "540000"}))
		detail2, err := impl.AdminProductDetail(ctx, spuId)
		t.AssertNil(err)
		t.Assert(len(detail2.SaleRestrictCodes), 2)

		// 下架 → 删除（软删）
		t.AssertNil(impl.AdminProductStatus(ctx, spuId, 0))
		t.AssertNil(impl.AdminProductDelete(ctx, spuId))

		// 清理（**库存行必须先删**: product_sku 一删, 按 SKU 定位库存的子查询就永远空集 → 静默泄漏,
		// 泄漏行 total=0 恒满足预警条件, 累积到 100 行会挤爆 InventoryWarnings 的分页窗口）
		_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id IN (SELECT id FROM product_sku WHERE spu_id=?)", spuId)
		_, _ = g.DB().Exec(ctx, "DELETE ps FROM product_sku ps JOIN product_spu s ON ps.spu_id=s.id WHERE s.name LIKE 'T-测试%'")
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_sku WHERE spu_id IN (?, ?)", sku1, sku2)
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE id IN (?, ?)", spuId, 0)
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_category WHERE id=?", catId)
	})
}

// TestSkuPriceRefreshOnEdit 编辑 SKU 价格时冗余刷新（research D1）。
func TestSkuPriceRefreshOnEdit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		impl := NewProductLogic()
		catId, err := impl.AdminCategoryCreate(ctx, model.CategoryInput{ParentId: 0, Name: "T-价格分类", Level: 3})
		t.AssertNil(err)
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM product_category WHERE id=?", catId) }()
		spuId, _, err := impl.AdminProductCreate(ctx, model.SpuInput{Name: "T-价格SPU", CategoryId: catId, Images: []string{"x"}})
		t.AssertNil(err)
		defer func() {
			// 库存行先删（同上: 后删 SKU 才能定位到库存; 顺序颠倒即静默泄漏）
			_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id IN (SELECT id FROM product_sku WHERE spu_id=?)", spuId)
			_, _ = g.DB().Exec(ctx, "DELETE FROM product_sku WHERE spu_id=?", spuId)
			_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE id=?", spuId)
		}()

		skuA, err := impl.AdminSkuCreate(ctx, spuId, model.SkuInput{Specs: map[string]string{"尺码": "M"}, Price: "10.00"})
		t.AssertNil(err)
		skuB, err := impl.AdminSkuCreate(ctx, spuId, model.SkuInput{Specs: map[string]string{"尺码": "L"}, Price: "30.00"})
		t.AssertNil(err)

		// 编辑 skuA 价格 10 → 50: 冗余 min 变 30
		t.AssertNil(impl.AdminSkuUpdate(ctx, skuA, model.SkuInput{Price: "50.00"}))
		rec, err := g.DB().GetOne(ctx, "SELECT price_min, price_max FROM product_spu WHERE id=?", spuId)
		t.AssertNil(err)
		t.Assert(rec["price_min"].String(), "30.00")
		t.Assert(rec["price_max"].String(), "50.00")

		// 删除 skuB: min/max 均 50
		t.AssertNil(impl.AdminSkuDelete(ctx, skuB))
		rec, err = g.DB().GetOne(ctx, "SELECT price_min, price_max FROM product_spu WHERE id=?", spuId)
		t.AssertNil(err)
		fmtCheck(rec, "50.00")
		_ = skuB
	})
}

func fmtCheck(rec gdb.Record, min string) {
	if rec["price_min"].String() != min {
		t2 := &testing.T{}
		t2.Errorf("price_min = %s, want %s", rec["price_min"].String(), min)
	}
	if rec["price_max"].String() != min {
		t2 := &testing.T{}
		t2.Errorf("price_max = %s, want %s", rec["price_max"].String(), min)
	}
}

// errCode 提取错误中的契约码（兼容 gerror 包装）。
func errCode(err error) int {
	var ge *gerror.Error
	if errors.As(err, &ge) {
		return ge.Code().Code()
	}
	return -1
}
