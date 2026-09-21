// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserRelation is the golang structure of table user_relation for DAO operations like Where/Data.
type UserRelation struct {
	g.Meta      `orm:"table:user_relation, do:true"`
	Id          any         // 关系ID
	UserId      any         // 用户ID(一人至多一条关系)
	InviterId   any         // 直接上级(邀请人);二级=上级的上级,两次单列查询解析;无祖父列/路径列(ADR-0003)
	BindChannel any         // 绑定方式:1分享链接 2邀请码
	BindTime    *gtime.Time // 绑定时间
	Locked      any         // 锁定:0保护期内(可换绑) 1已锁定(绑定窗口期满,防抢人纠纷)
	LockTime    *gtime.Time // 锁定时间
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
}
