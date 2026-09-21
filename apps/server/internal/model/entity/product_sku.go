// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ProductSku is the golang structure for table product_sku.
type ProductSku struct {
	Id        uint64      `json:"id"        orm:"id"         ` // SKU ID
	SkuNo     string      `json:"skuNo"     orm:"sku_no"     ` // SKU业务编码(全局唯一)
	SpuId     uint64      `json:"spuId"     orm:"spu_id"     ` // 所属SPU ID
	Name      string      `json:"name"      orm:"name"       ` // SKU名称=SPU名+规格串(冗余生成,展示用)
	Specs     string      `json:"specs"     orm:"specs"      ` // 规格组合快照:{"颜色":"黑","尺码":"M"}(同SPU内组合唯一由应用层保证)
	SpecsHash string      `json:"specsHash" orm:"specs_hash" ` // 规格组合哈希(同SPU内组合唯一;应用序列化稳定时生效)
	Price     float64     `json:"price"     orm:"price"      ` // 现售价(元)
	LinePrice float64     `json:"linePrice" orm:"line_price" ` // 划线价/原价(营销展示,可为空)
	CostPrice float64     `json:"costPrice" orm:"cost_price" ` // 成本价(后台可见,毛利分析预留)
	Image     string      `json:"image"     orm:"image"      ` // SKU图URL(空则用SPU主图)
	Weight    float64     `json:"weight"    orm:"weight"     ` // 重量(克,运费按重计费时使用)
	Barcode   string      `json:"barcode"   orm:"barcode"    ` // 条形码(预留)
	Sort      int         `json:"sort"      orm:"sort"       ` // 展示排序
	Status    int         `json:"status"    orm:"status"     ` // 状态:1启用 0禁用(断码/失效单品)
	Deleted   int         `json:"deleted"   orm:"deleted"    ` // 软删除:0否 1是
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` // 更新时间
}
