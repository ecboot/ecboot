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

// TestAddressUpdateKeepsDefaultAndFields C1 回归: 编辑（尤其默认地址）不得丢默认、不得清空未传字段。
func TestAddressUpdateKeepsDefaultAndFields(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900003010"
		defer cleanupMember(ctx, t, phone)
		uid := seedMember(ctx, t, phone, "编辑测试", 0)
		defer cleanupAddresses(ctx, t, uid)

		// 默认地址（完整字段）
		id := seedAddress(ctx, t, uid, "原名", true)

		// 仅改详址（未传 isDefault 与其他字段）
		t.AssertNil(AddressUpdate(ctx, uid, id, model.AddressInput{DetailAddress: "新详址"}))
		rec, err := g.DB().GetOne(ctx,
			"SELECT receiver_name, receiver_phone, province_code, district_code, detail_address, is_default FROM user_address WHERE id=?", id)
		t.AssertNil(err)
		t.Assert(rec["detail_address"].String(), "新详址") // 已更新
		t.Assert(rec["is_default"].Int(), 1)            // C1: 默认保持（此前被清零）
		t.Assert(rec["receiver_name"].String(), "原名")   // C1: 未传字段不被清空
		t.Assert(rec["receiver_phone"].String(), "13800000000")
		t.Assert(rec["province_code"].String(), "330000")
		t.Assert(rec["district_code"].String(), "330106")

		// 显式传 isDefault=true → 清其他（唯一）并置本条
		id2 := seedAddress(ctx, t, uid, "第二条", false)
		t.AssertNil(AddressUpdate(ctx, uid, id2, model.AddressInput{DetailAddress: "第二条改", IsDefault: true}))
		n, err := g.DB().GetValue(ctx,
			"SELECT COUNT(*) FROM user_address WHERE user_id=? AND is_default=1 AND deleted=0", uid)
		t.AssertNil(err)
		t.Assert(n.Int(), 1)
		def, err := g.DB().GetValue(ctx, "SELECT id FROM user_address WHERE user_id=? AND is_default=1", uid)
		t.AssertNil(err)
		t.Assert(def.Int64(), id2)
	})
}

// TestAddressGetForOrderRawPhone I1 回归: 内部取址返回**原始号码**（脱敏仅会员边界）。
func TestAddressGetForOrderRawPhone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		const phone = "13900003011"
		defer cleanupMember(ctx, t, phone)
		uid := seedMember(ctx, t, phone, "取号测试", 0)
		defer cleanupAddresses(ctx, t, uid)

		id := seedAddress(ctx, t, uid, "取号", false)
		got, err := AddressGetForOrder(ctx, uid, id)
		t.AssertNil(err)
		t.Assert(got.ReceiverPhone, "13800000000") // 原始（下单快照需真实号码）

		// 而会员可见的列表仍是脱敏
		list, err := AddressList(ctx, uid)
		t.AssertNil(err)
		t.Assert(list[0].ReceiverPhone, "138****0000")
	})
}
