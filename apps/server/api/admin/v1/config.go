package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 系统配置列表
	AdminConfigListReq struct {
		g.Meta `path:"/configs" method:"GET" summary:"系统配置列表"`
	}
	AdminConfigItem struct {
		Code        string `json:"code" dc:"配置编码"`
		Value       string `json:"value" dc:"配置值"`
		ValueType   int    `json:"valueType" dc:"1整数 2小数 3字符串 4布尔 5JSON"`
		Name        string `json:"name" dc:"名称"`
		Description string `json:"description" dc:"说明(含代码默认值)"`
		Status      int    `json:"status" dc:"1启用 0停用(回退代码默认)"`
	}
	AdminConfigListRes struct {
		List []AdminConfigItem `json:"list"`
	}

	// 修改配置
	// 权限: system:config:update
	AdminConfigUpdateReq struct {
		g.Meta `path:"/configs/{code}" method:"PUT" summary:"修改系统配置"`
		Code   string `json:"code" v:"required" dc:"配置编码"`
		Value  string `json:"value" v:"required" dc:"新值"`
		Status int    `json:"status" dc:"状态"`
	}
	AdminConfigUpdateRes struct {
		Success bool `json:"success"`
	}
)
