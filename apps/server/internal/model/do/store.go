// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Store is the golang structure of table store for DAO operations like Where/Data.
type Store struct {
	g.Meta        `orm:"table:store, do:true"`
	Id            any         // 门店ID
	StoreNo       any         // 门店编码(全局唯一)
	Name          any         // 门店名称
	ProvinceCode  any         // 省级行政区划代码(GB/T 2260)
	CityCode      any         // 市级行政区划代码
	DistrictCode  any         // 区县级行政区划代码
	DetailAddress any         // 详细地址(街道门牌)
	Longitude     any         // 经度(附近门店检索;GCJ-02坐标系)
	Latitude      any         // 纬度(同上)
	BusinessHours any         // 营业时间(如 09:00-21:00)
	ContactPhone  any         // 门店电话
	PickupEnabled any         // 自提开关:0否 1是(下单自提点候选)
	Sort          any         // 排序,越小越靠前
	Status        any         // 状态:1营业 2歇业(歇业=展示但不可自提/核销)
	Deleted       any         // 软删除:0否 1是
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
}
