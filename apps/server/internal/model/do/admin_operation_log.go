// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminOperationLog is the golang structure of table admin_operation_log for DAO operations like Where/Data.
type AdminOperationLog struct {
	g.Meta        `orm:"table:admin_operation_log, do:true"`
	Id            any         // 日志ID
	AdminId       any         // 操作账号ID
	Username      any         // 操作人用户名(冗余快照,账号删除后仍可读)
	Module        any         // 业务模块(如 商品/订单/售后/权限)
	Operation     any         // 操作(如 上下架/改价/退款审核/角色授权)
	Method        any         // HTTP方法(GET/POST/PUT/DELETE)
	RequestUri    any         // 请求路径
	RequestParams any         // 请求参数(敏感字段脱敏后留档)
	ResultStatus  any         // 结果:1成功 0失败
	ErrorMsg      any         // 失败原因
	Ip            any         // 来源IP
	CostMs        any         // 耗时(毫秒)
	CreatedAt     *gtime.Time // 操作时间
}
