// product.go 商品目录域——表: product_category / product_brand / product_spu / product_sku（V2/V23/V31）。
// 规则: 三级分类（有商品禁删）; 两级上下架（可售=SPU上架 AND SKU启用 AND 库存>0）;
// SKU 规格组合唯一=生成列 specs_hash（V23 语义级, 键序无关）; 限售=SPU 黑名单区划码（下单校验/列表 ES 过滤）。
package shop

import "context"

// IProductLogic 商品目录。
type IProductLogic interface {
	// ---- 游客/会员浏览（C 端） ----
	CategoryTree(ctx context.Context) ([]CategoryNode, error)
	Brands(ctx context.Context, page PageQuery) (*PageResult[BrandItem], error)
	Products(ctx context.Context, q ProductQuery) (*PageResult[ProductCard], error)
	// Detail 详情: SPU+SKU 列表（含可售态）+评价汇总; 软触足迹（会员）。
	Detail(ctx context.Context, spuId int64, viewerUserId int64) (*ProductDetail, error)
	Search(ctx context.Context, keyword string, q ProductQuery) (*PageResult[ProductCard], error)
	ProductReviews(ctx context.Context, spuId int64, score int, page PageQuery) (*PageResult[ReviewCard], error)

	// ---- 后台管理 ----
	AdminCategoryCreate(ctx context.Context, in CategoryInput) (int64, error)
	AdminCategoryUpdate(ctx context.Context, id int64, in CategoryInput) error
	AdminCategoryDelete(ctx context.Context, id int64) error // 有商品引用拒绝
	AdminBrandCreate(ctx context.Context, in BrandInput) (int64, error)
	AdminBrandUpdate(ctx context.Context, id int64, in BrandInput) error
	AdminBrandDelete(ctx context.Context, id int64) error
	AdminProductList(ctx context.Context, q AdminProductQuery) (*PageResult[AdminProductItem], error)
	AdminProductCreate(ctx context.Context, in SpuInput) (spuId int64, spuNo string, err error)
	AdminProductDetail(ctx context.Context, spuId int64) (*AdminProductDetail, error) // 含成本价
	AdminProductUpdate(ctx context.Context, spuId int64, in SpuInput) error
	AdminProductDelete(ctx context.Context, spuId int64) error
	// AdminProductStatus 上下架（无启用 SKU 禁上架, 30005）。
	AdminProductStatus(ctx context.Context, spuId int64, status int) error
	AdminProductRestrict(ctx context.Context, spuId int64, codes []string) error
	AdminSkuCreate(ctx context.Context, spuId int64, in SkuInput) (int64, error) // 规格组合冲突 30006
	AdminSkuUpdate(ctx context.Context, skuId int64, in SkuInput) error
	AdminSkuDelete(ctx context.Context, skuId int64) error
	AdminSkuStatus(ctx context.Context, skuId int64, status int) error
}

type CategoryNode struct {
	Id       int64          `json:"id"`
	Name     string         `json:"name"`
	Icon     string         `json:"icon"`
	Children []CategoryNode `json:"children"`
}

type CategoryInput struct {
	ParentId int64
	Name     string
	Icon     string
	Level    int
	Sort     int
}

type BrandItem struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type BrandInput struct {
	Name        string
	Logo        string
	Description string
	Sort        int
}

type ProductQuery struct {
	CategoryId int64
	BrandId    int64
	Sort       int // 0综合 1销量 2价格 3上新
	PriceMin   string
	PriceMax   string
	PageQuery
}

type ProductCard struct {
	SpuId      int64  `json:"spuId"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	PriceRange string `json:"priceRange"`
	SaleCount  int    `json:"saleCount"`
}

type ProductDetail struct {
	SpuId           int64             `json:"spuId"`
	SpuNo           string            `json:"spuNo"`
	Name            string            `json:"name"`
	SubTitle        string            `json:"subTitle"`
	Images          []string          `json:"images"`
	VideoUrl        string            `json:"videoUrl"`
	Description     string            `json:"description"`
	SpecDefinitions []map[string]any  `json:"specDefinitions"`
	Attributes      map[string]string `json:"attributes"`
	FreightSummary  string            `json:"freightSummary"`
	Skus            []SkuCard         `json:"skus"`
	ReviewSummary   ReviewSummary     `json:"reviewSummary"`
}

type SkuCard struct {
	SkuId     int64             `json:"skuId"`
	SkuNo     string            `json:"skuNo"`
	Specs     map[string]string `json:"specs"`
	Price     string            `json:"price"`
	LinePrice string            `json:"linePrice"`
	Sellable  bool              `json:"sellable"`
}

type ReviewSummary struct {
	Avg          string         `json:"avg"`
	Distribution map[string]int `json:"distribution"`
	Total        int            `json:"total"`
}

type ReviewCard struct {
	ReviewId  int64             `json:"reviewId"`
	User      string            `json:"user" dc:"匿名脱敏"`
	Score     int               `json:"score"`
	Content   string            `json:"content"`
	Images    []string          `json:"images"`
	Specs     map[string]string `json:"specs"`
	Reply     string            `json:"reply"`
	Extra     string            `json:"extra"`
	CreatedAt string            `json:"createdAt"`
}

type AdminProductQuery struct {
	Status     int
	CategoryId int64
	Keyword    string
	PageQuery
}

type AdminProductItem struct {
	SpuId      int64  `json:"spuId"`
	SpuNo      string `json:"spuNo"`
	Name       string `json:"name"`
	CategoryId int64  `json:"categoryId"`
	Status     int    `json:"status"`
	SaleCount  int    `json:"saleCount"`
	CreatedAt  string `json:"createdAt"`
}

type SpuInput struct {
	Name              string
	SubTitle          string
	CategoryId        int64
	BrandId           int64
	FreightTemplateId int64
	Images            []string
	VideoUrl          string
	Description       string
	SpecDefinitions   []map[string]any
	Attributes        map[string]string
}

type AdminProductDetail struct {
	SpuId             int64             `json:"spuId"`
	SpuNo             string            `json:"spuNo"`
	Name              string            `json:"name"`
	CategoryId        int64             `json:"categoryId"`
	BrandId           int64             `json:"brandId"`
	FreightTemplateId int64             `json:"freightTemplateId"`
	Images            []string          `json:"images"`
	Description       string            `json:"description"`
	SpecDefinitions   []map[string]any  `json:"specDefinitions"`
	Attributes        map[string]string `json:"attributes"`
	SaleRestrictCodes []string          `json:"saleRestrictCodes"`
	Status            int               `json:"status"`
	Skus              []AdminSkuDetail  `json:"skus"`
}

type AdminSkuDetail struct {
	SkuId     int64             `json:"skuId"`
	SkuNo     string            `json:"skuNo"`
	Specs     map[string]string `json:"specs"`
	Price     string            `json:"price"`
	LinePrice string            `json:"linePrice"`
	CostPrice string            `json:"costPrice" dc:"管理可见"`
	Status    int               `json:"status"`
}

type SkuInput struct {
	Specs     map[string]string
	Price     string
	LinePrice string
	CostPrice string
	Image     string
	Weight    string
	Barcode   string
}
