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

// couponValidDesc 组装有效期描述（列表/详情共用口径）。
func couponValidDesc(validType int, startAt, endAt gtime.Time, days int) string {
	switch validType {
	case 1:
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
		ValidDesc: couponValidDesc(r[cols.ValidType].Int(),
			*r[cols.ValidStartAt].GTime(), *r[cols.ValidEndAt].GTime(), r[cols.ValidDays].Int()),
		Status: r[cols.Status].Int(),
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

// AdminUpdate 修改券模板（部分字段; Status 显式启停——do 结构 omitempty 语义下
// 零值会被吞, 故 Status>0 或显式置停都以**字段白名单**写入, 沿批次 02 全量覆盖护栏先例）。
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
	// 启停: api 契约 Status 1启用 0停发; 0 为合法目标值, 不能被 omitempty 吞——
	// 用"<0 表示未传"不可行（JSON 无三态）, 故约定: 显式传 0 或 1 都写入, 未传(缺省 0)时
	// 与"停发"语义重合——管理端前端的 Update 必传 status（契约 dc 注明）。此处照单写入。
	data[cols.Status] = in.Status
	if len(data) == 0 {
		return nil
	}
	res, err := dao.Coupon.Ctx(ctx).Where(cols.Id, id).Where(cols.Deleted, 0).Data(data).Update()
	if err != nil {
		return gerror.Wrap(err, "修改券模板失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		// 幂等: 值未变化也视为成功（与批次 02 修改语义一致, 不报错）
		return nil
	}
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
