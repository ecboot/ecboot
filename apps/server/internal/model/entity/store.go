// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Store is the golang structure for table store.
type Store struct {
	Id            uint64      `json:"id"            orm:"id"             ` // 门店ID
	StoreNo       string      `json:"storeNo"       orm:"store_no"       ` // 门店编码(全局唯一)
	Name          string      `json:"name"          orm:"name"           ` // 门店名称
	ProvinceCode  string      `json:"provinceCode"  orm:"province_code"  ` // 省级行政区划代码(GB/T 2260)
	CityCode      string      `json:"cityCode"      orm:"city_code"      ` // 市级行政区划代码
	DistrictCode  string      `json:"districtCode"  orm:"district_code"  ` // 区县级行政区划代码
	DetailAddress string      `json:"detailAddress" orm:"detail_address" ` // 详细地址(街道门牌)
	Longitude     float64     `json:"longitude"     orm:"longitude"      ` // 经度(附近门店检索;GCJ-02坐标系)
	Latitude      float64     `json:"latitude"      orm:"latitude"       ` // 纬度(同上)
	BusinessHours string      `json:"businessHours" orm:"business_hours" ` // 营业时间(如 09:00-21:00)
	ContactPhone  string      `json:"contactPhone"  orm:"contact_phone"  ` // 门店电话
	PickupEnabled int         `json:"pickupEnabled" orm:"pickup_enabled" ` // 自提开关:0否 1是(下单自提点候选)
	Sort          int         `json:"sort"          orm:"sort"           ` // 排序,越小越靠前
	Status        int         `json:"status"        orm:"status"         ` // 状态:1营业 2歇业(歇业=展示但不可自提/核销)
	Deleted       int         `json:"deleted"       orm:"deleted"        ` // 软删除:0否 1是
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     ` // 创建时间
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     ` // 更新时间
}
