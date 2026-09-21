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
	if in.Status == 0 {
		enabled := 1
		out.Status = &enabled // 创建默认启用（零值覆盖防护——010 评审 C1: 缺此防护则落库 status=0 致 C 端不可见）
	}
	return out
}

func doCategoryUpdate(in model.CategoryInput) *do.ProductCategory {
	out := &do.ProductCategory{}
	_ = gconv.Struct(in, out)
	omitEmptyStrings(out)
	// 010 评审 C2: parent_id/level 不在 api Update 契约内（客户端无法携带）——必须显式跳过,
	// 否则每次"改名"都会把它们清零（分类被搬根 + level=0 且被禁用, C 端树消失）。
	out.ParentId = nil
	out.Level = nil
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
	// 010 评审 I2: 关联 ID 的零值来自"未传"（api 为空串 → optID 归一为 0）, 跳过以防静默清空关联；
	// 局限: 无法经本接口"解除关联"（如需, 后续以显式语义解决——宁可不能解除, 不可静默破坏）。
	if in.CategoryId == 0 {
		out.CategoryId = nil
	}
	if in.BrandId == 0 {
		out.BrandId = nil
	}
	if in.FreightTemplateId == 0 {
		out.FreightTemplateId = nil
	}
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
