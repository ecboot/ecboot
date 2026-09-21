// profile_impl.go 会员资料与登录记录（011-member-center; 接口契约见 profile.go IProfileLogic）。
// 约定: 手机号脱敏在服务层完成（明文不出参/不出日志——个保法）;
// 等级按 user_level_rule 门槛匹配成长值（取最大满足门槛者）;
// 资料修改只更新非零字段（防零值覆盖, 同批次 04 教训; 局限: 无法将性别改回"未知"）。
package user

import (
	"context"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/model"
	"ecboot/internal/model/do"
	"ecboot/internal/model/entity"
)

// maskPhone 手机号脱敏（前 3 + **** + 后 4）。
func maskPhone(p string) string {
	if p == "" {
		return ""
	}
	if len(p) < 7 {
		return "****" // 评审 Minor8: 短值兜底全掩码（原样返回等于不脱敏）
	}
	return p[:3] + "****" + p[len(p)-4:]
}

// maskIP IP 脱敏（IPv4 保留前两段; 其他形态原样返回）。
func maskIP(ip string) string {
	if ip == "" {
		return ""
	}
	parts := strings.Split(ip, ".")
	if len(parts) == 4 {
		return parts[0] + "." + parts[1] + ".*.*"
	}
	// 评审 Minor8: IPv6/非四段值兜底全掩码（原样返回等于不脱敏）
	return "****"
}

// matchLevel 按成长值匹配等级（取最大满足门槛者; 无规则或低于最低门槛 → 0/空）。
func matchLevel(ctx context.Context, growth int) (int64, string) {
	recs, err := dao.UserLevelRule.Ctx(ctx).
		Where(dao.UserLevelRule.Columns().Deleted, 0).
		Where(dao.UserLevelRule.Columns().Status, 1).
		OrderDesc(dao.UserLevelRule.Columns().GrowthThreshold).
		All()
	if err != nil {
		return 0, ""
	}
	for _, r := range recs {
		if growth >= r["growth_threshold"].Int() {
			return r["id"].Int64(), r["name"].String()
		}
	}
	return 0, ""
}

// ProfileDetail 会员资料（FR-001）: 脱敏手机号 + 等级名 + 成长值。
func ProfileDetail(ctx context.Context, userId int64) (*model.ProfileDetail, error) {
	cols := dao.User.Columns()
	rec, err := dao.User.Ctx(ctx).
		Where(cols.Id, userId).
		Where(cols.Deleted, 0).
		One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询会员失败")
	}
	if rec.IsEmpty() {
		return nil, errcode.New(errcode.CodeNotFound, "会员不存在")
	}
	var u entity.User
	if err = rec.Struct(&u); err != nil {
		return nil, gerror.Wrap(err, "解析会员失败")
	}

	// 手机号: 密文解密后脱敏（解密失败不阻断资料查询）
	masked := ""
	if u.Phone != "" {
		if p, derr := PhoneCipherFor(ctx).Decrypt(u.Phone); derr == nil {
			masked = maskPhone(p)
		}
	}
	levelId, levelName := matchLevel(ctx, int(u.GrowthValue))

	return &model.ProfileDetail{
		UserId:      strconv.FormatInt(userId, 10),
		Nickname:    u.Nickname,
		Avatar:      u.Avatar,
		Gender:      u.Gender,
		PhoneMasked: masked,
		HasPassword: u.PasswordHash != "",
		Level:       levelId,
		LevelName:   levelName,
		GrowthValue: int(u.GrowthValue),
	}, nil
}

// ProfileUpdate 资料修改（FR-002）: 只更新非零字段（未传=不变）。
func ProfileUpdate(ctx context.Context, userId int64, in model.ProfileUpdateInput) error {
	cols := dao.User.Columns()
	data := do.User{}
	if in.Nickname != "" {
		data.Nickname = in.Nickname
	}
	if in.Avatar != "" {
		data.Avatar = in.Avatar
	}
	if in.Gender != 0 {
		data.Gender = in.Gender
	}
	res, err := dao.User.Ctx(ctx).
		Where(cols.Id, userId).
		Where(cols.Deleted, 0).
		Data(data).
		Update()
	if err != nil {
		return gerror.Wrap(err, "修改资料失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		cnt, cerr := dao.User.Ctx(ctx).Where(cols.Id, userId).Where(cols.Deleted, 0).Count()
		if cerr != nil {
			return gerror.Wrap(cerr, "查询会员失败")
		}
		if cnt == 0 {
			return errcode.New(errcode.CodeNotFound, "会员不存在")
		}
	}
	return nil
}

// LoginLogs 近 30 天登录记录（FR-003）: 时间倒序 + 分页; IP 脱敏。
func LoginLogs(ctx context.Context, userId int64, page model.PageReq) (*model.PageResult[model.UserLoginLogItem], error) {
	page = page.Normalized()
	cols := dao.UserLoginLog.Columns()
	m := dao.UserLoginLog.Ctx(ctx).
		Where(cols.UserId, userId).
		Where("created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)")
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计登录记录失败")
	}
	recs, err := m.OrderDesc(cols.CreatedAt).OrderDesc(cols.Id).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询登录记录失败")
	}
	list := make([]model.UserLoginLogItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, model.UserLoginLogItem{
			Channel:   r["login_channel"].Int(),
			Status:    r["login_status"].Int(),
			Ip:        maskIP(r["ip"].String()),
			UserAgent: r["user_agent"].String(),
			CreatedAt: r["created_at"].String(),
		})
	}
	return &model.PageResult[model.UserLoginLogItem]{List: list, Total: int64(total)}, nil
}

// GrowthAdd 成长值累加（内部方法; 只增不减）。
func GrowthAdd(ctx context.Context, userId int64, value int) error {
	if value <= 0 {
		return nil
	}
	cols := dao.User.Columns()
	_, err := dao.User.Ctx(ctx).
		Where(cols.Id, userId).
		Where(cols.Deleted, 0).
		Data(do.User{GrowthValue: gdb.Raw("growth_value + (" + strconv.Itoa(value) + ")")}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "累加成长值失败")
	}
	return LevelRecalc(ctx, userId)
}

// LevelRecalc 依门槛重算等级（内部方法; 成长值变更后调用）。
func LevelRecalc(ctx context.Context, userId int64) error {
	cols := dao.User.Columns()
	rec, err := dao.User.Ctx(ctx).Fields(cols.GrowthValue).
		Where(cols.Id, userId).Where(cols.Deleted, 0).One()
	if err != nil {
		return gerror.Wrap(err, "查询成长值失败")
	}
	if rec.IsEmpty() {
		return nil
	}
	levelId, _ := matchLevel(ctx, rec[cols.GrowthValue].Int())
	var data do.User
	if levelId > 0 {
		data.Level = levelId
	}
	if _, err = dao.User.Ctx(ctx).Where(cols.Id, userId).Data(data).Update(); err != nil {
		return gerror.Wrap(err, "重算等级失败")
	}
	return nil
}
