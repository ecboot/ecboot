// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ProductReview is the golang structure for table product_review.
type ProductReview struct {
	Id           uint64      `json:"id"           orm:"id"            ` // 评价ID
	OrderItemId  uint64      `json:"orderItemId"  orm:"order_item_id" ` // 订单项ID(一项一评)
	OrderNo      string      `json:"orderNo"      orm:"order_no"      ` // 订单号(冗余,免联查)
	UserId       uint64      `json:"userId"       orm:"user_id"       ` // 评价用户ID
	SpuId        uint64      `json:"spuId"        orm:"spu_id"        ` // SPU ID(聚合统计用)
	SkuId        uint64      `json:"skuId"        orm:"sku_id"        ` // SKU ID
	SpuName      string      `json:"spuName"      orm:"spu_name"      ` // SPU名称(下单时快照,商品软删后自洽展示)
	SkuSpecs     string      `json:"skuSpecs"     orm:"sku_specs"     ` // 规格组合(下单时快照)
	Score        int         `json:"score"        orm:"score"         ` // 评分:1-5星
	Content      string      `json:"content"      orm:"content"       ` // 评价内容
	Images       string      `json:"images"       orm:"images"        ` // 评价图片URL数组
	IsAnonymous  int         `json:"isAnonymous"  orm:"is_anonymous"  ` // 匿名:0否 1是
	AuditStatus  int         `json:"auditStatus"  orm:"audit_status"  ` // 审核状态:0待审核 1通过 2驳回
	ExtraContent string      `json:"extraContent" orm:"extra_content" ` // 追评内容(仅一次,90天内)
	ExtraTime    *gtime.Time `json:"extraTime"    orm:"extra_time"    ` // 追评时间
	ReplyContent string      `json:"replyContent" orm:"reply_content" ` // 商家回复(仅一次)
	ReplyTime    *gtime.Time `json:"replyTime"    orm:"reply_time"    ` // 回复时间
	Deleted      int         `json:"deleted"      orm:"deleted"       ` // 软删除:0否 1是(用户删除自己的评价)
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` // 创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` // 更新时间
}
