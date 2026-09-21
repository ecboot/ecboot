// do_builders.go do 对象构造 helper（分层契约 §二: Data 一律 do 对象）。
// gconv.Struct 按 json/dc 标签映射, 零值字段由 do 指针语义自动跳过。
package shop

import (
	"github.com/gogf/gf/v2/util/gconv"

	"ecboot/internal/model"
	"ecboot/internal/model/do"
)

func doCategory(in model.CategoryInput) *do.ProductCategory {
	out := &do.ProductCategory{}
	_ = gconv.Struct(in, out)
	omitEmptyStrings(out)
	return out
}

func doCategoryUpdate(in model.CategoryInput) *do.ProductCategory {
	out := &do.ProductCategory{}
	_ = gconv.Struct(in, out)
	omitEmptyStrings(out)
	return out
}

func doBrand(in model.BrandInput) *do.ProductBrand {
	out := &do.ProductBrand{}
	_ = gconv.Struct(in, out)
	if in.Status == 0 {
		enabled := 1
		out.Status = &enabled // 创建默认启用（零值覆盖防护）
	}
	return out
}

func doBrandUpdate(in model.BrandInput) *do.ProductBrand {
	out := &do.ProductBrand{}
	_ = gconv.Struct(in, out)
	omitEmptyStrings(out)
	return out
}

func doSpuCreate(in model.SpuInput, spuNo string) *do.ProductSpu {
	out := &do.ProductSpu{}
	_ = gconv.Struct(in, out)
	omitEmptyStrings(out)
	out.SpuNo = &spuNo
	return out
}

func doSpuUpdate(in model.SpuInput) *do.ProductSpu {
	out := &do.ProductSpu{}
	_ = gconv.Struct(in, out)
	omitEmptyStrings(out)
	return out
}

func doSkuCreate(spuId int64, in model.SkuInput, skuNo string) *do.ProductSku {
	out := &do.ProductSku{}
	_ = gconv.Struct(in, out)
	omitEmptyStrings(out)
	out.SpuId = &spuId
	out.SkuNo = &skuNo
	return out
}

func doSkuUpdate(in model.SkuInput) *do.ProductSku {
	out := &do.ProductSku{}
	_ = gconv.Struct(in, out)
	omitEmptyStrings(out)
	return out
}
