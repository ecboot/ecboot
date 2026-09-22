// activity_impl.go IActivityLogic 实现（016-marketing-admin 批次 10）。
// 满减/拼团/秒杀/砍价/助力五类的管理面（契约映射见 specs/016-marketing-admin/contracts/）。
// 防线（plan D1/D2）: 唯一键 1062 转业务码; 全量替换走事务; 秒杀已售保护事务内复核;
// 金额一律 money.FromYuanString 校验（拒绝 >2 位小数/非正数）。
package shop

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/library/money"
	"ecboot/internal/model"
)

// ActivityLogicImpl IActivityLogic 实现。
type ActivityLogicImpl struct{}

func NewActivityLogic() *ActivityLogicImpl { return &ActivityLogicImpl{} }

// ---------- 通用校验与装配 ----------

// activityTimeCheck 活动时间窗。M8（评审修复）: 部分更新语义——非空侧才校验格式,
// 双非空才比较顺序（原实现传其一时另一侧按空解析报"时间格式非法", 误导且 Update 契约声明可选）。
func activityTimeCheck(start, end string) error {
	var s, e *gtime.Time
	if start != "" {
		if s = gtime.New(start); s == nil || s.Timestamp() == 0 {
			return errcode.New(errcode.CodeInvalidParam, "开始时间格式非法")
		}
	}
	if end != "" {
		if e = gtime.New(end); e == nil || e.Timestamp() == 0 {
			return errcode.New(errcode.CodeInvalidParam, "结束时间格式非法")
		}
	}
	if s != nil && e != nil && !e.After(s) {
		return errcode.New(errcode.CodeInvalidParam, "活动开始必须早于结束")
	}
	return nil
}

// moneyYuanCheck 金额校验（合法十进制、>0、≤2 位小数）→ 返回分。
func moneyYuanCheck(v string, field string) (int64, error) {
	fen, err := money.FromYuanString(v)
	if err != nil || fen <= 0 {
		return 0, errcode.New(errcode.CodeInvalidParam, field+"金额非法")
	}
	return fen, nil
}

// spuExistsCheck SPU 存在且未删（M13: 拼团/砍价创建的归属校验）。
func spuExistsCheck(ctx context.Context, spuId int64) error {
	n, err := dao.ProductSpu.Ctx(ctx).
		Where(dao.ProductSpu.Columns().Id, spuId).
		Where(dao.ProductSpu.Columns().Deleted, 0).Count()
	if err != nil {
		return gerror.Wrap(err, "查询SPU失败")
	}
	if n == 0 {
		return errcode.New(errcode.CodeProductNotFound, "SPU不存在或已删除")
	}
	return nil
}

// skuExistsCheck SKU 存在且未删。
func skuExistsCheck(ctx context.Context, skuId int64) error {
	n, err := dao.ProductSku.Ctx(ctx).
		Where(dao.ProductSku.Columns().Id, skuId).
		Where(dao.ProductSku.Columns().Deleted, 0).Count()
	if err != nil {
		return gerror.Wrap(err, "查询SKU失败")
	}
	if n == 0 {
		return errcode.New(errcode.CodeProductNotFound, "SKU不存在或已删除")
	}
	return nil
}

// activityRowOf 表行 → 列表 DTO（五类共用形态, D8 同源）。
// 时间列名参数化: 拼团表是 valid_start_at/valid_end_at, 其余表是 start_time/end_time（勘察实证）。
// 玩法字段（group_size/per_limit/reward_type/required_count）按列名防御式读取——
// 缺列的表返回零值（D3-⑧: 各类型列表项出参只在该类型上有意义）。
func activityRowOf(r gdb.Record, spuIdCol, tsCol, teCol string) model.PromotionActivityItem {
	return model.PromotionActivityItem{
		Id:            r["id"].Int64(),
		Name:          r["name"].String(),
		SpuId:         r[spuIdCol].Int64(),
		StartTime:     r[tsCol].String(),
		EndTime:       r[teCol].String(),
		Status:        r["status"].Int(),
		GroupSize:     r["group_size"].Int(),
		PerLimit:      r["per_limit"].Int(),
		RewardType:    r["reward_type"].Int(),
		RequiredCount: r["required_count"].Int(),
	}
}

const (
	tsColStd = "start_time" // 通用活动表
	teColStd = "end_time"
	tsColGrp = "valid_start_at" // 拼团表（000017 列名不同）
	teColGrp = "valid_end_at"
)

// activityCreate 通用活动创建（表名 + 额外表数据）。
func activityCreate(ctx context.Context, table string, data g.Map) (int64, error) {
	res, err := g.DB().Model(table).Ctx(ctx).Data(data).Insert()
	if err != nil {
		return 0, gerror.Wrap(err, "创建活动失败")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, gerror.Wrap(err, "读取活动ID失败")
	}
	return id, nil
}

// activityList 通用活动列表（status 0=全部; 分页; deleted=0; timeCols=[起列,止列]）。
func activityList(ctx context.Context, table, spuIdCol string, timeCols [2]string, status int, page model.PageReq) (*model.PageResult[model.PromotionActivityItem], error) {
	page = page.Normalized()
	m := g.DB().Model(table).Ctx(ctx).Where("deleted", 0)
	if status > 0 {
		m = m.Where("status", status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计活动失败")
	}
	recs, err := m.OrderDesc("id").Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询活动失败")
	}
	list := make([]model.PromotionActivityItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, activityRowOf(r, spuIdCol, timeCols[0], timeCols[1]))
	}
	return &model.PageResult[model.PromotionActivityItem]{List: list, Total: int64(total)}, nil
}

// activityDelete 通用软删（判行数; 不存在 → 50003）。
func activityDelete(ctx context.Context, table string, id int64) error {
	res, err := g.DB().Model(table).Ctx(ctx).
		Where("id", id).Where("deleted", 0).
		Data(g.Map{"deleted": 1, "status": 0}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "删除活动失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeActivityNotFound, "活动不存在")
	}
	return nil
}

// activityTimeUpdate 通用时间窗/名称/状态修改。
// I3（评审修复）: affected=0 有两种语义——"不存在"与"同值无变化"（MySQL 不计同值行）,
// 直接判行数会把同值更新误报 50003（评审探针 P2 实证）→ 前置存在性 Count, 更新本身幂等成功。
func activityTimeUpdate(ctx context.Context, table string, id int64, in model.ActivityTimeInput) error {
	if in.StartTime != "" || in.EndTime != "" {
		if err := activityTimeCheck(in.StartTime, in.EndTime); err != nil {
			return err
		}
	}
	if err := activityMustExist(ctx, table, id); err != nil {
		return err
	}
	data := g.Map{}
	if in.Name != "" {
		data["name"] = in.Name
	}
	if in.StartTime != "" {
		data["start_time"] = gtime.New(in.StartTime)
	}
	if in.EndTime != "" {
		data["end_time"] = gtime.New(in.EndTime)
	}
	if in.Status != nil { // I6: 三态——nil 不修改启停
		data["status"] = *in.Status
	}
	if len(data) == 0 {
		return nil
	}
	_, err := g.DB().Model(table).Ctx(ctx).Where("id", id).Where("deleted", 0).Data(data).Update()
	return gerror.Wrap(err, "修改活动失败")
}

// activityMustExist 活动存在且未删（SetItems 前置）。
func activityMustExist(ctx context.Context, table string, id int64) error {
	n, err := g.DB().Model(table).Ctx(ctx).Where("id", id).Where("deleted", 0).Count()
	if err != nil {
		return gerror.Wrap(err, "查询活动失败")
	}
	if n == 0 {
		return errcode.New(errcode.CodeActivityNotFound, "活动不存在")
	}
	return nil
}

// ---------- 满减（档位+范围嵌套） ----------

func (i *ActivityLogicImpl) FullReductionList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.PromotionActivityItem], error) {
	return activityList(ctx, "promotion_activity", "id", [2]string{tsColStd, teColStd}, status, page) // 满减无 spu 列 → 占位 0
}

func (i *ActivityLogicImpl) FullReductionCreate(ctx context.Context, in model.PromotionActivityInput) (int64, error) {
	if err := activityTimeCheck(in.StartTime, in.EndTime); err != nil {
		return 0, err
	}
	if len(in.Ladders) == 0 {
		return 0, errcode.New(errcode.CodeInvalidParam, "满减活动至少一个档位")
	}
	if err := fullReductionCheck(ctx, in); err != nil {
		return 0, err
	}
	var id int64
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, e := tx.Model("promotion_activity").Ctx(ctx).Data(g.Map{
			"name":       in.Name,
			"start_time": gtime.New(in.StartTime),
			"end_time":   gtime.New(in.EndTime),
			"status":     1,
		}).Insert()
		if e != nil {
			return gerror.Wrap(e, "创建满减活动失败")
		}
		id, e = res.LastInsertId()
		if e != nil {
			return gerror.Wrap(e, "读取活动ID失败")
		}
		return fullReductionSaveChildren(ctx, tx, id, in)
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (i *ActivityLogicImpl) FullReductionDetail(ctx context.Context, id int64) (*model.PromotionActivityDetail, error) {
	a, err := g.DB().Model("promotion_activity").Ctx(ctx).Where("id", id).Where("deleted", 0).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询满减活动失败")
	}
	if a.IsEmpty() {
		return nil, errcode.New(errcode.CodeActivityNotFound, "满减活动不存在")
	}
	ladders, scopes, err := fullReductionLoadChildren(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.PromotionActivityDetail{
		Id:        a["id"].Int64(),
		Name:      a["name"].String(),
		StartTime: a["start_time"].String(),
		EndTime:   a["end_time"].String(),
		Status:    a["status"].Int(),
		Ladders:   ladders,
		Scopes:    scopes,
	}, nil
}

func (i *ActivityLogicImpl) FullReductionUpdate(ctx context.Context, id int64, in model.PromotionActivityInput) error {
	if err := activityTimeCheck(in.StartTime, in.EndTime); err != nil {
		return err
	}
	if err := fullReductionCheck(ctx, in); err != nil {
		return err
	}
	// 存在性前置检查（同值 UPDATE 的 affected=0 是"无变化"不是"不存在"——MySQL 不计同值行,
	// 批次 02 修改语义先例）: 先 Count, Update 本身不判行数（幂等成功）
	if err := activityMustExist(ctx, "promotion_activity", id); err != nil {
		return err
	}
	// 时间字段: 全量更新语义（Create/Update 同一表单提交, api 两字段都随单提交）
	data := g.Map{
		"name":       in.Name,
		"start_time": gtime.New(in.StartTime),
		"end_time":   gtime.New(in.EndTime),
	}
	if in.Status != nil { // I6: 三态——nil 不修改启停
		data["status"] = *in.Status
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := tx.Model("promotion_activity").Ctx(ctx).Where("id", id).Where("deleted", 0).Data(data).Update(); e != nil {
			return gerror.Wrap(e, "修改满减活动失败")
		}
		// 全量替换: 先删后插（D2）
		if _, e := tx.Model("promotion_activity_ladder").Ctx(ctx).Where("activity_id", id).Delete(); e != nil {
			return gerror.Wrap(e, "清空档位失败")
		}
		if _, e := tx.Model("promotion_activity_scope").Ctx(ctx).Where("activity_id", id).Delete(); e != nil {
			return gerror.Wrap(e, "清空范围失败")
		}
		return fullReductionSaveChildren(ctx, tx, id, in)
	})
}

func (i *ActivityLogicImpl) FullReductionDelete(ctx context.Context, id int64) error {
	return activityDelete(ctx, "promotion_activity", id)
}

// fullReductionCheck 档位与范围的业务校验（门槛唯一 50008 / 三型 target 规则）。
func fullReductionCheck(ctx context.Context, in model.PromotionActivityInput) error {
	seen := map[string]bool{}
	for _, l := range in.Ladders {
		th, err := moneyYuanCheck(l.Threshold, "门槛")
		if err != nil {
			return err
		}
		dc, err := moneyYuanCheck(l.Discount, "优惠")
		if err != nil {
			return err
		}
		// I5（评审修复）: "减"不得超过"满"——配置面拦住负应付（满100减200 经计价直通
		// payFen=total-promotion 无下限 → 负应付订单; 计价封顶属 012 文件, 已挂账）
		if dc > th {
			return errcode.New(errcode.CodeInvalidParam, "优惠金额不得超过门槛金额")
		}
		if seen[l.Threshold] {
			return errcode.New(errcode.CodeLadderDup, "档位门槛重复")
		}
		seen[l.Threshold] = true
	}
	for _, s := range in.Scopes {
		switch s.ScopeType {
		case 1:
			if s.TargetId != 0 {
				return errcode.New(errcode.CodeInvalidParam, "全场范围不携带目标ID")
			}
		case 2, 3:
			if s.TargetId <= 0 {
				return errcode.New(errcode.CodeInvalidParam, "分类/商品范围必须指定目标ID")
			}
		default:
			return errcode.New(errcode.CodeInvalidParam, "范围类型非法")
		}
	}
	return nil
}

// fullReductionSaveChildren 档位+范围落库（事务内; 1062 → 50008）。
func fullReductionSaveChildren(ctx context.Context, tx gdb.TX, id int64, in model.PromotionActivityInput) error {
	for _, l := range in.Ladders {
		if _, e := tx.Model("promotion_activity_ladder").Ctx(ctx).Data(g.Map{
			"activity_id":      id,
			"threshold_amount": l.Threshold,
			"discount_amount":  l.Discount,
		}).Insert(); e != nil {
			if isDupKey(e) {
				return errcode.New(errcode.CodeLadderDup, "档位门槛重复")
			}
			return gerror.Wrap(e, "写入档位失败")
		}
	}
	for _, s := range in.Scopes {
		target := any(nil)
		if s.ScopeType != 1 {
			target = s.TargetId
		}
		if _, e := tx.Model("promotion_activity_scope").Ctx(ctx).Data(g.Map{
			"activity_id": id,
			"scope_type":  s.ScopeType,
			"target_id":   target,
		}).Insert(); e != nil {
			return gerror.Wrap(e, "写入范围失败")
		}
	}
	return nil
}

// fullReductionLoadChildren 档位（门槛升序）+范围（D7: 查询口径与 C 端 scopeDescOf 同源）。
func fullReductionLoadChildren(ctx context.Context, id int64) ([]model.PromotionLadder, []model.PromotionScope, error) {
	lrecs, err := g.DB().Model("promotion_activity_ladder").Ctx(ctx).
		Where("activity_id", id).OrderAsc("threshold_amount").All()
	if err != nil {
		return nil, nil, gerror.Wrap(err, "查询档位失败")
	}
	ladders := make([]model.PromotionLadder, 0, len(lrecs))
	for _, r := range lrecs {
		ladders = append(ladders, model.PromotionLadder{
			Threshold: r["threshold_amount"].String(), Discount: r["discount_amount"].String(),
		})
	}
	srecs, err := g.DB().Model("promotion_activity_scope").Ctx(ctx).Where("activity_id", id).All()
	if err != nil {
		return nil, nil, gerror.Wrap(err, "查询范围失败")
	}
	scopes := make([]model.PromotionScope, 0, len(srecs))
	for _, r := range srecs {
		scopes = append(scopes, model.PromotionScope{
			ScopeType: r["scope_type"].Int(), TargetId: r["target_id"].Int64(),
		})
	}
	return ladders, scopes, nil
}

// ---------- 拼团 ----------

func (i *ActivityLogicImpl) GroupBuyList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.PromotionActivityItem], error) {
	return activityList(ctx, "group_buy_activity", "spu_id", [2]string{tsColGrp, teColGrp}, status, page)
}

func (i *ActivityLogicImpl) GroupBuyCreate(ctx context.Context, in model.GroupBuyInput) (int64, error) {
	if in.GroupSize < 2 {
		return 0, errcode.New(errcode.CodeInvalidParam, "成团人数必须大于 1")
	}
	if err := activityTimeCheck(in.StartTime, in.EndTime); err != nil {
		return 0, err
	}
	if err := spuExistsCheck(ctx, in.SpuId); err != nil { // M13: 不挂悬空 SPU
		return 0, err
	}
	// 拼团表时间列为 valid_start_at/valid_end_at（000017, 与其余活动表不同名——勘察实证）
	return activityCreate(ctx, "group_buy_activity", g.Map{
		"name":           in.Name,
		"spu_id":         in.SpuId,
		"group_size":     in.GroupSize,
		"per_limit":      in.PerLimit,
		"valid_start_at": gtime.New(in.StartTime),
		"valid_end_at":   gtime.New(in.EndTime),
		"status":         1,
	})
}

func (i *ActivityLogicImpl) GroupBuyDetail(ctx context.Context, id int64) (*model.PromotionActivityDetail, error) {
	return activityDetailSimple(ctx, "group_buy_activity", id)
}

func (i *ActivityLogicImpl) GroupBuyUpdate(ctx context.Context, id int64, in model.GroupBuyInput) error {
	if in.GroupSize != 0 && in.GroupSize < 2 {
		return errcode.New(errcode.CodeInvalidParam, "成团人数必须大于 1")
	}
	if in.StartTime != "" || in.EndTime != "" {
		if err := activityTimeCheck(in.StartTime, in.EndTime); err != nil {
			return err
		}
	}
	if err := activityMustExist(ctx, "group_buy_activity", id); err != nil { // I3: 同值更新不得误报不存在
		return err
	}
	data := g.Map{}
	if in.Name != "" {
		data["name"] = in.Name
	}
	if in.GroupSize != 0 {
		data["group_size"] = in.GroupSize
	}
	if in.PerLimit != 0 {
		data["per_limit"] = in.PerLimit
	}
	if in.StartTime != "" {
		data["valid_start_at"] = gtime.New(in.StartTime)
	}
	if in.EndTime != "" {
		data["valid_end_at"] = gtime.New(in.EndTime)
	}
	if in.Status != nil { // I6: 三态
		data["status"] = *in.Status
	}
	if len(data) == 0 {
		return nil
	}
	_, err := g.DB().Model("group_buy_activity").Ctx(ctx).Where("id", id).Where("deleted", 0).Data(data).Update()
	return gerror.Wrap(err, "修改拼团活动失败")
}

func (i *ActivityLogicImpl) GroupBuyDelete(ctx context.Context, id int64) error {
	return activityDelete(ctx, "group_buy_activity", id)
}

// GroupBuySetItems 场次商品全量替换（成团价必填; D2 事务; 1062 → 50002）。
func (i *ActivityLogicImpl) GroupBuySetItems(ctx context.Context, id int64, items []model.ActivitySkuInput) error {
	if err := activityMustExist(ctx, "group_buy_activity", id); err != nil {
		return err
	}
	seen := map[int64]bool{}
	for _, it := range items {
		if err := skuExistsCheck(ctx, it.SkuId); err != nil {
			return err
		}
		if _, err := moneyYuanCheck(it.GroupPrice, "成团价"); err != nil {
			return err
		}
		if seen[it.SkuId] {
			return errcode.New(errcode.CodeActivityInvalid, "场次商品SKU重复")
		}
		seen[it.SkuId] = true
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := tx.Model("group_buy_item").Ctx(ctx).Where("activity_id", id).Delete(); e != nil {
			return gerror.Wrap(e, "清空场次商品失败")
		}
		for _, it := range items {
			if _, e := tx.Model("group_buy_item").Ctx(ctx).Data(g.Map{
				"activity_id": id,
				"sku_id":      it.SkuId,
				"group_price": it.GroupPrice,
			}).Insert(); e != nil {
				if isDupKey(e) {
					return errcode.New(errcode.CodeActivityInvalid, "场次商品SKU重复")
				}
				return gerror.Wrap(e, "写入场次商品失败")
			}
		}
		return nil
	})
}

// ---------- 秒杀 ----------

func (i *ActivityLogicImpl) FlashSaleList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.PromotionActivityItem], error) {
	return activityList(ctx, "flash_sale_activity", "id", [2]string{tsColStd, teColStd}, status, page)
}

func (i *ActivityLogicImpl) FlashSaleCreate(ctx context.Context, in model.ActivityTimeInput) (int64, error) {
	if err := activityTimeCheck(in.StartTime, in.EndTime); err != nil {
		return 0, err
	}
	return activityCreate(ctx, "flash_sale_activity", g.Map{
		"name":       in.Name,
		"start_time": gtime.New(in.StartTime),
		"end_time":   gtime.New(in.EndTime),
		"status":     1,
	})
}

func (i *ActivityLogicImpl) FlashSaleDetail(ctx context.Context, id int64) (*model.PromotionActivityDetail, error) {
	d, err := activityDetailSimple(ctx, "flash_sale_activity", id)
	if err != nil {
		return nil, err
	}
	// 场次商品（只读装配; api 无读端点, D4 死契约消费点）
	recs, err := g.DB().Model("flash_sale_item").Ctx(ctx).Where("activity_id", id).OrderAsc("id").All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询场次商品失败")
	}
	for _, r := range recs {
		d.Items = append(d.Items, model.PromotionActivitySku{
			SkuId: r["sku_id"].Int64(), FlashPrice: r["flash_price"].String(),
			StockCount: r["stock_count"].Int(), PerLimit: r["per_limit"].Int(),
		})
	}
	return d, nil
}

func (i *ActivityLogicImpl) FlashSaleUpdate(ctx context.Context, id int64, in model.ActivityTimeInput) error {
	return activityTimeUpdate(ctx, "flash_sale_activity", id, in)
}

func (i *ActivityLogicImpl) FlashSaleDelete(ctx context.Context, id int64) error {
	return activityDelete(ctx, "flash_sale_activity", id)
}

// FlashSaleSetItems 秒杀场次商品全量替换（**已售保护**: 被订单项引用且 sold_count>0 的行
// 不得移除——活动库存账悬空面; D2 事务 + 事务内锁定读复核）。
func (i *ActivityLogicImpl) FlashSaleSetItems(ctx context.Context, id int64, items []model.ActivitySkuInput) error {
	if err := activityMustExist(ctx, "flash_sale_activity", id); err != nil {
		return err
	}
	seen := map[int64]bool{}
	for _, it := range items {
		if err := skuExistsCheck(ctx, it.SkuId); err != nil {
			return err
		}
		if _, err := moneyYuanCheck(it.FlashPrice, "秒杀价"); err != nil {
			return err
		}
		if it.StockCount < 1 {
			return errcode.New(errcode.CodeInvalidParam, "活动限量必须大于 0")
		}
		if it.PerLimit < 0 {
			return errcode.New(errcode.CodeInvalidParam, "每人限购非法")
		}
		if seen[it.SkuId] {
			return errcode.New(errcode.CodeActivityInvalid, "场次商品SKU重复")
		}
		seen[it.SkuId] = true
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// C2（评审修复）: **diff 语义, 不再全删重插**——原实现 LockUpdate 只拦"移除已售行",
		// 随后的全删+重插把**保留的已售行**也删行重建 → sold_count 清零、
		// trade_order_item.flash_sale_item_id 悬空（限购累计/取消回补/售后回补全部失配,
		// 活动库存账清零可超卖——恰是本防线要防的账实悬空, 评审探针 P3 实证）。
		// 改为: 保留行**原位 UPDATE**（只改价格/限量/限购, sold_count 天然不动）,
		// 删除提交集合外的行, 新 SKU INSERT。锁定读整集合复核防并发下单竞态。
		old, e := tx.Model("flash_sale_item").Ctx(ctx).
			Where("activity_id", id).LockUpdate().All()
		if e != nil {
			return gerror.Wrap(e, "查询场次商品失败")
		}
		bySku := map[int64]model.ActivitySkuInput{}
		for _, it := range items {
			bySku[it.SkuId] = it
		}
		// 已售保护: sold_count>0 的行不得移除（必须在提交集合里）
		for _, r := range old {
			if r["sold_count"].Int() > 0 {
				if _, ok := bySku[r["sku_id"].Int64()]; !ok {
					return errcode.New(errcode.CodeActivityInvalid, "场次商品已被抢购, 不可移除")
				}
			}
		}
		for _, r := range old {
			skuId := r["sku_id"].Int64()
			it, keep := bySku[skuId]
			if !keep {
				if _, e = tx.Model("flash_sale_item").Ctx(ctx).Where("id", r["id"].Int64()).Delete(); e != nil {
					return gerror.Wrap(e, "移除场次商品失败")
				}
				continue
			}
			// 保留行: 原位更新配置（sold_count 不在写入集 → 守恒）
			perLimit := it.PerLimit
			if perLimit == 0 {
				perLimit = 1 // 缺省 1
			}
			if _, e = tx.Model("flash_sale_item").Ctx(ctx).Where("id", r["id"].Int64()).Data(g.Map{
				"flash_price": it.FlashPrice,
				"stock_count": it.StockCount,
				"per_limit":   perLimit,
			}).Update(); e != nil {
				return gerror.Wrap(e, "更新场次商品失败")
			}
			delete(bySku, skuId)
		}
		// 剩余 = 新增 SKU
		for _, it := range bySku {
			perLimit := it.PerLimit
			if perLimit == 0 {
				perLimit = 1
			}
			if _, e = tx.Model("flash_sale_item").Ctx(ctx).Data(g.Map{
				"activity_id": id,
				"sku_id":      it.SkuId,
				"flash_price": it.FlashPrice,
				"stock_count": it.StockCount,
				"per_limit":   perLimit,
			}).Insert(); e != nil {
				if isDupKey(e) {
					return errcode.New(errcode.CodeActivityInvalid, "场次商品SKU重复")
				}
				return gerror.Wrap(e, "写入场次商品失败")
			}
		}
		return nil
	})
}

// ---------- 砍价 ----------

func (i *ActivityLogicImpl) BargainList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.PromotionActivityItem], error) {
	return activityList(ctx, "bargain_activity", "spu_id", [2]string{tsColStd, teColStd}, status, page)
}

func (i *ActivityLogicImpl) BargainCreate(ctx context.Context, in model.BargainActivityInput) (int64, error) {
	if err := activityTimeCheck(in.StartTime, in.EndTime); err != nil {
		return 0, err
	}
	if err := spuExistsCheck(ctx, in.SpuId); err != nil { // M13
		return 0, err
	}
	return activityCreate(ctx, "bargain_activity", g.Map{
		"name":       in.Name,
		"spu_id":     in.SpuId,
		"start_time": gtime.New(in.StartTime),
		"end_time":   gtime.New(in.EndTime),
		"status":     1,
	})
}

func (i *ActivityLogicImpl) BargainDetail(ctx context.Context, id int64) (*model.PromotionActivityDetail, error) {
	d, err := activityDetailSimple(ctx, "bargain_activity", id)
	if err != nil {
		return nil, err
	}
	recs, err := g.DB().Model("bargain_item").Ctx(ctx).Where("activity_id", id).OrderAsc("id").All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询场次商品失败")
	}
	for _, r := range recs {
		d.Items = append(d.Items, model.PromotionActivitySku{
			SkuId: r["sku_id"].Int64(), OriginalPrice: r["original_price"].String(),
			FloorPrice: r["floor_price"].String(), MaxCutCount: r["max_cut_count"].Int(),
		})
	}
	return d, nil
}

func (i *ActivityLogicImpl) BargainUpdate(ctx context.Context, id int64, in model.BargainActivityInput) error {
	if in.StartTime != "" || in.EndTime != "" {
		if err := activityTimeCheck(in.StartTime, in.EndTime); err != nil {
			return err
		}
	}
	if err := activityMustExist(ctx, "bargain_activity", id); err != nil { // I3
		return err
	}
	data := g.Map{}
	if in.Name != "" {
		data["name"] = in.Name
	}
	if in.StartTime != "" {
		data["start_time"] = gtime.New(in.StartTime)
	}
	if in.EndTime != "" {
		data["end_time"] = gtime.New(in.EndTime)
	}
	if in.Status != nil { // I6: 三态
		data["status"] = *in.Status
	}
	if len(data) == 0 {
		return nil
	}
	_, err := g.DB().Model("bargain_activity").Ctx(ctx).Where("id", id).Where("deleted", 0).Data(data).Update()
	return gerror.Wrap(err, "修改砍价活动失败")
}

func (i *ActivityLogicImpl) BargainDelete(ctx context.Context, id int64) error {
	return activityDelete(ctx, "bargain_activity", id)
}

// BargainSetItems 场次商品全量替换（original>floor>0; maxCutCount≥0; D2 事务; 1062 → 50002）。
func (i *ActivityLogicImpl) BargainSetItems(ctx context.Context, id int64, items []model.ActivitySkuInput) error {
	if err := activityMustExist(ctx, "bargain_activity", id); err != nil {
		return err
	}
	seen := map[int64]bool{}
	for _, it := range items {
		if err := skuExistsCheck(ctx, it.SkuId); err != nil {
			return err
		}
		orig, err := moneyYuanCheck(it.OriginalPrice, "起始价")
		if err != nil {
			return err
		}
		floor, err := moneyYuanCheck(it.FloorPrice, "底价")
		if err != nil {
			return err
		}
		if floor >= orig {
			return errcode.New(errcode.CodeInvalidParam, "底价必须低于起始价")
		}
		if it.MaxCutCount < 0 {
			return errcode.New(errcode.CodeInvalidParam, "最大刀数非法")
		}
		if seen[it.SkuId] {
			return errcode.New(errcode.CodeActivityInvalid, "场次商品SKU重复")
		}
		seen[it.SkuId] = true
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := tx.Model("bargain_item").Ctx(ctx).Where("activity_id", id).Delete(); e != nil {
			return gerror.Wrap(e, "清空场次商品失败")
		}
		for _, it := range items {
			if _, e := tx.Model("bargain_item").Ctx(ctx).Data(g.Map{
				"activity_id":    id,
				"sku_id":         it.SkuId,
				"original_price": it.OriginalPrice,
				"floor_price":    it.FloorPrice,
				"max_cut_count":  it.MaxCutCount,
				"config":         "{}",
			}).Insert(); e != nil {
				if isDupKey(e) {
					return errcode.New(errcode.CodeActivityInvalid, "场次商品SKU重复")
				}
				return gerror.Wrap(e, "写入场次商品失败")
			}
		}
		return nil
	})
}

// ---------- 助力 ----------

func (i *ActivityLogicImpl) AssistList(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.PromotionActivityItem], error) {
	out, err := activityList(ctx, "assist_activity", "id", [2]string{tsColStd, teColStd}, status, page) // 助力无 spu 列
	if err != nil {
		return nil, err
	}
	// RewardDesc 装配（D5/D7）: 积分型从 config JSON; 券型一次 IN 查询批量取名（避免 N+1）
	couponIds := []int64{}
	// activityList 未带 reward_ref/config 列 → 单独补查本页活动的这两列（一次 IN）
	ids := make([]int64, 0, len(out.List))
	for _, it := range out.List {
		ids = append(ids, it.Id)
	}
	if len(ids) == 0 {
		return out, nil
	}
	extra, err := g.DB().Model("assist_activity").Ctx(ctx).
		Fields("id", "reward_type", "reward_ref", "required_count", "config").
		WhereIn("id", ids).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询助力奖励失败")
	}
	type rewardInfo struct {
		ref    int64
		rtype  int
		reqCnt int
		points int
	}
	info := map[int64]rewardInfo{}
	for _, r := range extra {
		pt := 0
		if r["reward_type"].Int() == 2 {
			if m := g.NewVar(r["config"].String()).Map(); m != nil {
				pt = g.NewVar(m["pointAmount"]).Int()
			}
		}
		info[r["id"].Int64()] = rewardInfo{ref: r["reward_ref"].Int64(), rtype: r["reward_type"].Int(),
			reqCnt: r["required_count"].Int(), points: pt}
		if r["reward_type"].Int() == 1 {
			couponIds = append(couponIds, r["reward_ref"].Int64())
		}
	}
	names := map[int64]string{}
	if len(couponIds) > 0 {
		coups, e := dao.Coupon.Ctx(ctx).
			Fields(dao.Coupon.Columns().Id, dao.Coupon.Columns().Name).
			WhereIn(dao.Coupon.Columns().Id, couponIds).All()
		if e != nil {
			return nil, gerror.Wrap(e, "查询奖励券失败")
		}
		for _, c := range coups {
			names[c[dao.Coupon.Columns().Id].Int64()] = c[dao.Coupon.Columns().Name].String()
		}
	}
	for i := range out.List {
		ri, ok := info[out.List[i].Id]
		if !ok {
			continue
		}
		out.List[i].RewardType = ri.rtype
		out.List[i].RequiredCount = ri.reqCnt
		if ri.rtype == 1 {
			name := names[ri.ref]
			if name == "" {
				name = "奖励券已删除" // M13: 券被删后不出现空引号
			}
			out.List[i].RewardDesc = "邀 " + intToString(ri.reqCnt) + " 人得券「" + name + "」"
		} else {
			out.List[i].RewardDesc = "邀 " + intToString(ri.reqCnt) + " 人得 " + intToString(ri.points) + " 积分"
		}
	}
	return out, nil
}

// AssistCreate 创建助力活动（D5: 券→reward_ref 校验存在; 积分→config JSON; D6: perLimit 如实存储）。
func (i *ActivityLogicImpl) AssistCreate(ctx context.Context, in model.AssistActivityInput) (int64, error) {
	if err := assistInputCheck(ctx, in); err != nil {
		return 0, err
	}
	if err := activityTimeCheck(in.StartTime, in.EndTime); err != nil {
		return 0, err
	}
	data := g.Map{
		"name":           in.Name,
		"reward_type":    in.RewardType,
		"reward_ref":     0,
		"required_count": in.RequiredCount,
		"per_limit":      in.PerLimit,
		"start_time":     gtime.New(in.StartTime),
		"end_time":       gtime.New(in.EndTime),
		"status":         1,
		"config":         "{}",
	}
	if in.RewardType == 1 {
		data["reward_ref"] = in.RewardRef
	} else {
		data["config"] = `{"pointAmount":` + intToString(in.PointAmount) + `}`
	}
	return activityCreate(ctx, "assist_activity", data)
}

func (i *ActivityLogicImpl) AssistDetail(ctx context.Context, id int64) (*model.PromotionActivityDetail, error) {
	return activityDetailSimple(ctx, "assist_activity", id)
}

func (i *ActivityLogicImpl) AssistUpdate(ctx context.Context, id int64, in model.AssistActivityInput) error {
	if in.StartTime != "" || in.EndTime != "" {
		if err := activityTimeCheck(in.StartTime, in.EndTime); err != nil {
			return err
		}
	}
	if in.RequiredCount != 0 && in.RequiredCount < 1 {
		return errcode.New(errcode.CodeInvalidParam, "所需人数必须大于 0")
	}
	if err := activityMustExist(ctx, "assist_activity", id); err != nil { // I3
		return err
	}
	data := g.Map{}
	if in.Name != "" {
		data["name"] = in.Name
	}
	if in.RequiredCount != 0 {
		data["required_count"] = in.RequiredCount
	}
	if in.PerLimit != 0 {
		data["per_limit"] = in.PerLimit // D6: >1 如实存储（uk_activity_user 兜底"每人一次"）
	}
	if in.StartTime != "" {
		data["start_time"] = gtime.New(in.StartTime)
	}
	if in.EndTime != "" {
		data["end_time"] = gtime.New(in.EndTime)
	}
	if in.Status != nil { // I6: 三态
		data["status"] = *in.Status
	}
	if len(data) == 0 {
		return nil
	}
	_, err := g.DB().Model("assist_activity").Ctx(ctx).Where("id", id).Where("deleted", 0).Data(data).Update()
	return gerror.Wrap(err, "修改助力活动失败")
}

func (i *ActivityLogicImpl) AssistDelete(ctx context.Context, id int64) error {
	return activityDelete(ctx, "assist_activity", id)
}

// assistInputCheck 助力创建校验（奖励两型/所需人数/时间窗）。
func assistInputCheck(ctx context.Context, in model.AssistActivityInput) error {
	if in.RewardType != 1 && in.RewardType != 2 {
		return errcode.New(errcode.CodeInvalidParam, "奖励类型非法")
	}
	if in.RequiredCount < 1 {
		return errcode.New(errcode.CodeInvalidParam, "所需人数必须大于 0")
	}
	if in.RewardType == 1 {
		if in.RewardRef <= 0 {
			return errcode.New(errcode.CodeInvalidParam, "券奖励必须指定券模板")
		}
		n, err := dao.Coupon.Ctx(ctx).
			Where(dao.Coupon.Columns().Id, in.RewardRef).
			Where(dao.Coupon.Columns().Deleted, 0).Count()
		if err != nil {
			return gerror.Wrap(err, "查询券模板失败")
		}
		if n == 0 {
			return errcode.New(errcode.CodeActivityNotFound, "奖励券不存在")
		}
	} else if in.PointAmount < 1 {
		return errcode.New(errcode.CodeInvalidParam, "积分奖励必须大于 0")
	}
	return nil
}

// ---------- 共用小件 ----------

// activityDetailSimple 活动基础详情（无场次商品装配; D4: 三类 Detail 的公共部分）。
func activityDetailSimple(ctx context.Context, table string, id int64) (*model.PromotionActivityDetail, error) {
	a, err := g.DB().Model(table).Ctx(ctx).Where("id", id).Where("deleted", 0).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询活动失败")
	}
	if a.IsEmpty() {
		return nil, errcode.New(errcode.CodeActivityNotFound, "活动不存在")
	}
	return &model.PromotionActivityDetail{
		Id:        a["id"].Int64(),
		Name:      a["name"].String(),
		StartTime: a["start_time"].String(),
		EndTime:   a["end_time"].String(),
		Status:    a["status"].Int(),
	}, nil
}

// intToString 极简 itoa（包内已有 itoa 为 int 版; 此处规避再次包装）。
func intToString(n int) string {
	return g.NewVar(n).String()
}
