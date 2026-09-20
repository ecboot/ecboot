package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 门店列表（区县筛选或经纬度附近检索, 距离排序）—— contracts/common-api.md
	StoreListReq struct {
		g.Meta       `path:"/stores" method:"GET" summary:"门店列表"`
		DistrictCode string  `json:"districtCode" dc:"区县区划码(与经纬度二选一)"`
		Longitude    float64 `json:"longitude" dc:"经度(附近检索)"`
		Latitude     float64 `json:"latitude" dc:"纬度"`
		RadiusKm     int     `json:"radiusKm" dc:"半径公里,默认10" d:"10"`
		PageReq
	}
	StoreItem struct {
		Id            string `json:"id" dc:"门店ID"`
		StoreNo       string `json:"storeNo" dc:"门店编码"`
		Name          string `json:"name" dc:"门店名称"`
		DetailAddress string `json:"detailAddress" dc:"详细地址"`
		DistanceM     int64  `json:"distanceM" dc:"距离米(附近检索时返回)"`
		BusinessHours string `json:"businessHours" dc:"营业时间"`
		PickupEnabled bool   `json:"pickupEnabled" dc:"是否支持自提"`
		Status        int    `json:"status" dc:"状态:1营业 2歇业"`
	}
	StoreListRes struct {
		PageRes
		List []StoreItem `json:"list"`
	}

	// 门店详情
	StoreDetailReq struct {
		g.Meta `path:"/stores/{id}" method:"GET" summary:"门店详情"`
		Id     string `json:"id" v:"required" dc:"门店ID"`
	}
	StoreDetailRes struct {
		Id            string  `json:"id"`
		StoreNo       string  `json:"storeNo" dc:"门店编码"`
		Name          string  `json:"name" dc:"门店名称"`
		ProvinceCode  string  `json:"provinceCode" dc:"省区划码"`
		CityCode      string  `json:"cityCode" dc:"市区划码"`
		DistrictCode  string  `json:"districtCode" dc:"区县区划码"`
		DetailAddress string  `json:"detailAddress" dc:"详细地址"`
		Longitude     float64 `json:"longitude" dc:"经度"`
		Latitude      float64 `json:"latitude" dc:"纬度"`
		BusinessHours string  `json:"businessHours" dc:"营业时间"`
		ContactPhone  string  `json:"contactPhone" dc:"门店电话"`
		PickupEnabled bool    `json:"pickupEnabled" dc:"是否支持自提"`
		Status        int     `json:"status" dc:"状态"`
	}
)
