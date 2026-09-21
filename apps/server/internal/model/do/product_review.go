// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ProductReview is the golang structure of table product_review for DAO operations like Where/Data.
type ProductReview struct {
	g.Meta       `orm:"table:product_review, do:true"`
	Id           any         // 评价ID
	OrderItemId  any         // 订单项ID(一项一评)
	OrderNo      any         // 订单号(冗余,免联查)
	UserId       any         // 评价用户ID
	SpuId        any         // SPU ID(聚合统计用)
	SkuId        any         // SKU ID
	SpuName      any         // SPU名称(下单时快照,商品软删后自洽展示)
	SkuSpecs     any         // 规格组合(下单时快照)
	Score        any         // 评分:1-5星
	Content      any         // 评价内容
	Images       any         // 评价图片URL数组
	IsAnonymous  any         // 匿名:0否 1是
	AuditStatus  any         // 审核状态:0待审核 1通过 2驳回
	ExtraContent any         // 追评内容(仅一次,90天内)
	ExtraTime    *gtime.Time // 追评时间
	ReplyContent any         // 商家回复(仅一次)
	ReplyTime    *gtime.Time // 回复时间
	Deleted      any         // 软删除:0否 1是(用户删除自己的评价)
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
}
