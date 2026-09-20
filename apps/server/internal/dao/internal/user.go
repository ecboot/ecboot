// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserDao is the data access object for the table user.
type UserDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserColumns defines and stores column names for the table user.
type UserColumns struct {
	Id              string // 用户ID
	Nickname        string // 昵称
	Avatar          string // 头像URL(OSS/CDN,只存链接)
	Phone           string // 手机号密文(加密存储,算法由应用层确定)
	PhoneHash       string // 手机号哈希(加盐SHA-256,登录精确检索与唯一性)
	PasswordHash    string // 密码哈希(bcrypt);微信首次登录未设密码时为NULL
	WxOpenid        string // 微信openid(小程序登录凭证,唯一)
	WxUnionid       string // 微信unionid(开放平台多应用打通预留)
	ShareCode       string // 个人推广码(唯一;NULL=未生成,唯一索引不参与)
	Gender          string // 性别:0未知 1男 2女
	GrowthValue     string // 成长值(只增不减,累计列)
	Level           string // 会员等级(user_level_rule.id,未启体系为NULL)
	LastLoginAt     string // 最后登录时间(休眠分级依据;登录链路一次命中,免login_log聚合)
	LastActiveAt    string // 最后活跃时间(登录/下单/浏览任一;休眠判定以此为准)
	Status          string // 账号状态:1正常 2禁用
	RegisterChannel string // 注册渠道:1微信小程序 2H5 3后台创建
	Deleted         string // 软删除:0否 1是(删除后手机号仍占用,复用走后台改绑)
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
}

// userColumns holds the columns for the table user.
var userColumns = UserColumns{
	Id:              "id",
	Nickname:        "nickname",
	Avatar:          "avatar",
	Phone:           "phone",
	PhoneHash:       "phone_hash",
	PasswordHash:    "password_hash",
	WxOpenid:        "wx_openid",
	WxUnionid:       "wx_unionid",
	ShareCode:       "share_code",
	Gender:          "gender",
	GrowthValue:     "growth_value",
	Level:           "level",
	LastLoginAt:     "last_login_at",
	LastActiveAt:    "last_active_at",
	Status:          "status",
	RegisterChannel: "register_channel",
	Deleted:         "deleted",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewUserDao creates and returns a new DAO object for table data access.
func NewUserDao(handlers ...gdb.ModelHandler) *UserDao {
	return &UserDao{
		group:    "default",
		table:    "user",
		columns:  userColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserDao) Columns() UserColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
