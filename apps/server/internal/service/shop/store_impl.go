// store_impl.go 门店域实现（接口契约见 misc.go IStoreLogic）。
// 语义: 游客列表仅营业、详情歇业可见（表注释"歇业=展示但不可自提/核销", 交易侧裁决归 012+）;
// 附近检索 = 包围盒预筛（idx_location）+ Haversine 计算列（research D1）; 编码 sonyflake（D2）。
package shop

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/library/idgen"
	"ecboot/internal/model"
	"ecboot/internal/model/do"
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
	page := q.Normalized()
	nearby := q.Longitude != 0 && q.Latitude != 0

	if nearby {
		// 经纬度域校验（评审 M4: 越界/NaN 会落 DB 错误而非参数错误）
		if math.IsNaN(q.Longitude) || math.IsInf(q.Longitude, 0) || math.Abs(q.Longitude) > 180 {
			return nil, errcode.New(errcode.CodeInvalidParam, "经度须在 -180~180")
		}
		if math.IsNaN(q.Latitude) || math.IsInf(q.Latitude, 0) || math.Abs(q.Latitude) > 90 {
			return nil, errcode.New(errcode.CodeInvalidParam, "纬度须在 -90~90")
		}
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
		Fields(storeBaseCols+", "+expr+" AS distance_m").
		Where(dao.Store.Columns().Status, 1).
		Where(dao.Store.Columns().Deleted, 0).
		Where("latitude BETWEEN ? AND ?", lat-deltaLat, lat+deltaLat).
		Where("longitude BETWEEN ? AND ?", lng-deltaLng, lng+deltaLng).
		Having("distance_m <= ?", float64(radiusKm)*1000).
		Order("distance_m ASC, id ASC").
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
	// 评审 I3: 超大 page 使 (page-1)*pageSize 溢出为负, 必须双向钳制（公开端点抗畸形入参）
	if start < 0 || start > len(all) {
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

// ---- 管理面（US2） ----

// isValidAreaCode 区划码 6 位数字（GB/T 2260）。
func isValidAreaCode(s string) bool {
	if len(s) != 6 {
		return false
	}
	for i := 0; i < 6; i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// coordValue 坐标值: 0 视为未录入/清空。
// gdb 对 do 结构体自动启用 OmitNilData（nil 在 SET 中被丢弃, Fields 白名单也救不回,
// 评审 I2 实证）——故以 gdb.Raw("NULL") 显式置空: 非 nil 值不被 omit, 且为原生 SQL 片段。
func coordValue(v float64) any {
	if v == 0 {
		return gdb.Raw("NULL")
	}
	return v
}

// AdminCreate 创建门店（FR-007/008）: 服务端生成全局唯一编码（sonyflake, research D2）;
// 唯一键冲突重试 ≤3（概率极低, 不暴露给调用方）。
func AdminCreate(ctx context.Context, in model.StoreInput) (int64, error) {
	if in.Name == "" {
		return 0, errcode.New(errcode.CodeInvalidParam, "门店名称必填")
	}
	if in.DetailAddress == "" {
		return 0, errcode.New(errcode.CodeInvalidParam, "详细地址必填")
	}
	if !isValidAreaCode(in.ProvinceCode) || !isValidAreaCode(in.CityCode) || !isValidAreaCode(in.DistrictCode) {
		return 0, errcode.New(errcode.CodeInvalidParam, "省/市/区县区划码须为6位数字")
	}
	status := in.Status
	if status <= 0 {
		status = 1
	}
	pickup := 0
	if in.PickupEnabled {
		pickup = 1
	}
	for attempt := 0; attempt < 3; attempt++ {
		next, err := idgen.NextID()
		if err != nil {
			return 0, gerror.Wrap(err, "生成门店编码失败")
		}
		res, err := dao.Store.Ctx(ctx).Data(do.Store{
			StoreNo:       "ST" + strconv.FormatInt(next, 10),
			Name:          in.Name,
			ProvinceCode:  in.ProvinceCode,
			CityCode:      in.CityCode,
			DistrictCode:  in.DistrictCode,
			DetailAddress: in.DetailAddress,
			Longitude:     coordValue(in.Longitude),
			Latitude:      coordValue(in.Latitude),
			BusinessHours: in.BusinessHours,
			ContactPhone:  in.ContactPhone,
			PickupEnabled: pickup,
			Status:        status,
		}).Insert()
		if err == nil {
			id, ierr := res.LastInsertId()
			if ierr != nil {
				return 0, gerror.Wrap(ierr, "读取门店ID失败")
			}
			return id, nil
		}
		if !strings.Contains(err.Error(), "Duplicate entry") {
			return 0, gerror.Wrap(err, "创建门店失败")
		}
	}
	return 0, errcode.New(errcode.CodeSystemError, "门店编码生成冲突, 请重试")
}

// AdminUpdate 修改门店（FR-009）: 全量覆盖语义（后台表单提交完整档案——显式 Fields 白名单
// 保证零值可写: 自提关闭/歇业态不被 omitempty 吞掉）; 目标不存在返 10006。
func AdminUpdate(ctx context.Context, id int64, in model.StoreInput) error {
	// 全量覆盖必填护栏（评审 I1: api 契约无法保证调用方全量提交, service 兜底防静默清空档案）
	if in.Name == "" {
		return errcode.New(errcode.CodeInvalidParam, "门店名称必填")
	}
	if in.DetailAddress == "" {
		return errcode.New(errcode.CodeInvalidParam, "详细地址必填")
	}
	if !isValidAreaCode(in.ProvinceCode) || !isValidAreaCode(in.CityCode) || !isValidAreaCode(in.DistrictCode) {
		return errcode.New(errcode.CodeInvalidParam, "省/市/区县区划码须为6位数字")
	}
	// 状态白名单（评审 C1: 状态机仅 1营业/2歇业, 防域外值落库致门店静默消失）
	if in.Status != 1 && in.Status != 2 {
		return errcode.New(errcode.CodeInvalidParam, "门店状态须为1营业或2歇业")
	}
	pickup := 0
	if in.PickupEnabled {
		pickup = 1
	}
	cols := dao.Store.Columns()
	res, err := dao.Store.Ctx(ctx).
		Where(cols.Id, id).
		Where(cols.Deleted, 0).
		Data(do.Store{
			Name:          in.Name,
			ProvinceCode:  in.ProvinceCode,
			CityCode:      in.CityCode,
			DistrictCode:  in.DistrictCode,
			DetailAddress: in.DetailAddress,
			Longitude:     coordValue(in.Longitude),
			Latitude:      coordValue(in.Latitude),
			BusinessHours: in.BusinessHours,
			ContactPhone:  in.ContactPhone,
			PickupEnabled: pickup,
			Status:        in.Status,
		}).
		Fields(cols.Name, cols.ProvinceCode, cols.CityCode, cols.DistrictCode, cols.DetailAddress,
			cols.Longitude, cols.Latitude, cols.BusinessHours, cols.ContactPhone,
			cols.PickupEnabled, cols.Status).
		Update()
	if err != nil {
		return gerror.Wrap(err, "修改门店失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		// 目标不存在或档案无变化（两者均不构成成功语义, 以存在性优先判定）
		if cnt, cerr := dao.Store.Ctx(ctx).
			Where(cols.Id, id).Where(cols.Deleted, 0).Count(); cerr == nil && cnt == 0 {
			return errcode.New(errcode.CodeNotFound, "门店不存在")
		}
	}
	return nil
}

// AdminDelete 软删门店（FR-010）: 后台与游客均不可见; 目标不存在返 10006。
func AdminDelete(ctx context.Context, id int64) error {
	cols := dao.Store.Columns()
	res, err := dao.Store.Ctx(ctx).
		Where(cols.Id, id).
		Where(cols.Deleted, 0).
		Data(do.Store{Deleted: 1, Status: 2}).
		Fields(cols.Deleted, cols.Status).
		Update()
	if err != nil {
		return gerror.Wrap(err, "删除门店失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeNotFound, "门店不存在")
	}
	return nil
}

// AdminList 后台门店列表（FR-011）: 状态筛选 + 名称/编码关键词模糊 + 分页。
func AdminList(ctx context.Context, status int, keyword string, page model.PageReq) (*model.PageResult[model.StoreItem], error) {
	page = page.Normalized()
	m := dao.Store.Ctx(ctx).Where(dao.Store.Columns().Deleted, 0)
	if status > 0 {
		m = m.Where(dao.Store.Columns().Status, status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		n, no := dao.Store.Columns().Name, dao.Store.Columns().StoreNo
		m = m.Where(n+" LIKE ? OR "+no+" LIKE ?", like, like)
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

// AdminDetail 后台门店详情（FR-012）: 字段口径与游客详情一致（后台含歇业店）。
func AdminDetail(ctx context.Context, storeId int64) (*model.StoreItem, error) {
	return PublicDetail(ctx, storeId)
}
