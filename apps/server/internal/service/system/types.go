// types.go 治理域服务层公共类型。
package system

// PageQuery 服务层分页查询参数。
type PageQuery struct {
	Page     int
	PageSize int
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

// PageResult 分页结果（泛型容器）。
type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
}
