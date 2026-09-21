package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/service/user"
)

// AddressCreate 新增收货地址（isDefault=1 时同事务清其他默认）
func (c *ControllerV1) AddressCreate(ctx context.Context, req *v1.AddressCreateReq) (res *v1.AddressCreateRes, err error) {
	id, err := user.AddressCreate(ctx, middleware.CtxUserIdFrom(ctx), addrInputFromReq(
		req.ReceiverName, req.ReceiverPhone, req.Province, req.City, req.District,
		req.DetailAddress, req.ProvinceCode, req.CityCode, req.DistrictCode, req.IsDefault))
	if err != nil {
		return nil, err
	}
	return &v1.AddressCreateRes{Id: fmtID(id)}, nil
}
