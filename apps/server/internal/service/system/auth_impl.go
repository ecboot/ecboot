// auth_impl.go 后台账号认证实现（接口契约见 rbac.go IAdminAuthLogic）。
// 规则: 登录成功/失败均写 admin_login_log（只追加, research D8 不锁定;
// login_status 1成功 2密码错 3禁用或不存在）; 密码 bcrypt 不可逆（FR-007）;
// 会话走 security 库 admin 渠道（research D1）; 验证码携带即强校验（D7）。
package system

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"golang.org/x/crypto/bcrypt"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/library/captcha"
	"ecboot/internal/library/security"
	"ecboot/internal/model"
	"ecboot/internal/model/do"
	"ecboot/internal/model/entity"
)

// ActiveAdmin 账号存在且启用且未软删（FR-008: 中间件每请求校验, 禁用/软删后会话立即不可用）。
func ActiveAdmin(ctx context.Context, adminId int64) (bool, error) {
	cnt, err := dao.AdminUser.Ctx(ctx).
		Where(dao.AdminUser.Columns().Id, adminId).
		Where(dao.AdminUser.Columns().Status, 1).
		Where(dao.AdminUser.Columns().Deleted, 0).
		Count()
	if err != nil {
		return false, gerror.Wrap(err, "查询后台账号失败")
	}
	return cnt > 0, nil
}

// AdminLogin 后台登录（FR-001~003）: 校验账号态与密码, 全程审计, 成功发放 admin 会话。
func AdminLogin(ctx context.Context, username, password, captchaKey, captchaCode, ip, userAgent string) (*model.AdminLoginResult, error) {
	// 验证码: 携带即强校验（research D7）
	if captchaKey != "" {
		ok, err := captcha.Verify(ctx, captchaKey, captchaCode)
		if err != nil {
			return nil, gerror.Wrap(err, "校验验证码失败")
		}
		if !ok {
			return nil, errcode.New(errcode.CodeCaptchaError, "验证码错误或已过期")
		}
	}

	// 账号查询（未软删）
	rec, err := dao.AdminUser.Ctx(ctx).
		Where(dao.AdminUser.Columns().Username, username).
		Where(dao.AdminUser.Columns().Deleted, 0).
		One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询后台账号失败")
	}
	if rec.IsEmpty() {
		_ = writeLoginLog(ctx, username, nil, 3, ip, userAgent)
		return nil, errcode.New(errcode.CodeAdminBadCredential, "用户名或密码错误")
	}
	var admin entity.AdminUser
	if err = rec.Struct(&admin); err != nil {
		return nil, gerror.Wrap(err, "解析后台账号失败")
	}

	// 禁用账号拒绝（不比对密码, 不泄露状态细节于消息）
	if admin.Status != 1 {
		_ = writeLoginLog(ctx, username, int64(admin.Id), 3, ip, userAgent)
		return nil, errcode.New(errcode.CodeAdminDisabled, "账号已禁用")
	}

	// 密码校验（bcrypt）
	if bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)) != nil {
		_ = writeLoginLog(ctx, username, int64(admin.Id), 2, ip, userAgent)
		return nil, errcode.New(errcode.CodeAdminBadCredential, "用户名或密码错误")
	}

	// 成功: 审计 + 最后登录时间 + 发放会话
	_ = writeLoginLog(ctx, username, int64(admin.Id), 1, ip, userAgent)
	if _, err = dao.AdminUser.Ctx(ctx).
		Where(dao.AdminUser.Columns().Id, admin.Id).
		Data(do.AdminUser{LastLoginTime: gtime.Now()}).
		Update(); err != nil {
		return nil, gerror.Wrap(err, "更新最后登录时间失败")
	}
	token, refreshToken, err := security.NewSessionManagerFromConfig(ctx, "admin").Create(ctx, int64(admin.Id))
	if err != nil {
		return nil, gerror.Wrap(err, "创建会话失败")
	}
	return &model.AdminLoginResult{
		AdminId:      int64(admin.Id),
		Token:        token,
		RefreshToken: refreshToken,
		RealName:     admin.RealName,
		IsSuper:      admin.IsSuper == 1,
	}, nil
}

// writeLoginLog 登录审计只追加（失败不阻断登录响应, 审计缺失仅记日志）。
func writeLoginLog(ctx context.Context, username string, adminId any, status int, ip, userAgent string) error {
	entry := do.AdminLoginLog{
		Username:    username,
		LoginStatus: status,
		Ip:          ip,
		UserAgent:   userAgent,
	}
	if adminId != nil {
		entry.AdminId = adminId
	}
	if _, err := dao.AdminLoginLog.Ctx(ctx).Data(entry).Insert(); err != nil {
		g.Log().Errorf(ctx, "登录审计写入失败 username=%s status=%d err=%v", username, status, err)
		return err
	}
	return nil
}

// ChangePassword 改密（FR-007）: 旧密码校验 + bcrypt 新哈希。
func ChangePassword(ctx context.Context, adminId int64, oldPassword, newPassword string) error {
	rec, err := dao.AdminUser.Ctx(ctx).
		Where(dao.AdminUser.Columns().Id, adminId).
		Where(dao.AdminUser.Columns().Deleted, 0).
		One()
	if err != nil {
		return gerror.Wrap(err, "查询后台账号失败")
	}
	if rec.IsEmpty() {
		return errcode.New(errcode.CodeAdminBadCredential, "账号不存在")
	}
	var admin entity.AdminUser
	if err = rec.Struct(&admin); err != nil {
		return gerror.Wrap(err, "解析后台账号失败")
	}
	if bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(oldPassword)) != nil {
		return errcode.New(errcode.CodeOldPasswordWrong, "原密码错误")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return gerror.Wrap(err, "生成密码哈希失败")
	}
	_, err = dao.AdminUser.Ctx(ctx).
		Where(dao.AdminUser.Columns().Id, adminId).
		Data(do.AdminUser{PasswordHash: string(hash)}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "更新密码失败")
	}
	return nil
}

// Profile 个人信息（FR-006）: 账号 + 启用角色编码列表。
func Profile(ctx context.Context, adminId int64) (*model.AdminProfile, error) {
	rec, err := dao.AdminUser.Ctx(ctx).
		Where(dao.AdminUser.Columns().Id, adminId).
		Where(dao.AdminUser.Columns().Deleted, 0).
		One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询后台账号失败")
	}
	if rec.IsEmpty() {
		return nil, errcode.New(errcode.CodeAdminBadCredential, "账号不存在")
	}
	var admin entity.AdminUser
	if err = rec.Struct(&admin); err != nil {
		return nil, gerror.Wrap(err, "解析后台账号失败")
	}
	roles, err := adminRoleCodes(ctx, adminId)
	if err != nil {
		return nil, err
	}
	return &model.AdminProfile{
		Username: admin.Username,
		RealName: admin.RealName,
		Roles:    roles,
	}, nil
}

// adminRoleCodes 账号关联的启用角色编码（admin_user_role → admin_role）。
func adminRoleCodes(ctx context.Context, adminId int64) ([]string, error) {
	roleIds, err := dao.AdminUserRole.Ctx(ctx).
		Fields(dao.AdminUserRole.Columns().RoleId).
		Where(dao.AdminUserRole.Columns().AdminId, adminId).
		Array()
	if err != nil {
		return nil, gerror.Wrap(err, "查询账号角色失败")
	}
	if len(roleIds) == 0 {
		return []string{}, nil
	}
	codes, err := dao.AdminRole.Ctx(ctx).
		Fields(dao.AdminRole.Columns().Code).
		Where(dao.AdminRole.Columns().Deleted, 0).
		Where(dao.AdminRole.Columns().Status, 1).
		WhereIn(dao.AdminRole.Columns().Id, roleIds).
		Array()
	if err != nil {
		return nil, gerror.Wrap(err, "查询角色编码失败")
	}
	out := make([]string, 0, len(codes))
	for _, c := range codes {
		out = append(out, c.String())
	}
	return out, nil
}
