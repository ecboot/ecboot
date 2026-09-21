package user

import "ecboot/internal/model"

// addrInputFromReq 地址请求 → 服务入参（创建/修改共用）。
func addrInputFromReq(receiverName, receiverPhone, province, city, district, detailAddress,
	provinceCode, cityCode, districtCode string, isDefault bool) model.AddressInput {
	return model.AddressInput{
		ReceiverName: receiverName, ReceiverPhone: receiverPhone,
		Province: province, City: city, District: district, DetailAddress: detailAddress,
		ProvinceCode: provinceCode, CityCode: cityCode, DistrictCode: districtCode,
		IsDefault: isDefault,
	}
}
