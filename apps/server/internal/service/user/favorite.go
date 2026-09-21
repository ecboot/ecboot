// favorite_footprint.go 收藏与足迹——表: user_favorite / user_footprint（V13）。
// 规则: 双表均 (user_id, spu_id) 唯一; 收藏取消=软删、再收藏=复活（ADR-0002 严格唯一）;
// 足迹重复浏览=更新 last_view_at 与 view_count（UPSERT），保留期 90 天物理清理（idx last_view_at 扫描）。
package user

import (
	"context"

	"ecboot/internal/model"
)

// IFavoriteLogic 收藏。
type IFavoriteLogic interface {
	List(ctx context.Context, userId int64, page model.PageReq) (*model.PageResult[model.FavoriteItem], error) // 含实时价态
	Add(ctx context.Context, userId, spuId int64) error                                                        // 复活语义
	Remove(ctx context.Context, userId, spuId int64) error                                                     // 软删
}

// IFootprintLogic 足迹。
type IFootprintLogic interface {
	List(ctx context.Context, userId int64, page model.PageReq) (*model.PageResult[model.FootprintItem], error)
	Clear(ctx context.Context, userId int64) error
	// Record 浏览上报（重复浏览=UPSERT 更新 last_view_at/view_count; 由商品详情接口调用）。
	Record(ctx context.Context, userId, spuId int64) error
	// CleanExpired 90 天清理（定时任务入口）。
	CleanExpired(ctx context.Context) (int64, error)
}
