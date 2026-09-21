// Package money 金额运算工具：内部以「分」为 int64 运算（AGENTS 金额红线），
// 对外（API DTO）为元 string 两位小数。禁止 float 参与任何金额运算。
package money

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// FromYuanString 元字符串（两位小数）→ 分。
func FromYuanString(yuan string) (int64, error) {
	s := strings.TrimSpace(yuan)
	if s == "" {
		return 0, errors.New("金额为空")
	}
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}
	parts := strings.Split(s, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("金额格式非法: %s", yuan)
	}
	intPart, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("金额整数部分非法: %s", yuan)
	}
	fen := intPart * 100
	if len(parts) == 2 {
		dec := parts[1]
		if len(dec) == 1 {
			dec += "0"
		}
		if len(dec) != 2 {
			return 0, fmt.Errorf("金额小数部分非法: %s", yuan)
		}
		decPart, err := strconv.ParseInt(dec, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("金额小数部分非法: %s", yuan)
		}
		if decPart >= 100 {
			return 0, fmt.Errorf("金额小数部分非法: %s", yuan)
		}
		fen += decPart
	}
	if neg {
		fen = -fen
	}
	return fen, nil
}

// ToYuanString 分 → 元字符串（两位小数）。
func ToYuanString(fen int64) string {
	sign := ""
	if fen < 0 {
		sign = "-"
		fen = -fen
	}
	return fmt.Sprintf("%s%d.%02d", sign, fen/100, fen%100)
}

// MulQty 单价（分）× 数量。
func MulQty(unitFen int64, qty int64) int64 { return unitFen * qty }

// AllocateProRata 按行金额比例分摊 total（分），尾差记最后一行（分摊恒等式保证）。
// 输入行金额与返回分摊一一对应；全 0 行时返回全 0。
func AllocateProRata(total int64, lineAmounts []int64) []int64 {
	out := make([]int64, len(lineAmounts))
	var sum int64
	for _, a := range lineAmounts {
		sum += a
	}
	if sum == 0 {
		return out
	}
	var allocated int64
	for i, a := range lineAmounts {
		if i == len(lineAmounts)-1 {
			out[i] = total - allocated // 尾差记末行
			continue
		}
		// 四舍五入到分
		v := total * a / sum
		out[i] = v
		allocated += v
	}
	return out
}
