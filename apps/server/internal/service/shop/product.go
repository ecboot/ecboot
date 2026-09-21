// product.go 商品目录域——表: product_category / product_brand / product_spu / product_sku（V2/V23/V31）。
// 规则: 三级分类（有商品禁删）; 两级上下架（可售=SPU上架 AND SKU启用 AND 库存>0）;
// SKU 规格组合唯一=生成列 specs_hash（V23 语义级, 键序无关）; 限售=SPU 黑名单区划码（下单校验/列表 ES 过滤）。
package shop

import (
	"context"

	"ecboot/internal/model"
)

// IProductLogic 商品目录。
type IProductLogic interface {
	// ---- 游客/会员浏览（C 端） ----
	CategoryTree(ctx context.Context) ([]model.CategoryNode, error)
	Brands(ctx context.Context, page model.PageReq) (*model.PageResult[model.BrandItem], error)
	Products(ctx context.Context, q model.ProductQuery) (*model.PageResult[model.ProductCard], error)
	// Detail 详情: SPU+SKU 列表（含可售态）+评价汇总; 软触足迹（会员）。
	Detail(ctx context.Context, spuId int64, viewerUserId int64) (*model.ProductDetail, error)
	Search(ctx context.Context, keyword string, q model.ProductQuery) (*model.PageResult[model.ProductCard], error)
	ProductReviews(ctx context.Context, spuId int64, score int, page model.PageReq) (*model.PageResult[model.ReviewCard], error)

	// ---- 后台管理 ----
	AdminCategoryCreate(ctx context.Context, in model.CategoryInput) (int64, error)
	AdminCategoryUpdate(ctx context.Context, id int64, in model.CategoryInput) error
	AdminCategoryDelete(ctx context.Context, id int64) error // 有商品引用拒绝
	AdminBrandCreate(ctx context.Context, in model.BrandInput) (int64, error)
	AdminBrandUpdate(ctx context.Context, id int64, in model.BrandInput) error
	AdminBrandDelete(ctx context.Context, id int64) error
	AdminProductList(ctx context.Context, q model.AdminProductQuery) (*model.PageResult[model.AdminProductItem], error)
	AdminProductCreate(ctx context.Context, in model.SpuInput) (spuId int64, spuNo string, err error)
	AdminProductDetailView(ctx context.Context, spuId int64) (*model.AdminProductDetailView, error) // 含成本价
	AdminProductUpdate(ctx context.Context, spuId int64, in model.SpuInput) error
	AdminProductDelete(ctx context.Context, spuId int64) error
	// AdminProductStatus 上下架（无启用 SKU 禁上架, 30005）。
	AdminProductStatus(ctx context.Context, spuId int64, status int) error
	AdminProductRestrict(ctx context.Context, spuId int64, codes []string) error
	AdminSkuCreate(ctx context.Context, spuId int64, in model.SkuInput) (int64, error) // 规格组合冲突 30006
	AdminSkuUpdate(ctx context.Context, skuId int64, in model.SkuInput) error
	AdminSkuDelete(ctx context.Context, skuId int64) error
	AdminSkuStatus(ctx context.Context, skuId int64, status int) error
}
