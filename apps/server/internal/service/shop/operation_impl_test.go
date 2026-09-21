package shop

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// ---- 测试数据（自建+清理; init 基座见 trade_test.go） ----

// seedLogistics 建测试物流公司（用后 cleanupLogistics; 按 code 双路清理）。
func seedLogistics(ctx context.Context, t *gtest.T, code string, status int) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM `logistics_company` WHERE code=?", code)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO `logistics_company`(code,name,tracking_rule,sort,status) VALUES(?,?,?,?,?)",
		code, code+"-name", "", 0, status)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

func cleanupLogistics(ctx context.Context, t *gtest.T, codes ...string) {
	for _, c := range codes {
		_, _ = g.DB().Exec(ctx, "DELETE FROM `logistics_company` WHERE code=?", c)
	}
}

// ---- US1 物流公司字典 ----

// TestLogisticsCrud 物流公司 CRUD: 编码唯一 + 编码不可改 + 停用保留（FR-001~006）。
func TestLogisticsCrud(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const code = "T_SF"
		defer cleanupLogistics(ctx, t, code, "T_NEWCODE")

		// 创建成功
		id, err := LogisticsCreate(ctx, model.LogisticsCompanyInput{
			Code: code, Name: "顺丰速运", TrackingRule: "SF+12位数字", Status: 1,
		})
		t.AssertNil(err)
		t.AssertGT(id, 0)

		// 编码重复 → 40012
		_, err = LogisticsCreate(ctx, model.LogisticsCompanyInput{Code: code, Name: "另一家", Status: 1})
		t.Assert(errCode(err), errcode.CodeLogisticsCodeTaken)

		// 详情
		d, err := LogisticsDetail(ctx, id)
		t.AssertNil(err)
		t.Assert(d.Code, code)
		t.Assert(d.Name, "顺丰速运")
		t.Assert(d.TrackingRule, "SF+12位数字")
		t.Assert(d.Status, 1)

		// 修改: 名称/规则/状态生效; 编码不可改（入参携带亦被忽略, FR-003）
		t.AssertNil(LogisticsUpdate(ctx, id, model.LogisticsCompanyInput{
			Code: "T_NEWCODE", Name: "顺丰", TrackingRule: "SF+13位", Status: 0,
		}))
		d, err = LogisticsDetail(ctx, id)
		t.AssertNil(err)
		t.Assert(d.Name, "顺丰")
		t.Assert(d.TrackingRule, "SF+13位")
		t.Assert(d.Status, 0)
		t.Assert(d.Code, code) // 编码未变

		// 停用保留: 不筛选列表仍含（FR-006）
		res, err := LogisticsList(ctx, 0, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		found := false
		for _, it := range res.List {
			if it.Id == id {
				found = true
				t.Assert(it.Status, 0)
			}
		}
		t.Assert(found, true)

		// status=1 筛选不含停用项
		res, err = LogisticsList(ctx, 1, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		for _, it := range res.List {
			t.AssertNE(it.Id, id)
		}

		// 详情不存在 → 10006
		_, err = LogisticsDetail(ctx, 999999999)
		t.Assert(errCode(err), errcode.CodeNotFound)

		// 修改不存在 → 10006
		err = LogisticsUpdate(ctx, 999999999, model.LogisticsCompanyInput{Name: "x", Status: 1})
		t.Assert(errCode(err), errcode.CodeNotFound)

		// 软删 → 列表与详情均不可见
		t.AssertNil(LogisticsDelete(ctx, id))
		_, err = LogisticsDetail(ctx, id)
		t.Assert(errCode(err), errcode.CodeNotFound)
		res, err = LogisticsList(ctx, 0, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		for _, it := range res.List {
			t.AssertNE(it.Id, id)
		}

		// 重复删除 → 10006
		err = LogisticsDelete(ctx, id)
		t.Assert(errCode(err), errcode.CodeNotFound)
	})
}

// TestLogisticsListPaging 列表筛选与分页（FR-001）。
func TestLogisticsListPaging(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			c1 = "T_lst1"
			c2 = "T_lst2"
		)
		defer cleanupLogistics(ctx, t, c1, c2)
		seedLogistics(ctx, t, c1, 1)
		seedLogistics(ctx, t, c2, 1)

		// 分页: pageSize=1 → 单行但总数>=2
		res, err := LogisticsList(ctx, 1, model.PageReq{Page: 1, PageSize: 1})
		t.AssertNil(err)
		t.Assert(len(res.List), 1)
		t.AssertGE(res.Total, 2)
	})
}

// ---- US2 装修管理 ----

// seedBanner 建测试轮播; startMin/endMin 相对当前分钟偏移（nil=该端为 NULL 长期）。
func seedBanner(ctx context.Context, t *gtest.T, image string, position int, startMin, endMin *int, status int) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM `operation_banner` WHERE image_url=?", image)
	var sVal, eVal any
	if startMin != nil {
		sVal = gtime.Now().Add(time.Duration(*startMin) * time.Minute)
	}
	if endMin != nil {
		eVal = gtime.Now().Add(time.Duration(*endMin) * time.Minute)
	}
	res, err := g.DB().Exec(ctx,
		"INSERT INTO `operation_banner`(position,image_url,link_url,sort,start_time,end_time,status) VALUES(?,?,'',0,?,?,?)",
		position, image, sVal, eVal, status)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

func cleanupBanner(ctx context.Context, t *gtest.T, images ...string) {
	for _, img := range images {
		_, _ = g.DB().Exec(ctx, "DELETE FROM `operation_banner` WHERE image_url=?", img)
	}
}

func seedFloor(ctx context.Context, t *gtest.T, title string, floorType, status int) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM `operation_floor` WHERE title=?", title)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO `operation_floor`(floor_type,title,config,sort,status) VALUES(?,?,NULL,0,?)",
		floorType, title, status)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

func cleanupFloor(ctx context.Context, t *gtest.T, titles ...string) {
	for _, ti := range titles {
		_, _ = g.DB().Exec(ctx, "DELETE FROM `operation_floor` WHERE title=?", ti)
	}
}

// TestBannerCrud 轮播管理: 创建/位置校验/列表筛选/修改含时段/软删（FR-007/009）。
func TestBannerCrud(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			imgA = "t_ban_a"
			imgB = "t_ban_b"
			imgX = "t_ban_x"
		)
		defer cleanupBanner(ctx, t, imgA, imgB, imgX)

		id, err := BannerCreate(ctx, model.OperBannerInput{Position: 1, ImageUrl: imgA, Status: 1})
		t.AssertNil(err)
		t.AssertGT(id, 0)

		// 位置域外 → 10001（service 兜底, api 层另有 in:1,2）
		_, err = BannerCreate(ctx, model.OperBannerInput{Position: 9, ImageUrl: imgX, Status: 1})
		t.Assert(errCode(err), errcode.CodeInvalidParam)
		// 图片必填
		_, err = BannerCreate(ctx, model.OperBannerInput{Position: 1, ImageUrl: "", Status: 1})
		t.Assert(errCode(err), errcode.CodeInvalidParam)

		// 位置 2 另建一条
		_, err = BannerCreate(ctx, model.OperBannerInput{Position: 2, ImageUrl: imgB, Status: 1})
		t.AssertNil(err)

		// 位置筛选: position=1 仅含 A
		res, err := BannerList(ctx, 1, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		foundA := false
		for _, it := range res.List {
			t.Assert(it.Position, 1)
			if it.Id == id {
				foundA = true
			}
		}
		t.Assert(foundA, true)

		// 修改: 时段 + 启停 + 链接 + 排序（传 position=2 亦不改位——api 无该字段）
		t.AssertNil(BannerUpdate(ctx, id, model.OperBannerInput{
			Position: 2, ImageUrl: imgA, LinkUrl: "/pages/x", Sort: 5,
			StartTime: "2026-01-01T00:00:00+08:00", EndTime: "2027-01-01T00:00:00+08:00", Status: 0,
		}))
		res, err = BannerList(ctx, 1, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		for _, it := range res.List {
			if it.Id == id {
				t.Assert(it.LinkUrl, "/pages/x")
				t.Assert(it.Sort, 5)
				t.Assert(it.Status, 0)
				t.Assert(it.StartTime != "", true)
				t.Assert(it.EndTime != "", true)
				t.Assert(it.Position, 1) // 位置不可改（传入 2 亦保持原值）
			}
		}

		// 不存在 → 10006
		err = BannerUpdate(ctx, 999999999, model.OperBannerInput{Position: 1, ImageUrl: "x", Status: 1})
		t.Assert(errCode(err), errcode.CodeNotFound)

		// 软删 → 管理端不可见; 重复删除 10006
		t.AssertNil(BannerDelete(ctx, id))
		res, err = BannerList(ctx, 1, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		for _, it := range res.List {
			t.AssertNE(it.Id, id)
		}
		err = BannerDelete(ctx, id)
		t.Assert(errCode(err), errcode.CodeNotFound)
	})
}

// TestFloorCrud 楼层管理: 三型/config JSON/修改/软删（FR-008/009）。
func TestFloorCrud(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const (
			fA = "t_flr_a"
			fB = "t_flr_b"
			fX = "t_flr_x"
		)
		defer cleanupFloor(ctx, t, fA, fB, fX)

		// 金刚区（config 含入口数组）
		id, err := FloorCreate(ctx, model.OperFloorInput{
			FloorType: 1, Title: fA, Status: 1,
			Config: map[string]any{"entries": []any{map[string]any{"text": "入口1", "link": "/p/1"}}},
		})
		t.AssertNil(err)
		t.AssertGT(id, 0)

		// 类型域外 → 10001
		_, err = FloorCreate(ctx, model.OperFloorInput{FloorType: 9, Title: fX, Status: 1})
		t.Assert(errCode(err), errcode.CodeInvalidParam)

		// 商品楼层
		idB, err := FloorCreate(ctx, model.OperFloorInput{
			FloorType: 2, Title: fB, Status: 1,
			Config: map[string]any{"spuIds": []any{"1"}},
		})
		t.AssertNil(err)
		t.AssertGT(idB, 0)

		// 列表含两者（分页返回 config 反序列化）
		res, err := FloorList(ctx, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		seen := map[int64]bool{}
		for _, it := range res.List {
			seen[it.Id] = true
			if it.Id == id {
				t.Assert(it.FloorType, 1)
				t.Assert(it.Title, fA)
				t.Assert(it.Config != nil, true)
			}
		}
		t.Assert(seen[id], true)
		t.Assert(seen[idB], true)

		// 修改 config 与启停（传不同 floorType 亦不改型——api 入参无该字段）
		t.AssertNil(FloorUpdate(ctx, id, model.OperFloorInput{
			FloorType: 3, Title: fA + "_2", Status: 0,
			Config: map[string]any{"entries": []any{}},
		}))
		res, err = FloorList(ctx, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		for _, it := range res.List {
			if it.Id == id {
				t.Assert(it.Title, fA+"_2")
				t.Assert(it.Status, 0)
				t.Assert(it.FloorType, 1) // 类型不可改（传入 3 亦保持原值）
			}
		}

		// 不存在 → 10006
		err = FloorUpdate(ctx, 999999999, model.OperFloorInput{FloorType: 1, Title: "x", Status: 1})
		t.Assert(errCode(err), errcode.CodeNotFound)

		// 软删 → 不可见; 重复删除 10006
		t.AssertNil(FloorDelete(ctx, idB))
		res, err = FloorList(ctx, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		for _, it := range res.List {
			t.AssertNE(it.Id, idB)
		}
		err = FloorDelete(ctx, idB)
		t.Assert(errCode(err), errcode.CodeNotFound)
	})
}
