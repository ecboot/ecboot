// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ProductSpu is the golang structure for table product_spu.
type ProductSpu struct {
	Id                uint64      `json:"id"                orm:"id"                  ` // SPU ID
	SpuNo             string      `json:"spuNo"             orm:"spu_no"              ` // SPU业务编码(全局唯一,对外展示用)
	Name              string      `json:"name"              orm:"name"                ` // 商品名称(SPU级)
	SubTitle          string      `json:"subTitle"          orm:"sub_title"           ` // 副标题/卖点
	CategoryId        uint64      `json:"categoryId"        orm:"category_id"         ` // 所属分类ID(三级分类)
	BrandId           uint64      `json:"brandId"           orm:"brand_id"            ` // 品牌ID,可为空(无品牌商品)
	SellerId          uint64      `json:"sellerId"          orm:"seller_id"           ` // 所属商家(1=自营;多商家预留)
	FreightTemplateId uint64      `json:"freightTemplateId" orm:"freight_template_id" ` // 运费模板ID(NULL=包邮)
	Images            string      `json:"images"            orm:"images"              ` // 主图+轮播图URL数组,按序存储
	Description       string      `json:"description"       orm:"description"         ` // 图文详情(富文本)
	VideoUrl          string      `json:"videoUrl"          orm:"video_url"           ` // 主图视频URL(预留)
	SpecDefinitions   string      `json:"specDefinitions"   orm:"spec_definitions"    ` // 规格定义:[{"name":"颜色","values":["黑","白"]}]
	Attributes        string      `json:"attributes"        orm:"attributes"          ` // 非销售属性键值对(材质/产地等)
	SaleRestrictCodes string      `json:"saleRestrictCodes" orm:"sale_restrict_codes" ` // 禁售省级区划代码列表(黑名单,NULL=全国可售;下单时地址省代码∈列表即拒绝;列表页过滤由ES承载)
	Status            int         `json:"status"            orm:"status"              ` // 上架状态:0下架 1上架(默认下架)
	SaleCount         uint        `json:"saleCount"         orm:"sale_count"          ` // 累计销量(异步冗余,支付成功事件累加,排序用)
	PriceMin          float64     `json:"priceMin"          orm:"price_min"           ` // 启用SKU最低价(冗余,SKU变更时刷新;全部禁用为NULL)
	PriceMax          float64     `json:"priceMax"          orm:"price_max"           ` // 启用SKU最高价(冗余)
	Deleted           int         `json:"deleted"           orm:"deleted"             ` // 软删除:0否 1是(删除后历史订单不受影响,依赖订单项快照)
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"          ` // 创建时间
	UpdatedAt         *gtime.Time `json:"updatedAt"         orm:"updated_at"          ` // 更新时间
}
