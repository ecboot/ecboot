package model

type PartnerOptions struct {
	Type        string `json:"type" dc:"表单类型，如：input，select，radio"`
	Name        string `json:"name" dc:"配置项，如：ruleId"`
	Val         string `json:"val" dc:"配置值，如：xxxx"`
	Description string `json:"description" dc:"配置项描述，如：规则ID"`
	Placeholder string `json:"placeholder" dc:"配置项提示性内容"`
	Options     []struct {
		Name string `json:"name" dc:"选项文本"`
		Val  string `json:"val" dc:"选项值"`
	} `json:"options" dc:"配置项可选项"`
}

type QueryParams struct {
	ProviderId string `json:"providerId" validate:"required" v:"required" dc:"合作方appid标识"` // 合作方标识
	Sign       string `json:"sign" validate:"required" v:"required" dc:"签名"`               // 签名
	Timestamp  int64  `json:"timestamp" validate:"required" v:"required" dc:"接口调用时的毫秒时间戳"` // 接口调用时的毫秒时间戳
	TraceId    string `json:"traceId" validate:"required" v:"required" dc:"请求唯一标识"`        // 请求唯一标识
	SignMethod string `json:"signMethod" validate:"required" v:"required" dc:"签名方法，支持MD5"` // 请求唯一标识
}

type FileTypeDTO struct {
	Url  string `json:"url" dc:"附件连接"`
	Type string `json:"type" dc:"附件类型"`
}

type (
	PageRequest struct {
		PageNum  int
		PageSize int
	}

	PageResult[T any] struct {
		Records  []*T `json:"records"`
		PageNum  int  `json:"pageNum"`
		PageSize int  `json:"pageSize"`
		Total    int  `json:"total"`
	}
)
