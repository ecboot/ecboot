// freight.go 运费域——表: freight_template / freight_rule（V5/V23/V31 seller 维度）。
// 规则: 满额包邮按订单商品金额 ≥ free_threshold（free_exclude_codes 省份例外 V23）;
// 计费量 Q: 按件=总件数; 按重=Σ(weight*qty), 缺 weight 按 1 件折算;
// 区域匹配: 收货省代码(GB/T 2260, V23 区划码) ∈ rule.region_codes, 无命中用空 regions 兜底;
// 运费 = first_fee + Q>first_unit ? CEIL((Q-first_unit)/continue_unit)*continue_fee : 0。
package shop

import (
	"context"

	"ecboot/internal/model"
)

// IFreightLogic 运费。
type IFreightLogic interface {
	// Templates 模板列表（后台）。
	Templates(ctx context.Context, page model.PageReq) (*model.PageResult[model.FreightTemplateItem], error)
	Create(ctx context.Context, in model.FreightTemplateInput) (int64, error)
	Update(ctx context.Context, id int64, in model.FreightTemplateInput) error
	Delete(ctx context.Context, id int64) error

	// Calculate 订单运费计算（下单/试算调用; 收货省代码+SKU 明细 → 运费）。
	Calculate(ctx context.Context, templateId int64, provinceCode string, items []model.FreightCalcItem, goodsAmount string) (string, error)
}
