// marketing_impl.go 营销 C 端公开面实现（015-marketing-c 批次 09）。
// 口径见 marketing.go 的接口注释（只出有效活动 / 按结束时间升序 / 首页与列表共用查询）。
package shop

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"ecboot/internal/dao"
	"ecboot/internal/model"
)

// indexEntryLimit 首页各类活动入口的条数上限（research D8）。
const indexEntryLimit = 3

// MarketingLogicImpl IMarketingLogic 实现。
type MarketingLogicImpl struct{}

func NewMarketingLogic() *MarketingLogicImpl { return &MarketingLogicImpl{} }

// PublicFlashSales 秒杀公开列表（FR-001）: 进行中 + **预告**（已配置未开始），已结束不返回。
// 排序按活动结束时间升序（最快结束的先展示）——与首页入口共用。
func (i *MarketingLogicImpl) PublicFlashSales(
	ctx context.Context, page model.PageReq,
) (*model.PageResult[model.PublicFlashSaleItem], error) {
	page = page.Normalized()
	cols := dao.FlashSaleActivity.Columns()
	base := func() *gdb.Model {
		return dao.FlashSaleActivity.Ctx(ctx).
			Where(cols.Status, 1).
			Where(cols.Deleted, 0).
			Where(cols.EndTime + " > NOW()") // 未结束 = 进行中(已开始) 或 预告(未开始)
	}
	total, err := base().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计秒杀场次失败")
	}
	acts, err := base().OrderAsc(cols.EndTime).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询秒杀场次失败")
	}
	items, err := flashItemsByActivities(ctx, acts)
	if err != nil {
		return nil, err
	}
	now := gtime.Now()
	list := make([]model.PublicFlashSaleItem, 0, len(acts))
	for _, a := range acts {
		id := a[cols.Id].Int64()
		list = append(list, model.PublicFlashSaleItem{
			ActivityId: id,
			Name:       a[cols.Name].String(),
			StartTime:  a[cols.StartTime].String(),
			EndTime:    a[cols.EndTime].String(),
			Upcoming:   a[cols.StartTime].GTime() != nil && a[cols.StartTime].GTime().After(now),
			Items:      items[id],
		})
	}
	return &model.PageResult[model.PublicFlashSaleItem]{List: list, Total: int64(total)}, nil
}

// flashItemsByActivities 批量取场次商品（一次 IN 查询, 避免 N+1）→ activityId → 摘要列表。
func flashItemsByActivities(
	ctx context.Context, acts gdb.Result,
) (map[int64][]model.ActivitySkuBrief, error) {
	out := map[int64][]model.ActivitySkuBrief{}
	if len(acts) == 0 {
		return out, nil
	}
	cols := dao.FlashSaleActivity.Columns()
	ids := make([]int64, 0, len(acts))
	for _, a := range acts {
		ids = append(ids, a[cols.Id].Int64())
	}
	icols := dao.FlashSaleItem.Columns()
	recs, err := dao.FlashSaleItem.Ctx(ctx).WhereIn(icols.ActivityId, ids).OrderAsc(icols.Id).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询场次商品失败")
	}
	for _, r := range recs {
		remain := r[icols.StockCount].Int() - r[icols.SoldCount].Int()
		if remain < 0 {
			remain = 0
		}
		aid := r[icols.ActivityId].Int64()
		out[aid] = append(out[aid], model.ActivitySkuBrief{
			SkuId:       r[icols.SkuId].Int64(),
			Price:       r[icols.FlashPrice].String(),
			StockRemain: remain,
			PerLimit:    r[icols.PerLimit].Int(),
		})
	}
	return out, nil
}
