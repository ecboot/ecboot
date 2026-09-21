// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ProductSpu is the golang structure of table product_spu for DAO operations like Where/Data.
type ProductSpu struct {
	g.Meta            `orm:"table:product_spu, do:true"`
	Id                any         // SPU ID
	SpuNo             any         // SPU业务编码(全局唯一,对外展示用)
	Name              any         // 商品名称(SPU级)
	SubTitle          any         // 副标题/卖点
	CategoryId        any         // 所属分类ID(三级分类)
	BrandId           any         // 品牌ID,可为空(无品牌商品)
	FreightTemplateId any         // 运费模板ID(NULL=包邮)
	Images            any         // 主图+轮播图URL数组,按序存储
	Description       any         // 图文详情(富文本)
	VideoUrl          any         // 主图视频URL(预留)
	SpecDefinitions   any         // 规格定义:[{"name":"颜色","values":["黑","白"]}]
	Attributes        any         // 非销售属性键值对(材质/产地等)
	SaleRestrictCodes any         // 禁售省级区划代码列表(黑名单,NULL=全国可售;下单时地址省代码∈列表即拒绝;列表页过滤由ES承载)
	Status            any         // 上架状态:0下架 1上架(默认下架)
	SaleCount         any         // 累计销量(异步冗余,支付成功事件累加,排序用)
	PriceMin          any         // 启用SKU最低价(冗余,SKU变更时刷新;全部禁用为NULL)
	PriceMax          any         // 启用SKU最高价(冗余)
	Deleted           any         // 软删除:0否 1是(删除后历史订单不受影响,依赖订单项快照)
	CreatedAt         *gtime.Time // 创建时间
	UpdatedAt         *gtime.Time // 更新时间
}
