package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 三级分类树（禁用不返回）
	CategoryTreeReq struct {
		g.Meta `path:"/categories" method:"GET" summary:"商品分类树"`
	}
	CategoryNode struct {
		Id       string         `json:"id" dc:"分类ID"`
		Name     string         `json:"name" dc:"分类名称"`
		Icon     string         `json:"icon" dc:"图标"`
		Children []CategoryNode `json:"children" dc:"子分类"`
	}
	CategoryTreeRes struct {
		Tree []CategoryNode `json:"tree"`
	}

	// 品牌列表（启用）
	BrandListReq struct {
		g.Meta `path:"/brands" method:"GET" summary:"品牌列表"`
		PageReq
	}
	BrandItem struct {
		Id   string `json:"id"`
		Name string `json:"name" dc:"品牌名称"`
		Logo string `json:"logo" dc:"Logo"`
	}
	BrandListRes struct {
		PageRes
		List []BrandItem `json:"list"`
	}

	// 商品列表（在售; 筛选与排序）
	ProductListReq struct {
		g.Meta     `path:"/products" method:"GET" summary:"商品列表"`
		CategoryId string `json:"categoryId" dc:"三级分类ID"`
		BrandId    string `json:"brandId" dc:"品牌ID"`
		Sort       int    `json:"sort" dc:"排序:0综合 1销量 2价格 3上新" d:"0"`
		PriceMin   string `json:"priceMin" dc:"价格区间下限(元)"`
		PriceMax   string `json:"priceMax" dc:"价格区间上限(元)"`
		PageReq
	}
	ProductItem struct {
		SpuId      string `json:"spuId"`
		SpuName    string `json:"spuName" dc:"商品名称"`
		Image      string `json:"image" dc:"主图"`
		PriceRange string `json:"priceRange" dc:"价格区间(如 99.00-129.00)"`
		SaleCount  int    `json:"saleCount" dc:"销量"`
	}
	ProductListRes struct {
		PageRes
		List []ProductItem `json:"list"`
	}

	// 商品详情（SKU/规格/运费概要/评价汇总; 不含成本价）
	ProductDetailReq struct {
		g.Meta `path:"/products/{spuId}" method:"GET" summary:"商品详情"`
		SpuId  string `json:"spuId" v:"required" dc:"SPU ID"`
	}
	SkuItem struct {
		SkuId     string            `json:"skuId"`
		SkuNo     string            `json:"skuNo" dc:"SKU编码"`
		Specs     map[string]string `json:"specs" dc:"规格组合"`
		Price     string            `json:"price" dc:"售价(元)"`
		LinePrice string            `json:"linePrice" dc:"划线价"`
		Sellable  bool              `json:"sellable" dc:"是否可售(SPU上架+SKU启用+有库存)"`
	}
	ReviewSummary struct {
		Avg          string         `json:"avg" dc:"平均分"`
		Distribution map[string]int `json:"distribution" dc:"星级分布{1:n..5:n}"`
		Total        int            `json:"total" dc:"评价总数"`
	}
	ProductDetailRes struct {
		SpuId           string            `json:"spuId"`
		SpuNo           string            `json:"spuNo" dc:"SPU编码"`
		Name            string            `json:"name" dc:"商品名称"`
		SubTitle        string            `json:"subTitle" dc:"副标题"`
		Images          []string          `json:"images" dc:"图集"`
		VideoUrl        string            `json:"videoUrl" dc:"主图视频"`
		Description     string            `json:"description" dc:"图文详情"`
		SpecDefinitions []map[string]any  `json:"specDefinitions" dc:"规格定义"`
		Attributes      map[string]string `json:"attributes" dc:"非销售属性"`
		FreightSummary  string            `json:"freightSummary" dc:"运费概要描述"`
		Skus            []SkuItem         `json:"skus"`
		ReviewSummary   ReviewSummary     `json:"reviewSummary"`
	}

	// 关键词搜索
	ProductSearchReq struct {
		g.Meta     `path:"/search" method:"GET" summary:"商品搜索"`
		Keyword    string `json:"keyword" v:"required" dc:"关键词"`
		CategoryId string `json:"categoryId" dc:"分类筛选"`
		BrandId    string `json:"brandId" dc:"品牌筛选"`
		Sort       int    `json:"sort" dc:"排序:0综合 1销量 2价格 3上新" d:"0"`
		PageReq
	}
	ProductSearchRes struct {
		PageRes
		List []ProductItem `json:"list"`
	}

	// 商品评价列表（审核通过）
	ProductReviewListReq struct {
		g.Meta `path:"/products/{spuId}/reviews" method:"GET" summary:"商品评价列表"`
		SpuId  string `json:"spuId" v:"required" dc:"SPU ID"`
		Score  int    `json:"score" dc:"星级筛选1-5"`
		PageReq
	}
	ReviewItem struct {
		ReviewId  string            `json:"reviewId"`
		User      string            `json:"user" dc:"评价人(匿名则脱敏)"`
		Score     int               `json:"score" dc:"评分"`
		Content   string            `json:"content" dc:"评价内容"`
		Images    []string          `json:"images" dc:"评价图片"`
		Specs     map[string]string `json:"specs" dc:"购买规格"`
		Reply     string            `json:"reply" dc:"商家回复"`
		Extra     string            `json:"extra" dc:"追评内容"`
		CreatedAt string            `json:"createdAt" dc:"评价时间"`
	}
	ProductReviewListRes struct {
		PageRes
		Summary ReviewSummary `json:"summary" dc:"评价汇总"`
		List    []ReviewItem  `json:"list"`
	}
)
