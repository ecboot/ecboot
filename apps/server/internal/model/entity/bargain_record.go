// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// BargainRecord is the golang structure for table bargain_record.
type BargainRecord struct {
	Id           uint64      `json:"id"           orm:"id"            ` // 砍价单ID
	BargainNo    string      `json:"bargainNo"    orm:"bargain_no"    ` // 砍价单号(全局唯一)
	ItemId       uint64      `json:"itemId"       orm:"item_id"       ` // 砍价场次商品ID
	UserId       uint64      `json:"userId"       orm:"user_id"       ` // 发起人用户ID
	CurrentPrice float64     `json:"currentPrice" orm:"current_price" ` // 当前价(逐刀递减,>=floor_price由应用+条件更新保证)
	CutCount     uint        `json:"cutCount"     orm:"cut_count"     ` // 已砍刀数
	Status       int         `json:"status"       orm:"status"        ` // 状态:1砍价中 2到底价待下单 3已下单 4超时失败 5已取消
	ExpireTime   *gtime.Time `json:"expireTime"   orm:"expire_time"   ` // 砍价截止时间(超时扫描)
	SuccessTime  *gtime.Time `json:"successTime"  orm:"success_time"  ` // 到底价时间
	OrderNo      string      `json:"orderNo"      orm:"order_no"      ` // 成交订单号(下单后回填; NULL=未下单, NULL不参与唯一约束)
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` // 创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` // 更新时间
}
