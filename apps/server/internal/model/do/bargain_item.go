// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// BargainItem is the golang structure of table bargain_item for DAO operations like Where/Data.
type BargainItem struct {
	g.Meta        `orm:"table:bargain_item, do:true"`
	Id            any         // 砍价场次商品ID
	ActivityId    any         // 活动ID
	SkuId         any         // SKU ID
	OriginalPrice any         // 起始价(=发起时售价口径)
	FloorPrice    any         // 底价(砍到底价即可下单;floor<=original应用校验)
	MaxCutCount   any         // 最大帮砍刀数(0=不限,金额收敛到底价)
	Config        any         // 玩法参数(随机递减算法边界/首刀上限等,应用层约定)
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
}
