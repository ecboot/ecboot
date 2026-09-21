package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// AddressList 收货地址列表（电话脱敏）
func (c *ControllerV1) AddressList(ctx context.Context, req *v1.AddressListReq) (res *v1.AddressListRes, err error) {
	list, err := user.AddressList(ctx, middleware.CtxUserIdFrom(ctx))
	if err != nil {
		return nil, err
	}
	res = &v1.AddressListRes{List: make([]v1.AddressDetail, 0, len(list))}
	for _, it := range list {
		res.List = append(res.List, v1.AddressDetail{
			Id: fmtID(it.Id), ReceiverName: it.ReceiverName, ReceiverPhone: it.ReceiverPhone,
			Province: it.Province, City: it.City, District: it.District,
			DetailAddress: it.DetailAddress, ProvinceCode: it.ProvinceCode,
			CityCode: it.CityCode, DistrictCode: it.DistrictCode, IsDefault: it.IsDefault,
		})
	}
	return res, nil
}
