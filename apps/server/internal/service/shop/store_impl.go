// store_impl.go 门店域实现（接口契约见 misc.go IStoreLogic）。
// 语义: 游客列表仅营业、详情歇业可见（表注释"歇业=展示但不可自提/核销", 交易侧裁决归 012+）;
// 附近检索 = 包围盒预筛（idx_location）+ Haversine 计算列（research D1）; 编码 sonyflake（D2）。
package shop

import (
	"context"
	"fmt"
	"math"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// storeBaseCols 门店列表/详情基础列（与 storeItemFromRecord 映射一致）。
const storeBaseCols = "id,store_no,name,province_code,city_code,district_code,detail_address," +
	"longitude,latitude,business_hours,contact_phone,pickup_enabled,status"

// haversineDistanceSQL 球面距离（米）计算列; 参数顺序: 检索点纬度×2、检索点经度。
// 注: gdb 的 Fields 为原样 SQL 片段（不做参数绑定）, 故检索点以 %.6f 数值字面量内联——
// 值为 float64 格式化产物, 无注入面（经纬度参数本身是数字类型, 非用户字符串）。
const haversineDistanceSQL = `6371000 * 2 * ASIN(SQRT(` +
	`POW(SIN((RADIANS(latitude) - RADIANS(%.6f)) / 2), 2) + ` +
	`COS(RADIANS(%.6f)) * COS(RADIANS(latitude)) * ` +
	`POW(SIN((RADIANS(longitude) - RADIANS(%.6f)) / 2), 2)))`

// normalizeRadiusKm 半径归一（默认 10; 上限 100; 域外回退, spec 边界）。
func normalizeRadiusKm(km int) int {
	if km <= 0 {
		return 10
	}
	if km > 100 {
		return 100
	}
	return km
}

// storeItemFromRecord 行 → DTO（手写映射: 列名 snake_case 与 DTO json tag 不同名）。
func storeItemFromRecord(r gdb.Record) model.StoreItem {
	it := model.StoreItem{
		Id:            r["id"].Int64(),
		StoreNo:       r["store_no"].String(),
		Name:          r["name"].String(),
		ProvinceCode:  r["province_code"].String(),
		CityCode:      r["city_code"].String(),
		DistrictCode:  r["district_code"].String(),
		DetailAddress: r["detail_address"].String(),
		Longitude:     r["longitude"].Float64(),
		Latitude:      r["latitude"].Float64(),
		BusinessHours: r["business_hours"].String(),
		ContactPhone:  r["contact_phone"].String(),
		PickupEnabled: r["pickup_enabled"].Bool(),
		Status:        r["status"].Int(),
	}
	if v, ok := r["distance_m"]; ok {
		it.DistanceM = v.Int64()
	}
	return it
}

// PublicList 游客门店（FR-001~005）: 区县筛选或经纬度附近检索（经纬度优先）, 仅营业。
func PublicList(ctx context.Context, q model.StoreQuery) (*model.PageResult[model.StoreItem], error) {
	page := q.PageReq.Normalized()
	nearby := q.Longitude != 0 && q.Latitude != 0

	if nearby {
		return publicListNearby(ctx, q, page)
	}

	m := dao.Store.Ctx(ctx).
		Where(dao.Store.Columns().Status, 1).
		Where(dao.Store.Columns().Deleted, 0)
	if q.DistrictCode != "" {
		m = m.Where(dao.Store.Columns().DistrictCode, q.DistrictCode)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计门店失败")
	}
	recs, err := m.Fields(storeBaseCols).
		Order("sort ASC, id ASC").
		Page(page.Page, page.PageSize).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询门店失败")
	}
	list := make([]model.StoreItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, storeItemFromRecord(r))
	}
	return &model.PageResult[model.StoreItem]{List: list, Total: int64(total)}, nil
}

// publicListNearby 附近检索: 包围盒预筛（命中 idx_location）→ Haversine 计算列 →
// HAVING 半径过滤 → 距离升序。无坐标门店被包围盒条件天然排除（FR-004）。
// 分页在 Go 侧切页: 附近命中量为门店总量同阶（数百）, 免于 Count/Having 组合在 ORM 上的歧义。
func publicListNearby(ctx context.Context, q model.StoreQuery, page model.PageReq) (*model.PageResult[model.StoreItem], error) {
	radiusKm := normalizeRadiusKm(q.RadiusKm)
	lng, lat := q.Longitude, q.Latitude

	// 包围盒（1 纬度 ≈ 111km; 经度随纬度收缩, 极区收敛保护）
	cosLat := math.Cos(lat * math.Pi / 180)
	if cosLat < 0.01 {
		cosLat = 0.01
	}
	deltaLat := float64(radiusKm) / 111.0
	deltaLng := float64(radiusKm) / (111.0 * cosLat)

	expr := fmt.Sprintf(haversineDistanceSQL, lat, lat, lng)
	recs, err := dao.Store.Ctx(ctx).
		Fields(storeBaseCols + ", " + expr + " AS distance_m").
		Where(dao.Store.Columns().Status, 1).
		Where(dao.Store.Columns().Deleted, 0).
		Where("latitude BETWEEN ? AND ?", lat-deltaLat, lat+deltaLat).
		Where("longitude BETWEEN ? AND ?", lng-deltaLng, lng+deltaLng).
		Having("distance_m <= ?", float64(radiusKm)*1000).
		Order("distance_m ASC").
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "附近门店检索失败")
	}

	all := make([]model.StoreItem, 0, len(recs))
	for _, r := range recs {
		all = append(all, storeItemFromRecord(r))
	}
	total := int64(len(all))
	start := (page.Page - 1) * page.PageSize
	if start > len(all) {
		start = len(all)
	}
	end := start + page.PageSize
	if end > len(all) {
		end = len(all)
	}
	return &model.PageResult[model.StoreItem]{List: all[start:end], Total: total}, nil
}

// PublicDetail 游客门店详情（FR-006）: 歇业店可见; 不存在/已删返 10006。
func PublicDetail(ctx context.Context, storeId int64) (*model.StoreItem, error) {
	rec, err := dao.Store.Ctx(ctx).
		Fields(storeBaseCols).
		Where(dao.Store.Columns().Id, storeId).
		Where(dao.Store.Columns().Deleted, 0).
		One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询门店失败")
	}
	if rec.IsEmpty() {
		return nil, errcode.New(errcode.CodeNotFound, "门店不存在")
	}
	it := storeItemFromRecord(rec)
	return &it, nil
}
