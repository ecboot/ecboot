// browse_impl.go C 端浏览实现（公开; 会员查看详情软触足迹 FR-006）。
// 口径（research D3）: 可售 = SPU上架 AND SKU启用 AND 可售>0; 价格区间取启用 SKU 冗余列。
package shop

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

func (i *ProductLogicImpl) CategoryTree(ctx context.Context) ([]model.CategoryNode, error) {
	all, err := dao.ProductCategory.Ctx(ctx).
		Where(dao.ProductCategory.Columns().Deleted, 0).
		Where(dao.ProductCategory.Columns().Status, 1).
		Order("level asc, sort asc, id asc").All()
	if err != nil {
		return nil, err
	}
	// 中间态含 parentId, 组树后输出契约形态
	type node struct {
		model.CategoryNode
		ParentId int64
	}
	nodes := map[int64]*node{}
	for _, r := range all {
		id := r["id"].Int64()
		nodes[id] = &node{
			CategoryNode: model.CategoryNode{Id: id, Name: r["name"].String(), Icon: r["icon"].String()},
			ParentId:     r["parent_id"].Int64(),
		}
	}
	var walk func(pid int64) []model.CategoryNode
	walk = func(pid int64) []model.CategoryNode {
		var out []model.CategoryNode
		for id, n := range nodes {
			if n.ParentId == pid {
				cp := n.CategoryNode
				cp.Children = walk(id)
				out = append(out, cp)
			}
		}
		return out
	}
	return walk(0), nil
}

func (i *ProductLogicImpl) Brands(ctx context.Context, page model.PageReq) (*model.PageResult[model.BrandItem], error) {
	page = page.Normalized()
	m := dao.ProductBrand.Ctx(ctx).
		Where(dao.ProductBrand.Columns().Deleted, 0).
		Where(dao.ProductBrand.Columns().Status, 1)
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	all, err := m.Page(page.Page, page.PageSize).Order("sort asc, id asc").All()
	if err != nil {
		return nil, err
	}
	list := []model.BrandItem{}
	for _, r := range all {
		list = append(list, model.BrandItem{Id: r["id"].Int64(), Name: r["name"].String(), Logo: r["logo"].String()})
	}
	return &model.PageResult[model.BrandItem]{List: list, Total: int64(total)}, nil
}

// sellableMap 批量可售判定（available>0）。
func sellableMap(ctx context.Context, skuIds []int64) map[int64]bool {
	out := map[int64]bool{}
	if len(skuIds) == 0 {
		return out
	}
	all, err := dao.Inventory.Ctx(ctx).WhereIn(dao.Inventory.Columns().SkuId, skuIds).All()
	if err != nil {
		return out
	}
	for _, r := range all {
		if r["total"].Int()-r["locked"].Int() > 0 {
			out[r["sku_id"].Int64()] = true
		}
	}
	return out
}

func firstImage(imagesJSON string) string {
	var out []string
	if imagesJSON != "" {
		_ = json.Unmarshal([]byte(imagesJSON), &out)
	}
	if len(out) > 0 {
		return out[0]
	}
	return ""
}

func priceRangeText(min, max string) string {
	if min == "" || min == max {
		return max
	}
	return min + "-" + max
}

func (i *ProductLogicImpl) Products(ctx context.Context, q model.ProductQuery) (*model.PageResult[model.ProductCard], error) {
	page := q.Normalized()
	m := dao.ProductSpu.Ctx(ctx).
		Where(dao.ProductSpu.Columns().Deleted, 0).
		Where(dao.ProductSpu.Columns().Status, 1)
	if q.CategoryId > 0 {
		m = m.Where(dao.ProductSpu.Columns().CategoryId, q.CategoryId)
	}
	if q.BrandId > 0 {
		m = m.Where(dao.ProductSpu.Columns().BrandId, q.BrandId)
	}
	if q.PriceMin != "" {
		m = m.WhereGTE("price_min", q.PriceMin)
	}
	if q.PriceMax != "" {
		m = m.WhereLTE("price_max", q.PriceMax)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	orderBy := "sale_count desc, id asc"
	switch q.Sort {
	case 2:
		orderBy = "price_min asc, id asc"
	case 3:
		orderBy = "price_min desc, id asc"
	case 4:
		orderBy = "id desc"
	}
	all, err := m.Page(page.Page, page.PageSize).Order(orderBy).All()
	if err != nil {
		return nil, err
	}
	list := []model.ProductCard{}
	for _, r := range all {
		list = append(list, model.ProductCard{
			SpuId:      r["id"].Int64(),
			Name:       r["name"].String(),
			Image:      firstImage(r["images"].String()),
			PriceRange: priceRangeText(r["price_min"].String(), r["price_max"].String()),
			SaleCount:  r["sale_count"].Int(),
		})
	}
	return &model.PageResult[model.ProductCard]{List: list, Total: int64(total)}, nil
}

// ProductDetail C 端详情（FR-004: 启用 SKU+可售态+评价汇总+限售标记; 不含成本价; 会员软触足迹）。
func (i *ProductLogicImpl) ProductDetail(ctx context.Context, spuId int64, viewerUserId int64) (*model.ProductDetailView, error) {
	rec, err := dao.ProductSpu.Ctx(ctx).
		Where(dao.ProductSpu.Columns().Id, spuId).
		Where(dao.ProductSpu.Columns().Deleted, 0).
		Where(dao.ProductSpu.Columns().Status, 1).One()
	if err != nil {
		return nil, err
	}
	if rec.IsEmpty() {
		return nil, errcode.New(errcode.CodeProductNotFound, "商品不存在或已下架")
	}
	skus, err := dao.ProductSku.Ctx(ctx).
		Where(dao.ProductSku.Columns().SpuId, spuId).
		Where(dao.ProductSku.Columns().Deleted, 0).
		Where(dao.ProductSku.Columns().Status, 1).
		Order("sort asc, id asc").All()
	if err != nil {
		return nil, err
	}
	skuIds := []int64{}
	for _, r := range skus {
		skuIds = append(skuIds, r["id"].Int64())
	}
	sellable := sellableMap(ctx, skuIds)

	summary := model.ReviewSummary{Avg: "—", Distribution: map[string]int{}}
	v, err := g.DB().GetOne(ctx,
		"SELECT COALESCE(ROUND(AVG(score),1),0) AS avg_score, COUNT(*) AS total FROM product_review WHERE spu_id=? AND audit_status=1 AND deleted=0", spuId)
	if err == nil && !v.IsEmpty() {
		summary.Total = v["total"].Int()
		if summary.Total > 0 {
			summary.Avg = v["avg_score"].String()
		}
	}
	dist, err := g.DB().GetAll(ctx,
		"SELECT score, COUNT(*) AS c FROM product_review WHERE spu_id=? AND audit_status=1 AND deleted=0 GROUP BY score", spuId)
	if err == nil {
		for _, r := range dist {
			summary.Distribution[r["score"].String()] = r["c"].Int()
		}
	}

	view := &model.ProductDetailView{
		SpuId:           spuId,
		SpuNo:           rec["spu_no"].String(),
		Name:            rec["name"].String(),
		SubTitle:        rec["sub_title"].String(),
		Images:          spuImages(rec["images"].String()),
		VideoUrl:        rec["video_url"].String(),
		Description:     rec["description"].String(),
		SpecDefinitions: specDefs(rec["spec_definitions"].String()),
		Attributes:      attrMap(rec["attributes"].String()),
		SaleRestricted:  len(attrMapList(rec["sale_restrict_codes"].String())) > 0,
		ReviewSummary:   summary,
	}
	for _, r := range skus {
		view.Skus = append(view.Skus, model.SkuCard{
			SkuId:     r["id"].Int64(),
			SkuNo:     r["sku_no"].String(),
			Specs:     specsMap(r["specs"].String()),
			Price:     r["price"].String(),
			LinePrice: r["line_price"].String(),
			Sellable:  sellable[r["id"].Int64()],
		})
	}
	// 足迹软触（UPSERT; research D7）
	if viewerUserId > 0 {
		_, _ = g.DB().Exec(ctx, `
INSERT INTO user_footprint (user_id, spu_id, view_count, last_view_at)
VALUES (?, ?, 1, NOW())
ON DUPLICATE KEY UPDATE view_count=view_count+1, last_view_at=NOW()`, viewerUserId, spuId)
	}
	return view, nil
}

// Search 关键词搜索（LIKE 转义; FR-005 V1 数据库实现, ES 演进不改契约）。
func (i *ProductLogicImpl) Search(ctx context.Context, keyword string, q model.ProductQuery) (*model.PageResult[model.ProductCard], error) {
	page := q.Normalized()
	escaped := likeEscape(keyword)
	m := dao.ProductSpu.Ctx(ctx).
		Where(dao.ProductSpu.Columns().Deleted, 0).
		Where(dao.ProductSpu.Columns().Status, 1).
		WhereLike(dao.ProductSpu.Columns().Name, "%"+escaped+"%")
	if q.CategoryId > 0 {
		m = m.Where(dao.ProductSpu.Columns().CategoryId, q.CategoryId)
	}
	if q.BrandId > 0 {
		m = m.Where(dao.ProductSpu.Columns().BrandId, q.BrandId)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	orderBy := "sale_count desc, id asc"
	switch q.Sort {
	case 2:
		orderBy = "price_min asc, id asc"
	case 3:
		orderBy = "price_min desc, id asc"
	case 4:
		orderBy = "id desc"
	}
	all, err := m.Page(page.Page, page.PageSize).Order(orderBy).All()
	if err != nil {
		return nil, err
	}
	list := []model.ProductCard{}
	for _, r := range all {
		list = append(list, model.ProductCard{
			SpuId:      r["id"].Int64(),
			Name:       r["name"].String(),
			Image:      firstImage(r["images"].String()),
			PriceRange: priceRangeText(r["price_min"].String(), r["price_max"].String()),
			SaleCount:  r["sale_count"].Int(),
		})
	}
	return &model.PageResult[model.ProductCard]{List: list, Total: int64(total)}, nil
}

// ProductReviews 商品评价列表（审核通过; 星级筛选）。
func (i *ProductLogicImpl) ProductReviews(ctx context.Context, spuId int64, score int, page model.PageReq) (*model.PageResult[model.ReviewCard], error) {
	page = page.Normalized()
	m := dao.ProductReview.Ctx(ctx).
		Where(dao.ProductReview.Columns().SpuId, spuId).
		Where(dao.ProductReview.Columns().AuditStatus, 1).
		Where(dao.ProductReview.Columns().Deleted, 0)
	if score > 0 {
		m = m.Where(dao.ProductReview.Columns().Score, score)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	all, err := m.Page(page.Page, page.PageSize).Order("created_at desc, id desc").All()
	if err != nil {
		return nil, err
	}
	list := []model.ReviewCard{}
	for _, r := range all {
		images := []string{}
		if raw := r["images"].String(); raw != "" {
			_ = json.Unmarshal([]byte(raw), &images)
		}
		specs := map[string]string{}
		if raw := r["sku_specs"].String(); raw != "" {
			_ = json.Unmarshal([]byte(raw), &specs)
		}
		list = append(list, model.ReviewCard{
			ReviewId:  r["id"].Int64(),
			User:      "商城用户",
			Score:     r["score"].Int(),
			Content:   r["content"].String(),
			Images:    images,
			Specs:     specs,
			Reply:     r["reply_content"].String(),
			Extra:     r["extra_content"].String(),
			CreatedAt: r["created_at"].String(),
		})
	}
	return &model.PageResult[model.ReviewCard]{List: list, Total: int64(total)}, nil
}
