// address_impl.go 收货地址（011-member-center; 接口契约见 address.go IAddressLogic）。
// 约定: 归属校验——他人地址一律按"不存在"处理（10006, 不泄露存在性）;
// 默认唯一由"同事务清该用户其他默认 + 置本条"保证（创建/修改/设默认三路径统一走此语义）。
package user

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/model"
	"ecboot/internal/model/do"
)

// addressFromRecordRaw 行 → DTO（**电话不脱敏**——内部链路用: 下单快照需真实号码）。
func addressFromRecordRaw(r gdb.Record) model.AddressItem {
	return model.AddressItem{
		Id:            r["id"].Int64(),
		ReceiverName:  r["receiver_name"].String(),
		ReceiverPhone: r["receiver_phone"].String(),
		Province:      r["province"].String(),
		City:          r["city"].String(),
		District:      r["district"].String(),
		DetailAddress: r["detail_address"].String(),
		ProvinceCode:  r["province_code"].String(),
		CityCode:      r["city_code"].String(),
		DistrictCode:  r["district_code"].String(),
		IsDefault:     r["is_default"].Int() == 1,
	}
}

// addressFromRecord 行 → DTO（**会员可见边界**: 电话脱敏）。
func addressFromRecord(r gdb.Record) model.AddressItem {
	it := addressFromRecordRaw(r)
	it.ReceiverPhone = maskPhone(it.ReceiverPhone)
	return it
}

// addressUpdateDo 更新入参 → do（010 评审判例延伸: api 各字段可选, **仅非空才更新**;
// is_default 为 false 时**不动作**——"取消默认"经给其他地址设默认实现, 避免静默丢默认）。
func addressUpdateDo(in model.AddressInput) do.UserAddress {
	d := do.UserAddress{}
	if in.ReceiverName != "" {
		d.ReceiverName = in.ReceiverName
	}
	if in.ReceiverPhone != "" {
		d.ReceiverPhone = in.ReceiverPhone
	}
	if in.Province != "" {
		d.Province = in.Province
	}
	if in.City != "" {
		d.City = in.City
	}
	if in.District != "" {
		d.District = in.District
	}
	if in.DetailAddress != "" {
		d.DetailAddress = in.DetailAddress
	}
	if in.ProvinceCode != "" {
		d.ProvinceCode = in.ProvinceCode
	}
	if in.CityCode != "" {
		d.CityCode = in.CityCode
	}
	if in.DistrictCode != "" {
		d.DistrictCode = in.DistrictCode
	}
	if in.IsDefault {
		d.IsDefault = 1
	}
	return d
}

// addressInputDo 入参 → do（创建用: 全量写入, 含区划码列）。
func addressInputDo(in model.AddressInput, isDefault bool) do.UserAddress {
	d := 0
	if isDefault {
		d = 1
	}
	return do.UserAddress{
		ReceiverName:  in.ReceiverName,
		ReceiverPhone: in.ReceiverPhone,
		Province:      in.Province,
		City:          in.City,
		District:      in.District,
		DetailAddress: in.DetailAddress,
		ProvinceCode:  in.ProvinceCode,
		CityCode:      in.CityCode,
		DistrictCode:  in.DistrictCode,
		IsDefault:     d,
	}
}

// clearOtherDefault 清该用户其他默认标记（同事务内调用; 排除 exceptId）。
func clearOtherDefault(ctx context.Context, tx gdb.TX, userId, exceptId int64) error {
	cols := dao.UserAddress.Columns()
	_, err := dao.UserAddress.Ctx(ctx).
		Where(cols.UserId, userId).
		WhereNot(cols.Id, exceptId).
		Where(cols.IsDefault, 1).
		Data(do.UserAddress{IsDefault: 0}).
		Fields(cols.IsDefault).
		Update()
	return err
}

// ensureOwned 归属校验: 返回本人地址行, 他人/不存在 → 10006。
func ensureOwned(ctx context.Context, userId, addressId int64) (gdb.Record, error) {
	cols := dao.UserAddress.Columns()
	rec, err := dao.UserAddress.Ctx(ctx).
		Where(cols.Id, addressId).
		Where(cols.UserId, userId).
		Where(cols.Deleted, 0).
		One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询地址失败")
	}
	if rec.IsEmpty() {
		return nil, errcode.New(errcode.CodeNotFound, "地址不存在")
	}
	return rec, nil
}

// AddressList 地址列表（默认在前, 其余按 id 倒序）; 电话脱敏。
func AddressList(ctx context.Context, userId int64) ([]model.AddressItem, error) {
	cols := dao.UserAddress.Columns()
	recs, err := dao.UserAddress.Ctx(ctx).
		Where(cols.UserId, userId).
		Where(cols.Deleted, 0).
		OrderDesc(cols.IsDefault).
		OrderDesc(cols.Id).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询地址列表失败")
	}
	list := make([]model.AddressItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, addressFromRecord(r))
	}
	return list, nil
}

// AddressCreate 新增地址（is_default=1 时同事务清其他默认）。
func AddressCreate(ctx context.Context, userId int64, in model.AddressInput) (int64, error) {
	var id int64
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		d := addressInputDo(in, in.IsDefault)
		d.UserId = userId // 归属列（遗漏将触发 1364 无默认值）
		res, e := dao.UserAddress.Ctx(ctx).Data(d).Insert()
		if e != nil {
			return gerror.Wrap(e, "新增地址失败")
		}
		nid, e := res.LastInsertId()
		if e != nil {
			return gerror.Wrap(e, "读取新地址ID失败")
		}
		id = nid
		if in.IsDefault {
			return clearOtherDefault(ctx, tx, userId, id)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

// AddressUpdate 修改地址（归属校验; 传 is_default=1 时清其他）。
func AddressUpdate(ctx context.Context, userId, addressId int64, in model.AddressInput) error {
	if _, err := ensureOwned(ctx, userId, addressId); err != nil {
		return err
	}
	cols := dao.UserAddress.Columns()
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := dao.UserAddress.Ctx(ctx).
			Where(cols.Id, addressId).
			Where(cols.UserId, userId).
			Data(addressUpdateDo(in)).
			Update(); e != nil {
			return gerror.Wrap(e, "修改地址失败")
		}
		if in.IsDefault {
			return clearOtherDefault(ctx, tx, userId, addressId)
		}
		return nil
	})
	return err
}

// AddressDelete 删除地址（软删; 归属校验）。
func AddressDelete(ctx context.Context, userId, addressId int64) error {
	if _, err := ensureOwned(ctx, userId, addressId); err != nil {
		return err
	}
	cols := dao.UserAddress.Columns()
	_, err := dao.UserAddress.Ctx(ctx).
		Where(cols.Id, addressId).
		Where(cols.UserId, userId).
		Data(do.UserAddress{Deleted: 1}).
		Fields(cols.Deleted).
		Update()
	if err != nil {
		return gerror.Wrap(err, "删除地址失败")
	}
	return nil
}

// AddressSetDefault 设默认（归属校验 + 同事务清其他, 默认唯一）。
func AddressSetDefault(ctx context.Context, userId, addressId int64) error {
	if _, err := ensureOwned(ctx, userId, addressId); err != nil {
		return err
	}
	cols := dao.UserAddress.Columns()
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := dao.UserAddress.Ctx(ctx).
			Where(cols.Id, addressId).
			Where(cols.UserId, userId).
			Data(do.UserAddress{IsDefault: 1}).
			Fields(cols.IsDefault).
			Update(); e != nil {
			return gerror.Wrap(e, "设默认失败")
		}
		return clearOtherDefault(ctx, tx, userId, addressId)
	})
}

// AddressGetForOrder 下单取址（内部方法; 归属校验 + 含区划码——运费计算依据）。
func AddressGetForOrder(ctx context.Context, userId, addressId int64) (*model.AddressItem, error) {
	rec, err := ensureOwned(ctx, userId, addressId)
	if err != nil {
		return nil, err
	}
	// 内部链路（下单/运费）返回**原始号码**——脱敏仅发生在会员可见边界
	it := addressFromRecordRaw(rec)
	return &it, nil
}
