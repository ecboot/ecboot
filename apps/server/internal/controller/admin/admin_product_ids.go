package admin

// optID 可选 ID 字符串 → int64（空串=0）。商品域可选关联 ID（分类/品牌/运费模板）共用。
func optID(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	return parseID(s)
}

// optIDPair/optIDTriple 多可选 ID 批量转换（减少样板）。
func optIDPair(a, b string) (int64, int64, error) {
	x, err := optID(a)
	if err != nil {
		return 0, 0, err
	}
	y, err := optID(b)
	return x, y, err
}

func optIDTriple(a, b, c string) (int64, int64, int64, error) {
	x, err := optID(a)
	if err != nil {
		return 0, 0, 0, err
	}
	y, err := optID(b)
	if err != nil {
		return 0, 0, 0, err
	}
	z, err := optID(c)
	return x, y, z, err
}
