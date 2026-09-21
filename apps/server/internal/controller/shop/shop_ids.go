package shop

import "strconv"

// fmtID api 层 ID 一律 string（防 JS 精度）, 与 service int64 的统一转换。
// 各渠道包各自持有（四渠道目录互不引用, 见分层契约）。
// 注: 本渠道端点无 ID 入参（纯公开列表）, 故未定义 parseID（YAGNI）。
func fmtID(id int64) string {
	return strconv.FormatInt(id, 10)
}
