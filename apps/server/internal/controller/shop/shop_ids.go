package shop

import (
	"strconv"

	"ecboot/internal/errcode"
)

// parseID/fmtID api 层 ID 一律 string（防 JS 精度）, 与 service int64 的统一转换。
// 各渠道包各自持有（四渠道目录互不引用, 见分层契约）。
// 注: 010 起本渠道端点含 ID 入参（商品详情/筛选）, 故 parseID 启用。
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

// spuIdOf 可选的 spuId 查询参数（空/非法 → 0 = 不过滤; 满减列表用）。
func spuIdOf(s string) int64 {
	if s == "" {
		return 0
	}
	id, err := parseID(s)
	if err != nil {
		return 0
	}
	return id
}
