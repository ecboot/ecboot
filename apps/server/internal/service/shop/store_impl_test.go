package shop

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// ---- 门店测试数据（自建+清理; init 基座见 trade_test.go） ----

const (
	tSearchLng = 120.000000 // 检索点（杭州）
	tSearchLat = 30.000000
)

// seedStore 建测试门店, 返回 ID; lng/lat 同时为 0 表示未录入坐标（NULL, 不参与附近检索）。
func seedStore(ctx context.Context, t *gtest.T, name, districtCode string, lng, lat float64, status int) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM `store` WHERE name=?", name)
	var lngVal, latVal any
	if lng == 0 && lat == 0 {
		lngVal, latVal = nil, nil
	} else {
		lngVal, latVal = lng, lat
	}
	res, err := g.DB().Exec(ctx,
		"INSERT INTO `store`(store_no,name,province_code,city_code,district_code,detail_address,"+
			"longitude,latitude,business_hours,contact_phone,pickup_enabled,sort,status) "+
			"VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)",
		"ST-T-"+name, name, "330000", "330100", districtCode, "测试地址",
		lngVal, latVal, "09:00-21:00", "0571-00000000", 1, 0, status)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

// cleanupStore 清理测试门店。
func cleanupStore(ctx context.Context, t *gtest.T, names ...string) {
	for _, n := range names {
		_, _ = g.DB().Exec(ctx, "DELETE FROM `store` WHERE name=?", n)
	}
}

// TestPublicListDistrict 区县筛选仅营业（FR-001/002, US1 验收 1）。
func TestPublicListDistrict(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			nA = "t_st_dA"
			nB = "t_st_dB"
			nC = "t_st_dC"
		)
		defer cleanupStore(ctx, t, nA, nB, nC)
		seedStore(ctx, t, nA, "330108", tSearchLng, tSearchLat, 1)          // A 区营业
		seedStore(ctx, t, nB, "330110", tSearchLng, tSearchLat, 1)          // B 区营业
		seedStore(ctx, t, nC, "330108", tSearchLng, tSearchLat, 2)          // A 区歇业

		res, err := PublicList(ctx, model.StoreQuery{DistrictCode: "330108"})
		t.AssertNil(err)
		t.Assert(len(res.List), 1)
		t.Assert(res.List[0].Name, nA)
		t.Assert(res.List[0].Status, 1)
	})
}

// TestPublicListNearby 附近检索: 距离升序 + distanceM + 半径过滤 + 无坐标排除（FR-003/004, 验收 2/3/4）。
func TestPublicListNearby(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			nNear   = "t_st_near"   // ~500m
			nMid    = "t_st_mid"    // ~3km
			nFar    = "t_st_far"    // ~30km
			nNoGeo  = "t_st_nogeo"  // 无坐标
			nClosed = "t_st_closed" // 歇业（近）
		)
		defer cleanupStore(ctx, t, nNear, nMid, nFar, nNoGeo, nClosed)
		nearId := seedStore(ctx, t, nNear, "330108", tSearchLng, tSearchLat+0.0045, 1)
		seedStore(ctx, t, nMid, "330108", tSearchLng, tSearchLat+0.027, 1)
		seedStore(ctx, t, nFar, "330108", tSearchLng, tSearchLat+0.27, 1)
		seedStore(ctx, t, nNoGeo, "330108", 0, 0, 1)
		seedStore(ctx, t, nClosed, "330108", tSearchLng, tSearchLat+0.001, 2)

		// 默认半径 10km: 近 + 中, 按距离升序, 无坐标与歇业不出现
		res, err := PublicList(ctx, model.StoreQuery{Longitude: tSearchLng, Latitude: tSearchLat})
		t.AssertNil(err)
		t.Assert(len(res.List), 2)
		t.Assert(res.List[0].Name, nNear)
		t.Assert(res.List[1].Name, nMid)
		t.Assert(res.List[0].DistanceM >= 400 && res.List[0].DistanceM <= 600, true)
		t.Assert(res.List[1].DistanceM >= 2800 && res.List[1].DistanceM <= 3200, true)
		t.Assert(res.List[0].Id, nearId)

		// 半径扩到 50km: 远店出现, 距离仍升序
		res, err = PublicList(ctx, model.StoreQuery{Longitude: tSearchLng, Latitude: tSearchLat, RadiusKm: 50})
		t.AssertNil(err)
		t.Assert(len(res.List), 3)
		t.Assert(res.List[2].Name, nFar)
		t.Assert(res.List[2].DistanceM >= 29000 && res.List[2].DistanceM <= 31000, true)

		// radiusKm 超上限回退（100km 口径内仍含远店）
		res, err = PublicList(ctx, model.StoreQuery{Longitude: tSearchLng, Latitude: tSearchLat, RadiusKm: 9999})
		t.AssertNil(err)
		t.Assert(len(res.List) >= 3, true)

		// 命中 0 家: 空列表非错误
		res, err = PublicList(ctx, model.StoreQuery{Longitude: 100.0, Latitude: 20.0, RadiusKm: 1})
		t.AssertNil(err)
		t.Assert(len(res.List), 0)
	})
}

// TestPublicListPrecedence 区县与经纬度同传按附近检索（FR-005, 验收 5）。
func TestPublicListPrecedence(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			nIn  = "t_st_pIn"  // 附近命中但区县不同
			nOut = "t_st_pOut" // 区县命中但超出半径
		)
		defer cleanupStore(ctx, t, nIn, nOut)
		seedStore(ctx, t, nIn, "330110", tSearchLng, tSearchLat+0.0045, 1)
		seedStore(ctx, t, nOut, "330108", tSearchLng, tSearchLat+0.27, 1)

		// 同传: 按附近检索（半径内只有 nIn, 即使区县不匹配）
		res, err := PublicList(ctx, model.StoreQuery{
			DistrictCode: "330108", Longitude: tSearchLng, Latitude: tSearchLat,
		})
		t.AssertNil(err)
		t.Assert(len(res.List), 1)
		t.Assert(res.List[0].Name, nIn)
	})
}

// TestPublicDetail 游客详情: 完整字段 + 歇业可见 + 不存在 10006（FR-006, 验收 6）。
func TestPublicDetail(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const n = "t_st_detail"
		defer cleanupStore(ctx, t, n)
		id := seedStore(ctx, t, n, "330108", tSearchLng, tSearchLat, 2) // 歇业

		d, err := PublicDetail(ctx, id)
		t.AssertNil(err)
		t.Assert(d.Name, n)
		t.Assert(d.Status, 2)
		t.Assert(d.PickupEnabled, true)
		t.Assert(d.DistrictCode, "330108")
		t.Assert(d.StoreNo, "ST-T-"+n)

		_, err = PublicDetail(ctx, 999999999)
		t.Assert(errCode(err), errcode.CodeNotFound)
	})
}
