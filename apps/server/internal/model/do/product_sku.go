// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ProductSku is the golang structure of table product_sku for DAO operations like Where/Data.
type ProductSku struct {
	g.Meta    `orm:"table:product_sku, do:true"`
	Id        any         // SKU ID
	SkuNo     any         // SKU业务编码(全局唯一)
	SpuId     any         // 所属SPU ID
	Name      any         // SKU名称=SPU名+规格串(冗余生成,展示用)
	Specs     any         // 规格组合快照:{"颜色":"黑","尺码":"M"}(同SPU内组合唯一由应用层保证)
	SpecsHash any         // 规格组合哈希(同SPU内组合唯一;应用序列化稳定时生效)
	Price     any         // 现售价(元)
	LinePrice any         // 划线价/原价(营销展示,可为空)
	CostPrice any         // 成本价(后台可见,毛利分析预留)
	Image     any         // SKU图URL(空则用SPU主图)
	Weight    any         // 重量(克,运费按重计费时使用)
	Barcode   any         // 条形码(预留)
	Sort      any         // 展示排序
	Status    any         // 状态:1启用 0禁用(断码/失效单品)
	Deleted   any         // 软删除:0否 1是
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
