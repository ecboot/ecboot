// Package base 跨渠道公共契约结构（contracts/conventions.md）。
// 约定：ID/金额对外一律 string（int64 与十进制精度安全）；时间 RFC3339；
// 分页请求 page 默认 1、pageSize 默认 10 上限 100；列表响应 total + list。
package base

// PageReq 分页请求嵌入结构。
type PageReq struct {
	Page     int `json:"page" dc:"页码,默认1" d:"1"`
	PageSize int `json:"pageSize" dc:"每页数量,默认10,上限100" d:"10"`
}

// PageRes 分页响应嵌入结构。
type PageRes struct {
	Total int64 `json:"total" dc:"总条数"`
}
