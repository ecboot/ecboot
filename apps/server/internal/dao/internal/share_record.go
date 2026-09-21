// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ShareRecordDao is the data access object for the table share_record.
type ShareRecordDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ShareRecordColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ShareRecordColumns defines and stores column names for the table share_record.
type ShareRecordColumns struct {
	Id           string // 分享记录ID
	SharerUserId string // 分享人用户ID
	SpuId        string // 分享对象SPU(NULL=整店/页面分享)
	ShareChannel string // 渠道:1小程序卡片 2海报图片 3复制链接/口令 4朋友圈/社群
	SceneValue   string // 小程序场景值(scene参数,落地追踪)
	ShareCode    string // 使用的推广码
	CreatedAt    string // 分享时间(归因窗口起点)
}

// shareRecordColumns holds the columns for the table share_record.
var shareRecordColumns = ShareRecordColumns{
	Id:           "id",
	SharerUserId: "sharer_user_id",
	SpuId:        "spu_id",
	ShareChannel: "share_channel",
	SceneValue:   "scene_value",
	ShareCode:    "share_code",
	CreatedAt:    "created_at",
}

// NewShareRecordDao creates and returns a new DAO object for table data access.
func NewShareRecordDao(handlers ...gdb.ModelHandler) *ShareRecordDao {
	return &ShareRecordDao{
		group:    "default",
		table:    "share_record",
		columns:  shareRecordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ShareRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ShareRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ShareRecordDao) Columns() ShareRecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ShareRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ShareRecordDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ShareRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
