// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BargainHelperDao is the data access object for the table bargain_helper.
type BargainHelperDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  BargainHelperColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// BargainHelperColumns defines and stores column names for the table bargain_helper.
type BargainHelperColumns struct {
	Id           string // 帮砍记录ID
	RecordId     string // 砍价单ID
	HelperUserId string // 帮砍人用户ID
	CutAmount    string // 本刀砍掉金额(与record条件更新同事务)
	CreatedAt    string // 帮砍时间
}

// bargainHelperColumns holds the columns for the table bargain_helper.
var bargainHelperColumns = BargainHelperColumns{
	Id:           "id",
	RecordId:     "record_id",
	HelperUserId: "helper_user_id",
	CutAmount:    "cut_amount",
	CreatedAt:    "created_at",
}

// NewBargainHelperDao creates and returns a new DAO object for table data access.
func NewBargainHelperDao(handlers ...gdb.ModelHandler) *BargainHelperDao {
	return &BargainHelperDao{
		group:    "default",
		table:    "bargain_helper",
		columns:  bargainHelperColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BargainHelperDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BargainHelperDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BargainHelperDao) Columns() BargainHelperColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BargainHelperDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BargainHelperDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *BargainHelperDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
