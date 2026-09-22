// review_impl_test.go 评价域（014-review / 批次 08）——4 端点 + 并发/边界。
// fixture: 一笔"已完成(40)"订单 + 单个订单项（含 spu_name/sku_specs 快照）+ 会员。
// 全部自清；"一项一评"依赖 uk_order_item 唯一索引兜底，故并发用例是真防线验证。
package shop

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// reviewFixture 评价测试数据。
type reviewFixture struct {
	UserId  int64
	OrderId int64
	OrderNo string
	ItemId  int64
	SpuId   int64
	SkuId   int64
	phone   string
}

// seedReviewFixture 建 fixture（已完成订单 + 一个订单项，含商品名/规格快照）。
func seedReviewFixture(ctx context.Context, t *gtest.T, phone, orderNo string, spuId, skuId int64) *reviewFixture {
	cleanupReviewByOrder(ctx, orderNo)
	uid := seedPointUser(ctx, t, phone)
	f := &reviewFixture{UserId: uid, OrderNo: orderNo, SpuId: spuId, SkuId: skuId, phone: phone}

	res, err := g.DB().Exec(ctx,
		"INSERT INTO trade_order(order_no,user_id,status,total_amount,promotion_amount,pay_amount,currency,"+
			"receiver_name,receiver_phone,receiver_province,receiver_city,receiver_detail,refund_status) "+
			"VALUES(?,?,40,10.00,0,10.00,'CNY','评价测试','13800000000','浙江省','杭州市','T路评价',0)",
		orderNo, uid)
	t.AssertNil(err)
	oid, _ := res.LastInsertId()
	f.OrderId = oid

	res, err = g.DB().Exec(ctx,
		"INSERT INTO trade_order_item(order_no,order_id,spu_id,sku_id,sku_no,spu_name,sku_name,sku_specs,"+
			"quantity,original_price,price,coupon_amount,full_reduction_amount,point_amount,promotion_amount,pay_amount) "+
			"VALUES(?,?,?,?,'TF-RV-SKU','评价商品','评价商品 默认','{\"颜色\":\"黑\"}',1,10.00,10.00,0,0,0,0,10.00)",
		orderNo, oid, spuId, skuId)
	t.AssertNil(err)
	iid, _ := res.LastInsertId()
	f.ItemId = iid
	return f
}

// cleanupReviewByOrder 按订单号清场（评价 → 订单项 → 订单）。
func cleanupReviewByOrder(ctx context.Context, orderNo string) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_review WHERE order_no=?", orderNo)
	_, _ = g.DB().Exec(ctx, "DELETE i FROM trade_order_item i JOIN trade_order o ON i.order_id=o.id WHERE o.order_no=?", orderNo)
	_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order WHERE order_no=?", orderNo)
}

func cleanupReviewFixture(ctx context.Context, t *gtest.T, f *reviewFixture) {
	cleanupReviewByOrder(ctx, f.OrderNo)
	cleanupPointUser(ctx, t, f.phone)
}

// seedReviewRow 直接造评价（绕开 Create 构造多状态数据）。
func seedReviewRow(
	ctx context.Context, t *gtest.T, f *reviewFixture, orderNo string, spuId, score, auditStatus, anonymous int, content string,
) int64 {
	res, err := g.DB().Exec(ctx,
		"INSERT INTO product_review(order_item_id,order_no,user_id,spu_id,sku_id,spu_name,sku_specs,score,content,"+
			"images,is_anonymous,audit_status) VALUES(?,?,?,?,?,'评价商品','{\"颜色\":\"黑\"}',?,?,'[]',?,?)",
		f.ItemId+int64(score)*1000+int64(auditStatus), orderNo, f.UserId, spuId, f.SkuId, score, content, anonymous, auditStatus)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

// reviewRow 取某订单下的评价（单条）。
func reviewRow(ctx context.Context, orderNo string) map[string]any {
	rec, err := g.DB().GetOne(ctx,
		"SELECT id, audit_status, score, content, is_anonymous, extra_content, extra_time "+
			"FROM product_review WHERE order_no=?", orderNo)
	if err != nil || rec.IsEmpty() {
		return nil
	}
	return rec.Map()
}

// ---------- US1: 提交评价 ----------

// TestReviewCreate 提交评价（FR-001~004）: 成功落库（快照 + V1 自动通过）; 重复/他人/未完成/越界被拒。
func TestReviewCreate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedReviewFixture(ctx, t, "RV-CRT-1", "T-RV-CRT1", 970000001, 970000101)
		defer cleanupReviewFixture(ctx, t, f)

		id, err := NewReviewLogic().Create(ctx, f.UserId, model.ReviewCreateInput{
			OrderItemId: f.ItemId, Score: 5, Content: "很好用", Images: []string{"http://img/a.jpg"}, IsAnonymous: false,
		})
		t.AssertNil(err)
		t.Assert(id > 0, true)

		rec, err := g.DB().GetOne(ctx,
			"SELECT order_no, user_id, spu_id, sku_id, spu_name, sku_specs, score, content, audit_status "+
				"FROM product_review WHERE id=?", id)
		t.AssertNil(err)
		t.Assert(rec["order_no"].String(), f.OrderNo)
		t.Assert(rec["user_id"].Int64(), f.UserId)
		t.Assert(rec["spu_id"].Int64(), f.SpuId)
		t.Assert(rec["spu_name"].String(), "评价商品") // 下单时快照
		t.Assert(rec["audit_status"].Int(), 1)     // V1 自动通过（用户裁定 B）
		t.Assert(rec["score"].Int(), 5)
		t.Assert(rec["content"].String(), "很好用")

		// 一项一评: 重复提交 → 40010, 且活库仍只有 1 条
		_, err = NewReviewLogic().Create(ctx, f.UserId, model.ReviewCreateInput{
			OrderItemId: f.ItemId, Score: 4, Content: "再来一条",
		})
		t.Assert(errCode(err), errcode.CodeAlreadyReviewed)
		n, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM product_review WHERE order_item_id=?", f.ItemId)
		t.AssertNil(err)
		t.Assert(n.Int(), 1)

		// 他人订单项 → 按不存在（40009）
		_, err = NewReviewLogic().Create(ctx, f.UserId+999, model.ReviewCreateInput{
			OrderItemId: f.ItemId, Score: 5,
		})
		t.Assert(errCode(err), errcode.CodeNotFound)

		// 评分越界（服务侧兜底）
		_, err = NewReviewLogic().Create(ctx, f.UserId, model.ReviewCreateInput{
			OrderItemId: f.ItemId, Score: 6,
		})
		t.Assert(errCode(err), errcode.CodeInvalidParam)
	})
}

// TestReviewCreateOrderNotFinished 未完成订单不可评（FR-001）。
func TestReviewCreateOrderNotFinished(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedReviewFixture(ctx, t, "RV-CRT-2", "T-RV-CRT2", 970000002, 970000102)
		defer cleanupReviewFixture(ctx, t, f)

		_, _ = g.DB().Exec(ctx, "UPDATE trade_order SET status=20 WHERE id=?", f.OrderId)
		_, err := NewReviewLogic().Create(ctx, f.UserId, model.ReviewCreateInput{
			OrderItemId: f.ItemId, Score: 5,
		})
		t.Assert(errCode(err), errcode.CodeStatusNotAllowed) // 订单状态不允许（40006）
	})
}

// TestReviewCreateConcurrent 并发提交只成功一次（SC-003）: uk_order_item 唯一索引兜底 + 1062 转业务码。
func TestReviewCreateConcurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		f := seedReviewFixture(ctx, t, "RV-CRT-3", "T-RV-CRT3", 970000003, 970000103)
		defer cleanupReviewFixture(ctx, t, f)

		const n = 8
		var wg sync.WaitGroup
		errs := make([]error, n)
		start := make(chan struct{})
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				<-start
				_, errs[idx] = NewReviewLogic().Create(ctx, f.UserId, model.ReviewCreateInput{
					OrderItemId: f.ItemId, Score: 4, Content: "并发",
				})
			}(i)
		}
		close(start)
		wg.Wait()

		ok := 0
		for _, e := range errs {
			if e == nil {
				ok++
			} else {
				t.Assert(errCode(e), errcode.CodeAlreadyReviewed) // 冲突必须转业务码, 不得裸抛 1062
			}
		}
		t.Assert(ok, 1)
		cnt, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM product_review WHERE order_item_id=?", f.ItemId)
		t.AssertNil(err)
		t.Assert(cnt.Int(), 1)
	})
}

// ---------- US2: 商品评价列表与汇总 ----------

// TestReviewProductList 列表+汇总（FR-007/008/009）: 只出审核通过; 星级筛选; 汇总同源同筛; 匿名脱敏。
func TestReviewProductList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const spuId = 970000011
		f := seedReviewFixture(ctx, t, "RV-LST-1", "T-RV-LST1", spuId, 970000111)
		defer cleanupReviewFixture(ctx, t, f)

		// 同 SPU 四条: 2 通过（5 星 / 3 星）+ 1 待审（4 星）+ 1 驳回（1 星）
		seedReviewRow(ctx, t, f, f.OrderNo, spuId, 5, reviewAuditPassed, 0, "通过5星")
		seedReviewRow(ctx, t, f, f.OrderNo, spuId, 3, reviewAuditPassed, 0, "通过3星")
		seedReviewRow(ctx, t, f, f.OrderNo, spuId, 4, reviewAuditPending, 0, "待审4星")
		seedReviewRow(ctx, t, f, f.OrderNo, spuId, 1, reviewAuditReject, 0, "驳回1星")

		summary, page, err := NewReviewLogic().ProductList(ctx, spuId, 0, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(page.Total, 2)     // 只统计审核通过
		t.Assert(len(page.List), 2) // 待审与驳回不出现
		t.Assert(summary.Total, 2)
		t.Assert(summary.Avg, "4.0") // (5+3)/2
		t.Assert(summary.Distribution["5"], 1)
		t.Assert(summary.Distribution["3"], 1)
		t.Assert(summary.Distribution["4"], 0)

		// 空集合口径与商品详情页一致（占位 "—"）
		s0, p0, err := NewReviewLogic().ProductList(ctx, 979999999, 0, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(p0.Total, 0)
		t.Assert(s0.Avg, "—")
		t.Assert(s0.Total, 0)

		// 星级筛选
		s5, p5, err := NewReviewLogic().ProductList(ctx, spuId, 5, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(p5.Total, 1)
		t.Assert(p5.List[0].Score, 5)
		t.Assert(s5.Total, 1)
		t.Assert(s5.Avg, "5.0")

		// 条目字段: 规格快照 + 脱敏评价人（非匿名） + 图片数组
		var hit *model.ReviewCard
		for idx := range page.List {
			if page.List[idx].Score == 5 {
				hit = &page.List[idx]
			}
		}
		t.Assert(hit != nil, true)
		t.Assert(hit.Specs["颜色"], "黑")    // 下单时规格快照
		t.Assert(hit.User, "积***")        // "积分测试"(4字) → 首字符 + 3 掩码       // 昵称"积分测试" → 首字符+掩码
		t.Assert(hit.Images != nil, true) // 空图片也给空数组而非 nil
	})
}

// TestReviewProductListAnonymous 匿名评价展示为固定占位（FR-009）。
func TestReviewProductListAnonymous(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const spuId = 970000012
		f := seedReviewFixture(ctx, t, "RV-LST-2", "T-RV-LST2", spuId, 970000112)
		defer cleanupReviewFixture(ctx, t, f)
		seedReviewRow(ctx, t, f, f.OrderNo, spuId, 5, reviewAuditPassed, 1, "匿名评价")

		_, page, err := NewReviewLogic().ProductList(ctx, spuId, 0, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(len(page.List), 1)
		t.Assert(page.List[0].User, "匿名用户")
	})
}

// ---------- US3: 我的评价 ----------

// TestReviewMyList 我的评价（FR-005/010）: 只返回本人（含全部审核状态）; 追评/回复透出; 分页。
func TestReviewMyList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const spuId = 970000021
		f := seedReviewFixture(ctx, t, "RV-MY-1", "T-RV-MY1", spuId, 970000121)
		defer cleanupReviewFixture(ctx, t, f)

		id1 := seedReviewRow(ctx, t, f, f.OrderNo, spuId, 5, reviewAuditPassed, 0, "我的通过")
		seedReviewRow(ctx, t, f, f.OrderNo, spuId, 2, reviewAuditPending, 0, "我的待审")
		// 带追评与回复的一条
		_, err := g.DB().Exec(ctx,
			"UPDATE product_review SET extra_content='用了一个月还很好', reply_content='感谢支持' WHERE id=?", id1)
		t.AssertNil(err)
		// 他人评价（不同 user_id, 用不同的 order_item_id 避开唯一键）
		_, err = g.DB().Exec(ctx,
			"INSERT INTO product_review(order_item_id,order_no,user_id,spu_id,sku_id,spu_name,sku_specs,score,"+
				"content,images,is_anonymous,audit_status) VALUES(?,?,?,?,?,'评价商品','{}',4,'他人的评价','[]',0,1)",
			f.ItemId+7777, "T-RV-OTHER", f.UserId+999, spuId, f.SkuId)
		t.AssertNil(err)
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM product_review WHERE order_no='T-RV-OTHER'") }()

		res, err := NewReviewLogic().MyList(ctx, f.UserId, model.PageReq{Page: 1, PageSize: 10})
		t.AssertNil(err)
		t.Assert(res.Total, 2) // 不含他人
		hit := map[int64]model.MyReviewItem{}
		for _, it := range res.List {
			hit[it.ReviewId] = it
		}
		t.Assert(len(hit), 2)
		t.Assert(hit[id1].AuditStatus, reviewAuditPassed)
		t.Assert(hit[id1].SpuName, "评价商品")
		t.Assert(hit[id1].Extra, "用了一个月还很好")
		t.Assert(hit[id1].Reply, "感谢支持")
		// 待审的那条也在（我的评价不筛审核状态）
		found := false
		for _, it := range res.List {
			if it.AuditStatus == reviewAuditPending {
				found = true
			}
		}
		t.Assert(found, true)

		// 分页（页大小 1 → 只回 1 条但 total 仍 2）
		pg, err := NewReviewLogic().MyList(ctx, f.UserId, model.PageReq{Page: 1, PageSize: 1})
		t.AssertNil(err)
		t.Assert(len(pg.List), 1)
		t.Assert(pg.Total, 2)
	})
}

// ---------- US4: 追加评价 ----------

// TestReviewExtra 追评（FR-006）: 成功一次; 二次/超期/他人/空内容被拒; 并发只成功一次。
func TestReviewExtra(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const spuId = 970000031
		f := seedReviewFixture(ctx, t, "RV-EXT-1", "T-RV-EXT1", spuId, 970000131)
		defer cleanupReviewFixture(ctx, t, f)

		rid, err := NewReviewLogic().Create(ctx, f.UserId, model.ReviewCreateInput{
			OrderItemId: f.ItemId, Score: 5, Content: "首评",
		})
		t.AssertNil(err)

		// ① 成功: 追评内容与时间落库
		t.AssertNil(NewReviewLogic().Extra(ctx, f.UserId, rid, "用了一个月", []string{"http://img/e.jpg"}))
		rec := reviewRow(ctx, f.OrderNo)
		t.Assert(rec["extra_content"].(string), "用了一个月")
		t.Assert(rec["extra_time"] != nil, true)

		// ② 二次追评 → 40011
		t.Assert(errCode(NewReviewLogic().Extra(ctx, f.UserId, rid, "再来一次", nil)), errcode.CodeExtraReviewDenied)

		// ③ 他人追评 → 按不存在（40009）
		t.Assert(errCode(NewReviewLogic().Extra(ctx, f.UserId+999, rid, "他人追评", nil)), errcode.CodeNotFound)

		// ④ 空内容 → 10001
		t.Assert(errCode(NewReviewLogic().Extra(ctx, f.UserId, rid, "  ", nil)), errcode.CodeInvalidParam)
	})
}

// TestReviewExtraExpired 超 90 天不可追评（FR-006）。
func TestReviewExtraExpired(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const spuId = 970000032
		f := seedReviewFixture(ctx, t, "RV-EXT-2", "T-RV-EXT2", spuId, 970000132)
		defer cleanupReviewFixture(ctx, t, f)

		rid := seedReviewRow(ctx, t, f, f.OrderNo, spuId, 5, reviewAuditPassed, 0, "老评价")
		_, _ = g.DB().Exec(ctx, "UPDATE product_review SET created_at=DATE_SUB(NOW(), INTERVAL 100 DAY) WHERE id=?", rid)
		t.Assert(errCode(NewReviewLogic().Extra(ctx, f.UserId, rid, "超期追评", nil)), errcode.CodeExtraReviewDenied)

		// 89 天前仍可追评（边界）
		rid2 := seedReviewRow(ctx, t, f, f.OrderNo, spuId, 4, reviewAuditPassed, 0, "边界评价") // 评审 I1: 必须用 f.OrderNo, 否则 cleanup 按 order_no 删不到 → 每次跑测试留 1 行
		_, _ = g.DB().Exec(ctx, "UPDATE product_review SET created_at=DATE_SUB(NOW(), INTERVAL 89 DAY) WHERE id=?", rid2)
		t.AssertNil(NewReviewLogic().Extra(ctx, f.UserId, rid2, "89天追评", nil))
	})
}

// TestReviewExtraConcurrent 并发追评只成功一次（条件更新防线）。
func TestReviewExtraConcurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const spuId = 970000033
		f := seedReviewFixture(ctx, t, "RV-EXT-3", "T-RV-EXT3", spuId, 970000133)
		defer cleanupReviewFixture(ctx, t, f)
		rid := seedReviewRow(ctx, t, f, f.OrderNo, spuId, 5, reviewAuditPassed, 0, "首评")

		const n = 6
		var wg sync.WaitGroup
		oks := make([]bool, n)
		start := make(chan struct{})
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				<-start
				oks[idx] = NewReviewLogic().Extra(ctx, f.UserId, rid, "并发追评", nil) == nil
			}(i)
		}
		close(start)
		wg.Wait()

		success := 0
		for _, b := range oks {
			if b {
				success++
			}
		}
		t.Assert(success, 1) // 条件更新: 只有一次能把 '' 写成内容
	})
}

// TestReviewEdgeCoverage 评审 Minor 的覆盖补齐: 90 天边界/软删不可见/空昵称占位/图片往返/排序/追评文案。
func TestReviewEdgeCoverage(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const spuId = 970000041
		f := seedReviewFixture(ctx, t, "RV-EDGE-1", "T-RV-EDGE1", spuId, 970000141)
		defer cleanupReviewFixture(ctx, t, f)

		// ① 恰好 90 天: 允许追评（`>=` 边界; 改成 `>` 会被这条抓住）
		rid := seedReviewRow(ctx, t, f, f.OrderNo, spuId, 5, reviewAuditPassed, 0, "边界90天")
		_, _ = g.DB().Exec(ctx, "UPDATE product_review SET created_at=DATE_SUB(NOW(), INTERVAL 90 DAY) WHERE id=?", rid)
		t.AssertNil(NewReviewLogic().Extra(ctx, f.UserId, rid, "恰好90天可追评", nil))

		// ② 软删评价不可见（对外列表与我的评价都要过滤 deleted=0）
		rid2 := seedReviewRow(ctx, t, f, f.OrderNo, spuId, 1, reviewAuditPassed, 0, "待软删")
		_, _ = g.DB().Exec(ctx, "UPDATE product_review SET deleted=1 WHERE id=?", rid2)
		_, page, err := NewReviewLogic().ProductList(ctx, spuId, 0, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		for _, it := range page.List {
			t.Assert(it.ReviewId != rid2, true)
		}
		mine, err := NewReviewLogic().MyList(ctx, f.UserId, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		for _, it := range mine.List {
			t.Assert(it.ReviewId != rid2, true)
		}

		// ③ 空昵称占位（脱敏 helper 的空值分支）
		_, _ = g.DB().Exec(ctx, "UPDATE `user` SET nickname='' WHERE id=?", f.UserId)
		_, page2, err := NewReviewLogic().ProductList(ctx, spuId, 5, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		t.Assert(page2.List[0].User, "用户****")
		_, _ = g.DB().Exec(ctx, "UPDATE `user` SET nickname='积分测试' WHERE id=?", f.UserId)

		// ④ 图片往返: 提交时写非空数组 → 列表读回同内容（写入侧与解析侧各钉一条）
		id3, err := NewReviewLogic().Create(ctx, f.UserId, model.ReviewCreateInput{
			OrderItemId: f.ItemId, Score: 4, Content: "带图", Images: []string{"http://img/x.jpg", "http://img/y.jpg"},
		})
		t.AssertNil(err)
		raw, err := g.DB().GetValue(ctx, "SELECT images FROM product_review WHERE id=?", id3)
		t.AssertNil(err)
		t.Assert(strings.Contains(raw.String(), "x.jpg"), true) // 写入侧
		_, page3, err := NewReviewLogic().ProductList(ctx, spuId, 4, model.PageReq{Page: 1, PageSize: 50})
		t.AssertNil(err)
		t.Assert(len(page3.List), 1)
		t.Assert(len(page3.List[0].Images), 2) // 解析侧
		t.Assert(page3.List[0].Images[0], "http://img/x.jpg")

		// ⑤ 我的评价按 id 倒序（最新在前）: 断言首页第一条是刚创建的
		mine2, err := NewReviewLogic().MyList(ctx, f.UserId, model.PageReq{Page: 1, PageSize: 1})
		t.AssertNil(err)
		t.Assert(mine2.List[0].ReviewId, id3)
	})
}

// TestReviewExtraMessages 追评 affected=0 的两条文案可区分（评审 M8: 交换两条消息的变异须被抓住）。
func TestReviewExtraMessages(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const spuId = 970000042
		f := seedReviewFixture(ctx, t, "RV-MSG-1", "T-RV-MSG1", spuId, 970000142)
		defer cleanupReviewFixture(ctx, t, f)

		rid := seedReviewRow(ctx, t, f, f.OrderNo, spuId, 5, reviewAuditPassed, 0, "文案用")
		// 已追评 → 文案含"已追评"
		_, _ = g.DB().Exec(ctx, "UPDATE product_review SET extra_content='已追评过' WHERE id=?", rid)
		err := NewReviewLogic().Extra(ctx, f.UserId, rid, "再来", nil)
		t.Assert(errCode(err), errcode.CodeExtraReviewDenied)
		t.Assert(strings.Contains(err.Error(), "已追评"), true)

		// 超期 → 文案含"90 天"
		rid2 := seedReviewRow(ctx, t, f, f.OrderNo, spuId, 4, reviewAuditPassed, 0, "超期用")
		_, _ = g.DB().Exec(ctx, "UPDATE product_review SET created_at=DATE_SUB(NOW(), INTERVAL 200 DAY) WHERE id=?", rid2)
		err = NewReviewLogic().Extra(ctx, f.UserId, rid2, "再来", nil)
		t.Assert(errCode(err), errcode.CodeExtraReviewDenied)
		t.Assert(strings.Contains(err.Error(), "90 天"), true)
	})
}
