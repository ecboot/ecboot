package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/model"
)

type (
	// 分类树（含禁用）
	AdminCategoryTreeReq struct {
		g.Meta `path:"/categories" method:"GET" summary:"分类树"`
	}
	AdminCategoryNode struct {
		Id       string              `json:"id"`
		ParentId string              `json:"parentId" dc:"父ID,0为根"`
		Name     string              `json:"name"`
		Icon     string              `json:"icon"`
		Level    int                 `json:"level" dc:"1/2/3"`
		Sort     int                 `json:"sort"`
		Status   int                 `json:"status" dc:"1启用 0禁用"`
		Children []AdminCategoryNode `json:"children"`
	}
	AdminCategoryTreeRes struct {
		Tree []AdminCategoryNode `json:"tree"`
	}

	// 权限: product:category:create
	AdminCategoryCreateReq struct {
		g.Meta   `path:"/categories" method:"POST" summary:"新增分类"`
		ParentId string `json:"parentId" dc:"父ID,0为根" d:"0"`
		Name     string `json:"name" v:"required" dc:"名称"`
		Icon     string `json:"icon" dc:"图标"`
		Level    int    `json:"level" v:"required|in:1,2,3" dc:"层级"`
		Sort     int    `json:"sort" dc:"排序"`
	}
	AdminCategoryCreateRes struct {
		Id string `json:"id"`
	}

	// 权限: product:category:update
	AdminCategoryUpdateReq struct {
		g.Meta `path:"/categories/{id}" method:"PUT" summary:"修改分类"`
		Id     string `json:"id" v:"required" dc:"分类ID"`
		Name   string `json:"name" dc:"名称"`
		Icon   string `json:"icon" dc:"图标"`
		Sort   int    `json:"sort" dc:"排序"`
		Status int    `json:"status" dc:"状态"`
	}
	AdminCategoryUpdateRes struct {
		Success bool `json:"success"`
	}

	// 权限: product:category:delete
	AdminCategoryDeleteReq struct {
		g.Meta `path:"/categories/{id}" method:"DELETE" summary:"删除分类(软删,有商品禁删)"`
		Id     string `json:"id" v:"required" dc:"分类ID"`
	}
	AdminCategoryDeleteRes struct {
		Success bool `json:"success"`
	}

	// 品牌列表
	AdminBrandListReq struct {
		g.Meta `path:"/brands" method:"GET" summary:"品牌列表"`
		Status int `json:"status" dc:"状态筛选"`
		model.PageReq
	}
	AdminBrandItem struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Logo        string `json:"logo"`
		Description string `json:"description"`
		Sort        int    `json:"sort"`
		Status      int    `json:"status"`
	}
	AdminBrandListRes struct {
		model.PageRes
		List []AdminBrandItem `json:"list"`
	}

	// 权限: product:brand:create
	AdminBrandCreateReq struct {
		g.Meta      `path:"/brands" method:"POST" summary:"新增品牌"`
		Name        string `json:"name" v:"required" dc:"名称"`
		Logo        string `json:"logo" dc:"Logo"`
		Description string `json:"description" dc:"简介"`
		Sort        int    `json:"sort" dc:"排序"`
	}
	AdminBrandCreateRes struct {
		Id string `json:"id"`
	}

	// 权限: product:brand:update
	AdminBrandUpdateReq struct {
		g.Meta      `path:"/brands/{id}" method:"PUT" summary:"修改品牌"`
		Id          string `json:"id" v:"required" dc:"品牌ID"`
		Name        string `json:"name" dc:"名称"`
		Logo        string `json:"logo" dc:"Logo"`
		Description string `json:"description" dc:"简介"`
		Sort        int    `json:"sort" dc:"排序"`
		Status      int    `json:"status" dc:"状态"`
	}
	AdminBrandUpdateRes struct {
		Success bool `json:"success"`
	}

	// 权限: product:brand:delete
	AdminBrandDeleteReq struct {
		g.Meta `path:"/brands/{id}" method:"DELETE" summary:"删除品牌(软删)"`
		Id     string `json:"id" v:"required" dc:"品牌ID"`
	}
	AdminBrandDeleteRes struct {
		Success bool `json:"success"`
	}

	// SPU 列表（全状态）
	AdminSpuListReq struct {
		g.Meta     `path:"/products" method:"GET" summary:"商品SPU列表"`
		Status     int    `json:"status" dc:"0下架 1上架(空=全部)"`
		CategoryId string `json:"categoryId" dc:"分类筛选"`
		Keyword    string `json:"keyword" dc:"名称/编码"`
		model.PageReq
	}
	AdminSpuItem struct {
		SpuId      string `json:"spuId"`
		SpuNo      string `json:"spuNo"`
		Name       string `json:"name"`
		CategoryId string `json:"categoryId"`
		BrandId    string `json:"brandId"`
		Status     int    `json:"status" dc:"0下架 1上架"`
		SaleCount  int    `json:"saleCount"`
		CreatedAt  string `json:"createdAt"`
	}
	AdminSpuListRes struct {
		model.PageRes
		List []AdminSpuItem `json:"list"`
	}

	// 创建 SPU（含规格定义/图文/运费模板）
	// 权限: product:spu:create
	AdminSpuCreateReq struct {
		g.Meta            `path:"/products" method:"POST" summary:"创建商品SPU"`
		Name              string            `json:"name" v:"required" dc:"商品名称"`
		SubTitle          string            `json:"subTitle" dc:"副标题"`
		CategoryId        string            `json:"categoryId" v:"required" dc:"三级分类"`
		BrandId           string            `json:"brandId" dc:"品牌"`
		FreightTemplateId string            `json:"freightTemplateId" dc:"运费模板(空=包邮)"`
		Images            []string          `json:"images" v:"required" dc:"图集"`
		VideoUrl          string            `json:"videoUrl" dc:"视频"`
		Description       string            `json:"description" dc:"图文详情"`
		SpecDefinitions   []map[string]any  `json:"specDefinitions" dc:"规格定义"`
		Attributes        map[string]string `json:"attributes" dc:"非销售属性"`
	}
	AdminSpuCreateRes struct {
		SpuId string `json:"spuId"`
		SpuNo string `json:"spuNo" dc:"SPU编码"`
	}

	AdminSpuDetailReq struct {
		g.Meta `path:"/products/{spuId}" method:"GET" summary:"商品详情(管理视角,含成本价)"`
		SpuId  string `json:"spuId" v:"required" dc:"SPU ID"`
	}
	AdminSkuAdminItem struct {
		SkuId     string            `json:"skuId"`
		SkuNo     string            `json:"skuNo"`
		Specs     map[string]string `json:"specs"`
		Price     string            `json:"price"`
		LinePrice string            `json:"linePrice"`
		CostPrice string            `json:"costPrice" dc:"成本价(管理可见)"`
		Weight    string            `json:"weight" dc:"重量克"`
		Barcode   string            `json:"barcode"`
		Status    int               `json:"status" dc:"1启用 0禁用"`
	}
	AdminSpuDetailRes struct {
		SpuId             string              `json:"spuId"`
		SpuNo             string              `json:"spuNo"`
		Name              string              `json:"name"`
		SubTitle          string              `json:"subTitle"`
		CategoryId        string              `json:"categoryId"`
		BrandId           string              `json:"brandId"`
		FreightTemplateId string              `json:"freightTemplateId"`
		Images            []string            `json:"images"`
		VideoUrl          string              `json:"videoUrl"`
		Description       string              `json:"description"`
		SpecDefinitions   []map[string]any    `json:"specDefinitions"`
		Attributes        map[string]string   `json:"attributes"`
		SaleRestrictCodes []string            `json:"saleRestrictCodes" dc:"限售省级代码(空=不限)"`
		Status            int                 `json:"status"`
		Skus              []AdminSkuAdminItem `json:"skus"`
	}

	// 权限: product:spu:update
	AdminSpuUpdateReq struct {
		g.Meta            `path:"/products/{spuId}" method:"PUT" summary:"修改商品SPU"`
		SpuId             string            `json:"spuId" v:"required" dc:"SPU ID"`
		Name              string            `json:"name" dc:"名称"`
		SubTitle          string            `json:"subTitle" dc:"副标题"`
		CategoryId        string            `json:"categoryId" dc:"分类"`
		BrandId           string            `json:"brandId" dc:"品牌"`
		FreightTemplateId string            `json:"freightTemplateId" dc:"运费模板"`
		Images            []string          `json:"images" dc:"图集"`
		VideoUrl          string            `json:"videoUrl" dc:"视频"`
		Description       string            `json:"description" dc:"图文详情"`
		SpecDefinitions   []map[string]any  `json:"specDefinitions" dc:"规格定义"`
		Attributes        map[string]string `json:"attributes" dc:"属性"`
	}
	AdminSpuUpdateRes struct {
		Success bool `json:"success"`
	}

	// 权限: product:spu:delete
	AdminSpuDeleteReq struct {
		g.Meta `path:"/products/{spuId}" method:"DELETE" summary:"删除商品(软删)"`
		SpuId  string `json:"spuId" v:"required" dc:"SPU ID"`
	}
	AdminSpuDeleteRes struct {
		Success bool `json:"success"`
	}

	// SPU 上下架（无启用 SKU 禁上架）
	// 权限: product:spu:update
	AdminSpuStatusReq struct {
		g.Meta `path:"/products/{spuId}/status" method:"POST" summary:"商品上下架"`
		SpuId  string `json:"spuId" v:"required" dc:"SPU ID"`
		Status int    `json:"status" v:"required|in:0,1" dc:"0下架 1上架"`
	}
	AdminSpuStatusRes struct {
		Success bool `json:"success"`
	}

	// 限售区域设置
	// 权限: product:spu:update
	AdminSpuRestrictReq struct {
		g.Meta            `path:"/products/{spuId}/restrict-codes" method:"PUT" summary:"限售区域设置"`
		SpuId             string   `json:"spuId" v:"required" dc:"SPU ID"`
		SaleRestrictCodes []string `json:"saleRestrictCodes" dc:"禁售省级代码列表(空数组=不限售)"`
	}
	AdminSpuRestrictRes struct {
		Success bool `json:"success"`
	}

	// 新增 SKU
	// 权限: product:sku:create
	AdminSkuCreateReq struct {
		g.Meta    `path:"/products/{spuId}/skus" method:"POST" summary:"新增SKU"`
		SpuId     string            `json:"spuId" v:"required" dc:"SPU ID"`
		Specs     map[string]string `json:"specs" v:"required" dc:"规格组合"`
		Price     string            `json:"price" v:"required" dc:"售价"`
		LinePrice string            `json:"linePrice" dc:"划线价"`
		CostPrice string            `json:"costPrice" dc:"成本价"`
		Image     string            `json:"image" dc:"SKU图"`
		Weight    string            `json:"weight" dc:"重量克"`
		Barcode   string            `json:"barcode" dc:"条码"`
	}
	AdminSkuCreateRes struct {
		SkuId string `json:"skuId"`
		SkuNo string `json:"skuNo" dc:"SKU编码"`
	}

	// 权限: product:sku:update
	AdminSkuUpdateReq struct {
		g.Meta    `path:"/skus/{skuId}" method:"PUT" summary:"修改SKU"`
		SkuId     string            `json:"skuId" v:"required" dc:"SKU ID"`
		Specs     map[string]string `json:"specs" dc:"规格组合"`
		Price     string            `json:"price" dc:"售价"`
		LinePrice string            `json:"linePrice" dc:"划线价"`
		CostPrice string            `json:"costPrice" dc:"成本价"`
		Image     string            `json:"image" dc:"SKU图"`
		Weight    string            `json:"weight" dc:"重量克"`
		Barcode   string            `json:"barcode" dc:"条码"`
	}
	AdminSkuUpdateRes struct {
		Success bool `json:"success"`
	}

	// 权限: product:sku:delete
	AdminSkuDeleteReq struct {
		g.Meta `path:"/skus/{skuId}" method:"DELETE" summary:"删除SKU(软删)"`
		SkuId  string `json:"skuId" v:"required" dc:"SKU ID"`
	}
	AdminSkuDeleteRes struct {
		Success bool `json:"success"`
	}

	// 权限: product:sku:update
	AdminSkuStatusReq struct {
		g.Meta `path:"/skus/{skuId}/status" method:"POST" summary:"SKU启停"`
		SkuId  string `json:"skuId" v:"required" dc:"SKU ID"`
		Status int    `json:"status" v:"required|in:0,1" dc:"1启用 0禁用"`
	}
	AdminSkuStatusRes struct {
		Success bool `json:"success"`
	}
)
