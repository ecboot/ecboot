package admin

import "strconv"

// parseID/fmtID api 层 ID 一律 string（防 JS 精度）, 与 service int64 的统一转换。
func parseID(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func fmtID(id int64) string {
	return strconv.FormatInt(id, 10)
}
