package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/service/user"
)

// DistRuleQuery 佣金比例查询（商品覆盖 > 分类默认）。
func (c *ControllerV1) DistRuleQuery(ctx context.Context, req *v1.DistRuleQueryReq) (res *v1.DistRuleQueryRes, err error) {
	spuId, err := parseID(req.SpuId)
	if err != nil {
		return nil, err
	}
	hits, err := user.NewDistributionLogic().RuleQuery(ctx, spuId)
	if err != nil {
		return nil, err
	}
	res = &v1.DistRuleQueryRes{List: make([]v1.DistRuleItem, 0, len(hits))}
	for _, h := range hits {
		res.List = append(res.List, v1.DistRuleItem{ScopeDesc: h.ScopeDesc, Level1Rate: h.Level1Rate, Level2Rate: h.Level2Rate})
	}
	return res, nil
}
