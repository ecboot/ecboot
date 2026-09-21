package user

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/model"
)

// TestFavoriteLifecycle 收藏（FR-007）: 实时价态 / 取消软删 / 再收藏复活（无重复行）。
func TestFavoriteLifecycle(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			phone = "13900004001"
			sfx   = "fav"
		)
		defer cleanupMember(ctx, t, phone)
		defer cleanupSpuForMember(ctx, t, sfx)
		uid := seedMember(ctx, t, phone, "收藏测试", 0)
		defer cleanupMemberCollections(ctx, t, uid)
		spuId := seedSpuForMember(ctx, t, sfx, "29.90", `["http://img/um.png"]`, 1)

		// 收藏（幂等）
		t.AssertNil(FavoriteAdd(ctx, uid, spuId))
		t.AssertNil(FavoriteAdd(ctx, uid, spuId)) // 重复收藏不报错（复活语义覆盖已在位场景）

		// 列表含实时价态
		res, err := FavoriteList(ctx, uid, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.Total, 1)
		t.Assert(res.List[0].SpuId, spuId)
		t.Assert(res.List[0].Name, "UM商品"+sfx)
		t.Assert(res.List[0].Price, "29.90")
		t.Assert(res.List[0].Sellable, true)
		t.Assert(res.List[0].Invalid, false)

		// 取消收藏 = 软删
		t.AssertNil(FavoriteRemove(ctx, uid, spuId))
		res, err = FavoriteList(ctx, uid, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.Total, 0)
		// 取消幂等
		t.AssertNil(FavoriteRemove(ctx, uid, spuId))

		// 再收藏 = 复活（复用原行, 行数仍为 1, 不违反唯一键）
		t.AssertNil(FavoriteAdd(ctx, uid, spuId))
		res, err = FavoriteList(ctx, uid, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.Total, 1)
		n, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM user_favorite WHERE user_id=? AND spu_id=?", uid, spuId)
		t.AssertNil(err)
		t.Assert(n.Int(), 1) // 复活而非新插行

		// 商品下架 → 标 Invalid（仍可见, 不静默消失——spec US3 验收 3）
		_, err = g.DB().Exec(ctx, "UPDATE product_spu SET status=0 WHERE id=?", spuId)
		t.AssertNil(err)
		res, err = FavoriteList(ctx, uid, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.Total, 1)
		t.Assert(res.List[0].Invalid, true)
		t.Assert(res.List[0].Sellable, false)
	})
}

// TestFootprintLifecycle 足迹（FR-008）: Record UPSERT / 列表倒序含次数 / 清空。
func TestFootprintLifecycle(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			phone = "13900004002"
			s1    = "fp1"
			s2    = "fp2"
		)
		defer cleanupMember(ctx, t, phone)
		defer cleanupSpuForMember(ctx, t, s1)
		defer cleanupSpuForMember(ctx, t, s2)
		uid := seedMember(ctx, t, phone, "足迹测试", 0)
		defer cleanupMemberCollections(ctx, t, uid)
		spu1 := seedSpuForMember(ctx, t, s1, "10.00", `["http://img/1.png"]`, 1)
		spu2 := seedSpuForMember(ctx, t, s2, "20.00", `["http://img/2.png"]`, 1)

		// 浏览上报（UPSERT: 重复浏览不新增行, 计数递增）
		t.AssertNil(FootprintRecord(ctx, uid, spu1))
		t.AssertNil(FootprintRecord(ctx, uid, spu2))
		t.AssertNil(FootprintRecord(ctx, uid, spu1)) // 再次浏览 spu1

		res, err := FootprintList(ctx, uid, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.Total, 2)
		// 最近浏览在前（spu1 刚被再看）
		t.Assert(res.List[0].SpuId, spu1)
		// 次数递增（表内 view_count=2）
		n, err := g.DB().GetValue(ctx,
			"SELECT view_count FROM user_footprint WHERE user_id=? AND spu_id=?", uid, spu1)
		t.AssertNil(err)
		t.Assert(n.Int(), 2)

		// 清空
		t.AssertNil(FootprintClear(ctx, uid))
		res, err = FootprintList(ctx, uid, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.Total, 0)
	})
}
