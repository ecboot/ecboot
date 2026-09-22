// marketing_lists_impl.go 营销公开列表（拼团/砍价/助力/满减）+ 首页聚合（015-marketing-c）；
// 秒杀列表在 marketing_impl.go。口径统一: 只出"启用 + 未删 + 在时间窗内", 按**活动结束时间升序**。
package shop

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

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

// PublicFullReductions 满减公开列表（FR-016）: 只出进行中; 含档位与适用范围摘要;
// spuId>0 时只出**范围命中**该商品的活动（I3/FR-018: 全场 / 商品直配 / 分类含该商品,
// 无 scope 行视为全场——与 scopeDescOf 口径一致）。活动为管理侧配置数据量级有界,
// 故全量取回→过滤→内存分页, 保证 total 与筛选语义一致。
func (i *MarketingLogicImpl) PublicFullReductions(
	ctx context.Context, spuId int64, page model.PageReq,
) (*model.PageResult[model.PublicFullReductionItem], error) {
	page = page.Normalized()
	cols := dao.PromotionActivity.Columns()
	acts, err := dao.PromotionActivity.Ctx(ctx).
		Where(cols.Status, 1).Where(cols.Deleted, 0).
		Where(cols.StartTime + " <= NOW()").Where(cols.EndTime + " >= NOW()").
		OrderAsc(cols.EndTime).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询满减活动失败")
	}

	// SpuId 范围命中过滤（批量取 scope 行, 判定三级范围; 无 scope 行 = 全场命中）
	if spuId > 0 {
		scols := dao.ProductSpu.Columns()
		cateId, err := dao.ProductSpu.Ctx(ctx).
			Where(scols.Id, spuId).Value(scols.CategoryId)
		if err != nil {
			return nil, gerror.Wrap(err, "查询商品分类失败")
		}
		hit, err := fullReductionScopeHits(ctx, spuId, cateId.Int64())
		if err != nil {
			return nil, err
		}
		filtered := acts[:0]
		for _, a := range acts {
			if hit[a[cols.Id].Int64()] {
				filtered = append(filtered, a)
			}
		}
		acts = filtered
	}

	// 内存分页（total 为过滤后语义）
	total := len(acts)
	start := (page.Page - 1) * page.PageSize
	if start > total {
		start = total
	}
	end := start + page.PageSize
	if end > total {
		end = total
	}
	acts = acts[start:end]

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
		scopeDesc, e := scopeDescOf(ctx, id)
		if e != nil {
			return nil, e // M9: 查询失败必须显式失败, 不得静默标"全场"
		}
		list = append(list, model.PublicFullReductionItem{
			ActivityId: id,
			Name:       a[cols.Name].String(),
			Ladders:    ladders,
			ScopeDesc:  scopeDesc,
			EndTime:    a[cols.EndTime].String(),
		})
	}
	return &model.PageResult[model.PublicFullReductionItem]{List: list, Total: int64(total)}, nil
}

// fullReductionScopeHits 计算满减活动的范围命中集合（批量一次 IN 查询）。
// 命中规则: 有全场行(1) / 商品行(3) target=spuId / 分类行(2) target=该商品分类;
// 完全没有 scope 行的活动视为全场（与 scopeDescOf 展示口径一致）。
func fullReductionScopeHits(ctx context.Context, spuId, cateId int64) (map[int64]bool, error) {
	ids := []int64{}
	actIds, err := dao.PromotionActivity.Ctx(ctx).
		Fields(dao.PromotionActivity.Columns().Id).
		Where(dao.PromotionActivity.Columns().Status, 1).
		Where(dao.PromotionActivity.Columns().Deleted, 0).
		Where(dao.PromotionActivity.Columns().StartTime + " <= NOW()").
		Where(dao.PromotionActivity.Columns().EndTime + " >= NOW()").All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询满减活动失败")
	}
	for _, a := range actIds {
		ids = append(ids, a[dao.PromotionActivity.Columns().Id].Int64())
	}
	out := map[int64]bool{}
	if len(ids) == 0 {
		return out, nil
	}
	scols := dao.PromotionActivityScope.Columns()
	recs, err := dao.PromotionActivityScope.Ctx(ctx).
		WhereIn(scols.ActivityId, ids).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询活动范围失败")
	}
	// 无 scope 行的活动先全部视为全场命中, 再被有 scope 行的活动覆盖
	hasScope := map[int64]bool{}
	for _, r := range recs {
		hasScope[r[scols.ActivityId].Int64()] = true
	}
	for _, id := range ids {
		if !hasScope[id] {
			out[id] = true
		}
	}
	for _, r := range recs {
		id := r[scols.ActivityId].Int64()
		st := r[scols.ScopeType].Int()
		target := r[scols.TargetId].Int64()
		if st == 1 || (st == 3 && target == spuId) || (st == 2 && cateId > 0 && target == cateId) {
			out[id] = true
		}
	}
	return out, nil
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

// spuBriefs 批量取商品名称与主图（一次 IN 查询, 避免 N+1; 主图口径统一走 firstImage——
// M10: 原先 gconv.Struct 重写了一遍"取首图", 与 browse/operation 的口径是两份实现）。
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
		out[r[cols.Id].Int64()] = spuBrief{Name: r[cols.Name].String(), Image: firstImage(r[cols.Images].String())}
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
// M9（评审修复）: 查询失败必须显式返回错误——原实现把任何错误都吞成"全场",
// 会给用户**错误的标签**（把"分类/指定商品"活动标成全场等于误导下单）。
func scopeDescOf(ctx context.Context, actId int64) (string, error) {
	cols := dao.PromotionActivityScope.Columns()
	recs, err := dao.PromotionActivityScope.Ctx(ctx).
		Where(cols.ActivityId, actId).Fields(cols.ScopeType).All()
	if err != nil {
		return "", gerror.Wrap(err, "查询活动范围失败")
	}
	if len(recs) == 0 {
		return "全场", nil
	}
	has := map[int]bool{}
	for _, r := range recs {
		has[r[cols.ScopeType].Int()] = true
	}
	switch {
	case has[1]:
		return "全场", nil
	case has[2]:
		return "分类", nil
	case has[3]:
		return "指定商品", nil
	}
	return "全场", nil
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

	// 各分块失败不阻断首页（"空数组而非报错"针对的是**无内容**; M9: 查询**失败**必须留下告警,
	// 否则 DB 故障时首页静默变空, 可用性掩盖了故障）
	if banners, err := PublicBanners(ctx, 1); err == nil {
		out.Banners = banners
	} else {
		g.Log().Warningf(ctx, "[首页聚合] 轮播分块失败: %v", err)
	}
	if floors, err := PublicFloors(ctx); err == nil {
		out.Floors = floors
	} else {
		g.Log().Warningf(ctx, "[首页聚合] 楼层分块失败: %v", err)
	}
	if fs, err := i.PublicFlashSales(ctx, entryPage); err == nil && fs != nil {
		out.FlashSales = fs.List
	} else if err != nil {
		g.Log().Warningf(ctx, "[首页聚合] 秒杀分块失败: %v", err)
	}
	if gb, err := i.PublicGroupBuys(ctx, entryPage); err == nil && gb != nil {
		out.GroupBuys = gb.List
	} else if err != nil {
		g.Log().Warningf(ctx, "[首页聚合] 拼团分块失败: %v", err)
	}
	if bg, err := i.PublicBargains(ctx, entryPage); err == nil && bg != nil {
		out.Bargains = bg.List
	} else if err != nil {
		g.Log().Warningf(ctx, "[首页聚合] 砍价分块失败: %v", err)
	}
	if as, err := i.PublicAssists(ctx, entryPage); err == nil && as != nil {
		out.Assists = as.List
	} else if err != nil {
		g.Log().Warningf(ctx, "[首页聚合] 助力分块失败: %v", err)
	}
	if fr, err := i.PublicFullReductions(ctx, 0, entryPage); err == nil && fr != nil {
		out.FullReductions = fr.List
	} else if err != nil {
		g.Log().Warningf(ctx, "[首页聚合] 满减分块失败: %v", err)
	}
	// 可领券: shop 域直读券表（**既有先例**: promotion_calc.go 早已直读 coupon/user_coupon 计价）。
	// 注（记账）: 此处口径是"**有哪些券可领**"（模板级: 启用未删 + 未领完）; 而 user 域的
	// `/user/coupons/available` 是"**我还能领哪些**"（会员级: 叠加个人限领判定）——两者回答不同问题,
	// 故不强行统一; 若后续要把会员级口径也搬上首页, 应走既有 `ICouponQuery` 端口而非在此重复实现。
	if cp, err := publicCouponBriefs(ctx, indexEntryLimit); err == nil {
		out.Coupons = cp
	} else {
		g.Log().Warningf(ctx, "[首页聚合] 可领券分块失败: %v", err)
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
