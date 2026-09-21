// freight.go 运费域——表: freight_template / freight_rule（V5/V23/V31 seller 维度）。
// 规则: 满额包邮按订单商品金额 ≥ free_threshold（free_exclude_codes 省份例外 V23）;
// 计费量 Q: 按件=总件数; 按重=Σ(weight*qty), 缺 weight 按 1 件折算;
// 区域匹配: 收货省代码(GB/T 2260, V23 区划码) ∈ rule.region_codes, 无命中用空 regions 兜底;
// 运费 = first_fee + Q>first_unit ? CEIL((Q-first_unit)/continue_unit)*continue_fee : 0。
package shop

import "context"

// IFreightLogic 运费。
type IFreightLogic interface {
	// Templates 模板列表（后台）。
	Templates(ctx context.Context, page PageQuery) (*PageResult[FreightTemplateItem], error)
	Create(ctx context.Context, in FreightTemplateInput) (int64, error)
	Update(ctx context.Context, id int64, in FreightTemplateInput) error
	Delete(ctx context.Context, id int64) error

	// Calculate 订单运费计算（下单/试算调用; 收货省代码+SKU 明细 → 运费）。
	Calculate(ctx context.Context, templateId int64, provinceCode string, items []FreightCalcItem, goodsAmount string) (string, error)
}

type FreightTemplateItem struct {
	Id               int64   `json:"id"`
	Name             string  `json:"name"`
	ChargeType       int     `json:"chargeType" dc:"1按件 2按重"`
	FreeThreshold    *string `json:"freeThreshold" dc:"满额包邮,空=不包邮"`
	FreeExcludeCodes string  `json:"freeExcludeCodes" dc:"不参与包邮的省级代码(JSON)"`
	Status           int     `json:"status"`
}

type FreightTemplateInput struct {
	Name             string
	ChargeType       int
	FreeThreshold    *string
	FreeExcludeCodes []string
	Rules            []FreightRuleInput
}

type FreightRuleInput struct {
	RegionCodes  []string `json:"regionCodes" dc:"省级代码; 空=全国兜底"`
	FirstUnit    int
	FirstFee     string
	ContinueUnit int
	ContinueFee  string
}

type FreightCalcItem struct {
	SkuId    int64
	Weight   string // 克（按重计费用）
	Quantity int
}
