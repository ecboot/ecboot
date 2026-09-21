// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserAddress is the golang structure of table user_address for DAO operations like Where/Data.
type UserAddress struct {
	g.Meta        `orm:"table:user_address, do:true"`
	Id            any         // 地址ID
	UserId        any         // 所属用户ID
	ReceiverName  any         // 收货人姓名
	ReceiverPhone any         // 收货人手机号
	Province      any         // 省
	ProvinceCode  any         // 省级行政区划代码(GB/T 2260六位,运费规则匹配口径;空=历史数据待补)
	City          any         // 市
	CityCode      any         // 市级行政区划代码
	District      any         // 区/县(直筒子市可为空)
	DistrictCode  any         // 区县级行政区划代码
	DetailAddress any         // 详细地址(街道门牌)
	IsDefault     any         // 默认地址:0否 1是(每用户至多一个,应用层保证)
	Deleted       any         // 软删除:0否 1是
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
}
