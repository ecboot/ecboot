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

// ---------- 分页（单一语义三件套, 2026-09-20 合并原 PageRequest/PageResult 双轨） ----------
// 约定：字段统一 page/pageSize（契约 JSON 与内部参数同名，零转换）；
// 默认与上限只在契约层声明（d/v 标签），仓储层做防御性兜底；
// ID/金额对外一律 string（int64 与十进制精度安全）；时间 RFC3339。
// 各渠道 api/{渠道}/v1 以类型别名引用（PageReq = model.PageReq）。

// PageReq 分页入参：嵌入 api 各列表 Req（契约结构），亦直接作为仓储查询参数。
type PageReq struct {
	Page     int `json:"page" dc:"页码,默认1" v:"min:1" d:"1"`
	PageSize int `json:"pageSize" dc:"每页数量,默认10,上限100" v:"min:1|max:100" d:"10"`
}

// PageRes 分页出参：嵌入 api 各列表 Res（total 之外的 list 字段由各接口自定义）。
type PageRes struct {
	Total int64 `json:"total" dc:"总条数"`
}

// PageResult 内部分页查询容器：repository 泛型返回（List + Total），
// 页参数由调用方持有的 PageReq 承载，不再重复回传。
type PageResult[T any] struct {
	List  []T   `json:"list" dc:"数据列表"`
	Total int64 `json:"total" dc:"总条数"`
}
