// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// FreightRule is the golang structure for table freight_rule.
type FreightRule struct {
	Id           uint64      `json:"id"           orm:"id"            ` // 规则ID
	TemplateId   uint64      `json:"templateId"   orm:"template_id"   ` // 所属模板ID
	RegionCodes  string      `json:"regionCodes"  orm:"region_codes"  ` // 适用省级行政区代码列表(如["310000","330000"];空数组=全国兜底规则,每模板至多一条)
	FirstUnit    uint        `json:"firstUnit"    orm:"first_unit"    ` // 首段额度:按件=首件数;按重=首重克数
	FirstFee     float64     `json:"firstFee"     orm:"first_fee"     ` // 首段运费
	ContinueUnit uint        `json:"continueUnit" orm:"continue_unit" ` // 续段步长:按件=每续N件;按重=每续N克
	ContinueFee  float64     `json:"continueFee"  orm:"continue_fee"  ` // 续段单价(每步长追加费用)
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` // 创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` // 更新时间
}
