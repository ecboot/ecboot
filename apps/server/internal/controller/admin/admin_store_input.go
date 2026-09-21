package admin

import "ecboot/internal/model"

// storeInputFromReq 门店请求 → 服务入参（创建/修改共用; 两渠道字段同构）。
func storeInputFromReq(name, province, city, district, addr, hours, phone string,
	lng, lat float64, pickup bool, status int) model.StoreInput {
	return model.StoreInput{
		Name: name, ProvinceCode: province, CityCode: city, DistrictCode: district,
		DetailAddress: addr, Longitude: lng, Latitude: lat,
		BusinessHours: hours, ContactPhone: phone, PickupEnabled: pickup, Status: status,
	}
}
