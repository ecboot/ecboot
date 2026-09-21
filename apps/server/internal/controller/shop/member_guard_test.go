// member_guard_test.go 会员端点未登录防御（I8）与购物车修改三态在 **controller 层** 的契约回归（I9）。
// 为什么测 controller: I8/I9 两处缺陷都只在 controller 层可见——service 层签名（userId int64 / checked *bool）
// 一直是对的, 出问题的是 controller 怎样取值（userId=0 哨兵直接下传、bool 包成恒非 nil 指针）。
// 这是本仓第一个 controller 层测试: 直接构造 ctx + 调端点方法, 不经 HTTP。
package shop

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/database/gdb"
	gredis "github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/api/shop/v1"
	"ecboot/internal/errcode"
	"ecboot/internal/middleware"

	"ecboot/internal/testutil"
)

func init() {
	// 与 main.go / 各 service 测试基座同一口径（009 评审 C2: 进程与库内 UTC 墙钟一致）
	time.Local = time.UTC
	_ = os.Setenv("ECBOOT_MOCK", "true")
	_ = gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			{Link: testutil.DSN()},
		},
	})
	gredis.SetConfig(&gredis.Config{Address: "127.0.0.1:6379", Db: 0})
}

// codeOf 提取业务错误码（与 service 测试的 errCode 同构）。
func codeOf(err error) int {
	var ge *gerror.Error
	if errors.As(err, &ge) {
		return ge.Code().Code()
	}
	return -1
}

// TestRequireMember 未登录防御（I8）: 哨兵 userId<=0 一律 10003, 绝不把 0 下传给"后台视角"分支。
// 背景: OrderDetail/Cancel 以 userId=0 语义为"后台视角不限归属"——C 端端点若漏挂鉴权,
// 缺此防御就会静默变成"任意订单可见/可取消"。
func TestRequireMember(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, v := range []int64{0, -1} {
			ctx := context.WithValue(context.Background(), middleware.CtxUserId, v)
			id, err := requireMember(ctx)
			t.Assert(id, int64(0))
			t.Assert(codeOf(err), errcode.CodeUnauthorized)
		}
		// 无凭证上下文（值缺失）
		_, err := requireMember(context.Background())
		t.Assert(codeOf(err), errcode.CodeUnauthorized)
		// 正常会员原样透传
		ctx := context.WithValue(context.Background(), middleware.CtxUserId, int64(7))
		id, err := requireMember(ctx)
		t.AssertNil(err)
		t.Assert(id, int64(7))
	})
}

// seedOrderCreateSku 下单所需的最小商品数据（SPU/SKU/库存）, 返回 skuId。
func seedOrderCreateSku(ctx context.Context, t *gtest.T) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id IN (SELECT id FROM product_sku WHERE sku_no='CTL-SKU-001')")
	_, _ = g.DB().Exec(ctx, "DELETE ps FROM product_sku ps JOIN product_spu s ON ps.spu_id=s.id WHERE s.name='CTL-商品'")
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE name='CTL-商品'")

	catId, err := g.DB().Model("product_category").Ctx(ctx).Data(g.Map{
		"parent_id": 0, "name": "CTL-分类", "level": 3, "status": 1,
	}).InsertAndGetId()
	t.AssertNil(err)
	spuId, err := g.DB().Model("product_spu").Ctx(ctx).Data(g.Map{
		"spu_no": "CTL-SPU-001", "name": "CTL-商品", "category_id": catId,
		"images": `["http://img/tf.png"]`, "status": 1,
	}).InsertAndGetId()
	t.AssertNil(err)
	skuId, err := g.DB().Model("product_sku").Ctx(ctx).Data(g.Map{
		"sku_no": "CTL-SKU-001", "spu_id": spuId, "name": "CTL-商品 黑",
		"specs": `{"颜色":"黑"}`, "price": "10.00", "status": 1,
	}).InsertAndGetId()
	t.AssertNil(err)
	_, err = g.DB().Exec(ctx, "INSERT INTO inventory(sku_id,total,locked) VALUES(?,10,0)", skuId)
	t.AssertNil(err)
	return skuId
}

func cleanupOrderCreateSku(ctx context.Context) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM inventory WHERE sku_id IN (SELECT id FROM product_sku WHERE sku_no='CTL-SKU-001')")
	_, _ = g.DB().Exec(ctx, "DELETE ps FROM product_sku ps JOIN product_spu s ON ps.spu_id=s.id WHERE s.name='CTL-商品'")
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE name='CTL-商品'")
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_category WHERE name='CTL-分类'")
}

// TestOrderCreateController 下单端点 controller→service 全链路（修复轮 C4a/C4b 的终端验证）。
// 背景: 批次 06 声称 create 端点"已连线", 但 service 层两处快照写错使该端点**必然 500**——
// 端点级验证只有走到 controller 才能覆盖请求映射 + 鉴权防御 + service 落库这一整条。
func TestOrderCreateController(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// uid 必须是 int64: 无类型常量会被存成 int, 而 CtxUserIdFrom 断言 int64 → 取不到值（退化为未登录）
		const uid int64 = 8804
		ctx := context.WithValue(context.Background(), middleware.CtxUserId, uid)
		_, _ = g.DB().Exec(ctx, "DELETE FROM user_address WHERE user_id=?", uid)
		_, _ = g.DB().Exec(ctx, "DELETE l FROM trade_order_log l JOIN trade_order o ON l.order_id=o.id WHERE o.user_id=?", uid)
		_, _ = g.DB().Exec(ctx, "DELETE i FROM trade_order_item i JOIN trade_order o ON i.order_id=o.id WHERE o.user_id=?", uid)
		_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order WHERE user_id=?", uid)
		defer func() {
			_, _ = g.DB().Exec(ctx, "DELETE FROM user_address WHERE user_id=?", uid)
			_, _ = g.DB().Exec(ctx, "DELETE l FROM trade_order_log l JOIN trade_order o ON l.order_id=o.id WHERE o.user_id=?", uid)
			_, _ = g.DB().Exec(ctx, "DELETE i FROM trade_order_item i JOIN trade_order o ON i.order_id=o.id WHERE o.user_id=?", uid)
			_, _ = g.DB().Exec(ctx, "DELETE FROM trade_order WHERE user_id=?", uid)
		}()
		skuId := seedOrderCreateSku(ctx, t)
		defer cleanupOrderCreateSku(ctx)

		addrRes, err := g.DB().Exec(ctx,
			"INSERT INTO user_address(user_id,receiver_name,receiver_phone,province,city,district,detail_address,is_default) "+
				"VALUES(?,'控制器收货人','13800000001','浙江省','杭州市','西湖区','T路2号',1)", uid)
		t.AssertNil(err)
		addrId, _ := addrRes.LastInsertId()

		c := &ControllerV1{}
		out, err := c.OrderCreate(ctx, &v1.OrderCreateReq{
			RequestToken: "T-CTL-TOKEN", AddressId: fmtID(addrId), SkuId: fmtID(skuId), Quantity: 2,
		})
		t.AssertNil(err)
		t.Assert(out.OrderNo != "", true)
		// 金额只钉恒等式（不写死"20.00"——满减取全局最早在架活动, 那是环境假设; 评审 Minor）
		hdr, err := g.DB().GetOne(ctx,
			"SELECT total_amount, promotion_amount, pay_amount FROM trade_order WHERE order_no=?", out.OrderNo)
		t.AssertNil(err)
		t.Assert(hdr["total_amount"].String(), "20.00")
		t.Assert(hdr["pay_amount"].Float64(), hdr["total_amount"].Float64()-hdr["promotion_amount"].Float64())
		t.Assert(out.PayAmount, hdr["pay_amount"].String())

		// 未登录 → 10003（不得到达 service）
		_, err = c.OrderCreate(context.Background(), &v1.OrderCreateReq{
			RequestToken: "T-CTL-TOKEN-2", AddressId: fmtID(addrId), SkuId: fmtID(skuId), Quantity: 1,
		})
		t.Assert(codeOf(err), errcode.CodeUnauthorized)
	})
}

// TestCartUpdateItemControllerThreeState controller 层三态（I9）: 请求不传 checked → 不改勾选。
// 原实现 `checked := req.Checked; &checked` 使指针恒非 nil, "只改数量"会顺手清掉勾选。
func TestCartUpdateItemControllerThreeState(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.WithValue(context.Background(), middleware.CtxUserId, int64(8802))
		const uid = 8802
		_, _ = g.DB().Exec(ctx, "DELETE FROM cart_item WHERE user_id=?", uid)
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM cart_item WHERE user_id=?", uid) }()

		res, err := g.DB().Exec(ctx,
			"INSERT INTO cart_item(user_id,sku_id,quantity,checked) VALUES(?,990000777,1,1)", uid)
		t.AssertNil(err)
		itemId, _ := res.LastInsertId()

		c := &ControllerV1{}
		// 只改数量（checked 缺省 nil）
		out, err := c.CartUpdateItem(ctx, &v1.CartUpdateItemReq{ItemId: fmtID(itemId), Quantity: 3})
		t.AssertNil(err)
		t.Assert(out.Success, true)
		rec, err := g.DB().GetOne(ctx, "SELECT quantity, checked FROM cart_item WHERE id=?", itemId)
		t.AssertNil(err)
		t.Assert(rec["quantity"].Int(), 3)
		t.Assert(rec["checked"].Int(), 1) // 勾选未被误清

		// 显式取消勾选
		no := false
		_, err = c.CartUpdateItem(ctx, &v1.CartUpdateItemReq{ItemId: fmtID(itemId), Checked: &no})
		t.AssertNil(err)
		v, err := g.DB().GetValue(ctx, "SELECT checked FROM cart_item WHERE id=?", itemId)
		t.AssertNil(err)
		t.Assert(v.Int(), 0)

		// 未登录 → 10003（I8 与 I9 的交汇: 无凭证不得触达购物车）
		_, err = c.CartUpdateItem(context.Background(), &v1.CartUpdateItemReq{ItemId: fmtID(itemId)})
		t.Assert(codeOf(err), errcode.CodeUnauthorized)
	})
}
