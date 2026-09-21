// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminOperationLog is the golang structure for table admin_operation_log.
type AdminOperationLog struct {
	Id            uint64      `json:"id"            orm:"id"             ` // 日志ID
	AdminId       uint64      `json:"adminId"       orm:"admin_id"       ` // 操作账号ID
	Username      string      `json:"username"      orm:"username"       ` // 操作人用户名(冗余快照,账号删除后仍可读)
	Module        string      `json:"module"        orm:"module"         ` // 业务模块(如 商品/订单/售后/权限)
	Operation     string      `json:"operation"     orm:"operation"      ` // 操作(如 上下架/改价/退款审核/角色授权)
	Method        string      `json:"method"        orm:"method"         ` // HTTP方法(GET/POST/PUT/DELETE)
	RequestUri    string      `json:"requestUri"    orm:"request_uri"    ` // 请求路径
	RequestParams string      `json:"requestParams" orm:"request_params" ` // 请求参数(敏感字段脱敏后留档)
	ResultStatus  int         `json:"resultStatus"  orm:"result_status"  ` // 结果:1成功 0失败
	ErrorMsg      string      `json:"errorMsg"      orm:"error_msg"      ` // 失败原因
	Ip            string      `json:"ip"            orm:"ip"             ` // 来源IP
	CostMs        uint        `json:"costMs"        orm:"cost_ms"        ` // 耗时(毫秒)
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     ` // 操作时间
}
