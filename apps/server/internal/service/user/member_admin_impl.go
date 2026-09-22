// member_admin_impl.go IMemberAdminLogic 实现（018 批次 12）。
// 防线: 禁用/启用条件更新判行数（同值幂等——批次 10 I3 教训）; 改绑 uk_phone_hash 1062→业务码。
package user

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/library/money"
	"ecboot/internal/model"
)

// MemberAdminLogicImpl IMemberAdminLogic 实现。
type MemberAdminLogicImpl struct{}

func NewMemberAdminLogic() *MemberAdminLogicImpl { return &MemberAdminLogicImpl{} }

var _ IMemberAdminLogic = (*MemberAdminLogicImpl)(nil)

// memberStatusOf user.status 直通——C1（评审修复）: 真实枚举是 **1正常 2禁用**（000001 列注释
// + auth.go 消费点 status==2 拒绝登录实证）, 原注释"0=禁用"系未查列注释的事实性错误,
// 导致禁用写 0 而登录查 2 → 治理语义为空。api 枚举与库一致, 无需转换。

func memberRowOf(r gdb.Record, phone string) AdminMemberItem {
	cols := dao.User.Columns()
	return AdminMemberItem{
		UserId:      r[cols.Id].Int64(),
		Nickname:    r[cols.Nickname].String(),
		Phone:       maskPhone(phone),
		Gender:      r[cols.Gender].Int(),
		Level:       r[cols.Level].Int64(),
		GrowthValue: r[cols.GrowthValue].Int(),
		Status:      r[cols.Status].Int(),
		CreatedAt:   r[cols.CreatedAt].String(),
	}
}

// memberPhone 解密手机号（复用 auth.go 的 phoneCipher() 单例; 解密失败回退空——不影响列表主流程）。
func memberPhone(ctx context.Context, encrypted string) string {
	if encrypted == "" {
		return ""
	}
	plain, err := phoneCipher().Decrypt(encrypted)
	if err != nil {
		g.Log().Warningf(ctx, "[会员管理] 手机号解密失败: %v", err)
		return ""
	}
	return plain
}

// AdminList 会员列表。
func (i *MemberAdminLogicImpl) AdminList(ctx context.Context, phone string, status int, keyword string, page model.PageReq) (*model.PageResult[AdminMemberItem], error) {
	page = page.Normalized()
	cols := dao.User.Columns()
	m := dao.User.Ctx(ctx).Where(cols.Deleted, 0)
	if phone != "" {
		m = m.Where(cols.PhoneHash, phoneCipher().Hash(phone)) // 精确检索走哈希（加密列不可 LIKE）
	}
	switch status {
	case 2:
		m = m.Where(cols.Status, 2) // C1: 库枚举 2=禁用（直通, 无 0 值转换）
	case 1:
		m = m.Where(cols.Status, 1)
	}
	if keyword != "" {
		m = m.WhereLike(cols.Nickname, "%"+strings.ReplaceAll(keyword, "%", "\\%")+"%")
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计会员失败")
	}
	recs, err := m.OrderDesc(cols.Id).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询会员失败")
	}
	list := make([]AdminMemberItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, memberRowOf(r, memberPhone(ctx, r[cols.Phone].String())))
	}
	return &model.PageResult[AdminMemberItem]{List: list, Total: int64(total)}, nil
}

// AdminDetail 会员详情（资产概要 + 订单统计）。
func (i *MemberAdminLogicImpl) AdminDetail(ctx context.Context, userId int64) (*AdminMemberItem, error) {
	// M5（评审修复）: 补 deleted=0（与列表/memberMustExist 口径一致, 软删会员详情不可达）
	r, err := dao.User.Ctx(ctx).
		Where(dao.User.Columns().Id, userId).
		Where(dao.User.Columns().Deleted, 0).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询会员失败")
	}
	if r.IsEmpty() {
		return nil, errcode.New(errcode.CodeNotFound, "会员不存在")
	}
	out := memberRowOf(r, memberPhone(ctx, r[dao.User.Columns().Phone].String()))

	// 资产概要（查不到/无行 → 零值占位）
	points, _ := g.DB().Model("point_account").Ctx(ctx).Where("user_id", userId).Value("balance")
	balance, _ := g.DB().Model("user_account").Ctx(ctx).Where("user_id", userId).Value("balance")
	coupons, _ := g.DB().Model("user_coupon").Ctx(ctx).
		Where("user_id", userId).Where("status", 1).Count()
	out.Assets = map[string]string{
		"points":  money.ToYuanString(points.Int64()),
		"balance": balance.String(),
		"coupons": money.ToYuanString(int64(coupons)),
	}
	// 订单统计
	stat := func(where string) int64 {
		v, e := g.DB().Model("trade_order").Ctx(ctx).
			Where("user_id", userId).Where(where).Count()
		if e != nil {
			return 0
		}
		return int64(v)
	}
	out.OrderStats = map[string]int64{
		"total":    stat("status >= 10"),
		"finished": stat("status = 40"),
		"refunded": stat("status = 90"),
	}
	return &out, nil
}

// Disable 禁用/启用（C1 修复: 库枚举 **2=禁用**（000001 列注释）, auth.go 以 status==2 拒绝登录
// ——治理语义即"禁用后无法登录"; 同值幂等; 审计 reason 由操作日志面承载）。
func (i *MemberAdminLogicImpl) Disable(ctx context.Context, userId int64, disable bool, reason string) error {
	target := 1 // 启用
	if disable {
		target = 2 // 禁用（登录侧 status==2 拒绝, 语义闭环）
	}
	if err := memberMustExist(ctx, userId); err != nil {
		return err
	}
	_, err := dao.User.Ctx(ctx).
		Where(dao.User.Columns().Id, userId).
		Data(g.Map{dao.User.Columns().Status: target}).Update()
	// 同值幂等: 不判行数（批次 02/10 修改语义先例）
	return gerror.Wrap(err, "变更会员状态失败")
}

// RebindPhone 改绑手机号（加密列三件套: phone 密文 + phone_hash 检索键; 唯一键兜底并发）。
func (i *MemberAdminLogicImpl) RebindPhone(ctx context.Context, userId int64, newPhone string) error {
	if len(newPhone) != 11 || !strings.HasPrefix(newPhone, "1") {
		return errcode.New(errcode.CodeInvalidParam, "手机号格式非法")
	}
	if err := memberMustExist(ctx, userId); err != nil {
		return err
	}
	pc := phoneCipher()
	enc, err := pc.Encrypt(newPhone)
	if err != nil {
		return gerror.Wrap(err, "加密手机号失败")
	}
	_, err = dao.User.Ctx(ctx).
		Where(dao.User.Columns().Id, userId).
		Data(g.Map{
			dao.User.Columns().Phone:     enc,
			dao.User.Columns().PhoneHash: pc.Hash(newPhone),
		}).Update()
	if err != nil {
		if isDupKeyUser(err) { // uk_phone_hash: 新号已被占用
			return errcode.New(errcode.CodeInvalidParam, "新手机号已被占用")
		}
		return gerror.Wrap(err, "改绑手机号失败")
	}
	return nil
}

// memberMustExist 会员存在且未删。
func memberMustExist(ctx context.Context, userId int64) error {
	n, err := dao.User.Ctx(ctx).
		Where(dao.User.Columns().Id, userId).
		Where(dao.User.Columns().Deleted, 0).Count()
	if err != nil {
		return gerror.Wrap(err, "查询会员失败")
	}
	if n == 0 {
		return errcode.New(errcode.CodeNotFound, "会员不存在")
	}
	return nil
}
