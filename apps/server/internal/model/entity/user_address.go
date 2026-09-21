// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAddress is the golang structure for table user_address.
type UserAddress struct {
	Id            uint64      `json:"id"            orm:"id"             ` // 地址ID
	UserId        uint64      `json:"userId"        orm:"user_id"        ` // 所属用户ID
	ReceiverName  string      `json:"receiverName"  orm:"receiver_name"  ` // 收货人姓名
	ReceiverPhone string      `json:"receiverPhone" orm:"receiver_phone" ` // 收货人手机号
	Province      string      `json:"province"      orm:"province"       ` // 省
	ProvinceCode  string      `json:"provinceCode"  orm:"province_code"  ` // 省级行政区划代码(GB/T 2260六位,运费规则匹配口径;空=历史数据待补)
	City          string      `json:"city"          orm:"city"           ` // 市
	CityCode      string      `json:"cityCode"      orm:"city_code"      ` // 市级行政区划代码
	District      string      `json:"district"      orm:"district"       ` // 区/县(直筒子市可为空)
	DistrictCode  string      `json:"districtCode"  orm:"district_code"  ` // 区县级行政区划代码
	DetailAddress string      `json:"detailAddress" orm:"detail_address" ` // 详细地址(街道门牌)
	IsDefault     int         `json:"isDefault"     orm:"is_default"     ` // 默认地址:0否 1是(每用户至多一个,应用层保证)
	Deleted       int         `json:"deleted"       orm:"deleted"        ` // 软删除:0否 1是
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     ` // 创建时间
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     ` // 更新时间
}
