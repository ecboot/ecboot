// types.go 用户域服务层公共类型（服务层内部形态, 与 api 层 model.PageReq/PageRes JSON 结构解耦）。
package user

// PageQuery 服务层分页查询参数（api 层 PageReq 转换而来, 归一化后传入）。
type PageQuery struct {
	Page     int // 默认 1
	PageSize int // 默认 10, 上限 100
}

// Normalized 归一化非法值。
func (q PageQuery) Normalized() PageQuery {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 || q.PageSize > 100 {
		q.PageSize = 10
	}
	return q
}

// Offset SQL 偏移。
func (q PageQuery) Offset() int { return (q.Normalized().Page - 1) * q.Normalized().PageSize }

// Limit SQL 行数。
func (q PageQuery) Limit() int { return q.Normalized().PageSize }

// PageResult 分页结果（泛型容器: List+Total, 字段与契约 list/total 对齐）。
type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
}
