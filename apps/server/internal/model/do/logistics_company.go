// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LogisticsCompany is the golang structure of table logistics_company for DAO operations like Where/Data.
type LogisticsCompany struct {
	g.Meta       `orm:"table:logistics_company, do:true"`
	Id           any         // 物流公司ID
	Code         any         // 编码(唯一,订单deliver_company存此值)
	Name         any         // 公司名称(如 顺丰速运)
	TrackingRule any         // 运单号校验规则描述(如 SF+12位数字)
	Sort         any         // 排序
	Status       any         // 状态:1启用 0停用
	Deleted      any         // 软删除:0否 1是
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
}
