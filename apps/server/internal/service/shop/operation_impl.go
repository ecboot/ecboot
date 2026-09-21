// operation_impl.go 运营装修实现（接口契约见 misc.go IOperationLogic）。
// 语义: 轮播投放时段（NULL=立即/长期, research D4）; 楼层 config 为不透明 JSON
// （商品楼层约定 {"spuIds":["<id>",...]}, D2）; 商品楼层装配摘要并剔除失效商品（D3）;
// 软删即从管理端与 C 端同时消失。
package shop

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/model"
	"ecboot/internal/model/do"
	"ecboot/internal/model/entity"
)

// ---- 编解码 helper ----

// parseRFC3339 RFC3339 → *gtime.Time（空=nil; 列类型为强类型时间指针, 不能承载 Raw）。
func parseRFC3339(s string) (*gtime.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, errcode.New(errcode.CodeInvalidParam, "时间格式须为 RFC3339")
	}
	return gtime.NewFromTime(t), nil
}

// clearTimeExpr 时段清空的 SET 片段（强类型时间列须经 Raw 显式置 NULL——
// gdb 对 do 自动 OmitNilData 会丢弃 nil, Fields 白名单亦救不回; 同批次 02 教训）。
func clearTimeExpr(startEmpty, endEmpty bool) gdb.Raw {
	parts := make([]string, 0, 2)
	if startEmpty {
		parts = append(parts, "start_time=NULL")
	}
	if endEmpty {
		parts = append(parts, "end_time=NULL")
	}
	return gdb.Raw(strings.Join(parts, ", "))
}

// timeText 库值 → RFC3339（空=零值）。
// 注: 必须用标准库 time.Time.Format——gtime.Format 接受的是 **gf 布局**（Y-m-d H:i:s / c）,
// 传 Go 布局会产出字面串（评审 C1 实证: "2006-01-02CST15:04:05Z07:00"）。
func timeText(t *gtime.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Time.Format(time.RFC3339)
}

// configValue map → JSON 库值; nil 写 NULL。
func configValue(c map[string]any) (any, error) {
	if c == nil {
		return gdb.Raw("NULL"), nil
	}
	b, err := json.Marshal(c)
	if err != nil {
		return nil, gerror.Wrap(err, "序列化楼层配置失败")
	}
	return string(b), nil
}

// configMap JSON → map（空/非法返回空对象, 展示层尽力而为）。
func configMap(raw string) map[string]any {
	out := map[string]any{}
	if raw == "" {
		return out
	}
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return map[string]any{}
	}
	return out
}

func bannerFromEntity(e entity.OperationBanner) model.OperBannerItem {
	return model.OperBannerItem{
		Id:        int64(e.Id),
		Position:  e.Position,
		ImageUrl:  e.ImageUrl,
		LinkUrl:   e.LinkUrl,
		Sort:      e.Sort,
		StartTime: timeText(e.StartTime),
		EndTime:   timeText(e.EndTime),
		Status:    e.Status,
	}
}

func floorFromEntity(e entity.OperationFloor) model.OperFloorItem {
	return model.OperFloorItem{
		Id:        int64(e.Id),
		FloorType: e.FloorType,
		Title:     e.Title,
		Config:    configMap(e.Config),
		Sort:      e.Sort,
		Status:    e.Status,
	}
}

// ---- 轮播管理（US2, FR-007） ----

// BannerList 轮播列表: position>0 筛选; 0=不筛选; sort,id 稳定序。
func BannerList(ctx context.Context, position int, page model.PageReq) (*model.PageResult[model.OperBannerItem], error) {
	page = page.Normalized()
	cols := dao.OperationBanner.Columns()
	m := dao.OperationBanner.Ctx(ctx).Where(cols.Deleted, 0)
	if position > 0 {
		m = m.Where(cols.Position, position)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计轮播失败")
	}
	recs, err := m.Fields(cols.Id, cols.Position, cols.ImageUrl, cols.LinkUrl, cols.Sort,
		cols.StartTime, cols.EndTime, cols.Status).
		Order("sort ASC, id ASC").Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询轮播失败")
	}
	list := make([]model.OperBannerItem, 0, len(recs))
	for _, r := range recs {
		var e entity.OperationBanner
		if err = r.Struct(&e); err != nil {
			return nil, gerror.Wrap(err, "解析轮播失败")
		}
		list = append(list, bannerFromEntity(e))
	}
	return &model.PageResult[model.OperBannerItem]{List: list, Total: int64(total)}, nil
}

// BannerCreate 新增轮播: 位置白名单 + 图片必填; 新建默认启用。
func BannerCreate(ctx context.Context, in model.OperBannerInput) (int64, error) {
	if in.Position != 1 && in.Position != 2 {
		return 0, errcode.New(errcode.CodeInvalidParam, "轮播位置须为1首页轮播或2首页弹窗")
	}
	if in.ImageUrl == "" {
		return 0, errcode.New(errcode.CodeInvalidParam, "图片必填")
	}
	start, err := parseRFC3339(in.StartTime)
	if err != nil {
		return 0, err
	}
	end, err := parseRFC3339(in.EndTime)
	if err != nil {
		return 0, err
	}
	status := in.Status
	if status == 0 {
		status = 1
	}
	res, err := dao.OperationBanner.Ctx(ctx).Data(do.OperationBanner{
		Position:  in.Position,
		ImageUrl:  in.ImageUrl,
		LinkUrl:   in.LinkUrl,
		Sort:      in.Sort,
		StartTime: start,
		EndTime:   end,
		Status:    status,
	}).Insert()
	if err != nil {
		return 0, gerror.Wrap(err, "创建轮播失败")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, gerror.Wrap(err, "读取新轮播ID失败")
	}
	return id, nil
}

// BannerUpdate 修改轮播（图片/链接/排序/时段/状态）; **位置创建后不可改**（api 入参无 position,
// 位置决定 C 端消费场景, 改位需删后重建）; 目标不存在返 10006。
func BannerUpdate(ctx context.Context, id int64, in model.OperBannerInput) error {
	if in.ImageUrl == "" {
		return errcode.New(errcode.CodeInvalidParam, "图片必填")
	}
	if in.Status != 0 && in.Status != 1 {
		return errcode.New(errcode.CodeInvalidParam, "状态须为1启用或0停用")
	}
	start, err := parseRFC3339(in.StartTime)
	if err != nil {
		return err
	}
	end, err := parseRFC3339(in.EndTime)
	if err != nil {
		return err
	}
	cols := dao.OperationBanner.Columns()
	data := do.OperationBanner{
		ImageUrl: in.ImageUrl, LinkUrl: in.LinkUrl, Sort: in.Sort, Status: in.Status,
	}
	fields := []any{cols.ImageUrl, cols.LinkUrl, cols.Sort, cols.Status}
	if start != nil {
		data.StartTime = start
		fields = append(fields, cols.StartTime)
	}
	if end != nil {
		data.EndTime = end
		fields = append(fields, cols.EndTime)
	}
	// 评审 I2: 主体更新与时段置空同事务（避免半更新——图片/状态已改而时段未清）
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, err := dao.OperationBanner.Ctx(ctx).
			Where(cols.Id, id).Where(cols.Deleted, 0).
			Data(data).Fields(fields...).
			Update()
		if err != nil {
			return gerror.Wrap(err, "修改轮播失败")
		}
		if n, _ := res.RowsAffected(); n == 0 {
			cnt, cerr := dao.OperationBanner.Ctx(ctx).
				Where(cols.Id, id).Where(cols.Deleted, 0).Count()
			if cerr != nil {
				return gerror.Wrap(cerr, "查询轮播失败")
			}
			if cnt == 0 {
				return errcode.New(errcode.CodeNotFound, "轮播不存在")
			}
		}
		// 时段清空（空串=清空该端 → 显式置 NULL, 强类型列须经 Raw）
		if in.StartTime == "" || in.EndTime == "" {
			if _, err = dao.OperationBanner.Ctx(ctx).
				Where(cols.Id, id).Where(cols.Deleted, 0).
				Data(clearTimeExpr(in.StartTime == "", in.EndTime == "")).
				Update(); err != nil {
				return gerror.Wrap(err, "清空投放时段失败")
			}
		}
		return nil
	})
}

// BannerDelete 软删轮播（FR-009）: 管理端与 C 端同时消失。
func BannerDelete(ctx context.Context, id int64) error {
	cols := dao.OperationBanner.Columns()
	res, err := dao.OperationBanner.Ctx(ctx).
		Where(cols.Id, id).Where(cols.Deleted, 0).
		Data(do.OperationBanner{Deleted: 1}).
		Fields(cols.Deleted).
		Update()
	if err != nil {
		return gerror.Wrap(err, "删除轮播失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeNotFound, "轮播不存在")
	}
	return nil
}

// ---- 楼层管理（US2, FR-008） ----

// FloorList 楼层列表: sort,id 稳定序, config 反序列化返回。
func FloorList(ctx context.Context, page model.PageReq) (*model.PageResult[model.OperFloorItem], error) {
	page = page.Normalized()
	cols := dao.OperationFloor.Columns()
	m := dao.OperationFloor.Ctx(ctx).Where(cols.Deleted, 0)
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计楼层失败")
	}
	recs, err := m.Fields(cols.Id, cols.FloorType, cols.Title, cols.Config, cols.Sort, cols.Status).
		Order("sort ASC, id ASC").Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询楼层失败")
	}
	list := make([]model.OperFloorItem, 0, len(recs))
	for _, r := range recs {
		var e entity.OperationFloor
		if err = r.Struct(&e); err != nil {
			return nil, gerror.Wrap(err, "解析楼层失败")
		}
		list = append(list, floorFromEntity(e))
	}
	return &model.PageResult[model.OperFloorItem]{List: list, Total: int64(total)}, nil
}

// FloorCreate 新增楼层: 类型白名单（三型）; 新建默认启用; config 为不透明 JSON。
func FloorCreate(ctx context.Context, in model.OperFloorInput) (int64, error) {
	if in.FloorType < 1 || in.FloorType > 3 {
		return 0, errcode.New(errcode.CodeInvalidParam, "楼层类型须为1金刚区2商品楼层3专题")
	}
	cfg, err := configValue(in.Config)
	if err != nil {
		return 0, err
	}
	status := in.Status
	if status == 0 {
		status = 1
	}
	res, err := dao.OperationFloor.Ctx(ctx).Data(do.OperationFloor{
		FloorType: in.FloorType, Title: in.Title, Config: cfg, Sort: in.Sort, Status: status,
	}).Insert()
	if err != nil {
		return 0, gerror.Wrap(err, "创建楼层失败")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, gerror.Wrap(err, "读取新楼层ID失败")
	}
	return id, nil
}

// FloorUpdate 修改楼层（标题/配置/排序/状态）; **楼层类型创建后不可改**（api 入参无 floorType,
// 类型决定 config 语义, 改型需删后重建）; 目标不存在返 10006。
func FloorUpdate(ctx context.Context, id int64, in model.OperFloorInput) error {
	if in.Status != 0 && in.Status != 1 {
		return errcode.New(errcode.CodeInvalidParam, "状态须为1启用或0停用")
	}
	cfg, err := configValue(in.Config)
	if err != nil {
		return err
	}
	cols := dao.OperationFloor.Columns()
	res, err := dao.OperationFloor.Ctx(ctx).
		Where(cols.Id, id).Where(cols.Deleted, 0).
		Data(do.OperationFloor{
			Title: in.Title, Config: cfg, Sort: in.Sort, Status: in.Status,
		}).
		Fields(cols.Title, cols.Config, cols.Sort, cols.Status).
		Update()
	if err != nil {
		return gerror.Wrap(err, "修改楼层失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		cnt, cerr := dao.OperationFloor.Ctx(ctx).
			Where(cols.Id, id).Where(cols.Deleted, 0).Count()
		if cerr != nil {
			return gerror.Wrap(cerr, "查询楼层失败")
		}
		if cnt == 0 {
			return errcode.New(errcode.CodeNotFound, "楼层不存在")
		}
	}
	return nil
}

// FloorDelete 软删楼层（FR-009）: 管理端与 C 端同时消失。
func FloorDelete(ctx context.Context, id int64) error {
	cols := dao.OperationFloor.Columns()
	res, err := dao.OperationFloor.Ctx(ctx).
		Where(cols.Id, id).Where(cols.Deleted, 0).
		Data(do.OperationFloor{Deleted: 1}).
		Fields(cols.Deleted).
		Update()
	if err != nil {
		return gerror.Wrap(err, "删除楼层失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeNotFound, "楼层不存在")
	}
	return nil
}

// ---- C 端展示（US3, FR-010~013） ----

// PublicBanners 在投轮播（FR-010）: 启用 + 投放时段内（NULL=立即/长期, research D4）, sort,id 序。
func PublicBanners(ctx context.Context, position int) ([]model.PublicBannerItem, error) {
	if position != 1 && position != 2 {
		return nil, errcode.New(errcode.CodeInvalidParam, "轮播位置须为1首页轮播或2首页弹窗")
	}
	cols := dao.OperationBanner.Columns()
	recs, err := dao.OperationBanner.Ctx(ctx).
		Fields(cols.Id, cols.ImageUrl, cols.LinkUrl).
		Where(cols.Deleted, 0).
		Where(cols.Status, 1).
		Where(cols.Position, position).
		Where("(start_time IS NULL OR start_time <= NOW())").
		Where("(end_time IS NULL OR end_time >= NOW())").
		Order("sort ASC, id ASC").
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询在投轮播失败")
	}
	list := make([]model.PublicBannerItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, model.PublicBannerItem{
			Id:       r["id"].Int64(),
			ImageUrl: r["image_url"].String(),
			LinkUrl:  r["link_url"].String(),
		})
	}
	return list, nil
}

// spuIdsFromConfig 解析商品楼层配置的商品 ID 列表（约定 {"spuIds":["<id>",...]}, research D2）;
// 容错: 缺失/类型不符/非法值一律跳过（展示层尽力而为, 不因脏配置报错）。
func spuIdsFromConfig(cfg map[string]any) []int64 {
	raw, ok := cfg["spuIds"]
	if !ok {
		return nil
	}
	arr, ok := raw.([]any)
	if !ok {
		return nil
	}
	out := make([]int64, 0, len(arr))
	for _, v := range arr {
		switch t := v.(type) {
		case string:
			if id, err := strconv.ParseInt(t, 10, 64); err == nil && id > 0 {
				out = append(out, id)
			}
		case float64: // JSON 数字形态容错
			if t > 0 {
				out = append(out, int64(t))
			}
		}
	}
	return out
}

// PublicFloors C 端楼层（FR-011/012）: 启用楼层按序; 商品楼层装配 SPU 摘要并剔除失效商品（research D3）。
func PublicFloors(ctx context.Context) ([]model.PublicFloorItem, error) {
	cols := dao.OperationFloor.Columns()
	recs, err := dao.OperationFloor.Ctx(ctx).
		Fields(cols.Id, cols.FloorType, cols.Title, cols.Config, cols.Sort).
		Where(cols.Deleted, 0).
		Where(cols.Status, 1).
		Order("sort ASC, id ASC").
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询楼层失败")
	}
	list := make([]model.PublicFloorItem, 0, len(recs))
	allIds := make([]int64, 0)
	for _, r := range recs {
		var e entity.OperationFloor
		if err = r.Struct(&e); err != nil {
			return nil, gerror.Wrap(err, "解析楼层失败")
		}
		cfg := configMap(e.Config)
		item := model.PublicFloorItem{
			FloorId:   int64(e.Id),
			FloorType: e.FloorType,
			Title:     e.Title,
			Config:    cfg,
			Products:  []model.FloorProductSummary{},
		}
		if e.FloorType == 2 {
			allIds = append(allIds, spuIdsFromConfig(cfg)...)
		}
		list = append(list, item)
	}

	// 批量装配: 一次查全部涉及 SPU（上架且未删）, 失效 ID 自然缺失 → 逐个剔除
	briefs := map[int64]model.FloorProductSummary{}
	if len(allIds) > 0 {
		spuCols := dao.ProductSpu.Columns()
		spuRecs, serr := dao.ProductSpu.Ctx(ctx).
			Fields(spuCols.Id, spuCols.Name, spuCols.Images, spuCols.PriceMin).
			Where(spuCols.Deleted, 0).
			Where(spuCols.Status, 1).
			WhereIn(spuCols.Id, allIds).
			All()
		if serr != nil {
			return nil, gerror.Wrap(serr, "装配楼层商品失败")
		}
		for _, r := range spuRecs {
			id := r["id"].Int64()
			briefs[id] = model.FloorProductSummary{
				SpuId: id,
				Name:  r["name"].String(),
				Image: firstImage(r["images"].String()),
				Price: r["price_min"].String(),
			}
		}
	}
	for i := range list {
		if list[i].FloorType != 2 {
			continue
		}
		for _, sid := range spuIdsFromConfig(list[i].Config) {
			if b, ok := briefs[sid]; ok {
				list[i].Products = append(list[i].Products, b)
			}
		}
	}
	return list, nil
}
