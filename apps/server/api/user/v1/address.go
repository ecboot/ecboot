package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	AddressDetail struct {
		Id            string `json:"id" dc:"地址ID"`
		ReceiverName  string `json:"receiverName" dc:"收货人"`
		ReceiverPhone string `json:"receiverPhone" dc:"手机号(脱敏)"`
		Province      string `json:"province" dc:"省"`
		City          string `json:"city" dc:"市"`
		District      string `json:"district" dc:"区县"`
		DetailAddress string `json:"detailAddress" dc:"详细地址"`
		ProvinceCode  string `json:"provinceCode" dc:"省区划码"`
		CityCode      string `json:"cityCode" dc:"市区划码"`
		DistrictCode  string `json:"districtCode" dc:"区县区划码"`
		IsDefault     bool   `json:"isDefault" dc:"默认地址"`
	}

	AddressListReq struct {
		g.Meta `path:"/addresses" method:"GET" summary:"收货地址列表"`
	}
	AddressListRes struct {
		List []AddressDetail `json:"list"`
	}

	AddressCreateReq struct {
		g.Meta        `path:"/addresses" method:"POST" summary:"新增收货地址"`
		ReceiverName  string `json:"receiverName" v:"required" dc:"收货人"`
		ReceiverPhone string `json:"receiverPhone" v:"required" dc:"手机号"`
		Province      string `json:"province" v:"required" dc:"省"`
		City          string `json:"city" v:"required" dc:"市"`
		District      string `json:"district" dc:"区县"`
		DetailAddress string `json:"detailAddress" v:"required" dc:"详细地址"`
		ProvinceCode  string `json:"provinceCode" v:"required" dc:"省区划码"`
		CityCode      string `json:"cityCode" v:"required" dc:"市区划码"`
		DistrictCode  string `json:"districtCode" dc:"区县区划码"`
		IsDefault     bool   `json:"isDefault" dc:"设为默认"`
	}
	AddressCreateRes struct {
		Id string `json:"id" dc:"地址ID"`
	}

	AddressUpdateReq struct {
		g.Meta        `path:"/addresses/{id}" method:"PUT" summary:"修改收货地址"`
		Id            string `json:"id" v:"required" dc:"地址ID"`
		ReceiverName  string `json:"receiverName" dc:"收货人"`
		ReceiverPhone string `json:"receiverPhone" dc:"手机号"`
		Province      string `json:"province" dc:"省"`
		City          string `json:"city" dc:"市"`
		District      string `json:"district" dc:"区县"`
		DetailAddress string `json:"detailAddress" dc:"详细地址"`
		ProvinceCode  string `json:"provinceCode" dc:"省区划码"`
		CityCode      string `json:"cityCode" dc:"市区划码"`
		DistrictCode  string `json:"districtCode" dc:"区县区划码"`
		IsDefault     bool   `json:"isDefault" dc:"设为默认"`
	}
	AddressUpdateRes struct {
		Success bool `json:"success"`
	}

	AddressDeleteReq struct {
		g.Meta `path:"/addresses/{id}" method:"DELETE" summary:"删除收货地址"`
		Id     string `json:"id" v:"required" dc:"地址ID"`
	}
	AddressDeleteRes struct {
		Success bool `json:"success"`
	}

	AddressSetDefaultReq struct {
		g.Meta `path:"/addresses/{id}/default" method:"PUT" summary:"设默认地址"`
		Id     string `json:"id" v:"required" dc:"地址ID"`
	}
	AddressSetDefaultRes struct {
		Success bool `json:"success"`
	}
)
