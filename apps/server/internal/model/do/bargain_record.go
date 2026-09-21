// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// BargainRecord is the golang structure of table bargain_record for DAO operations like Where/Data.
type BargainRecord struct {
	g.Meta       `orm:"table:bargain_record, do:true"`
	Id           any         // 砍价单ID
	BargainNo    any         // 砍价单号(全局唯一)
	ItemId       any         // 砍价场次商品ID
	UserId       any         // 发起人用户ID
	CurrentPrice any         // 当前价(逐刀递减,>=floor_price由应用+条件更新保证)
	CutCount     any         // 已砍刀数
	Status       any         // 状态:1砍价中 2到底价待下单 3已下单 4超时失败 5已取消
	ExpireTime   *gtime.Time // 砍价截止时间(超时扫描)
	SuccessTime  *gtime.Time // 到底价时间
	OrderNo      any         // 成交订单号(下单后回填,防重复成交)
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
}
