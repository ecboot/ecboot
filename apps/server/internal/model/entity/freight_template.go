// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// FreightTemplate is the golang structure for table freight_template.
type FreightTemplate struct {
	Id               uint64      `json:"id"               orm:"id"                 ` // 模板ID
	Name             string      `json:"name"             orm:"name"               ` // 模板名称(如:华东包邮-全国2件10元)
	ChargeType       int         `json:"chargeType"       orm:"charge_type"        ` // 计费方式:1按件数 2按重量
	FreeThreshold    float64     `json:"freeThreshold"    orm:"free_threshold"     ` // 满额包邮阈值(订单商品金额≥此值免运费,NULL=不包邮)
	FreeExcludeCodes string      `json:"freeExcludeCodes" orm:"free_exclude_codes" ` // 不参与满额包邮的省级区划代码列表(如["650000","540000"]);NULL=全国均可包邮
	Status           int         `json:"status"           orm:"status"             ` // 状态:1启用 0停用(停用后新订单不可用,已下单不受影响)
	Deleted          int         `json:"deleted"          orm:"deleted"            ` // 软删除:0否 1是
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"         ` // 创建时间
	UpdatedAt        *gtime.Time `json:"updatedAt"        orm:"updated_at"         ` // 更新时间
}
