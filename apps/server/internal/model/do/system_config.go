// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemConfig is the golang structure of table system_config for DAO operations like Where/Data.
type SystemConfig struct {
	g.Meta      `orm:"table:system_config, do:true"`
	Id          any         // 配置ID
	Code        any         // 配置编码(点分命名空间,如 distribution.level.threshold)
	Value       any         // 配置值(标量或JSON对象)
	ValueType   any         // 值类型:1整数 2小数 3字符串 4布尔 5JSON对象
	Name        any         // 配置名称(后台展示)
	Description any         // 配置说明(含代码默认值,便于回退排查)
	Status      any         // 状态:1启用 0停用(停用回退应用代码内置默认值)
	Deleted     any         // 软删除:0否 1是
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
}
