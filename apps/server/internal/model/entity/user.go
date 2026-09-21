// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure for table user.
type User struct {
	Id              uint64      `json:"id"              orm:"id"               ` // 用户ID
	Nickname        string      `json:"nickname"        orm:"nickname"         ` // 昵称
	Avatar          string      `json:"avatar"          orm:"avatar"           ` // 头像URL(OSS/CDN,只存链接)
	Phone           string      `json:"phone"           orm:"phone"            ` // 手机号密文(加密存储,算法由应用层确定)
	PhoneHash       string      `json:"phoneHash"       orm:"phone_hash"       ` // 手机号哈希(加盐SHA-256,登录精确检索与唯一性)
	PasswordHash    string      `json:"passwordHash"    orm:"password_hash"    ` // 密码哈希(bcrypt);微信首次登录未设密码时为NULL
	WxOpenid        string      `json:"wxOpenid"        orm:"wx_openid"        ` // 微信openid(小程序登录凭证,唯一)
	WxUnionid       string      `json:"wxUnionid"       orm:"wx_unionid"       ` // 微信unionid(开放平台多应用打通预留)
	ShareCode       string      `json:"shareCode"       orm:"share_code"       ` // 个人推广码(唯一;NULL=未生成,唯一索引不参与)
	Gender          int         `json:"gender"          orm:"gender"           ` // 性别:0未知 1男 2女
	GrowthValue     uint        `json:"growthValue"     orm:"growth_value"     ` // 成长值(只增不减,累计列)
	Level           uint64      `json:"level"           orm:"level"            ` // 会员等级(user_level_rule.id,未启体系为NULL)
	LastLoginAt     *gtime.Time `json:"lastLoginAt"     orm:"last_login_at"    ` // 最后登录时间(休眠分级依据;登录链路一次命中,免login_log聚合)
	LastActiveAt    *gtime.Time `json:"lastActiveAt"    orm:"last_active_at"   ` // 最后活跃时间(登录/下单/浏览任一;休眠判定以此为准)
	Status          int         `json:"status"          orm:"status"           ` // 账号状态:1正常 2禁用
	RegisterChannel int         `json:"registerChannel" orm:"register_channel" ` // 注册渠道:1微信小程序 2H5 3后台创建
	Deleted         int         `json:"deleted"         orm:"deleted"          ` // 软删除:0否 1是(删除后手机号仍占用,复用走后台改绑)
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       ` // 创建时间
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       ` // 更新时间
}
