// marketing.go 营销 C 端公开面（015-marketing-c 批次 09 新建）。
// 背景: IActivityLogic 是**管理面**接口（含创建/改配/场次商品维护）; C 端此前**无任何接口与 DTO**,
// 故按批次 03（IOperationLogic）先例在本批新建。
// 口径（与数据模型 research D8 一致）:
//   - 五类列表**只出有效活动**（启用 + 未删 + 在时间窗内）; 秒杀额外含"预告"（已配置未开始）。
//   - 统一**按活动结束时间升序**（最快结束的先展示）+ 分页。
//   - **首页入口与各列表共用同一查询函数**（只多一个 Limit）——避免"首页有而列表没有"的口径漂移
//     （批次 08 的"列表与汇总漂移"教训）。
package shop

import (
	"context"

	"ecboot/internal/model"
)

// IMarketingLogic 营销公开面。
type IMarketingLogic interface {
	PublicGroupBuys(ctx context.Context, page model.PageReq) (*model.PageResult[model.PublicGroupBuyItem], error)
	PublicFlashSales(ctx context.Context, page model.PageReq) (*model.PageResult[model.PublicFlashSaleItem], error)
	PublicBargains(ctx context.Context, page model.PageReq) (*model.PageResult[model.PublicBargainItem], error)
	PublicAssists(ctx context.Context, page model.PageReq) (*model.PageResult[model.PublicAssistItem], error)
	PublicFullReductions(ctx context.Context, page model.PageReq) (*model.PageResult[model.PublicFullReductionItem], error)
	// Index 首页聚合（FR-017）: 轮播 + 楼层 + 五类活动入口（各取前 N 条）+ 可领券;
	// **任一块为空返回空数组而非报错**（首页可用性优先）。
	Index(ctx context.Context) (*model.IndexAggregate, error)
}
