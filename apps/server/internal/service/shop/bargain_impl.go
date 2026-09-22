// bargain_impl.go 砍价动作实现（015-marketing-c 批次 09）。
// 规则见 bargain.go 的接口注释; 金额算法见 research D4（确定性均分 + 末刀补齐）。
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
	"ecboot/internal/model/do"
)

// BargainLogicImpl IBargainLogic 实现。
type BargainLogicImpl struct{}

func NewBargainLogic() *BargainLogicImpl { return &BargainLogicImpl{} }

// riskHitForPlay 玩法入口风控（FR-015）: 走端口, **未装配时放行 + 告警**（不阻塞玩法可用;
// 实现属批次 12 风控规则引擎）。迁移 000029 明写"砍价/助力是被刷重灾区——入口须挂风控"。
func riskHitForPlay(ctx context.Context, userId int64, ruleType int, scene string) (bool, error) {
	if RiskHit == nil {
		g.Log().Warningf(ctx, "风控端口未装配(批次12): 放行 %s user_id=%d", scene, userId)
		return false, nil
	}
	blocked, err := RiskHit.Hit(ctx, userId, ruleType, scene)
	if err != nil {
		return false, gerror.Wrap(err, "风控判定失败")
	}
	return blocked, nil
}

// perCutFen 每刀金额 = (起始价 − 底价) ÷ 最大刀数（**向下取整到分**, research D4）。
func perCutFen(origFen, floorFen int64, maxCuts int) int64 {
	if maxCuts <= 0 || origFen <= floorFen {
		return 0
	}
	return (origFen - floorFen) / int64(maxCuts)
}

// Launch 发起砍价（FR-007）: 校验场次商品与活动时间窗 → 建单（**发起即首刀**）。
func (i *BargainLogicImpl) Launch(
	ctx context.Context, userId, bargainItemId int64,
) (*model.BargainLaunchResult, error) {
	item, err := loadBargainItemForPlay(ctx, bargainItemId)
	if err != nil {
		return nil, err
	}
	icols := dao.BargainItem.Columns()
	origFen, err := yuanFen(item[icols.OriginalPrice].String())
	if err != nil {
		return nil, err
	}
	floorFen, err := yuanFen(item[icols.FloorPrice].String())
	if err != nil {
		return nil, err
	}
	maxCuts := item[icols.MaxCutCount].Int()
	per := perCutFen(origFen, floorFen, maxCuts)
	curFen := origFen - per
	if curFen < floorFen {
		curFen = floorFen
	}

	acols := dao.BargainActivity.Columns()
	act, err := dao.BargainActivity.Ctx(ctx).
		Where(acols.Id, item[icols.ActivityId].Int64()).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询砍价活动失败")
	}
	// 有效期: 与活动截止同源（活动结束时砍价自然失效）
	expire := act[acols.EndTime].GTime()

	res, err := dao.BargainRecord.Ctx(ctx).Data(do.BargainRecord{
		BargainNo:    "BG" + gtime.Now().Format("060102150405") + itoa(int(userId%10000)),
		ItemId:       bargainItemId,
		UserId:       userId,
		CurrentPrice: money.ToYuanString(curFen),
		CutCount:     1, // 发起即首刀
		Status:       1,
		ExpireTime:   expire,
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "创建砍价单失败")
	}
	recId, err := res.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "读取砍价单ID失败")
	}
	// 首刀记入帮砍流水（发起者本人; 与后续帮砍同构, 故进度列表自然含首刀）
	if _, e := dao.BargainHelper.Ctx(ctx).Data(do.BargainHelper{
		RecordId:     recId,
		HelperUserId: userId,
		CutAmount:    money.ToYuanString(per),
	}).Insert(); e != nil {
		return nil, gerror.Wrap(e, "记录首刀失败")
	}
	return &model.BargainLaunchResult{RecordId: recId, CurrentPrice: money.ToYuanString(curFen)}, nil
}

// Progress 砍价进度（FR-008）: 公开; 帮砍列表昵称脱敏; 超时**惰性判定**（不改表）。
func (i *BargainLogicImpl) Progress(ctx context.Context, recordId int64) (*model.BargainProgressView, error) {
	rec, item, err := loadBargainRecord(ctx, recordId)
	if err != nil {
		return nil, err
	}
	rcols, icols := dao.BargainRecord.Columns(), dao.BargainItem.Columns()
	helpers, err := bargainHelpers(ctx, recordId)
	if err != nil {
		return nil, err
	}
	status := rec[rcols.Status].Int()
	if status == 1 && rec[rcols.ExpireTime].GTime() != nil && !rec[rcols.ExpireTime].GTime().After(gtime.Now()) {
		status = 4 // 超时（惰性判定; 不改表, 由后续调度或查询口径统一）
	}
	return &model.BargainProgressView{
		RecordId:      recordId,
		SkuId:         item[icols.SkuId].Int64(),
		OriginalPrice: item[icols.OriginalPrice].String(),
		CurrentPrice:  rec[rcols.CurrentPrice].String(),
		FloorPrice:    item[icols.FloorPrice].String(),
		CutCount:      rec[rcols.CutCount].Int(),
		Status:        status,
		ExpireTime:    rec[rcols.ExpireTime].String(),
		OrderNo:       rec[rcols.OrderNo].String(),
		Helpers:       helpers,
	}, nil
}

// Cut 帮砍一刀（FR-009）: 一人一刀 / 本人不可 / 超时不可 / **不砍穿（末刀补齐到恰好底价）** / 风控。
func (i *BargainLogicImpl) Cut(
	ctx context.Context, userId, recordId int64,
) (*model.BargainCutResult, error) {
	if blocked, err := riskHitForPlay(ctx, userId, 2, "bargain_cut"); err != nil {
		return nil, err
	} else if blocked {
		return nil, errcode.New(errcode.CodeForbidden, "操作过于频繁, 请稍后再试")
	}

	rec, item, err := loadBargainRecord(ctx, recordId)
	if err != nil {
		return nil, err
	}
	rcols, icols := dao.BargainRecord.Columns(), dao.BargainItem.Columns()
	if rec[rcols.UserId].Int64() == userId {
		return nil, errcode.New(errcode.CodeActivityInvalid, "不能帮砍自己的砍价单")
	}
	if rec[rcols.Status].Int() != 1 {
		return nil, errcode.New(errcode.CodeActivityInvalid, "该砍价单已结束")
	}
	if rec[rcols.ExpireTime].GTime() != nil && !rec[rcols.ExpireTime].GTime().After(gtime.Now()) {
		return nil, errcode.New(errcode.CodeActivityInvalid, "该砍价单已超时")
	}

	origFen, err := yuanFen(item[icols.OriginalPrice].String())
	if err != nil {
		return nil, err
	}
	floorFen, err := yuanFen(item[icols.FloorPrice].String())
	if err != nil {
		return nil, err
	}
	curFen, err := yuanFen(rec[rcols.CurrentPrice].String())
	if err != nil {
		return nil, err
	}
	cutCount := rec[rcols.CutCount].Int()
	maxCuts := item[icols.MaxCutCount].Int()

	// 本刀金额: 常规取均分; **本次是最后一刀时补齐到恰好底价**（不砍穿、多刀之和 = 起始价 − 底价）
	cut := perCutFen(origFen, floorFen, maxCuts)
	if cutCount+1 >= maxCuts {
		cut = curFen - floorFen
	}
	if cut > curFen-floorFen {
		cut = curFen - floorFen
	}
	if cut < 0 {
		cut = 0
	}
	nextFen := curFen - cut

	// 一人一刀: 先写流水（唯一键 uk_record_helper 兜底并发, 1062 转业务码）
	if _, e := dao.BargainHelper.Ctx(ctx).Data(do.BargainHelper{
		RecordId:     recordId,
		HelperUserId: userId,
		CutAmount:    money.ToYuanString(cut),
	}).Insert(); e != nil {
		if isDupKey(e) {
			return nil, errcode.New(errcode.CodeAlreadyCut, "您已经帮砍过了")
		}
		return nil, gerror.Wrap(e, "记录帮砍失败")
	}

	// 条件更新当前价与刀数（防并发重入; 判行数——012/013/014 三轮评审确立）
	status := 1
	if nextFen <= floorFen {
		status = 2 // 到底价（可下单）
	}
	res, err := dao.BargainRecord.Ctx(ctx).
		Where(rcols.Id, recordId).
		Where(rcols.Status, 1).
		Data(g.Map{
			rcols.CurrentPrice: money.ToYuanString(nextFen),
			rcols.CutCount:     gdb.Raw("cut_count + 1"),
			rcols.Status:       status,
		}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "推进砍价失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return nil, errcode.New(errcode.CodeActivityInvalid, "该砍价单已结束")
	}
	return &model.BargainCutResult{
		CutAmount:    money.ToYuanString(cut),
		CurrentPrice: money.ToYuanString(nextFen),
		FloorReached: nextFen <= floorFen,
	}, nil
}

// loadBargainItemForPlay 取砍价场次商品并校验活动（启用未删 + 时间窗内）。
func loadBargainItemForPlay(ctx context.Context, itemId int64) (gdb.Record, error) {
	icols := dao.BargainItem.Columns()
	item, err := dao.BargainItem.Ctx(ctx).Where(icols.Id, itemId).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询砍价场次商品失败")
	}
	if item.IsEmpty() {
		return nil, errcode.New(errcode.CodeActivityNotFound, "砍价场次商品不存在")
	}
	acols := dao.BargainActivity.Columns()
	act, err := dao.BargainActivity.Ctx(ctx).
		Where(acols.Id, item[icols.ActivityId].Int64()).
		Where(acols.Status, 1).Where(acols.Deleted, 0).
		Where(acols.StartTime+" <= NOW()").Where(acols.EndTime+" >= NOW()").One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询砍价活动失败")
	}
	if act.IsEmpty() {
		return nil, errcode.New(errcode.CodeActivityInvalid, "砍价活动不在进行中")
	}
	return item, nil
}

// loadBargainRecord 取砍价单与其场次商品。
func loadBargainRecord(ctx context.Context, recordId int64) (gdb.Record, gdb.Record, error) {
	rcols := dao.BargainRecord.Columns()
	rec, err := dao.BargainRecord.Ctx(ctx).Where(rcols.Id, recordId).One()
	if err != nil {
		return nil, nil, gerror.Wrap(err, "查询砍价单失败")
	}
	if rec.IsEmpty() {
		return nil, nil, errcode.New(errcode.CodeActivityNotFound, "砍价单不存在")
	}
	icols := dao.BargainItem.Columns()
	item, err := dao.BargainItem.Ctx(ctx).Where(icols.Id, rec[rcols.ItemId].Int64()).One()
	if err != nil {
		return nil, nil, gerror.Wrap(err, "查询砍价场次商品失败")
	}
	if item.IsEmpty() {
		return nil, nil, errcode.New(errcode.CodeActivityNotFound, "砍价场次商品不存在")
	}
	return rec, item, nil
}

// bargainHelpers 帮砍列表（昵称脱敏; 一次 IN 查询取昵称, 避免 N+1）。
func bargainHelpers(ctx context.Context, recordId int64) ([]model.PlayHelperItem, error) {
	cols := dao.BargainHelper.Columns()
	recs, err := dao.BargainHelper.Ctx(ctx).Where(cols.RecordId, recordId).OrderAsc(cols.Id).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询帮砍记录失败")
	}
	nicks, err := helperNicknames(ctx, collectHelperUids(recs, cols.HelperUserId))
	if err != nil {
		return nil, err
	}
	out := make([]model.PlayHelperItem, 0, len(recs))
	for _, r := range recs {
		uid := r[cols.HelperUserId].Int64()
		out = append(out, model.PlayHelperItem{
			UserId:    uid,
			Nickname:  maskReviewUser(nicks[uid], 0), // 复用评价域的脱敏口径（shop 域内自带, 兄弟域不可 import）
			Amount:    r[cols.CutAmount].String(),
			CreatedAt: r[cols.CreatedAt].String(),
		})
	}
	return out, nil
}
