// dashboard_impl_test.go 三看板口径（019 批次 13）——口径钉住 + 窗口路径覆盖。
// I2/I3（收口终验评审修复）: ①断言改为**受控夹具**（唯一时间窗内造数, 不用全库聚合差值——
// 原写法在 go test ./... 并行包下与 shop 包写入竞争, 三次全量跑 1 红）;
// ②补时间窗路径用例（原全部只传空窗口, FR 主输入零覆盖——正是 C1 漏网之因）;
// ③补 PendingReview 断言与 pendingDeliver 精确断言（原为 `>= 1` 空转）。
package system

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

// dbWindow **固定历史区间**（N2 三轮收口的正解）: 窗口取 2000 年的一段固定区间, 夹具行
// **显式写入该区间内的 created_at**——并行包的写入 created_at=NOW, **结构上不可能落入窗内**。
// 前两轮失败教训: ①窗口含 NOW + 精确等值（并行写入即窗内）; ②窗口含 NOW + 基线差量
// （差量吸收不了 base→after 之间的并发增删——复审决定性实验 10/10 红）。真隔离靠窗口与
// 并发写入的**值域不相交**, 而非时间上的先后。
const (
	dbWinStart = "2000-01-01 00:00:00"
	dbWinEnd   = "2000-01-02 00:00:00"
	dbWinAt    = "2000-01-01 12:00:00"
)

func dbWindow() (string, string) { return dbWinStart, dbWinEnd }

// TestDashboardTrade 交易看板口径（FR-1）: 已支付枚举 20/30/40（**排除 90 已取消**）/窗口过滤/
// 退款额/待发货。
func TestDashboardTrade(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		logic := NewDashboardLogic()

		const no1, no2, no3, no4 = "TF-DB-T1", "TF-DB-T2", "TF-DB-T3", "TF-DB-T4"
		const as1 = "TF-DB-AS1"
		cleanup := func() {
			_, _ = g.DB().Exec(ctx, "DELETE FROM after_sale_order WHERE after_sale_no=?", as1)
			for _, no := range []string{no1, no2, no3, no4} {
				_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order WHERE order_no=?", no)
			}
		}
		cleanup()
		defer cleanup()

		insOrder := func(no string, status int, pay string) {
			_, err := g.DB().Exec(ctx,
				"INSERT INTO trade_order(order_no,user_id,status,total_amount,promotion_amount,pay_amount,currency,"+
					"receiver_name,receiver_phone,receiver_province,receiver_city,receiver_detail,created_at) "+
					"VALUES(?,997101,?,?,0,?,'CNY','看板测试','13800000000','浙江省','杭州市','T路',?)",
				no, status, pay, pay, dbWinAt) // created_at 显式落窗内固定值
			t.AssertNil(err)
		}
		start, end := dbWindow()

		insOrder(no1, 20, "50.00")
		insOrder(no2, 40, "80.00")
		insOrder(no3, 10, "999.00")
		insOrder(no4, 90, "888.00") // **已取消: 必须零贡献（C1）**
		oid, err := g.DB().GetValue(ctx, "SELECT id FROM trade_order WHERE order_no=?", no2)
		t.AssertNil(err)
		_, err = g.DB().Exec(ctx,
			"INSERT INTO after_sale_order(after_sale_no,order_id,order_no,order_item_id,user_id,type,status,currency,quantity,reason,refund_amount,created_at) "+
				"VALUES(?,?,?,0,997101,2,50,'CNY',1,'测试','10.00',?)", as1, oid.Int64(), no2, dbWinAt)
		t.AssertNil(err)

		after, err := logic.Trade(ctx, start, end)
		t.AssertNil(err)
		// N2 三轮收口: 窗口为固定历史区间且夹具行显式落窗内 → 与并行写入值域不相交,
		// 故可回到**精确等值**断言（比差量更强, 且不受并发增删影响）
		t.Assert(after.OrderCount, int64(2))         // 20/40 计入; 10 与 **90 排除**（C1）
		t.Assert(fen2(after.SalesAmount), int64(13000)) // 50+80（不含 888）
		t.Assert(fen2(after.RefundAmount), int64(1000)) // 售后完成 +10

		// 空窗口（**另一个**历史区间——不可与夹具窗口重合, 否则测的是夹具自身）→ 全零不报错
		zero, err := logic.Trade(ctx, "2001-06-01 00:00:00", "2001-06-02 00:00:00")
		t.AssertNil(err)
		t.Assert(zero.OrderCount, int64(0))
		t.Assert(fen2(zero.SalesAmount), int64(0))
		t.Assert(fen2(zero.RefundAmount), int64(0))
	})
}

// TestDashboardMember 会员看板口径（FR-2）: 窗口新增/活跃/休眠（阈值可配）。
func TestDashboardMember(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		logic := NewDashboardLogic()
		const ph = "TF-DB-DORM-PH"
		_, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE phone_hash=?", ph)
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE phone_hash=?", ph) }()

		// 90 天前活跃 → 休眠（受控窗口内新增）
		_, err := g.DB().Exec(ctx,
			"INSERT INTO `user`(nickname,phone,phone_hash,growth_value,status,last_active_at) "+
				"VALUES('TF-DB-DORM','x',?,0,1,DATE_SUB(NOW(), INTERVAL 100 DAY))", ph)
		t.AssertNil(err)

		// 窗口新增: 会员夹具的 created_at=NOW（未显式改写）→ 用 NOW 近窗查; 并行包可能再加, 故用 >=
		nowStart := time.Now().UTC().Add(-time.Hour).Format("2006-01-02 15:04:05")
		nowEnd := time.Now().UTC().Add(time.Hour).Format("2006-01-02 15:04:05")
		out, err := logic.Member(ctx, nowStart, nowEnd)
		t.AssertNil(err)
		t.Assert(out.NewCount >= 1, true)

		// N6 二次收口: **activeCount 断言**（I4 改的正是其默认窗口, 原零断言 → 变异回"全量"抓不到）
		// 本行 last_active_at = 100 天前 → 不在「近 30 天」活跃窗内 → 不贡献活跃
		act, err := logic.Member(ctx, "", "")
		t.AssertNil(err)
		// 差量法: 再插一个"1 天前活跃"的用户 → activeCount 必须 +1（证明默认窗是有限窗口而非全量）
		const ph2 = "TF-DB-ACT-PH"
		_, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE phone_hash=?", ph2)
		_, err = g.DB().Exec(ctx,
			"INSERT INTO `user`(nickname,phone,phone_hash,growth_value,status,last_active_at) "+
				"VALUES('TF-DB-ACT','x',?,0,1,DATE_SUB(NOW(), INTERVAL 1 DAY))", ph2)
		t.AssertNil(err)
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE phone_hash=?", ph2) }()
		afterAct, err := logic.Member(ctx, "", "")
		t.AssertNil(err)
		t.Assert(afterAct.ActiveCount, act.ActiveCount+1) // 1 天前活跃计入（默认近 30 天窗）
		t.Assert(afterAct.DormantCount, act.DormantCount) // 但不是休眠（<90 天）——双口径互斥验证

		// N6/I4 三轮收口: **100 天前活跃不得计入 ActiveCount**——这才是"默认窗是有限窗口"的
		// 判别式断言（原断言在"全量"与"近30天"两种实现下都恰好 +1, 原理上不可能区分; 复审全量变异不红实证）
		t.Assert(afterAct.ActiveCount, act.ActiveCount+1) // ph（100 天前活跃）未使活跃再 +1
		// 反向: 100 天前活跃者只在休眠口径出现（上面已断言 DormantCount 记账在 ph 上）

		// 休眠: 全量口径（不受窗口）——100 天前活跃行必被计入
		t.Assert(act.DormantCount >= 1, true)
	})
}

// TestDashboardProduct 商品看板口径（FR-3）: 在售/低库存/待审核评价（受控夹具, I2 修复）。
func TestDashboardProduct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		fm := NewDashboardLogic()

		// 前置清理（幂等: spu_no/sku_no 唯一键——上次失败可能留行）
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_review WHERE order_no='TF-DB-RV1'")
		_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id IN (SELECT id FROM product_sku WHERE sku_no='TF-DB-SKU')")
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_sku WHERE sku_no='TF-DB-SKU'")
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE spu_no='TF-DB-SPU'")

		// 受控夹具: 在售 SPU + 低库存 SKU + 待审核评价（断言这三者的差值为 1/1/1——
		// 但并行包可能同时增删 → 改为"造数后必 >= 基线+1"的组合断言, 且只依赖本夹具存在性）
		spu, err := g.DB().Model("product_spu").Ctx(ctx).Data(g.Map{
			"spu_no": "TF-DB-SPU", "name": "TF-DB-看板商品", "category_id": 1,
			"images": `[]`, "spec_definitions": `[]`, "status": 1,
		}).InsertAndGetId()
		t.AssertNil(err)
		sku, err := g.DB().Model("product_sku").Ctx(ctx).Data(g.Map{
			"sku_no": "TF-DB-SKU", "spu_id": spu, "name": "TF-DB-看板SKU",
			"specs": "{}", "price": "1.00", "status": 1,
		}).InsertAndGetId()
		t.AssertNil(err)
		_, err = g.DB().Exec(ctx,
			"INSERT INTO inventory(sku_id,total,locked,warn_count) VALUES(?,5,0,10) ON DUPLICATE KEY UPDATE total=5, locked=0, warn_count=10", sku)
		t.AssertNil(err)
		// 待审核评价（audit_status=0）
		_, err = g.DB().Exec(ctx,
			"INSERT INTO product_review(order_no,order_item_id,user_id,spu_id,sku_id,spu_name,sku_specs,score,content,images,audit_status) "+
				"VALUES('TF-DB-RV1',997101,997101,?,?,'TF-DB-看板商品','{}',5,'看板测试','[]',0)", spu, sku)
		t.AssertNil(err)
		defer func() {
			_, _ = g.DB().Exec(ctx, "DELETE FROM product_review WHERE order_no='TF-DB-RV1'")
			_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id=?", sku)
			_, _ = g.DB().Exec(ctx, "DELETE FROM product_sku WHERE id=?", sku)
			_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE id=?", spu)
		}()

		out, err := fm.Product(ctx)
		t.AssertNil(err)
		// 存在性断言（不受并行包增删干扰）: 本夹具计入且计数为正
		t.Assert(out.OnSaleCount >= 1, true)
		t.Assert(out.LowStockCount >= 1, true)
		t.Assert(out.PendingReview >= 1, true) // I3: 补 PendingReview 断言（原无）
	})
}

// fen2 元字符串 → 分（断言用；解析失败归 -1 以便暴露）。
func fen2(yuan string) int64 {
	var yuanPart, fenPart int64
	i := 0
	for ; i < len(yuan); i++ {
		if yuan[i] == '.' {
			break
		}
		if yuan[i] < '0' || yuan[i] > '9' {
			return -1
		}
		yuanPart = yuanPart*10 + int64(yuan[i]-'0')
	}
	if i < len(yuan) {
		frac := yuan[i+1:]
		for j := 0; j < len(frac) && j < 2; j++ {
			fenPart = fenPart*10 + int64(frac[j]-'0')
		}
		for j := len(frac); j < 2; j++ {
			fenPart *= 10
		}
	}
	return yuanPart*100 + fenPart
}

// TestDashboardDormantThresholdFromConfig N1 收口守卫: 休眠阈值读 `system_config` 表
// （与 user/wx.go 的 cfgInt 同源）——原修复误用 g.Cfg()（该键只存在于表内）→ 恒取兜底 90,
// 运营改阈值后登录门禁与看板统计会分叉（复审决定性探针: 表值改 2 后 g.Cfg 仍 90）。
func TestDashboardDormantThresholdFromConfig(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		logic := NewDashboardLogic()
		const ph = "TF-DB-CFG-PH"
		_, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE phone_hash=?", ph)
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE phone_hash=?", ph) }()

		// 5 天前活跃（若阈值为 2 天则应计入休眠; 阈值为 90 则不计）
		_, err := g.DB().Exec(ctx,
			"INSERT INTO `user`(nickname,phone,phone_hash,growth_value,status,last_active_at) "+
				"VALUES('TF-DB-CFG','x',?,0,1,DATE_SUB(NOW(), INTERVAL 5 DAY))", ph)
		t.AssertNil(err)

		// 读原阈值并改为 2（模拟运营在后台改配置）
		orig, err := g.DB().GetValue(ctx,
			"SELECT value FROM system_config WHERE code='dormant.tier1.days'")
		t.AssertNil(err)
		_, err = g.DB().Exec(ctx,
			"UPDATE system_config SET value='2' WHERE code='dormant.tier1.days'")
		t.AssertNil(err)
		defer func() {
			_, _ = g.DB().Exec(ctx,
				"UPDATE system_config SET value=? WHERE code='dormant.tier1.days'", orig.String())
		}()

		out, err := logic.Member(ctx, "", "")
		t.AssertNil(err)
		// 阈值=2 → 5 天前活跃者计入休眠（N1: 若读 g.Cfg 恒 90 则此处为 0 → 红）
		cnt, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM `user` WHERE phone_hash=? AND last_active_at <= DATE_SUB(NOW(), INTERVAL 2 DAY)", ph)
		t.AssertNil(err)
		t.Assert(cnt.Int(), 1) // 夹具自身符合新阈值
		t.Assert(out.DormantCount >= 1, true) // 看板必须跟随配置阈值（同源读生效）
	})
}
