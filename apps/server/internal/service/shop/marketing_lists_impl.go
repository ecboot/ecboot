// marketing_lists_impl.go 营销公开列表（拼团/砍价/助力/满减）+ 首页聚合（015-marketing-c）；
// 秒杀列表在 marketing_impl.go。口径统一: 只出"启用 + 未删 + 在时间窗内", 按**活动结束时间升序**。
package shop

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"

	"ecboot/internal/dao"
	"ecboot/internal/model"
)

// PublicGroupBuys 拼团公开列表（FR-018）: 只出进行中; 含成团价与成团人数。
func (i *MarketingLogicImpl) PublicGroupBuys(
	ctx context.Context, page model.PageReq,
) (*model.PageResult[model.PublicGroupBuyItem], error) {
	page = page.Normalized()
	cols := dao.GroupBuyActivity.Columns()
	base := func() *gdb.Model {
		return dao.GroupBuyActivity.Ctx(ctx).
			Where(cols.Status, 1).Where(cols.Deleted, 0).
			Where(cols.ValidStartAt + " <= NOW()").Where(cols.ValidEndAt + " >= NOW()")
	}
	total, err := base().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计拼团活动失败")
	}
	acts, err := base().OrderAsc(cols.ValidEndAt).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询拼团活动失败")
	}
	spus, err := spuBriefs(ctx, collectIDs(acts, cols.SpuId))
	if err != nil {
		return nil, err
	}
	icols := dao.GroupBuyItem.Columns()
	list := make([]model.PublicGroupBuyItem, 0, len(acts))
	for _, a := range acts {
		id := a[cols.Id].Int64()
		spuId := a[cols.SpuId].Int64()
		briefs, e := skuBriefsOf(ctx, dao.GroupBuyItem.Table(), icols.ActivityId, icols.GroupPrice, id)
		if e != nil {
			return nil, e
		}
		list = append(list, model.PublicGroupBuyItem{
			ActivityId: id,
			Name:       a[cols.Name].String(),
			SpuId:      spuId,
			SpuName:    spus[spuId].Name,
			Image:      spus[spuId].Image,
			GroupSize:  a[cols.GroupSize].Int(),
			Items:      briefs,
			EndTime:    a[cols.ValidEndAt].String(),
		})
	}
	return &model.PageResult[model.PublicGroupBuyItem]{List: list, Total: int64(total)}, nil
}

// PublicBargains 砍价公开列表（FR-006）: 只出进行中; 含起始价/底价/最大刀数。
func (i *MarketingLogicImpl) PublicBargains(
	ctx context.Context, page model.PageReq,
) (*model.PageResult[model.PublicBargainItem], error) {
	page = page.Normalized()
	cols := dao.BargainActivity.Columns()
	base := func() *gdb.Model {
		return dao.BargainActivity.Ctx(ctx).
			Where(cols.Status, 1).Where(cols.Deleted, 0).
			Where(cols.StartTime + " <= NOW()").Where(cols.EndTime + " >= NOW()")
	}
	total, err := base().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计砍价活动失败")
	}
	acts, err := base().OrderAsc(cols.EndTime).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询砍价活动失败")
	}
	spus, err := spuBriefs(ctx, collectIDs(acts, cols.SpuId))
	if err != nil {
		return nil, err
	}
	list := make([]model.PublicBargainItem, 0, len(acts))
	for _, a := range acts {
		id := a[cols.Id].Int64()
		spuId := a[cols.SpuId].Int64()
		briefs, e := bargainBriefs(ctx, id)
		if e != nil {
			return nil, e
		}
		list = append(list, model.PublicBargainItem{
			ActivityId: id,
			Name:       a[cols.Name].String(),
			SpuId:      spuId,
			SpuName:    spus[spuId].Name,
			Image:      spus[spuId].Image,
			Items:      briefs,
			EndTime:    a[cols.EndTime].String(),
		})
	}
	return &model.PageResult[model.PublicBargainItem]{List: list, Total: int64(total)}, nil
}

// PublicAssists 助力公开列表（FR-011）: 只出进行中; 含所需人数/限发起次数/奖励类型。
func (i *MarketingLogicImpl) PublicAssists(
	ctx context.Context, page model.PageReq,
) (*model.PageResult[model.PublicAssistItem], error) {
	page = page.Normalized()
	cols := dao.AssistActivity.Columns()
	base := func() *gdb.Model {
		return dao.AssistActivity.Ctx(ctx).
			Where(cols.Status, 1).Where(cols.Deleted, 0).
			Where(cols.StartTime + " <= NOW()").Where(cols.EndTime + " >= NOW()")
	}
	total, err := base().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计助力活动失败")
	}
	acts, err := base().OrderAsc(cols.EndTime).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询助力活动失败")
	}
	list := make([]model.PublicAssistItem, 0, len(acts))
	for _, a := range acts {
		list = append(list, model.PublicAssistItem{
			ActivityId:    a[cols.Id].Int64(),
			Name:          a[cols.Name].String(),
			RequiredCount: a[cols.RequiredCount].Int(),
			PerLimit:      a[cols.PerLimit].Int(),
			RewardType:    a[cols.RewardType].Int(),
			EndTime:       a[cols.EndTime].String(),
		})
	}
	return &model.PageResult[model.PublicAssistItem]{List: list, Total: int64(total)}, nil
}

// PublicFullReductions 满减公开列表（FR-016）: 只出进行中; 含档位与适用范围摘要。
func (i *MarketingLogicImpl) PublicFullReductions(
	ctx context.Context, page model.PageReq,
) (*model.PageResult[model.PublicFullReductionItem], error) {
	page = page.Normalized()
	cols := dao.PromotionActivity.Columns()
	base := func() *gdb.Model {
		return dao.PromotionActivity.Ctx(ctx).
			Where(cols.Status, 1).Where(cols.Deleted, 0).
			Where(cols.StartTime + " <= NOW()").Where(cols.EndTime + " >= NOW()")
	}
	total, err := base().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计满减活动失败")
	}
	acts, err := base().OrderAsc(cols.EndTime).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询满减活动失败")
	}
	lcols := dao.PromotionActivityLadder.Columns()
	list := make([]model.PublicFullReductionItem, 0, len(acts))
	for _, a := range acts {
		id := a[cols.Id].Int64()
		recs, e := dao.PromotionActivityLadder.Ctx(ctx).
			Where(lcols.ActivityId, id).OrderAsc(lcols.ThresholdAmount).All()
		if e != nil {
			return nil, gerror.Wrap(e, "查询满减档位失败")
		}
		ladders := make([]model.LadderBrief, 0, len(recs))
		for _, r := range recs {
			ladders = append(ladders, model.LadderBrief{
				Threshold: r[lcols.ThresholdAmount].String(),
				Discount:  r[lcols.DiscountAmount].String(),
			})
		}
		list = append(list, model.PublicFullReductionItem{
			ActivityId: id,
			Name:       a[cols.Name].String(),
			Ladders:    ladders,
			ScopeDesc:  scopeDescOf(ctx, id),
			EndTime:    a[cols.EndTime].String(),
		})
	}
	return &model.PageResult[model.PublicFullReductionItem]{List: list, Total: int64(total)}, nil
}

// spuBrief 商品摘要（名称 + 主图）。
type spuBrief struct {
	Name  string
	Image string
}

// collectIDs 取指定列的非空 id 集合（去重）。
func collectIDs(recs gdb.Result, col string) []int64 {
	seen := map[int64]bool{}
	out := []int64{}
	for _, r := range recs {
		id := r[col].Int64()
		if id > 0 && !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	return out
}

// spuBriefs 批量取商品名称与主图（一次 IN 查询, 避免 N+1; 主图取 images 数组首元素）。
func spuBriefs(ctx context.Context, spuIds []int64) (map[int64]spuBrief, error) {
	out := map[int64]spuBrief{}
	if len(spuIds) == 0 {
		return out, nil
	}
	cols := dao.ProductSpu.Columns()
	recs, err := dao.ProductSpu.Ctx(ctx).WhereIn(cols.Id, spuIds).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询商品摘要失败")
	}
	for _, r := range recs {
		imgs := []string{}
		_ = gconv.Struct(r[cols.Images].String(), &imgs)
		img := ""
		if len(imgs) > 0 {
			img = imgs[0]
		}
		out[r[cols.Id].Int64()] = spuBrief{Name: r[cols.Name].String(), Image: img}
	}
	return out, nil
}

// skuBriefsOf 场次商品摘要（表名/关联列/价格列参数化: 拼团与秒杀共用形态）。
func skuBriefsOf(ctx context.Context, table, actCol, priceCol string, actId int64) ([]model.ActivitySkuBrief, error) {
	recs, err := g.DB().Model(table).Ctx(ctx).Where(actCol, actId).OrderAsc("id").All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询场次商品失败")
	}
	out := make([]model.ActivitySkuBrief, 0, len(recs))
	for _, r := range recs {
		out = append(out, model.ActivitySkuBrief{
			SkuId: r["sku_id"].Int64(),
			Price: r[priceCol].String(),
		})
	}
	return out, nil
}

// bargainBriefs 砍价场次商品摘要（起始价/底价/最大刀数）。
func bargainBriefs(ctx context.Context, actId int64) ([]model.ActivitySkuBrief, error) {
	cols := dao.BargainItem.Columns()
	recs, err := dao.BargainItem.Ctx(ctx).Where(cols.ActivityId, actId).OrderAsc(cols.Id).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询砍价场次商品失败")
	}
	out := make([]model.ActivitySkuBrief, 0, len(recs))
	for _, r := range recs {
		out = append(out, model.ActivitySkuBrief{
			ItemId:        r[cols.Id].Int64(), // 砍价发起要用**场次商品 id**（非 sku id）
			SkuId:         r[cols.SkuId].Int64(),
			Price:         r[cols.OriginalPrice].String(),
			OriginalPrice: r[cols.OriginalPrice].String(),
			FloorPrice:    r[cols.FloorPrice].String(),
			MaxCutCount:   r[cols.MaxCutCount].Int(),
		})
	}
	return out, nil
}

// scopeDescOf 满减适用范围摘要（全场/分类/商品）。
func scopeDescOf(ctx context.Context, actId int64) string {
	cols := dao.PromotionActivityScope.Columns()
	recs, err := dao.PromotionActivityScope.Ctx(ctx).
		Where(cols.ActivityId, actId).Fields(cols.ScopeType).All()
	if err != nil || len(recs) == 0 {
		return "全场"
	}
	has := map[int]bool{}
	for _, r := range recs {
		has[r[cols.ScopeType].Int()] = true
	}
	switch {
	case has[1]:
		return "全场"
	case has[2]:
		return "分类"
	case has[3]:
		return "指定商品"
	}
	return "全场"
}

// Index 首页聚合（FR-017, 用户裁定 D2）: 轮播 + 楼层 + 五类活动入口（各取前 N）+ 可领券。
// **任一块为空返回空数组而非报错**（首页可用性优先）; 活动入口与各列表**共用同一查询函数**（只多一个 Limit）。
func (i *MarketingLogicImpl) Index(ctx context.Context) (*model.IndexAggregate, error) {
	out := &model.IndexAggregate{
		Banners:        []model.PublicBannerItem{},
		Floors:         []model.PublicFloorItem{},
		FlashSales:     []model.PublicFlashSaleItem{},
		GroupBuys:      []model.PublicGroupBuyItem{},
		Bargains:       []model.PublicBargainItem{},
		Assists:        []model.PublicAssistItem{},
		FullReductions: []model.PublicFullReductionItem{},
		Coupons:        []model.CouponTemplate{},
	}
	entryPage := model.PageReq{Page: 1, PageSize: indexEntryLimit}

	// 轮播与楼层（批次 03 装修域）: 失败不阻断首页其余分块
	if banners, err := PublicBanners(ctx, 1); err == nil {
		out.Banners = banners
	}
	if floors, err := PublicFloors(ctx); err == nil {
		out.Floors = floors
	}
	if fs, err := i.PublicFlashSales(ctx, entryPage); err == nil && fs != nil {
		out.FlashSales = fs.List
	}
	if gb, err := i.PublicGroupBuys(ctx, entryPage); err == nil && gb != nil {
		out.GroupBuys = gb.List
	}
	if bg, err := i.PublicBargains(ctx, entryPage); err == nil && bg != nil {
		out.Bargains = bg.List
	}
	if as, err := i.PublicAssists(ctx, entryPage); err == nil && as != nil {
		out.Assists = as.List
	}
	if fr, err := i.PublicFullReductions(ctx, entryPage); err == nil && fr != nil {
		out.FullReductions = fr.List
	}
	// 可领券: shop 域直读券表（**既有先例**: promotion_calc.go 早已直读 coupon/user_coupon 计价）。
	// 注（记账）: 此处口径是"**有哪些券可领**"（模板级: 启用未删 + 未领完）; 而 user 域的
	// `/user/coupons/available` 是"**我还能领哪些**"（会员级: 叠加个人限领判定）——两者回答不同问题,
	// 故不强行统一; 若后续要把会员级口径也搬上首页, 应走既有 `ICouponQuery` 端口而非在此重复实现。
	if cp, err := publicCouponBriefs(ctx, indexEntryLimit); err == nil {
		out.Coupons = cp
	}
	return out, nil
}

// publicCouponBriefs 首页可领券（模板级口径: 启用 + 未删 + 未领完; 上限 limit 条）。
func publicCouponBriefs(ctx context.Context, limit int) ([]model.CouponTemplate, error) {
	cols := dao.Coupon.Columns()
	recs, err := dao.Coupon.Ctx(ctx).
		Where(cols.Status, 1).
		Where(cols.Deleted, 0).
		Where("total_count = 0 OR received_count < total_count"). // total_count=0 视为不限量
		OrderDesc(cols.Id).Limit(limit).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询可领券失败")
	}
	list := make([]model.CouponTemplate, 0, len(recs))
	for _, r := range recs {
		list = append(list, model.CouponTemplate{
			Id:         r[cols.Id].Int64(),
			Name:       r[cols.Name].String(),
			Type:       r[cols.Type].Int(),
			Threshold:  r[cols.ThresholdAmount].String(),
			Discount:   r[cols.DiscountAmount].String(),
			TotalCount: r[cols.TotalCount].Int(),
			Received:   r[cols.ReceivedCount].Int(),
			PerLimit:   r[cols.PerLimit].Int(),
			ValidDesc:  couponValidDescOf(r),
			Status:     r[cols.Status].Int(),
		})
	}
	return list, nil
}

// couponValidDescOf 券有效期描述（模板级; 口径与 user 域的展示一致: 固定区间 / 领取后 N 天）。
func couponValidDescOf(r gdb.Record) string {
	cols := dao.Coupon.Columns()
	if r[cols.ValidType].Int() == 2 {
		return "领取后 " + r[cols.ValidDays].String() + " 天有效"
	}
	return r[cols.ValidStartAt].String() + " 至 " + r[cols.ValidEndAt].String()
}
