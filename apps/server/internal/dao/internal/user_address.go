// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserAddressDao is the data access object for the table user_address.
type UserAddressDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserAddressColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserAddressColumns defines and stores column names for the table user_address.
type UserAddressColumns struct {
	Id            string // 地址ID
	UserId        string // 所属用户ID
	ReceiverName  string // 收货人姓名
	ReceiverPhone string // 收货人手机号
	Province      string // 省
	ProvinceCode  string // 省级行政区划代码(GB/T 2260六位,运费规则匹配口径;空=历史数据待补)
	City          string // 市
	CityCode      string // 市级行政区划代码
	District      string // 区/县(直筒子市可为空)
	DistrictCode  string // 区县级行政区划代码
	DetailAddress string // 详细地址(街道门牌)
	IsDefault     string // 默认地址:0否 1是(每用户至多一个,应用层保证)
	Deleted       string // 软删除:0否 1是
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
}

// userAddressColumns holds the columns for the table user_address.
var userAddressColumns = UserAddressColumns{
	Id:            "id",
	UserId:        "user_id",
	ReceiverName:  "receiver_name",
	ReceiverPhone: "receiver_phone",
	Province:      "province",
	ProvinceCode:  "province_code",
	City:          "city",
	CityCode:      "city_code",
	District:      "district",
	DistrictCode:  "district_code",
	DetailAddress: "detail_address",
	IsDefault:     "is_default",
	Deleted:       "deleted",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewUserAddressDao creates and returns a new DAO object for table data access.
func NewUserAddressDao(handlers ...gdb.ModelHandler) *UserAddressDao {
	return &UserAddressDao{
		group:    "default",
		table:    "user_address",
		columns:  userAddressColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserAddressDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserAddressDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserAddressDao) Columns() UserAddressColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserAddressDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserAddressDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserAddressDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
