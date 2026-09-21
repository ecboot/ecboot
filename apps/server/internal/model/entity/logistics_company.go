// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LogisticsCompany is the golang structure for table logistics_company.
type LogisticsCompany struct {
	Id           uint64      `json:"id"           orm:"id"            ` // 物流公司ID
	Code         string      `json:"code"         orm:"code"          ` // 编码(唯一,订单deliver_company存此值)
	Name         string      `json:"name"         orm:"name"          ` // 公司名称(如 顺丰速运)
	TrackingRule string      `json:"trackingRule" orm:"tracking_rule" ` // 运单号校验规则描述(如 SF+12位数字)
	Sort         int         `json:"sort"         orm:"sort"          ` // 排序
	Status       int         `json:"status"       orm:"status"        ` // 状态:1启用 0停用
	Deleted      int         `json:"deleted"      orm:"deleted"       ` // 软删除:0否 1是
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` // 创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` // 更新时间
}
