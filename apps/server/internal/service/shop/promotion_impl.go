// promotion_impl.go ICouponLogic 实现（016-marketing-admin 批次 10）。
// 规则见 promotion.go 接口注释; 契约映射见 specs/016-marketing-admin/contracts/。
// 防线（plan D1）: 唯一键/条件更新判行数; 1062 转业务码; 金额走 money 库。
package shop

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/library/money"
	"ecboot/internal/model"
	"ecboot/internal/model/do"
)

// CouponLogicImpl ICouponLogic 实现。
type CouponLogicImpl struct{}

func NewCouponLogic() *CouponLogicImpl { return &CouponLogicImpl{} }

// 接口满足性编译断言（批次 10 评审 M7: 此前 PublicList 无实现者、CouponLogicImpl 不满足
// 自家接口也无人发现——编译期钉住, 不再靠文字承诺）
var (
	_ ICouponLogic   = (*CouponLogicImpl)(nil)
	_ IActivityLogic = (*ActivityLogicImpl)(nil)
)

// couponValidDesc 组装有效期描述（列表/详情共用口径）。
// C1（评审修复）: valid_start_at/valid_end_at 为 NULL DEFAULT NULL（000009）——
// validType=2（领取后 N 天）的券这两列恒 NULL, 原实现对 GTime() 直接解引用 → 列表/详情 panic。
func couponValidDesc(validType int, startAt, endAt *gtime.Time, days int) string {
	switch validType {
	case 1:
		if startAt == nil || endAt == nil {
			return ""
		}
		return startAt.Format("Y-m-d") + " ~ " + endAt.Format("Y-m-d")
	case 2:
		return fmt.Sprintf("领取后 %d 天内有效", days)
	}
	return ""
}

// couponInputCheck 创建侧校验（api v 标签之外的业务条件, FR-1/契约校验矩阵）。
func couponInputCheck(in model.CouponInput) error {
	if in.Type != 1 && in.Type != 2 {
		return errcode.New(errcode.CodeInvalidParam, "券类型非法")
	}
	if in.Discount != "" {
		d, err := money.FromYuanString(in.Discount)
		if err != nil || d <= 0 {
			return errcode.New(errcode.CodeInvalidParam, "抵扣金额非法")
		}
	}
	if in.Type == 1 {
		if in.Threshold == "" {
			return errcode.New(errcode.CodeInvalidParam, "满减券必须设置使用门槛")
		}
		if _, err := money.FromYuanString(in.Threshold); err != nil {
			return errcode.New(errcode.CodeInvalidParam, "门槛金额非法")
		}
	}
	switch in.ValidType {
	case 1:
		if in.ValidStartAt == "" || in.ValidEndAt == "" {
			return errcode.New(errcode.CodeInvalidParam, "固定区间有效期必须填写起止时间")
		}
		s, e := gtime.New(in.ValidStartAt), gtime.New(in.ValidEndAt)
		if s == nil || e == nil {
			return errcode.New(errcode.CodeInvalidParam, "时间格式非法")
		}
		if !e.After(s) {
			return errcode.New(errcode.CodeInvalidParam, "有效期开始必须早于结束")
		}
	case 2:
		if in.ValidDays < 1 {
			return errcode.New(errcode.CodeInvalidParam, "领取后有效天数必须大于 0")
		}
	default:
		return errcode.New(errcode.CodeInvalidParam, "有效期方式非法")
	}
	return nil
}

// AdminList 券模板列表（status 0=全部; 含已领数; 分页）。
func (i *CouponLogicImpl) AdminList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.CouponTemplate], error) {
	page = page.Normalized()
	cols := dao.Coupon.Columns()
	m := dao.Coupon.Ctx(ctx).Where(cols.Deleted, 0)
	if status > 0 {
		m = m.Where(cols.Status, status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计券模板失败")
	}
	recs, err := m.OrderDesc(cols.Id).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询券模板失败")
	}
	list := make([]model.CouponTemplate, 0, len(recs))
	for _, r := range recs {
		list = append(list, couponTemplateOf(r))
	}
	return &model.PageResult[model.CouponTemplate]{List: list, Total: int64(total)}, nil
}

// couponTemplateOf 行→DTO（列表/详情同源, D8 口径统一）。
func couponTemplateOf(r gdb.Record) model.CouponTemplate {
	cols := dao.Coupon.Columns()
	// NULL 日期守卫（C1）: validType=2 的券 valid_start_at/valid_end_at 恒 NULL
	var vs, ve *gtime.Time
	if v := r[cols.ValidStartAt]; !v.IsNil() {
		vs = v.GTime()
	}
	if v := r[cols.ValidEndAt]; !v.IsNil() {
		ve = v.GTime()
	}
	return model.CouponTemplate{
		Id:         r[cols.Id].Int64(),
		Name:       r[cols.Name].String(),
		Type:       r[cols.Type].Int(),
		Threshold:  r[cols.ThresholdAmount].String(),
		Discount:   r[cols.DiscountAmount].String(),
		TotalCount: r[cols.TotalCount].Int(),
		Received:   r[cols.ReceivedCount].Int(),
		PerLimit:   r[cols.PerLimit].Int(),
		ValidType:  r[cols.ValidType].Int(), // D3-②
		ValidDesc:  couponValidDesc(r[cols.ValidType].Int(), vs, ve, r[cols.ValidDays].Int()),
		Status:     r[cols.Status].Int(),
	}
}

// AdminDetail 券模板详情（D3-①）。
func (i *CouponLogicImpl) AdminDetail(ctx context.Context, id int64) (*model.CouponTemplate, error) {
	cols := dao.Coupon.Columns()
	r, err := dao.Coupon.Ctx(ctx).Where(cols.Id, id).Where(cols.Deleted, 0).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询券模板失败")
	}
	if r.IsEmpty() {
		return nil, errcode.New(errcode.CodeNotFound, "券模板不存在")
	}
	out := couponTemplateOf(r)
	return &out, nil
}

// AdminCreate 创建券模板（业务条件校验见 couponInputCheck）。
func (i *CouponLogicImpl) AdminCreate(ctx context.Context, in model.CouponInput) (int64, error) {
	if err := couponInputCheck(in); err != nil {
		return 0, err
	}
	data := do.Coupon{
		Name:         in.Name,
		Type:         in.Type,
		TotalCount:   in.TotalCount,
		PerLimit:     in.PerLimit,
		ValidType:    in.ValidType,
		Status:       1,
		DiscountRate: nil,
	}
	if in.Type == 1 {
		data.ThresholdAmount = in.Threshold
	} else {
		data.ThresholdAmount = "0.00" // 无门槛
	}
	data.DiscountAmount = in.Discount
	switch in.ValidType {
	case 1:
		data.ValidStartAt = gtime.New(in.ValidStartAt)
		data.ValidEndAt = gtime.New(in.ValidEndAt)
	case 2:
		data.ValidDays = in.ValidDays
	}
	res, err := dao.Coupon.Ctx(ctx).Data(data).Insert()
	if err != nil {
		return 0, gerror.Wrap(err, "创建券模板失败")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, gerror.Wrap(err, "读取券模板ID失败")
	}
	return id, nil
}

// AdminUpdate 修改券模板（部分字段; I6 三态: Status *int nil=不修改启停, 0=停发, 1=启用）。
func (i *CouponLogicImpl) AdminUpdate(ctx context.Context, id int64, in model.CouponInput) error {
	cols := dao.Coupon.Columns()
	cur, err := dao.Coupon.Ctx(ctx).Where(cols.Id, id).Where(cols.Deleted, 0).One()
	if err != nil {
		return gerror.Wrap(err, "查询券模板失败")
	}
	if cur.IsEmpty() {
		return errcode.New(errcode.CodeNotFound, "券模板不存在")
	}
	data := g.Map{}
	if in.Name != "" {
		data[cols.Name] = in.Name
	}
	if in.Threshold != "" {
		if _, e := money.FromYuanString(in.Threshold); e != nil {
			return errcode.New(errcode.CodeInvalidParam, "门槛金额非法")
		}
		data[cols.ThresholdAmount] = in.Threshold
	}
	if in.Discount != "" {
		if d, e := money.FromYuanString(in.Discount); e != nil || d <= 0 {
			return errcode.New(errcode.CodeInvalidParam, "抵扣金额非法")
		}
		data[cols.DiscountAmount] = in.Discount
	}
	if in.TotalCount > 0 {
		data[cols.TotalCount] = in.TotalCount
	}
	if in.PerLimit > 0 {
		data[cols.PerLimit] = in.PerLimit
	}
	// I6（评审修复）: status 三态——*int nil=不修改启停, 0=停发, 1=启用。
	// 原实现无条件写 int status, "只改名不传 status"的部分更新会**静默停发**（C 端可领券消失
	// 且返回 success）; "JSON 无三态"的辩解不成立——*int 即三态。
	if in.Status != nil {
		data[cols.Status] = *in.Status
	}
	if len(data) == 0 {
		return nil
	}
	_, err = dao.Coupon.Ctx(ctx).Where(cols.Id, id).Where(cols.Deleted, 0).Data(data).Update()
	if err != nil {
		return gerror.Wrap(err, "修改券模板失败")
	}
	// 同值无变化（affected=0）视为幂等成功（批次 02 修改语义先例）
	return nil
}

// AdminDelete 软删（既有 user_coupon 不受影响——列注释"停用不影响已领取的券", 软删同口径）。
func (i *CouponLogicImpl) AdminDelete(ctx context.Context, id int64) error {
	cols := dao.Coupon.Columns()
	res, err := dao.Coupon.Ctx(ctx).
		Where(cols.Id, id).Where(cols.Deleted, 0).
		Data(g.Map{cols.Deleted: 1, cols.Status: 0}). // 删除即停发（C 端口径双保险）
		Update()
	if err != nil {
		return gerror.Wrap(err, "删除券模板失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeNotFound, "券模板不存在")
	}
	return nil
}

// PublicList 公开可领列表（模板级口径, 与首页 publicCouponBriefs 同源: 启用未删且未领完;
// M7 补实现——此前该方法是死契约。注: "canReceive 由会员态补充"属消费端逻辑（个人限领判定,
// user 域 /user/coupons/available 口径）, 本方法只回答"有哪些券正在发放"）。
func (i *CouponLogicImpl) PublicList(ctx context.Context, page model.PageReq) (*model.PageResult[model.CouponTemplate], error) {
	page = page.Normalized()
	cols := dao.Coupon.Columns()
	base := func() *gdb.Model {
		return dao.Coupon.Ctx(ctx).
			Where(cols.Status, 1).Where(cols.Deleted, 0).
			Where("total_count = 0 OR received_count < total_count") // total_count=0 视为不限量
	}
	total, err := base().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计可领券失败")
	}
	recs, err := base().OrderDesc(cols.Id).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询可领券失败")
	}
	list := make([]model.CouponTemplate, 0, len(recs))
	for _, r := range recs {
		list = append(list, couponTemplateOf(r))
	}
	return &model.PageResult[model.CouponTemplate]{List: list, Total: int64(total)}, nil
}

// AdminRecords 券领取/使用记录分页（FR-1 下钻）。
func (i *CouponLogicImpl) AdminRecords(ctx context.Context, couponId int64, page model.PageReq) (*model.PageResult[model.CouponRecordItem], error) {
	page = page.Normalized()
	cols := dao.Coupon.Columns()
	exists, err := dao.Coupon.Ctx(ctx).Where(cols.Id, couponId).Count()
	if err != nil {
		return nil, gerror.Wrap(err, "查询券模板失败")
	}
	if exists == 0 {
		return nil, errcode.New(errcode.CodeNotFound, "券模板不存在")
	}
	ucols := dao.UserCoupon.Columns()
	m := dao.UserCoupon.Ctx(ctx).Where(ucols.CouponId, couponId)
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计券记录失败")
	}
	recs, err := m.OrderDesc(ucols.Id).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询券记录失败")
	}
	list := make([]model.CouponRecordItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, model.CouponRecordItem{
			UserCouponId: r[ucols.Id].Int64(),
			UserId:       r[ucols.UserId].Int64(),
			Status:       r[ucols.Status].Int(),
			OrderNo:      r[ucols.OrderNo].String(),
			CreatedAt:    r[ucols.CreatedAt].String(),
		})
	}
	return &model.PageResult[model.CouponRecordItem]{List: list, Total: int64(total)}, nil
}
