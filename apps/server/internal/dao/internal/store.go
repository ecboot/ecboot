// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// StoreDao is the data access object for the table store.
type StoreDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  StoreColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// StoreColumns defines and stores column names for the table store.
type StoreColumns struct {
	Id            string // 门店ID
	StoreNo       string // 门店编码(全局唯一)
	Name          string // 门店名称
	ProvinceCode  string // 省级行政区划代码(GB/T 2260)
	CityCode      string // 市级行政区划代码
	DistrictCode  string // 区县级行政区划代码
	DetailAddress string // 详细地址(街道门牌)
	Longitude     string // 经度(附近门店检索;GCJ-02坐标系)
	Latitude      string // 纬度(同上)
	BusinessHours string // 营业时间(如 09:00-21:00)
	ContactPhone  string // 门店电话
	PickupEnabled string // 自提开关:0否 1是(下单自提点候选)
	Sort          string // 排序,越小越靠前
	Status        string // 状态:1营业 2歇业(歇业=展示但不可自提/核销)
	Deleted       string // 软删除:0否 1是
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
}

// storeColumns holds the columns for the table store.
var storeColumns = StoreColumns{
	Id:            "id",
	StoreNo:       "store_no",
	Name:          "name",
	ProvinceCode:  "province_code",
	CityCode:      "city_code",
	DistrictCode:  "district_code",
	DetailAddress: "detail_address",
	Longitude:     "longitude",
	Latitude:      "latitude",
	BusinessHours: "business_hours",
	ContactPhone:  "contact_phone",
	PickupEnabled: "pickup_enabled",
	Sort:          "sort",
	Status:        "status",
	Deleted:       "deleted",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewStoreDao creates and returns a new DAO object for table data access.
func NewStoreDao(handlers ...gdb.ModelHandler) *StoreDao {
	return &StoreDao{
		group:    "default",
		table:    "store",
		columns:  storeColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *StoreDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *StoreDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *StoreDao) Columns() StoreColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *StoreDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *StoreDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *StoreDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
