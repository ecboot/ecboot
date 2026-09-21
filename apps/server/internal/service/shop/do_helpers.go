package shop

import (
	"reflect"
)

// omitEmptyStrings 遍历 do 结构体的 any 字段, 把空串值置 nil
// （nil = gf ORM 跳过该列; 空串写 DECIMAL/数值列会 1366——实测根因）。
func omitEmptyStrings(v any) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() || rv.Elem().Kind() != reflect.Struct {
		return
	}
	re := rv.Elem()
	for i := 0; i < re.NumField(); i++ {
		f := re.Field(i)
		if !f.CanSet() || f.Kind() != reflect.Interface {
			continue
		}
		if s, ok := f.Interface().(string); ok && s == "" {
			f.Set(reflect.Zero(f.Type()))
		}
	}
}
