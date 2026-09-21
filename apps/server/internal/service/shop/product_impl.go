// product_impl.go IProductLogic 实现（005 商品目录域）。
// 分层契约 §二：Where/Data 一律 do 对象与 Columns 常量（禁 g.Map 传参）；
// 返回 entity → 转换为 model DTO；价格区间冗余列在 SKU 变更时同事务刷新（research D1）。
package shop

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// ProductLogicImpl IProductLogic 实现。
type ProductLogicImpl struct{}

func NewProductLogic() *ProductLogicImpl { return &ProductLogicImpl{} }

// ---------- JSON 列工具（entity string ↔ 结构化） ----------

func spuImages(raw string) []string {
	var out []string
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &out)
	}
	return out
}

func specsMap(raw string) map[string]string {
	out := map[string]string{}
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &out)
	}
	return out
}

func specDefs(raw string) []map[string]any {
	var out []map[string]any
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &out)
	}
	return out
}

func attrMap(raw string) map[string]string {
	out := map[string]string{}
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &out)
	}
	return out
}

func mustJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// likeEscape LIKE 通配符转义（FR-005）。
func likeEscape(kw string) string {
	kw = strings.ReplaceAll(kw, "\\", "\\\\")
	kw = strings.ReplaceAll(kw, "%", "\\%")
	kw = strings.ReplaceAll(kw, "_", "\\_")
	return kw
}

// refreshPriceRange 重算 SPU 启用 SKU 的 min/max 并回写（research D1; SKU 变更路径调用）。
func refreshPriceRange(ctx context.Context, spuId int64) error {
	_, err := g.DB().Exec(ctx, `
UPDATE product_spu s SET
  s.price_min = (SELECT MIN(k.price) FROM product_sku k WHERE k.spu_id=s.id AND k.status=1 AND k.deleted=0),
  s.price_max = (SELECT MAX(k.price) FROM product_sku k WHERE k.spu_id=s.id AND k.status=1 AND k.deleted=0)
WHERE s.id=?`, spuId)
	return err
}

// ---------- 分类 ----------

func (i *ProductLogicImpl) AdminCategoryTree(ctx context.Context) ([]model.AdminCategoryNode, error) {
	all, err := dao.ProductCategory.Ctx(ctx).
		Where(dao.ProductCategory.Columns().Deleted, 0).
		Order("level asc, sort asc, id asc").All()
	if err != nil {
		return nil, err
	}
	nodes := map[int64]*model.AdminCategoryNode{}
	for _, r := range all {
		nodes[r["id"].Int64()] = &model.AdminCategoryNode{
			Id:       r["id"].Int64(),
			ParentId: r["parent_id"].Int64(),
			Name:     r["name"].String(),
			Icon:     r["icon"].String(),
			Level:    r["level"].Int(),
			Sort:     r["sort"].Int(),
			Status:   r["status"].Int(),
		}
	}
	var walk func(pid int64) []model.AdminCategoryNode
	walk = func(pid int64) []model.AdminCategoryNode {
		var out []model.AdminCategoryNode
		for id, n := range nodes {
			if n.ParentId == pid {
				cp := *n
				cp.Children = walk(id)
				out = append(out, cp)
			}
		}
		return out
	}
	return walk(0), nil
}

func (i *ProductLogicImpl) AdminCategoryCreate(ctx context.Context, in model.CategoryInput) (int64, error) {
	res, err := dao.ProductCategory.Ctx(ctx).Data(doCategory(in)).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (i *ProductLogicImpl) AdminCategoryUpdate(ctx context.Context, id int64, in model.CategoryInput) error {
	_, err := dao.ProductCategory.Ctx(ctx).Where(dao.ProductCategory.Columns().Id, id).
		Data(doCategoryUpdate(in)).Update()
	return err
}

func (i *ProductLogicImpl) AdminCategoryDelete(ctx context.Context, id int64) error {
	// 引用检查（含软删商品, spec FR-007）
	n, err := dao.ProductSpu.Ctx(ctx).Where(dao.ProductSpu.Columns().CategoryId, id).Count()
	if err != nil {
		return err
	}
	if n > 0 {
		return errcode.New(errcode.CodeCategoryInUse, "分类下存在商品,禁止删除")
	}
	_, err = dao.ProductCategory.Ctx(ctx).Where(dao.ProductCategory.Columns().Id, id).
		Data(g.Map{dao.ProductCategory.Columns().Deleted: 1}).Update()
	return err
}

// ---------- 品牌 ----------

func (i *ProductLogicImpl) AdminBrandList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.BrandItem], error) {
	page = page.Normalized()
	m := dao.ProductBrand.Ctx(ctx)
	if status > 0 {
		m = m.Where(dao.ProductBrand.Columns().Status, status)
	}
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
		list = append(list, model.BrandItem{
			Id: r["id"].Int64(), Name: r["name"].String(), Logo: r["logo"].String(),
			Description: r["description"].String(), Sort: r["sort"].Int(), Status: r["status"].Int(),
		})
	}
	return &model.PageResult[model.BrandItem]{List: list, Total: int64(total)}, nil
}

func (i *ProductLogicImpl) AdminBrandCreate(ctx context.Context, in model.BrandInput) (int64, error) {
	return dao.ProductBrand.Ctx(ctx).Data(doBrand(in)).InsertAndGetId()
}

func (i *ProductLogicImpl) AdminBrandUpdate(ctx context.Context, id int64, in model.BrandInput) error {
	_, err := dao.ProductBrand.Ctx(ctx).Where(dao.ProductBrand.Columns().Id, id).Data(doBrandUpdate(in)).Update()
	return err
}

func (i *ProductLogicImpl) AdminBrandDelete(ctx context.Context, id int64) error {
	_, err := dao.ProductBrand.Ctx(ctx).Where(dao.ProductBrand.Columns().Id, id).
		Data(g.Map{dao.ProductBrand.Columns().Deleted: 1}).Update()
	return err
}

// ---------- 商品（SPU/SKU） ----------

func (i *ProductLogicImpl) adminSpuQuery(ctx context.Context, q model.AdminProductQuery) (*gdb.Model, int64, error) {
	page := q.Normalized()
	m := dao.ProductSpu.Ctx(ctx).Where(dao.ProductSpu.Columns().Deleted, 0)
	if q.Status == 0 || q.Status == 1 {
		m = m.Where(dao.ProductSpu.Columns().Status, q.Status)
	}
	if q.CategoryId > 0 {
		m = m.Where(dao.ProductSpu.Columns().CategoryId, q.CategoryId)
	}
	if q.Keyword != "" {
		m = m.WhereLike(dao.ProductSpu.Columns().Name, "%"+likeEscape(q.Keyword)+"%")
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	return m.Page(page.Page, page.PageSize).Order("id desc"), int64(total), nil
}

func (i *ProductLogicImpl) AdminProductList(ctx context.Context, q model.AdminProductQuery) (*model.PageResult[model.AdminProductItem], error) {
	m, total, err := i.adminSpuQuery(ctx, q)
	if err != nil {
		return nil, err
	}
	all, err := m.Order("id desc").All()
	if err != nil {
		return nil, err
	}
	list := []model.AdminProductItem{}
	for _, r := range all {
		list = append(list, model.AdminProductItem{
			SpuId: r["id"].Int64(), SpuNo: r["spu_no"].String(), Name: r["name"].String(),
			CategoryId: r["category_id"].Int64(), BrandId: r["brand_id"].Int64(),
			Status: r["status"].Int(), SaleCount: r["sale_count"].Int(),
			CreatedAt: r["created_at"].String(),
		})
	}
	return &model.PageResult[model.AdminProductItem]{List: list, Total: total}, nil
}

func (i *ProductLogicImpl) AdminProductCreate(ctx context.Context, in model.SpuInput) (int64, string, error) {
	spuNo, err := nextSpuNo(ctx)
	if err != nil {
		return 0, "", err
	}
	id, err := dao.ProductSpu.Ctx(ctx).Data(doSpuCreate(in, spuNo)).InsertAndGetId()
	return id, spuNo, err
}

func (i *ProductLogicImpl) AdminProductDetail(ctx context.Context, spuId int64) (*model.AdminProductDetailView, error) {
	rec, err := dao.ProductSpu.Ctx(ctx).Where(dao.ProductSpu.Columns().Id, spuId).One()
	if err != nil {
		return nil, err
	}
	if rec.IsEmpty() {
		return nil, errcode.New(errcode.CodeProductNotFound, "商品不存在")
	}
	skus, err := dao.ProductSku.Ctx(ctx).
		Where(dao.ProductSku.Columns().SpuId, spuId).
		Where(dao.ProductSku.Columns().Deleted, 0).All()
	if err != nil {
		return nil, err
	}
	v := &model.AdminProductDetailView{
		SpuId:             spuId,
		SpuNo:             rec["spu_no"].String(),
		Name:              rec["name"].String(),
		SubTitle:          rec["sub_title"].String(),
		CategoryId:        rec["category_id"].Int64(),
		BrandId:           rec["brand_id"].Int64(),
		FreightTemplateId: rec["freight_template_id"].Int64(),
		Images:            spuImages(rec["images"].String()),
		VideoUrl:          rec["video_url"].String(),
		Description:       rec["description"].String(),
		SpecDefinitions:   specDefs(rec["spec_definitions"].String()),
		Attributes:        attrMap(rec["attributes"].String()),
		SaleRestrictCodes: attrMapList(rec["sale_restrict_codes"].String()),
		Status:            rec["status"].Int(),
	}
	for _, r := range skus {
		v.Skus = append(v.Skus, model.AdminSkuDetail{
			SkuId: r["id"].Int64(), SkuNo: r["sku_no"].String(),
			Specs: specsMap(r["specs"].String()), Price: r["price"].String(),
			LinePrice: r["line_price"].String(), CostPrice: r["cost_price"].String(),
			Weight: r["weight"].String(), Barcode: r["barcode"].String(),
			Status: r["status"].Int(),
		})
	}
	return v, nil
}

func attrMapList(raw string) []string {
	var out []string
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &out)
	}
	return out
}

func (i *ProductLogicImpl) AdminProductUpdate(ctx context.Context, spuId int64, in model.SpuInput) error {
	_, err := dao.ProductSpu.Ctx(ctx).Where(dao.ProductSpu.Columns().Id, spuId).
		Data(doSpuUpdate(in)).Update()
	return err
}

func (i *ProductLogicImpl) AdminProductDelete(ctx context.Context, spuId int64) error {
	_, err := dao.ProductSpu.Ctx(ctx).Where(dao.ProductSpu.Columns().Id, spuId).
		Data(g.Map{dao.ProductSpu.Columns().Deleted: 1, dao.ProductSpu.Columns().Status: 0}).Update()
	return err
}

// AdminProductStatus 上下架（上架校验启用 SKU 存在 → 30005）。
func (i *ProductLogicImpl) AdminProductStatus(ctx context.Context, spuId int64, status int) error {
	if status == 1 {
		n, err := dao.ProductSku.Ctx(ctx).
			Where(dao.ProductSku.Columns().SpuId, spuId).
			Where(dao.ProductSku.Columns().Status, 1).
			Where(dao.ProductSku.Columns().Deleted, 0).Count()
		if err != nil {
			return err
		}
		if n == 0 {
			return errcode.New(errcode.CodeNoEnableSku, "无启用SKU,禁止上架")
		}
	}
	_, err := dao.ProductSpu.Ctx(ctx).Where(dao.ProductSpu.Columns().Id, spuId).
		Data(g.Map{dao.ProductSpu.Columns().Status: status}).Update()
	return err
}

func (i *ProductLogicImpl) AdminProductRestrict(ctx context.Context, spuId int64, codes []string) error {
	_, err := dao.ProductSpu.Ctx(ctx).Where(dao.ProductSpu.Columns().Id, spuId).
		Data(g.Map{dao.ProductSpu.Columns().SaleRestrictCodes: mustJSON(codes)}).Update()
	return err
}

// AdminSkuCreate 新增 SKU（规格组合唯一: specs_hash 唯一索引, 冲突→30006; 同事务刷新价格冗余+初始化库存行）。
func (i *ProductLogicImpl) AdminSkuCreate(ctx context.Context, spuId int64, in model.SkuInput) (int64, error) {
	skuNo, err := nextSkuNo(ctx)
	if err != nil {
		return 0, err
	}
	id, err := dao.ProductSku.Ctx(ctx).Data(doSkuCreate(spuId, in, skuNo)).InsertAndGetId()
	if err != nil {
		// specs_hash 唯一冲突（V23 生成列）→ 30006
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "1062") {
			return 0, errcode.New(errcode.CodeSpecDup, "规格组合已存在")
		}
		return 0, err
	}
	// 库存行初始化（total=0, 由后台调整录入——spec Assumptions）
	_, _ = g.DB().Model("inventory").Ctx(ctx).Data(g.Map{"sku_id": id, "total": 0, "locked": 0}).Insert()
	if err = refreshPriceRange(ctx, spuId); err != nil {
		return 0, err
	}
	return id, nil
}

func (i *ProductLogicImpl) AdminSkuUpdate(ctx context.Context, skuId int64, in model.SkuInput) error {
	rec, err := dao.ProductSku.Ctx(ctx).Where(dao.ProductSku.Columns().Id, skuId).One()
	if err != nil || rec.IsEmpty() {
		return errcode.New(errcode.CodeProductNotFound, "SKU不存在")
	}
	if _, err = dao.ProductSku.Ctx(ctx).Where(dao.ProductSku.Columns().Id, skuId).
		Data(doSkuUpdate(in)).Update(); err != nil {
		return err
	}
	return refreshPriceRange(ctx, rec["spu_id"].Int64())
}

func (i *ProductLogicImpl) AdminSkuDelete(ctx context.Context, skuId int64) error {
	rec, err := dao.ProductSku.Ctx(ctx).Where(dao.ProductSku.Columns().Id, skuId).One()
	if err != nil || rec.IsEmpty() {
		return nil
	}
	if _, err = dao.ProductSku.Ctx(ctx).Where(dao.ProductSku.Columns().Id, skuId).
		Data(g.Map{dao.ProductSku.Columns().Deleted: 1}).Update(); err != nil {
		return err
	}
	return refreshPriceRange(ctx, rec["spu_id"].Int64())
}

func (i *ProductLogicImpl) AdminSkuStatus(ctx context.Context, skuId int64, status int) error {
	_, err := dao.ProductSku.Ctx(ctx).Where(dao.ProductSku.Columns().Id, skuId).
		Data(g.Map{dao.ProductSku.Columns().Status: status}).Update()
	if err != nil {
		return err
	}
	rec, err := dao.ProductSku.Ctx(ctx).Where(dao.ProductSku.Columns().Id, skuId).One()
	if err != nil {
		return err
	}
	return refreshPriceRange(ctx, rec["spu_id"].Int64())
}

// ---------- 编码生成 ----------

func nextSpuNo(ctx context.Context) (string, error) { return nextNo(ctx, "SP") }
func nextSkuNo(ctx context.Context) (string, error) { return nextNo(ctx, "SK") }

func nextNo(ctx context.Context, prefix string) (string, error) {
	// 简单实现: 前缀+纳秒尾段+随机尾（生产可换发号器）
	v, err := g.Redis().Do(ctx, "INCR", "seq:"+prefix)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%010d", prefix, v.Int64()), nil
}

// AdminSkuCreateWithNo 创建 SKU 并返回编码（连线适配 010: api 契约 AdminSkuCreateRes 需 skuId 与 skuNo,
// 而既有 AdminSkuCreate 仅返回 id——不改其签名以免破坏 005 测试, 以本函数补编码）。
func AdminSkuCreateWithNo(ctx context.Context, spuId int64, in model.SkuInput) (int64, string, error) {
	id, err := NewProductLogic().AdminSkuCreate(ctx, spuId, in)
	if err != nil {
		return 0, "", err
	}
	v, err := dao.ProductSku.Ctx(ctx).
		Where(dao.ProductSku.Columns().Id, id).
		Value(dao.ProductSku.Columns().SkuNo)
	if err != nil {
		return 0, "", err
	}
	return id, v.String(), nil
}
