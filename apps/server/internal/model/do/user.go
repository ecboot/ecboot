// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure of table user for DAO operations like Where/Data.
type User struct {
	g.Meta          `orm:"table:user, do:true"`
	Id              any         // 用户ID
	Nickname        any         // 昵称
	Avatar          any         // 头像URL(OSS/CDN,只存链接)
	Phone           any         // 手机号密文(加密存储,算法由应用层确定)
	PhoneHash       any         // 手机号哈希(加盐SHA-256,登录精确检索与唯一性)
	PasswordHash    any         // 密码哈希(bcrypt);微信首次登录未设密码时为NULL
	WxOpenid        any         // 微信openid(小程序登录凭证,唯一)
	WxUnionid       any         // 微信unionid(开放平台多应用打通预留)
	ShareCode       any         // 个人推广码(唯一;NULL=未生成,唯一索引不参与)
	Gender          any         // 性别:0未知 1男 2女
	GrowthValue     any         // 成长值(只增不减,累计列)
	Level           any         // 会员等级(user_level_rule.id,未启体系为NULL)
	LastLoginAt     *gtime.Time // 最后登录时间(休眠分级依据;登录链路一次命中,免login_log聚合)
	LastActiveAt    *gtime.Time // 最后活跃时间(登录/下单/浏览任一;休眠判定以此为准)
	Status          any         // 账号状态:1正常 2禁用
	RegisterChannel any         // 注册渠道:1微信小程序 2H5 3后台创建
	Deleted         any         // 软删除:0否 1是(删除后手机号仍占用,复用走后台改绑)
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
}
