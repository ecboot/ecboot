package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/model"
)

type (
	// 门店列表（管理）
	AdminStoreListReq struct {
		g.Meta  `path:"/stores" method:"GET" summary:"门店列表"`
		Status  int    `json:"status" dc:"状态筛选"`
		Keyword string `json:"keyword" dc:"名称/编码"`
		model.PageReq
	}
	AdminStoreItem struct {
		Id            string `json:"id"`
		StoreNo       string `json:"storeNo"`
		Name          string `json:"name"`
		ProvinceCode  string `json:"provinceCode" dc:"省区划码"`
		CityCode      string `json:"cityCode" dc:"市区划码"`
		DistrictCode  string `json:"districtCode" dc:"区县区划码"`
		DetailAddress string `json:"detailAddress"`
		BusinessHours string `json:"businessHours"`
		ContactPhone  string `json:"contactPhone"`
		PickupEnabled bool   `json:"pickupEnabled" dc:"自提"`
		Status        int    `json:"status" dc:"1营业 2歇业"`
	}
	AdminStoreListRes struct {
		model.PageRes
		List []AdminStoreItem `json:"list"`
	}

	// 权限: store:manage:create
	AdminStoreCreateReq struct {
		g.Meta        `path:"/stores" method:"POST" summary:"新增门店"`
		Name          string  `json:"name" v:"required" dc:"名称"`
		ProvinceCode  string  `json:"provinceCode" v:"required" dc:"省区划码"`
		CityCode      string  `json:"cityCode" v:"required" dc:"市区划码"`
		DistrictCode  string  `json:"districtCode" v:"required" dc:"区县区划码"`
		DetailAddress string  `json:"detailAddress" v:"required" dc:"详细地址"`
		Longitude     float64 `json:"longitude" dc:"经度"`
		Latitude      float64 `json:"latitude" dc:"纬度"`
		BusinessHours string  `json:"businessHours" dc:"营业时间"`
		ContactPhone  string  `json:"contactPhone" dc:"电话"`
		PickupEnabled bool    `json:"pickupEnabled" dc:"自提开关"`
	}
	AdminStoreCreateRes struct {
		Id string `json:"id"`
	}

	// 权限: store:manage:update
	AdminStoreUpdateReq struct {
		g.Meta        `path:"/stores/{id}" method:"PUT" summary:"修改门店"`
		Id            string  `json:"id" v:"required" dc:"门店ID"`
		Name          string  `json:"name" dc:"名称"`
		ProvinceCode  string  `json:"provinceCode" dc:"省区划码"`
		CityCode      string  `json:"cityCode" dc:"市区划码"`
		DistrictCode  string  `json:"districtCode" dc:"区县区划码"`
		DetailAddress string  `json:"detailAddress" dc:"详细地址"`
		Longitude     float64 `json:"longitude" dc:"经度"`
		Latitude      float64 `json:"latitude" dc:"纬度"`
		BusinessHours string  `json:"businessHours" dc:"营业时间"`
		ContactPhone  string  `json:"contactPhone" dc:"电话"`
		PickupEnabled bool    `json:"pickupEnabled" dc:"自提开关"`
		Status        int     `json:"status" dc:"1营业 2歇业"`
	}
	AdminStoreUpdateRes struct {
		Success bool `json:"success"`
	}

	// 权限: store:manage:delete
	AdminStoreDeleteReq struct {
		g.Meta `path:"/stores/{id}" method:"DELETE" summary:"删除门店(软删)"`
		Id     string `json:"id" v:"required" dc:"门店ID"`
	}
	AdminStoreDeleteRes struct {
		Success bool `json:"success"`
	}

	// 门店详情（管理）
	AdminStoreDetailReq struct {
		g.Meta `path:"/stores/{id}" method:"GET" summary:"门店详情"`
		Id     string `json:"id" v:"required" dc:"门店ID"`
	}
	AdminStoreDetailRes struct {
		Id            string  `json:"id"`
		StoreNo       string  `json:"storeNo"`
		Name          string  `json:"name"`
		ProvinceCode  string  `json:"provinceCode"`
		CityCode      string  `json:"cityCode"`
		DistrictCode  string  `json:"districtCode"`
		DetailAddress string  `json:"detailAddress"`
		Longitude     float64 `json:"longitude"`
		Latitude      float64 `json:"latitude"`
		BusinessHours string  `json:"businessHours"`
		ContactPhone  string  `json:"contactPhone"`
		PickupEnabled bool    `json:"pickupEnabled"`
		Status        int     `json:"status"`
	}
)
