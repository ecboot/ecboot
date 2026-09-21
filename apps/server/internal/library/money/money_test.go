package money

import "testing"

func TestFromYuanString(t *testing.T) {
	cases := []struct {
		in   string
		want int64
		err  bool
	}{
		{"99.00", 9900, false},
		{"0.10", 10, false},
		{"-5.50", -550, false},
		{"100", 10000, false},
		{"0.1", 10, false},
		{"", 0, true},
		{"1.234", 0, true},
		{"abc", 0, true},
		{"1.999", 0, true},
	}
	for _, c := range cases {
		got, err := FromYuanString(c.in)
		if c.err {
			if err == nil {
				t.Errorf("FromYuanString(%q) 期望报错, 实际 %d", c.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("FromYuanString(%q) 意外报错: %v", c.in, err)
		}
		if got != c.want {
			t.Errorf("FromYuanString(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestToYuanString(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{9900, "99.00"},
		{-550, "-5.50"},
		{5, "0.05"},
		{0, "0.00"},
	}
	for _, c := range cases {
		if got := ToYuanString(c.in); got != c.want {
			t.Errorf("ToYuanString(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestAllocateProRataRemainderToLast(t *testing.T) {
	// 100 分摊到 [33,33,34] 比例行: 尾差记末行
	got := AllocateProRata(100, []int64{3333, 3333, 3334})
	var sum int64
	for _, v := range got {
		sum += v
	}
	if sum != 100 {
		t.Errorf("分摊合计 = %d, want 100", sum)
	}
	if got[2] != 34 {
		t.Errorf("末行 = %d, want 34(尾差)", got[2])
	}
	// 全 0 行
	if got := AllocateProRata(100, []int64{0, 0}); len(got) != 2 {
		t.Error("全 0 行应返回等长 0 分摊")
	}
}
