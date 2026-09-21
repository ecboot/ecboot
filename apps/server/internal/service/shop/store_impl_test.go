package shop

import (
	"context"
	"strings"
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
	_, _ = g.DB().Exec(ctx, "DELETE FROM `store` WHERE name=? OR store_no=?", name, "ST-T-"+name)
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

// cleanupStore 清理测试门店（按 name 与测试编码双路——测试可能改名, name 不可靠）。
func cleanupStore(ctx context.Context, t *gtest.T, names ...string) {
	for _, n := range names {
		_, _ = g.DB().Exec(ctx, "DELETE FROM `store` WHERE name=? OR store_no=?", n, "ST-T-"+n)
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
		seedStore(ctx, t, nA, "330108", tSearchLng, tSearchLat, 1) // A 区营业
		seedStore(ctx, t, nB, "330110", tSearchLng, tSearchLat, 1) // B 区营业
		seedStore(ctx, t, nC, "330108", tSearchLng, tSearchLat, 2) // A 区歇业

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

// ---- US2 管理面 ----

// TestAdminCreate 创建门店: 系统生成唯一编码 + 区划码校验（FR-007/008, 验收 1/2）。
func TestAdminCreate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			n1 = "t_adm_cr1"
			n2 = "t_adm_cr2"
		)
		defer cleanupStore(ctx, t, n1, n2)

		id, err := AdminCreate(ctx, model.StoreInput{
			Name: n1, ProvinceCode: "330000", CityCode: "330100", DistrictCode: "330108",
			DetailAddress: "测试路1号", Longitude: 120.123456, Latitude: 30.123456, PickupEnabled: true,
		})
		t.AssertNil(err)
		t.AssertGT(id, 0)

		d, err := AdminDetail(ctx, id)
		t.AssertNil(err)
		t.Assert(strings.HasPrefix(d.StoreNo, "ST"), true)
		t.Assert(d.Name, n1)
		t.Assert(d.PickupEnabled, true)
		t.Assert(d.Longitude, 120.123456)
		t.Assert(d.DistrictCode, "330108")

		// 编码全局唯一（两次创建不同）
		id2, err := AdminCreate(ctx, model.StoreInput{
			Name: n2, ProvinceCode: "330000", CityCode: "330100", DistrictCode: "330108",
			DetailAddress: "测试路2号",
		})
		t.AssertNil(err)
		d2, err := AdminDetail(ctx, id2)
		t.AssertNil(err)
		t.AssertNE(d.StoreNo, d2.StoreNo)

		// 区划码非 6 位数字 → 10001
		_, err = AdminCreate(ctx, model.StoreInput{
			Name: "t_adm_bad", ProvinceCode: "33", CityCode: "330100", DistrictCode: "330108",
			DetailAddress: "x",
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)
		_, err = AdminCreate(ctx, model.StoreInput{
			Name: "t_adm_bad", ProvinceCode: "330000", CityCode: "33010X", DistrictCode: "330108",
			DetailAddress: "x",
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)
	})
}

// TestAdminUpdateDelete 修改（含歇业切换）/软删/不存在语义（FR-009/010, 验收 3/5）。
func TestAdminUpdateDelete(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const n = "t_adm_ud"
		defer cleanupStore(ctx, t, n)
		id := seedStore(ctx, t, n, "330108", tSearchLng, tSearchLat, 1)

		// 改名 + 歇业（全量覆盖语义: 提交完整档案）
		t.AssertNil(AdminUpdate(ctx, id, model.StoreInput{
			Name: n + "_x", ProvinceCode: "330000", CityCode: "330100", DistrictCode: "330108",
			DetailAddress: "测试地址", Longitude: tSearchLng, Latitude: tSearchLat,
			PickupEnabled: true, Status: 2,
		}))
		d, err := AdminDetail(ctx, id)
		t.AssertNil(err)
		t.Assert(d.Name, n+"_x")
		t.Assert(d.Status, 2)
		_ = d

		// 关闭自提: 零值写入场景（do 的 omitempty 需显式 Fields 白名单才不吞零值）
		t.AssertNil(AdminUpdate(ctx, id, model.StoreInput{
			Name: n + "_x", ProvinceCode: "330000", CityCode: "330100", DistrictCode: "330108",
			DetailAddress: "改后地址", Longitude: tSearchLng, Latitude: tSearchLat,
			PickupEnabled: false, Status: 2,
		}))
		dc, err := AdminDetail(ctx, id)
		t.AssertNil(err)
		t.Assert(dc.PickupEnabled, false)
		t.Assert(dc.DetailAddress, "改后地址")

		// 歇业后游客列表不含
		res, err := PublicList(ctx, model.StoreQuery{DistrictCode: "330108"})
		t.AssertNil(err)
		for _, it := range res.List {
			t.AssertNE(it.Id, id)
		}

		// 不存在 → 10006（传完整档案以越过必填护栏, 直达存在性判定）
		err = AdminUpdate(ctx, 999999999, model.StoreInput{
			Name: "x", ProvinceCode: "330000", CityCode: "330100", DistrictCode: "330108",
			DetailAddress: "x", Status: 1,
		})
		t.Assert(errCode(err), errcode.CodeNotFound)

		// 软删 → 后台与游客均不可见
		t.AssertNil(AdminDelete(ctx, id))
		_, err = AdminDetail(ctx, id)
		t.Assert(errCode(err), errcode.CodeNotFound)
		_, err = PublicDetail(ctx, id)
		t.Assert(errCode(err), errcode.CodeNotFound)

		// 重复删除 → 10006
		err = AdminDelete(ctx, id)
		t.Assert(errCode(err), errcode.CodeNotFound)
	})
}

// TestAdminListFilters 后台列表: 状态 + 关键词（名称/编码）+ 分页（FR-011, 验收 4）。
func TestAdminListFilters(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			nA = "t_adm_lstA"
			nB = "t_adm_lstB"
		)
		defer cleanupStore(ctx, t, nA, nB)
		seedStore(ctx, t, nA, "330108", tSearchLng, tSearchLat, 1)
		idB := seedStore(ctx, t, nB, "330108", tSearchLng, tSearchLat, 2)

		// status=2 → 含 nB 不含 nA
		res, err := AdminList(ctx, 2, "", model.PageReq{Page: 1, PageSize: 20})
		t.AssertNil(err)
		foundB, foundA := false, false
		for _, it := range res.List {
			if it.Name == nB {
				foundB = true
				t.Assert(it.Status, 2)
			}
			if it.Name == nA {
				foundA = true
			}
		}
		t.Assert(foundB, true)
		t.Assert(foundA, false)

		// keyword 命中名称
		res, err = AdminList(ctx, 0, "t_adm_lst", model.PageReq{Page: 1, PageSize: 20})
		t.AssertNil(err)
		t.AssertGE(res.Total, 2)

		// keyword 命中编码（取 nB 的真实编码）
		db, err := AdminDetail(ctx, idB)
		t.AssertNil(err)
		res, err = AdminList(ctx, 0, db.StoreNo, model.PageReq{Page: 1, PageSize: 20})
		t.AssertNil(err)
		t.Assert(res.Total, 1)
		t.Assert(res.List[0].Id, idB)

		// 分页
		res, err = AdminList(ctx, 0, "t_adm_lst", model.PageReq{Page: 1, PageSize: 1})
		t.AssertNil(err)
		t.Assert(len(res.List), 1)
		t.Assert(res.Total, 2)
	})
}

// ---- 评审修复轮断言 ----

// TestAdminUpdateGuards 修改护栏（评审 C1/I1）: status 域外与全量覆盖必填缺失均拒绝。
func TestAdminUpdateGuards(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const n = "t_upd_guard"
		defer cleanupStore(ctx, t, n)
		id := seedStore(ctx, t, n, "330108", tSearchLng, tSearchLat, 1)

		full := model.StoreInput{
			Name: n, ProvinceCode: "330000", CityCode: "330100", DistrictCode: "330108",
			DetailAddress: "测试地址", Longitude: tSearchLng, Latitude: tSearchLat,
			PickupEnabled: true, Status: 1,
		}
		// status 域外（状态机仅 1/2）→ 10001
		bad := full
		bad.Status = 0
		t.Assert(errCode(AdminUpdate(ctx, id, bad)), errcode.CodeInvalidParam)
		// 全量覆盖必填缺失 → 10001（评审 I1 护栏）
		missing := full
		missing.Name = ""
		t.Assert(errCode(AdminUpdate(ctx, id, missing)), errcode.CodeInvalidParam)
		// 未落库: 状态与名称保持不变
		d, err := AdminDetail(ctx, id)
		t.AssertNil(err)
		t.Assert(d.Status, 1)
		t.Assert(d.Name, n)
	})
}

// TestAdminUpdateClearCoords 坐标清零（评审 I2: gdb 丢 nil, 需显式置空）。
func TestAdminUpdateClearCoords(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const n = "t_upd_coord"
		defer cleanupStore(ctx, t, n)
		id := seedStore(ctx, t, n, "330108", tSearchLng, tSearchLat+0.0045, 1)

		// 前置: 附近检索可命中
		res, err := PublicList(ctx, model.StoreQuery{Longitude: tSearchLng, Latitude: tSearchLat})
		t.AssertNil(err)
		hit := false
		for _, it := range res.List {
			if it.Id == id {
				hit = true
			}
		}
		t.Assert(hit, true)

		// 清空坐标（传 0）
		t.AssertNil(AdminUpdate(ctx, id, model.StoreInput{
			Name: n, ProvinceCode: "330000", CityCode: "330100", DistrictCode: "330108",
			DetailAddress: "测试地址", Longitude: 0, Latitude: 0, Status: 1,
		}))
		rec, err := g.DB().GetOne(ctx, "SELECT longitude, latitude FROM store WHERE id=?", id)
		t.AssertNil(err)
		t.Assert(rec["longitude"].IsNil(), true)
		t.Assert(rec["latitude"].IsNil(), true)

		// 清空后附近检索不再命中（无坐标排除）
		res, err = PublicList(ctx, model.StoreQuery{Longitude: tSearchLng, Latitude: tSearchLat})
		t.AssertNil(err)
		for _, it := range res.List {
			t.AssertNE(it.Id, id)
		}
	})
}

// TestPublicListPagingBounds 附近检索分页边界（评审 I3: 超大 page 不得 panic）。
func TestPublicListPagingBounds(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const n = "t_page_b"
		defer cleanupStore(ctx, t, n)
		seedStore(ctx, t, n, "330108", tSearchLng, tSearchLat, 1)

		// 溢出边界（(page-1)*pageSize 溢出为负）
		res, err := PublicList(ctx, model.StoreQuery{
			Longitude: tSearchLng, Latitude: tSearchLat,
			PageReq: model.PageReq{Page: 922337203685477582, PageSize: 10},
		})
		t.AssertNil(err)
		t.Assert(len(res.List), 0)

		// page 超范围
		res, err = PublicList(ctx, model.StoreQuery{
			Longitude: tSearchLng, Latitude: tSearchLat,
			PageReq: model.PageReq{Page: 9999, PageSize: 10},
		})
		t.AssertNil(err)
		t.Assert(len(res.List), 0)
	})
}

// TestPublicListCoordinatesGuard 经纬度域校验（评审 M4）: 非法值返 10001 而非 DB 错误。
func TestPublicListCoordinatesGuard(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		_, err := PublicList(ctx, model.StoreQuery{Longitude: 120, Latitude: 999})
		t.Assert(errCode(err), errcode.CodeInvalidParam)
		_, err = PublicList(ctx, model.StoreQuery{Longitude: 999, Latitude: 30})
		t.Assert(errCode(err), errcode.CodeInvalidParam)
	})
}

// TestPublicListOnlyOneCoord 经纬度只传其一视为未提供附近检索（spec 边界）。
func TestPublicListOnlyOneCoord(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const n = "t_one_coord"
		defer cleanupStore(ctx, t, n)
		seedStore(ctx, t, n, "330108", tSearchLng, tSearchLat, 1)

		// 只传经度 + 区县 → 回退区县筛选（非附近模式不带距离）
		res, err := PublicList(ctx, model.StoreQuery{DistrictCode: "330108", Longitude: tSearchLng})
		t.AssertNil(err)
		found := false
		for _, it := range res.List {
			if it.Name == n {
				found = true
				t.Assert(it.DistanceM, int64(0))
			}
		}
		t.Assert(found, true)
	})
}
