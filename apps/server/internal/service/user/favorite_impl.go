// favorite_impl.go 收藏与浏览足迹（011-member-center; 接口契约见 favorite.go）。
// 收藏: 取消=软删, 再收藏=**复活**（复用原行解软删, 不新插——uk(user_id,spu_id) 约束下会冲突）;
// 列表含实时价态（价格取 price_min, 与商品浏览同口径）与失效标注（下架/软删仍可见但标 Invalid）。
// 足迹: Record = UPSERT（重复浏览 view_count+1 并刷新 last_view_at）; CleanExpired 清 90 天前。
package user

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"ecboot/internal/dao"
	"ecboot/internal/model"
	"ecboot/internal/model/do"
)

// spuFirstImage 取 SPU 主图（images JSON 数组首项; 会员域自包含, 不复用 shop 包私有 helper）。
func spuFirstImage(raw string) string {
	if raw == "" {
		return ""
	}
	var arr []string
	if err := json.Unmarshal([]byte(raw), &arr); err != nil || len(arr) == 0 {
		return ""
	}
	return arr[0]
}

// FavoriteList 收藏列表（FR-007）: 含实时价态与失效标注。
func FavoriteList(ctx context.Context, userId int64, page model.PageReq) (*model.PageResult[model.FavoriteItem], error) {
	page = page.Normalized()
	fcols, scols := dao.UserFavorite.Columns(), dao.ProductSpu.Columns()
	m := dao.UserFavorite.Ctx(ctx).As("f").
		LeftJoin(dao.ProductSpu.Table()+" s", "s.id=f.spu_id").
		Where("f."+fcols.UserId, userId).
		Where("f."+fcols.Deleted, 0)
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计收藏失败")
	}
	recs, err := m.Fields("f."+fcols.SpuId+", s."+scols.Name+", s."+scols.Images+", s."+scols.PriceMin+
		", s."+scols.Status+", s."+scols.Deleted+" AS spu_deleted").
		OrderDesc("f."+fcols.Id).
		Page(page.Page, page.PageSize).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询收藏失败")
	}
	list := make([]model.FavoriteItem, 0, len(recs))
	for _, r := range recs {
		// 失效 = 商品不存在/已删/未上架（仍返回, 前端置灰——不静默消失）
		invalid := r["name"].String() == "" || r["status"].Int() != 1 || r["spu_deleted"].Int() == 1
		list = append(list, model.FavoriteItem{
			SpuId:    r[fcols.SpuId].Int64(),
			Name:     r["name"].String(),
			Image:    spuFirstImage(r["images"].String()),
			Price:    r["price_min"].String(),
			Sellable: !invalid,
			Invalid:  invalid,
		})
	}
	return &model.PageResult[model.FavoriteItem]{List: list, Total: int64(total)}, nil
}

// FavoriteAdd 收藏商品（FR-007）: 已存在（含软删）→ 复活; 不存在 → 新建; 全程幂等。
func FavoriteAdd(ctx context.Context, userId, spuId int64) error {
	cols := dao.UserFavorite.Columns()
	cnt, err := dao.UserFavorite.Ctx(ctx).
		Where(cols.UserId, userId).Where(cols.SpuId, spuId).Count()
	if err != nil {
		return gerror.Wrap(err, "查询收藏失败")
	}
	if cnt > 0 {
		// 复活: 解软删（已在位时 RowsAffected=0 亦无妨——幂等）
		if _, err = dao.UserFavorite.Ctx(ctx).
			Where(cols.UserId, userId).Where(cols.SpuId, spuId).
			Data(do.UserFavorite{Deleted: 0}).
			Fields(cols.Deleted).
			Update(); err != nil {
			return gerror.Wrap(err, "收藏失败")
		}
		return nil
	}
	_, err = dao.UserFavorite.Ctx(ctx).Data(do.UserFavorite{UserId: userId, SpuId: spuId}).Insert()
	if err != nil {
		// 并发下 uk 冲突视为已收藏（幂等）
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "1062") {
			return nil
		}
		return gerror.Wrap(err, "收藏失败")
	}
	return nil
}

// FavoriteRemove 取消收藏（FR-007）: 软删, 幂等。
func FavoriteRemove(ctx context.Context, userId, spuId int64) error {
	cols := dao.UserFavorite.Columns()
	_, err := dao.UserFavorite.Ctx(ctx).
		Where(cols.UserId, userId).
		Where(cols.SpuId, spuId).
		Where(cols.Deleted, 0).
		Data(do.UserFavorite{Deleted: 1}).
		Fields(cols.Deleted).
		Update()
	if err != nil {
		return gerror.Wrap(err, "取消收藏失败")
	}
	return nil
}

// FootprintList 足迹列表（FR-008）: 最近浏览倒序。
func FootprintList(ctx context.Context, userId int64, page model.PageReq) (*model.PageResult[model.FootprintItem], error) {
	page = page.Normalized()
	fcols, scols := dao.UserFootprint.Columns(), dao.ProductSpu.Columns()
	m := dao.UserFootprint.Ctx(ctx).As("fp").
		LeftJoin(dao.ProductSpu.Table()+" s", "s.id=fp.spu_id").
		Where("fp."+fcols.UserId, userId)
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计足迹失败")
	}
	recs, err := m.Fields("fp."+fcols.SpuId+", fp."+fcols.LastViewAt+", s."+scols.Name+", s."+scols.Images+", s."+scols.PriceMin).
		OrderDesc("fp."+fcols.LastViewAt).
		Page(page.Page, page.PageSize).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询足迹失败")
	}
	list := make([]model.FootprintItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, model.FootprintItem{
			SpuId:      r[fcols.SpuId].Int64(),
			Name:       r["name"].String(),
			Image:      spuFirstImage(r["images"].String()),
			Price:      r["price_min"].String(),
			LastViewAt: r[fcols.LastViewAt].String(),
		})
	}
	return &model.PageResult[model.FootprintItem]{List: list, Total: int64(total)}, nil
}

// FootprintClear 清空足迹（FR-008）。
func FootprintClear(ctx context.Context, userId int64) error {
	_, err := dao.UserFootprint.Ctx(ctx).
		Where(dao.UserFootprint.Columns().UserId, userId).
		Delete()
	if err != nil {
		return gerror.Wrap(err, "清空足迹失败")
	}
	return nil
}

// FootprintRecord 浏览上报（内部方法, FR-008）: UPSERT——重复浏览计数递增并刷新时间。
func FootprintRecord(ctx context.Context, userId, spuId int64) error {
	cols := dao.UserFootprint.Columns()
	cnt, err := dao.UserFootprint.Ctx(ctx).
		Where(cols.UserId, userId).Where(cols.SpuId, spuId).Count()
	if err != nil {
		return gerror.Wrap(err, "查询足迹失败")
	}
	now := gtime.Now()
	if cnt > 0 {
		_, err = dao.UserFootprint.Ctx(ctx).
			Where(cols.UserId, userId).Where(cols.SpuId, spuId).
			Data(do.UserFootprint{
				ViewCount:  gdb.Raw("view_count + 1"),
				LastViewAt: now,
			}).
			Fields(cols.ViewCount, cols.LastViewAt).
			Update()
	} else {
		_, err = dao.UserFootprint.Ctx(ctx).Data(do.UserFootprint{
			UserId: userId, SpuId: spuId, ViewCount: 1, LastViewAt: now,
		}).Insert()
	}
	if err != nil {
		return gerror.Wrap(err, "记录足迹失败")
	}
	return nil
}

// FootprintCleanExpired 清理 90 天前足迹（内部方法; 定时任务入口）。
func FootprintCleanExpired(ctx context.Context) (int64, error) {
	res, err := dao.UserFootprint.Ctx(ctx).
		Where("last_view_at < DATE_SUB(NOW(), INTERVAL 90 DAY)").
		Delete()
	if err != nil {
		return 0, gerror.Wrap(err, "清理足迹失败")
	}
	n, _ := res.RowsAffected()
	return n, nil
}
