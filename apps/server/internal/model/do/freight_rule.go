// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// FreightRule is the golang structure of table freight_rule for DAO operations like Where/Data.
type FreightRule struct {
	g.Meta       `orm:"table:freight_rule, do:true"`
	Id           any         // 规则ID
	TemplateId   any         // 所属模板ID
	RegionCodes  any         // 适用省级行政区代码列表(如["310000","330000"];空数组=全国兜底规则,每模板至多一条)
	FirstUnit    any         // 首段额度:按件=首件数;按重=首重克数
	FirstFee     any         // 首段运费
	ContinueUnit any         // 续段步长:按件=每续N件;按重=每续N克
	ContinueFee  any         // 续段单价(每步长追加费用)
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
}
