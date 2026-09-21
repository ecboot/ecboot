package user

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"

	"ecboot/internal/errcode"
	"ecboot/internal/model"
)

// ---- 地址 fixture ----

// seedAddress 建测试地址（用后 cleanupAddresses）。
func seedAddress(ctx context.Context, t *gtest.T, userId int64, name string, isDefault bool) int64 {
	def := 0
	if isDefault {
		def = 1
	}
	res, err := g.DB().Exec(ctx,
		"INSERT INTO user_address(user_id,receiver_name,receiver_phone,province,province_code,city,city_code,district,district_code,detail_address,is_default) "+
			"VALUES(?,?,'13800000000','浙江省','330000','杭州市','330100','西湖区','330106',?,?)",
		userId, name, "T路1号", def)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

func cleanupAddresses(ctx context.Context, t *gtest.T, userId int64) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM user_address WHERE user_id=?", userId)
}

// TestAddressCrud 地址 CRUD + 归属校验（FR-004, US2 验收 1/3/4/5）。
func TestAddressCrud(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900003001"
		defer cleanupMember(ctx, t, phone)
		uid := seedMember(ctx, t, phone, "地址测试", 0)
		defer cleanupAddresses(ctx, t, uid)

		// 新建（含区划码）
		id, err := AddressCreate(ctx, uid, model.AddressInput{
			ReceiverName: "张三", ReceiverPhone: "13800001111",
			Province: "浙江省", City: "杭州市", District: "西湖区",
			ProvinceCode: "330000", CityCode: "330100", DistrictCode: "330106",
			DetailAddress: "文一路1号",
		})
		t.AssertNil(err)
		t.AssertGT(id, 0)

		// 列表可见（电话脱敏）
		list, err := AddressList(ctx, uid)
		t.AssertNil(err)
		t.Assert(len(list), 1)
		t.Assert(list[0].ReceiverName, "张三")
		t.Assert(list[0].DistrictCode, "330106")
		t.Assert(list[0].ReceiverPhone, "138****1111") // 脱敏

		// 修改
		t.AssertNil(AddressUpdate(ctx, uid, id, model.AddressInput{
			ReceiverName: "李四", ReceiverPhone: "13800001111",
			Province: "浙江省", City: "杭州市", District: "西湖区",
			ProvinceCode: "330000", CityCode: "330100", DistrictCode: "330106",
			DetailAddress: "文二路2号",
		}))
		list, err = AddressList(ctx, uid)
		t.AssertNil(err)
		t.Assert(list[0].ReceiverName, "李四")
		t.Assert(list[0].DetailAddress, "文二路2号")

		// 越权: 他人（另一会员）改/删/设默认 → 10006
		other := seedMember(ctx, t, "13900003999", "他人", 0)
		defer cleanupMember(ctx, t, "13900003999")
		err = AddressUpdate(ctx, other, id, model.AddressInput{
			ReceiverName: "x", ReceiverPhone: "13800001111", Province: "p", City: "c",
			ProvinceCode: "330000", CityCode: "330100", DetailAddress: "d",
		})
		t.Assert(errCode(err), errcode.CodeNotFound)
		err = AddressDelete(ctx, other, id)
		t.Assert(errCode(err), errcode.CodeNotFound)
		err = AddressSetDefault(ctx, other, id)
		t.Assert(errCode(err), errcode.CodeNotFound)

		// 删除
		t.AssertNil(AddressDelete(ctx, uid, id))
		list, err = AddressList(ctx, uid)
		t.AssertNil(err)
		t.Assert(len(list), 0)
	})
}

// TestAddressSetDefault 设默认唯一（FR-005）: 同事务清其他, 不留双默认。
func TestAddressSetDefault(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900003002"
		defer cleanupMember(ctx, t, phone)
		uid := seedMember(ctx, t, phone, "默认测试", 0)
		defer cleanupAddresses(ctx, t, uid)

		a := seedAddress(ctx, t, uid, "A", true)
		b := seedAddress(ctx, t, uid, "B", false)

		// 设 B 默认 → A 的默认标记被清（无双默认）
		t.AssertNil(AddressSetDefault(ctx, uid, b))
		n, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM user_address WHERE user_id=? AND is_default=1 AND deleted=0", uid)
		t.AssertNil(err)
		t.Assert(n.Int(), 1) // 唯一默认
		def, err := g.DB().GetValue(ctx, "SELECT id FROM user_address WHERE user_id=? AND is_default=1", uid)
		t.AssertNil(err)
		t.Assert(def.Int64(), b)

		// A 已非默认；再设 A 默认 → 切回
		t.AssertNil(AddressSetDefault(ctx, uid, a))
		def, err = g.DB().GetValue(ctx, "SELECT id FROM user_address WHERE user_id=? AND is_default=1", uid)
		t.AssertNil(err)
		t.Assert(def.Int64(), a)
	})
}

// TestAddressGetForOrder 下单取址（FR-006 内部方法）: 归属校验 + 含区划码。
func TestAddressGetForOrder(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900003003"
		defer cleanupMember(ctx, t, phone)
		uid := seedMember(ctx, t, phone, "取址测试", 0)
		defer cleanupAddresses(ctx, t, uid)

		id := seedAddress(ctx, t, uid, "取址", false)
		got, err := AddressGetForOrder(ctx, uid, id)
		t.AssertNil(err)
		t.Assert(got.ProvinceCode, "330000")
		t.Assert(got.CityCode, "330100")
		t.Assert(got.DistrictCode, "330106")

		// 他人地址 → 10006
		other := seedMember(ctx, t, "13900003998", "他人2", 0)
		defer cleanupMember(ctx, t, "13900003998")
		_, err = AddressGetForOrder(ctx, other, id)
		t.Assert(errCode(err), errcode.CodeNotFound)
	})
}
