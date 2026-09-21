// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// FreightTemplate is the golang structure of table freight_template for DAO operations like Where/Data.
type FreightTemplate struct {
	g.Meta           `orm:"table:freight_template, do:true"`
	Id               any         // 模板ID
	Name             any         // 模板名称(如:华东包邮-全国2件10元)
	ChargeType       any         // 计费方式:1按件数 2按重量
	FreeThreshold    any         // 满额包邮阈值(订单商品金额≥此值免运费,NULL=不包邮)
	FreeExcludeCodes any         // 不参与满额包邮的省级区划代码列表(如["650000","540000"]);NULL=全国均可包邮
	Status           any         // 状态:1启用 0停用(停用后新订单不可用,已下单不受影响)
	Deleted          any         // 软删除:0否 1是
	CreatedAt        *gtime.Time // 创建时间
	UpdatedAt        *gtime.Time // 更新时间
}
