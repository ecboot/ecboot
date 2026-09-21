// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BargainRecordDao is the data access object for the table bargain_record.
type BargainRecordDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  BargainRecordColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// BargainRecordColumns defines and stores column names for the table bargain_record.
type BargainRecordColumns struct {
	Id           string // 砍价单ID
	BargainNo    string // 砍价单号(全局唯一)
	ItemId       string // 砍价场次商品ID
	UserId       string // 发起人用户ID
	CurrentPrice string // 当前价(逐刀递减,>=floor_price由应用+条件更新保证)
	CutCount     string // 已砍刀数
	Status       string // 状态:1砍价中 2到底价待下单 3已下单 4超时失败 5已取消
	ExpireTime   string // 砍价截止时间(超时扫描)
	SuccessTime  string // 到底价时间
	OrderNo      string // 成交订单号(下单后回填,防重复成交)
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
}

// bargainRecordColumns holds the columns for the table bargain_record.
var bargainRecordColumns = BargainRecordColumns{
	Id:           "id",
	BargainNo:    "bargain_no",
	ItemId:       "item_id",
	UserId:       "user_id",
	CurrentPrice: "current_price",
	CutCount:     "cut_count",
	Status:       "status",
	ExpireTime:   "expire_time",
	SuccessTime:  "success_time",
	OrderNo:      "order_no",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewBargainRecordDao creates and returns a new DAO object for table data access.
func NewBargainRecordDao(handlers ...gdb.ModelHandler) *BargainRecordDao {
	return &BargainRecordDao{
		group:    "default",
		table:    "bargain_record",
		columns:  bargainRecordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BargainRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BargainRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BargainRecordDao) Columns() BargainRecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BargainRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BargainRecordDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *BargainRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
