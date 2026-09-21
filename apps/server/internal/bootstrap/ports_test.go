// ports_test.go 跨域端口装配回归（I5: 券查询口）。
// 装配是"静默失败"的高发区——漏注入不会编译报错, 只会让 /cart/checkout 的 usableCoupons 恒为空,
// 因此这里既断言"注入了", 也走一遍真实链路（shop 端口 → 适配器 → user 域实现 → 库）。
package bootstrap

import (
	"context"
	"os"
	"testing"
	"time"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/service/shop"

	"ecboot/internal/testutil"
)

func init() {
	time.Local = time.UTC // 与 main.go 及各 service 测试基座同一口径
	_ = os.Setenv("ECBOOT_MOCK", "true")
	_ = gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			{Link: testutil.DSN()},
		},
	})
}

// TestCouponQueryPortWired 装配断言: init 必须完成注入。
func TestCouponQueryPortWired(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(shop.CouponQuery != nil, true)
	})
}

// TestCouponQueryPortE2E 可用券链路: 门槛内返回、门槛外过滤。
func TestCouponQueryPortE2E(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const uid = 8803
		_, _ = g.DB().Exec(ctx, "DELETE FROM coupon WHERE name='TF-端口满100减20'")
		couponId, err := g.DB().Model("coupon").Ctx(ctx).Data(g.Map{
			"name": "TF-端口满100减20", "type": 1, "threshold_amount": "100.00", "discount_amount": "20.00",
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
		defer func() { _, _ = g.DB().Exec(ctx, "DELETE FROM user_coupon WHERE id=?", ucId) }()

		// 150 元: 达门槛 → 命中
		got, err := shop.CouponQuery.UsableForOrder(ctx, uid, "150.00")
		t.AssertNil(err)
		t.Assert(len(got), 1)
		t.Assert(got[0].UserCouponId, ucId)
		t.Assert(got[0].Name, "TF-端口满100减20")
		t.Assert(got[0].Discount, "20.00")

		// 50 元: 未达门槛 → 过滤
		got, err = shop.CouponQuery.UsableForOrder(ctx, uid, "50.00")
		t.AssertNil(err)
		t.Assert(len(got), 0)
	})
}
