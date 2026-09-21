// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AdminPermissionDao is the data access object for the table admin_permission.
type AdminPermissionDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  AdminPermissionColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// AdminPermissionColumns defines and stores column names for the table admin_permission.
type AdminPermissionColumns struct {
	Id        string // 权限ID
	ParentId  string // 父权限ID,0为根(菜单树)
	Name      string // 权限名称(如:商品管理/SPU上架)
	Code      string // 权限编码(如 product:spu:create;接口权限与API路径对应)
	Type      string // 类型:1菜单 2按钮/操作 3接口
	Path      string // 前端路由地址(菜单)或API路径(接口)
	Sort      string // 同级排序
	Status    string // 状态:1启用 0禁用
	Deleted   string // 软删除:0否 1是
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// adminPermissionColumns holds the columns for the table admin_permission.
var adminPermissionColumns = AdminPermissionColumns{
	Id:        "id",
	ParentId:  "parent_id",
	Name:      "name",
	Code:      "code",
	Type:      "type",
	Path:      "path",
	Sort:      "sort",
	Status:    "status",
	Deleted:   "deleted",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewAdminPermissionDao creates and returns a new DAO object for table data access.
func NewAdminPermissionDao(handlers ...gdb.ModelHandler) *AdminPermissionDao {
	return &AdminPermissionDao{
		group:    "default",
		table:    "admin_permission",
		columns:  adminPermissionColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AdminPermissionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AdminPermissionDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AdminPermissionDao) Columns() AdminPermissionColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AdminPermissionDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AdminPermissionDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AdminPermissionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
