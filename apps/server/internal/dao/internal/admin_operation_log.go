// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AdminOperationLogDao is the data access object for the table admin_operation_log.
type AdminOperationLogDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  AdminOperationLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// AdminOperationLogColumns defines and stores column names for the table admin_operation_log.
type AdminOperationLogColumns struct {
	Id            string // 日志ID
	AdminId       string // 操作账号ID
	Username      string // 操作人用户名(冗余快照,账号删除后仍可读)
	Module        string // 业务模块(如 商品/订单/售后/权限)
	Operation     string // 操作(如 上下架/改价/退款审核/角色授权)
	Method        string // HTTP方法(GET/POST/PUT/DELETE)
	RequestUri    string // 请求路径
	RequestParams string // 请求参数(敏感字段脱敏后留档)
	ResultStatus  string // 结果:1成功 0失败
	ErrorMsg      string // 失败原因
	Ip            string // 来源IP
	CostMs        string // 耗时(毫秒)
	CreatedAt     string // 操作时间
}

// adminOperationLogColumns holds the columns for the table admin_operation_log.
var adminOperationLogColumns = AdminOperationLogColumns{
	Id:            "id",
	AdminId:       "admin_id",
	Username:      "username",
	Module:        "module",
	Operation:     "operation",
	Method:        "method",
	RequestUri:    "request_uri",
	RequestParams: "request_params",
	ResultStatus:  "result_status",
	ErrorMsg:      "error_msg",
	Ip:            "ip",
	CostMs:        "cost_ms",
	CreatedAt:     "created_at",
}

// NewAdminOperationLogDao creates and returns a new DAO object for table data access.
func NewAdminOperationLogDao(handlers ...gdb.ModelHandler) *AdminOperationLogDao {
	return &AdminOperationLogDao{
		group:    "default",
		table:    "admin_operation_log",
		columns:  adminOperationLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AdminOperationLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AdminOperationLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AdminOperationLogDao) Columns() AdminOperationLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AdminOperationLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AdminOperationLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AdminOperationLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
