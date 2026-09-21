// order_create_test.go 下单九步事务（012 修复轮补测）。
// 背景: 批次 06 对 `OrderLogicImpl.Create` **零测试覆盖**, 且收口冒烟 6 项恰好绕开下单端点,
// 致两处快照写入缺陷长期存在——① 订单项漏写 4 个非空列（sku_no/spu_name/sku_name/original_price）
// ② 状态流水用不存在的 `operator` 列（与 I6 同根因, 评审仅在 Cancel 中捕获）。
// 两处均使下单事务报错回滚, 即 C 端 `/shop/orders` 实际不可用。
package shop

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// seedUserAddress 建测试收货地址（下单快照来源）。
func seedUserAddress(ctx context.Context, t *gtest.T, userId int64) int64 {
	res, err := g.DB().Exec(ctx,
		"INSERT INTO user_address(user_id,receiver_name,receiver_phone,province,city,district,detail_address,is_default) "+
			"VALUES(?,?,?,?,?,?,?,1)",
		userId, "测试收货人", "13800000000", "浙江省", "杭州市", "西湖区", "T路1号")
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

// cleanupOrderCreate 清理下单相关数据（流水/订单项/订单/地址）。
func cleanupOrderCreate(ctx context.Context, userId int64) {
	_, _ = g.DB().Exec(ctx, "DELETE l FROM trade_order_log l JOIN trade_order o ON l.order_id=o.id WHERE o.user_id=?", userId)
	_, _ = g.DB().Exec(ctx, "DELETE i FROM trade_order_item i JOIN trade_order o ON i.order_id=o.id WHERE o.user_id=?", userId)
	_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order WHERE user_id=?", userId)
	_, _ = g.DB().Exec(ctx, "DELETE FROM user_address WHERE user_id=?", userId)
}

// TestOrderCreateSnapshot 下单成功（FR-001）: 订单头/订单项快照（非空列写全）/状态流水（操作者两列）/
// 库存锁定（只动 locked, 不动 total）/金额账本。
func TestOrderCreateSnapshot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "ORD-CRT-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		defer cleanupOrderCreate(ctx, uid)

		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		addrId := seedUserAddress(ctx, t, uid)

		out, err := NewOrderLogic().Create(ctx, uid, model.OrderCreateInput{
			RequestToken: "T-CRT-TOKEN-1", AddressId: addrId, SkuId: f.SkuId, Quantity: 2,
		})
		t.AssertNil(err)
		t.Assert(out.OrderNo != "", true)

		// 订单头: 状态 10 + 收货快照 + 幂等 token
		// 金额断言只钉**恒等式**与自造数据（总额=本 fixture 的 10.00×2）, 不写死"实付 20.00"——
		// 后者隐含"测试库无在架促销活动"这一环境假设（满减取全局最早的在架活动, 并行包可能插入）,
		// 评审 Minor 指出该假设未被强制。
		rec, err := g.DB().GetOne(ctx,
			"SELECT id, status, total_amount, promotion_amount, pay_amount, receiver_name, request_token "+
				"FROM trade_order WHERE order_no=?", out.OrderNo)
		t.AssertNil(err)
		t.Assert(rec["status"].Int(), 10)
		t.Assert(rec["total_amount"].String(), "20.00")
		t.Assert(rec["pay_amount"].Float64(),
			rec["total_amount"].Float64()-rec["promotion_amount"].Float64()) // 账本恒等式
		t.Assert(out.PayAmount, rec["pay_amount"].String()) // 出参与落库一致
		t.Assert(rec["receiver_name"].String(), "测试收货人")
		t.Assert(rec["request_token"].String(), "T-CRT-TOKEN-1")

		// 订单项快照: 4 个非空列（sku_no/spu_name/sku_name/original_price）必须写全
		it, err := g.DB().GetOne(ctx,
			"SELECT sku_no, spu_name, sku_name, original_price, price, quantity, promotion_amount, pay_amount "+
				"FROM trade_order_item WHERE order_id=?", rec["id"].Int64())
		t.AssertNil(err)
		t.Assert(it["sku_no"].String(), "TF-SKU-001")
		t.Assert(it["spu_name"].String(), "TF-商品")
		t.Assert(it["sku_name"].String(), "TF-商品 黑")
		t.Assert(it["original_price"].String(), "10.00")
		t.Assert(it["price"].String(), "10.00")
		t.Assert(it["quantity"].Int(), 2)
		t.Assert(it["pay_amount"].Float64(),
			it["price"].Float64()*2-it["promotion_amount"].Float64()) // 行账本恒等式

		// 状态流水: 建单 to_status=10; operator_type=2（表注释: 1系统 2用户 3管理员）
		lg, err := g.DB().GetOne(ctx,
			"SELECT to_status, operator_type, operator_id, remark FROM trade_order_log WHERE order_id=?",
			rec["id"].Int64())
		t.AssertNil(err)
		t.Assert(lg["to_status"].Int(), 10)
		t.Assert(lg["operator_type"].Int(), 2)
		t.Assert(lg["operator_id"].String(), "user:"+fmt.Sprint(uid))
		t.Assert(lg["remark"].String(), "订单创建")

		// 库存锁定: locked +2, total 不变
		inv, err := g.DB().GetOne(ctx, "SELECT total, locked FROM inventory WHERE sku_id=?", f.SkuId)
		t.AssertNil(err)
		t.Assert(inv["total"].Int(), 100)
		t.Assert(inv["locked"].Int(), 2)
	})
}

// TestOrderCreateAllocatesPromotionPerLine 行分摊三列（评审 Important）: 券/满减/积分的**行分摊列**
// 必须按各自构成分摊（表契约 000019/000020: "合计=订单头对应明细, 尾差记末行"）。
// 原实现只写聚合列 promotion_amount, 三列恒 0.00 → 与订单头"有券"矛盾, 而 trade_order_item 是
// 售后按行取数的依据（退款/退回按行分摊）→ 账实不符。
func TestOrderCreateAllocatesPromotionPerLine(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "ORD-ALLOC-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		defer cleanupOrderCreate(ctx, uid)
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM cart_item WHERE user_id=?", uid) }()

		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		addrId := seedUserAddress(ctx, t, uid)

		// 第二件商品（同 SPU, 价格同 10.00）——两行才能验证"尾差记末行"
		_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id IN (SELECT id FROM product_sku WHERE sku_no='TF-SKU-002')")
		_, _ = g.DB().Exec(ctx, "DELETE FROM product_sku WHERE sku_no='TF-SKU-002'")
		sku2Res, err := g.DB().Exec(ctx,
			"INSERT INTO product_sku(sku_no,spu_id,name,specs,price,status) VALUES('TF-SKU-002',?,?,?,?,1)",
			f.SpuId, "TF-商品 白", `{"颜色":"白"}`, "10.00")
		t.AssertNil(err)
		sku2, _ := sku2Res.LastInsertId()
		_, err = g.DB().Exec(ctx, "INSERT INTO inventory(sku_id,total,locked) VALUES(?,100,0)", sku2)
		t.AssertNil(err)
		defer func() {
			_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id=?", sku2)
			_, _ = g.DB().Exec(ctx, "DELETE FROM product_sku WHERE id=?", sku2)
		}()

		// 无门槛券: 满 0 减 5.00（避免依赖测试库是否有在架促销活动）
		_, _ = g.DB().Exec(ctx, "DELETE FROM coupon WHERE name='TF-分摊券'")
		couponId, err := g.DB().Model("coupon").Ctx(ctx).Data(g.Map{
			"name": "TF-分摊券", "type": 2, "threshold_amount": "0.00", "discount_amount": "5.00",
			"total_count": 0, "per_limit": 1, "valid_type": 1, "status": 1, "deleted": 0,
		}).InsertAndGetId()
		t.AssertNil(err)
		defer func() {
			_, _ = g.DB().Exec(ctx, "DELETE FROM user_coupon WHERE coupon_id=?", couponId)
			_, _ = g.DB().Exec(ctx, "DELETE FROM coupon WHERE id=?", couponId)
		}()
		ucId, err := g.DB().Model("user_coupon").Ctx(ctx).Data(g.Map{
			"user_id": uid, "coupon_id": couponId, "status": 1,
			"expire_time": gtime.Now().AddDate(0, 0, 7),
		}).InsertAndGetId()
		t.AssertNil(err)

		// 两行: 1 件（1000 分）+ 2 件（2000 分）; 券 500 分按行金额比例分摊 → 166 / 334（尾差记末行）
		item1, err := g.DB().Model("cart_item").Ctx(ctx).Data(g.Map{
			"user_id": uid, "sku_id": f.SkuId, "quantity": 1, "checked": 1,
		}).InsertAndGetId()
		t.AssertNil(err)
		item2, err := g.DB().Model("cart_item").Ctx(ctx).Data(g.Map{
			"user_id": uid, "sku_id": sku2, "quantity": 2, "checked": 1,
		}).InsertAndGetId()
		t.AssertNil(err)

		out, err := NewOrderLogic().Create(ctx, uid, model.OrderCreateInput{
			RequestToken: "T-ALLOC-1", AddressId: addrId,
			UserCouponId: ucId, CartItemIds: []int64{item1, item2},
		})
		t.AssertNil(err)
		t.Assert(out.PayAmount, "25.00") // 30.00 - 5.00

		// 订单头三列
		rec, err := g.DB().GetOne(ctx,
			"SELECT id, coupon_amount, promotion_amount, pay_amount FROM trade_order WHERE order_no=?", out.OrderNo)
		t.AssertNil(err)
		t.Assert(rec["coupon_amount"].String(), "5.00")
		t.Assert(rec["promotion_amount"].String(), "5.00")
		t.Assert(rec["pay_amount"].String(), "25.00")

		// 行分摊: 三列写全, 且合计等于订单头; 尾差落在末行
		items, err := g.DB().Model("trade_order_item").Ctx(ctx).
			Where("order_id", rec["id"].Int64()).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(items), 2)
		t.Assert(items[0]["coupon_amount"].String(), "1.66")
		t.Assert(items[1]["coupon_amount"].String(), "3.34")
		t.Assert(items[0]["full_reduction_amount"].String(), "0.00")
		t.Assert(items[0]["point_amount"].String(), "0.00")
		t.Assert(items[0]["promotion_amount"].String(), "1.66")
		t.Assert(items[0]["pay_amount"].String(), "8.34")
		// 合计恒等: 分摊列之和 == 订单头
		t.Assert(items[0]["coupon_amount"].Float64()+items[1]["coupon_amount"].Float64(), 5.0)
	})
}

// TestOrderCreateFlashSaleRejected 秒杀未接线 → 明确拒绝（评审 Important）。
// 原分支产出的单永远无法支付（不锁 inventory → 支付回调核销必未命中而回滚）, 且从不读 flash_price
// （按原价计费）。拒绝优于产出无法履约的单; 待批次 09 接通后删除该闸（见 PROGRESS §五 记账）。
func TestOrderCreateFlashSaleRejected(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "ORD-FS-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		defer cleanupOrderCreate(ctx, uid)

		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		addrId := seedUserAddress(ctx, t, uid)

		_, err := NewOrderLogic().Create(ctx, uid, model.OrderCreateInput{
			RequestToken: "T-FS-1", AddressId: addrId, SkuId: f.SkuId, Quantity: 1, FlashSaleItemId: 999,
		})
		t.Assert(errCode(err), errcode.CodeActivityInvalid)

		n, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM trade_order WHERE user_id=?", uid)
		t.AssertNil(err)
		t.Assert(n.Int(), 0)
		// 也未锁库存
		inv, err := g.DB().GetOne(ctx, "SELECT locked FROM inventory WHERE sku_id=?", f.SkuId)
		t.AssertNil(err)
		t.Assert(inv["locked"].Int(), 0)
	})
}

// TestOrderCreateNoTxLeak 事务不得泄漏（评审 Important）: 早退路径必须回滚并归还连接。
// 用**不会被取消**的 ctx（context.Background）连续触发"无有效购买项"早退, 观察 processlist 连接数:
// 修复前该路径漏给 err 赋值 → BEGIN 后既不提交也不回滚, 每调一次多挂一个连接（实测 15 次 → delta=15）。
// 真实 net/http 会因请求 ctx 取消而自愈, 故属"潜伏"缺陷; 但定时任务/后台 worker 这类长生命周期 ctx
// 的调用方会真实泄漏, 且一旦判空被挪到取锁之后即变持锁泄漏。此测试钉住"不依赖 ctx 取消"这一不变量。
func TestOrderCreateNoTxLeak(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "ORD-TXL-1"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		defer cleanupOrderCreate(ctx, uid)
		addrId := seedUserAddress(ctx, t, uid)

		conns := func() int {
			v, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM information_schema.processlist")
			t.AssertNil(err)
			return v.Int()
		}
		for i := 0; i < 3; i++ { // 预热连接池, 避免把池的惰性增长记成泄漏
			_ = conns()
		}
		before := conns()
		for i := 0; i < 15; i++ {
			_, err := NewOrderLogic().Create(ctx, uid, model.OrderCreateInput{
				RequestToken: fmt.Sprintf("T-TXL-%d", i), AddressId: addrId, // 无购买项 → 早退
			})
			t.AssertNE(err, nil)
		}
		// 阈值留 3 的余量: 同库被其他测试包并行使用, 连接数本身会小幅波动; 泄漏则是 +15。
		t.Assert(conns() <= before+3, true)
	})
}

// TestOrderCreateInsufficientStock 库存不足（FR-001）: 锁定条件更新未命中 → 40001, 且整单不落库
// （原实现不判 RowsAffected, 条件未命中仍继续下单 → 静默超卖）。
func TestOrderCreateInsufficientStock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const h1 = "ORD-CRT-2"
		defer cleanupPointUser(ctx, t, h1)
		uid := seedPointUser(ctx, t, h1)
		defer cleanupOrderCreate(ctx, uid)

		f := setupTradeFixture(t)
		defer cleanupFixture(ctx, f)
		addrId := seedUserAddress(ctx, t, uid)
		// 可售仅 1, 下单 2
		_, _ = g.DB().Exec(ctx, "UPDATE inventory SET total=1, locked=0 WHERE sku_id=?", f.SkuId)

		_, err := NewOrderLogic().Create(ctx, uid, model.OrderCreateInput{
			RequestToken: "T-CRT-TOKEN-2", AddressId: addrId, SkuId: f.SkuId, Quantity: 2,
		})
		t.Assert(errCode(err), errcode.CodeStockInsufficient)

		n, err := g.DB().GetValue(ctx, "SELECT COUNT(*) FROM trade_order WHERE user_id=?", uid)
		t.AssertNil(err)
		t.Assert(n.Int(), 0)
		inv, err := g.DB().GetOne(ctx, "SELECT total, locked FROM inventory WHERE sku_id=?", f.SkuId)
		t.AssertNil(err)
		t.Assert(inv["locked"].Int(), 0)
	})
}
