package common

import (
	"strconv"

	"ecboot/internal/errcode"
)

// parseID/fmtID api 层 ID 一律 string（防 JS 精度）, 与 service int64 的统一转换。
// 非法 ID 走参数错误（10001）。各渠道包各自持有（四渠道目录互不引用, 见分层契约）。
func parseID(s string) (int64, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, errcode.New(errcode.CodeInvalidParam, "ID格式错误")
	}
	return id, nil
}

func fmtID(id int64) string {
	return strconv.FormatInt(id, 10)
}
