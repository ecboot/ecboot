// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserRelation is the golang structure for table user_relation.
type UserRelation struct {
	Id          uint64      `json:"id"          orm:"id"           ` // 关系ID
	UserId      uint64      `json:"userId"      orm:"user_id"      ` // 用户ID(一人至多一条关系)
	InviterId   uint64      `json:"inviterId"   orm:"inviter_id"   ` // 直接上级(邀请人);二级=上级的上级,两次单列查询解析;无祖父列/路径列(ADR-0003)
	BindChannel int         `json:"bindChannel" orm:"bind_channel" ` // 绑定方式:1分享链接 2邀请码
	BindTime    *gtime.Time `json:"bindTime"    orm:"bind_time"    ` // 绑定时间
	Locked      int         `json:"locked"      orm:"locked"       ` // 锁定:0保护期内(可换绑) 1已锁定(绑定窗口期满,防抢人纠纷)
	LockTime    *gtime.Time `json:"lockTime"    orm:"lock_time"    ` // 锁定时间
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   ` // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   ` // 更新时间
}
