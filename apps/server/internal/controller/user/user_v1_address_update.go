package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// AddressUpdate 修改收货地址（归属校验）
func (c *ControllerV1) AddressUpdate(ctx context.Context, req *v1.AddressUpdateReq) (res *v1.AddressUpdateRes, err error) {
	id, err := parseID(req.Id)
	if err != nil {
		return nil, err
	}
	if err = user.AddressUpdate(ctx, middleware.CtxUserIdFrom(ctx), id, addrInputFromReq(
		req.ReceiverName, req.ReceiverPhone, req.Province, req.City, req.District,
		req.DetailAddress, req.ProvinceCode, req.CityCode, req.DistrictCode, req.IsDefault)); err != nil {
		return nil, err
	}
	return &v1.AddressUpdateRes{Success: true}, nil
}
