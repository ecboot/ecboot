// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemConfig is the golang structure for table system_config.
type SystemConfig struct {
	Id          uint64      `json:"id"          orm:"id"          ` // 配置ID
	Code        string      `json:"code"        orm:"code"        ` // 配置编码(点分命名空间,如 distribution.level.threshold)
	Value       string      `json:"value"       orm:"value"       ` // 配置值(标量或JSON对象)
	ValueType   int         `json:"valueType"   orm:"value_type"  ` // 值类型:1整数 2小数 3字符串 4布尔 5JSON对象
	Name        string      `json:"name"        orm:"name"        ` // 配置名称(后台展示)
	Description string      `json:"description" orm:"description" ` // 配置说明(含代码默认值,便于回退排查)
	Status      int         `json:"status"      orm:"status"      ` // 状态:1启用 0停用(停用回退应用代码内置默认值)
	Deleted     int         `json:"deleted"     orm:"deleted"     ` // 软删除:0否 1是
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  ` // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  ` // 更新时间
}
