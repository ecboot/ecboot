// logistics_impl.go 物流公司字典（接口契约见 misc.go ILogisticsLogic）。
// 语义: 编码唯一且**不可改**（订单 deliver_company 存编码值, 改码割裂历史）;
// 停用保留（列表/详情可见, 仅不再供新发货选择）; 软删后列表与详情均不可见。
package shop

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/model"
	"ecboot/internal/model/do"
	"ecboot/internal/model/entity"
)

func logisticsFromEntity(e entity.LogisticsCompany) model.LogisticsCompany {
	return model.LogisticsCompany{
		Id:           int64(e.Id),
		Code:         e.Code,
		Name:         e.Name,
		TrackingRule: e.TrackingRule,
		Status:       e.Status,
	}
}

// LogisticsList 物流公司列表（FR-001）: status>0 筛选; 0=不筛选（列表含停用——停用保留语义）。
func LogisticsList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.LogisticsCompany], error) {
	page = page.Normalized()
	cols := dao.LogisticsCompany.Columns()
	m := dao.LogisticsCompany.Ctx(ctx).Where(cols.Deleted, 0)
	if status > 0 {
		m = m.Where(cols.Status, status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计物流公司失败")
	}
	recs, err := m.Fields(cols.Id, cols.Code, cols.Name, cols.TrackingRule, cols.Status).
		Order("sort ASC, id ASC").
		Page(page.Page, page.PageSize).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询物流公司失败")
	}
	list := make([]model.LogisticsCompany, 0, len(recs))
	for _, r := range recs {
		var e entity.LogisticsCompany
		if err = r.Struct(&e); err != nil {
			return nil, gerror.Wrap(err, "解析物流公司失败")
		}
		list = append(list, logisticsFromEntity(e))
	}
	return &model.PageResult[model.LogisticsCompany]{List: list, Total: int64(total)}, nil
}

// LogisticsCreate 创建物流公司（FR-002）: 编码唯一; 状态默认启用（欲建即停用请创建后修改）。
func LogisticsCreate(ctx context.Context, in model.LogisticsCompanyInput) (int64, error) {
	if in.Code == "" {
		return 0, errcode.New(errcode.CodeInvalidParam, "编码必填")
	}
	if in.Name == "" {
		return 0, errcode.New(errcode.CodeInvalidParam, "名称必填")
	}
	cols := dao.LogisticsCompany.Columns()
	cnt, err := dao.LogisticsCompany.Ctx(ctx).Where(cols.Code, in.Code).Count()
	if err != nil {
		return 0, gerror.Wrap(err, "查询物流公司失败")
	}
	if cnt > 0 {
		return 0, errcode.New(errcode.CodeLogisticsCodeTaken, "物流公司编码已存在")
	}
	status := in.Status
	if status == 0 {
		status = 1
	}
	res, err := dao.LogisticsCompany.Ctx(ctx).Data(do.LogisticsCompany{
		Code:         in.Code,
		Name:         in.Name,
		TrackingRule: in.TrackingRule,
		Status:       status,
	}).Insert()
	if err != nil {
		return 0, gerror.Wrap(err, "创建物流公司失败")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, gerror.Wrap(err, "读取新记录ID失败")
	}
	return id, nil
}

// LogisticsUpdate 修改物流公司（FR-003）: 名称/校验规则/状态;
// **编码不在待更新字段中**——订单存编码值, 改码割裂历史（入参携带亦被忽略）。
func LogisticsUpdate(ctx context.Context, id int64, in model.LogisticsCompanyInput) error {
	if in.Name == "" {
		return errcode.New(errcode.CodeInvalidParam, "名称必填")
	}
	// 状态白名单（009 评审 I1: 域外值会落库致公司从发货选择静默消失）
	if in.Status != 0 && in.Status != 1 {
		return errcode.New(errcode.CodeInvalidParam, "状态须为1启用或0停用")
	}
	// 状态全量覆盖语义（与门店一致）: 0=停用 1=启用 严格照传——
	// 后台修改表单总是携带完整档案（含状态单选）, 不做零值归一以免"无法停用"。
	cols := dao.LogisticsCompany.Columns()
	res, err := dao.LogisticsCompany.Ctx(ctx).
		Where(cols.Id, id).
		Where(cols.Deleted, 0).
		Data(do.LogisticsCompany{Name: in.Name, TrackingRule: in.TrackingRule, Status: in.Status}).
		Fields(cols.Name, cols.TrackingRule, cols.Status).
		Update()
	if err != nil {
		return gerror.Wrap(err, "修改物流公司失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		cnt, cerr := dao.LogisticsCompany.Ctx(ctx).
			Where(cols.Id, id).Where(cols.Deleted, 0).Count()
		if cerr != nil {
			return gerror.Wrap(cerr, "查询物流公司失败")
		}
		if cnt == 0 {
			return errcode.New(errcode.CodeNotFound, "物流公司不存在")
		}
	}
	return nil
}

// LogisticsDelete 软删物流公司（FR-004）: 列表与详情均不可见。
func LogisticsDelete(ctx context.Context, id int64) error {
	cols := dao.LogisticsCompany.Columns()
	res, err := dao.LogisticsCompany.Ctx(ctx).
		Where(cols.Id, id).
		Where(cols.Deleted, 0).
		Data(do.LogisticsCompany{Deleted: 1}).
		Fields(cols.Deleted).
		Update()
	if err != nil {
		return gerror.Wrap(err, "删除物流公司失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeNotFound, "物流公司不存在")
	}
	return nil
}

// LogisticsDetail 物流公司详情（FR-005）: 停用可见; 不存在/已删返 10006。
func LogisticsDetail(ctx context.Context, id int64) (*model.LogisticsCompany, error) {
	cols := dao.LogisticsCompany.Columns()
	rec, err := dao.LogisticsCompany.Ctx(ctx).
		Fields(cols.Id, cols.Code, cols.Name, cols.TrackingRule, cols.Status).
		Where(cols.Id, id).
		Where(cols.Deleted, 0).
		One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询物流公司失败")
	}
	if rec.IsEmpty() {
		return nil, errcode.New(errcode.CodeNotFound, "物流公司不存在")
	}
	var e entity.LogisticsCompany
	if err = rec.Struct(&e); err != nil {
		return nil, gerror.Wrap(err, "解析物流公司失败")
	}
	item := logisticsFromEntity(e)
	return &item, nil
}
