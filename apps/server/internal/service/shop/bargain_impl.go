// bargain_impl.go 砍价动作实现（015-marketing-c 批次 09）。
// 规则见 bargain.go 的接口注释; 金额算法见 research D4（确定性均分 + 末刀补齐）。
package shop

import (
	"context"
	"strconv"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/library/idgen"
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

// nextBargainNo 生成砍价单号（sonyflake + 前缀, 与售后单号同风格）。
// 评审 C1 修复: 原"秒级时间戳 + uid%10000"同用户同秒两次发起必撞 uk_bargain_no。
func nextBargainNo() (string, error) {
	next, err := idgen.NextID()
	if err != nil {
		return "", gerror.Wrap(err, "生成砍价单号失败")
	}
	return "BG" + strconv.FormatInt(next, 10), nil
}

// Launch 发起砍价（FR-007）: 校验场次商品与活动时间窗 → 建单（**发起即首刀**）。
// C1（评审修复）: order_no **不写入**（迁移 000041 已改可空, 默认 NULL）——
// "下单后回填"本就是 NULL 语义, MySQL 唯一索引不约束 NULL, 未下单的砍价单可任意多条共存;
// 原默认值 '' 会让全库只能存在一张砍价单（第二人发起必撞 uk_order_no 1062）。
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

	// M1（评审修复）: 首刀就到底价（max_cut_count=1 / original==floor / 均分取整为 0）时,
	// 记录应直接处于"到底价待下单"（2）+ success_time, 与 Cut 的到底判定同口径——
	// 否则会出现"当前价=底价却砍价中", 好友再砍只能记一笔 0.00 才翻状态。
	status := 1
	if curFen <= floorFen {
		status = 2
	}

	bargainNo, err := nextBargainNo()
	if err != nil {
		return nil, err
	}
	// 砍价单与首刀流水**同一事务**（复审新发现-2）: 原先 record 先提交、流水后插——
	// 流水插入失败（非 1062）会残留 cut_count=1 而无流水行, 违反本轮守恒不变量
	// （Σ帮砍流水 == 起始价 − 当前价; M2 同款收编, Cut/Help 已是事务而 Launch 漏了）。
	var out *model.BargainLaunchResult
	if err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		data := do.BargainRecord{
			BargainNo:    bargainNo,
			ItemId:       bargainItemId,
			UserId:       userId,
			CurrentPrice: money.ToYuanString(curFen),
			CutCount:     1, // 发起即首刀
			Status:       status,
			ExpireTime:   expire,
			// OrderNo: 不写 → 默认 NULL（见函数头 C1 说明）
		}
		if status == 2 {
			data.SuccessTime = gtime.Now()
		}
		res, e := dao.BargainRecord.Ctx(ctx).Data(data).Insert()
		if e != nil {
			// 1062 → 业务码（uk_bargain_no 等唯一键兜底; sonyflake 后撞键概率已可忽略, 纯防御）
			if isDupKey(e) {
				return errcode.New(errcode.CodeTooFrequent, "发起过于频繁, 请稍后再试")
			}
			return gerror.Wrap(e, "创建砍价单失败")
		}
		recId, e := res.LastInsertId()
		if e != nil {
			return gerror.Wrap(e, "读取砍价单ID失败")
		}
		// 首刀记入帮砍流水（发起者本人; 与后续帮砍同构, 故进度列表自然含首刀）
		if _, e = dao.BargainHelper.Ctx(ctx).Data(do.BargainHelper{
			RecordId:     recId,
			HelperUserId: userId,
			CutAmount:    money.ToYuanString(per),
		}).Insert(); e != nil {
			return gerror.Wrap(e, "记录首刀失败")
		}
		out = &model.BargainLaunchResult{RecordId: recId, CurrentPrice: money.ToYuanString(curFen)}
		return nil
	}); err != nil {
		return nil, err
	}
	return out, nil
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
// C2（评审修复）: "读价→算刀→写价"必须串行——原条件更新守卫只有 status=1, 而中间刀不改状态,
// 并发全部命中 → current_price 被各自的陈旧读覆盖（丢失更新: 8 人砍 8.88 元、流水却记 39.96 元）。
// 故整段收进**事务 + LockUpdate 当前读**重算（013 I-1 结论: 事务内普通读被 InnoDB 读视图钉住,
// 不可替代当前读）; 帮砍流水与价格推进同生共死（M2: 原先流水先写、推进失败不回滚,
// 该会员会被永久记为"已帮砍"且进度里多一笔从未影响价格的刀）。
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
	maxCuts := item[icols.MaxCutCount].Int()

	var result *model.BargainCutResult
	if err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// ① 锁定读**当前价与刀数**（当前读, 并发帮砍在此串行; 再校验状态防读过期）
		locked, e := dao.BargainRecord.Ctx(ctx).
			Where(rcols.Id, recordId).LockUpdate().One()
		if e != nil {
			return gerror.Wrap(e, "锁定砍价单失败")
		}
		if locked.IsEmpty() || locked[rcols.Status].Int() != 1 {
			return errcode.New(errcode.CodeActivityInvalid, "该砍价单已结束")
		}
		curFen, e := yuanFen(locked[rcols.CurrentPrice].String())
		if e != nil {
			return e
		}
		cutCount := locked[rcols.CutCount].Int()

		// ② 本刀金额: 常规取均分; **本次是最后一刀时补齐到恰好底价**（不砍穿、多刀之和 = 起始价 − 底价）
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

		// ③ 一人一刀: 写流水（金额为本刀实砍额, 与 ④ 的价格推进同事务原子生效;
		//    唯一键 uk_record_helper 兜底并发, 1062 → 50004）
		if _, e = dao.BargainHelper.Ctx(ctx).Data(do.BargainHelper{
			RecordId:     recordId,
			HelperUserId: userId,
			CutAmount:    money.ToYuanString(cut),
		}).Insert(); e != nil {
			if isDupKey(e) {
				return errcode.New(errcode.CodeAlreadyCut, "您已经帮砍过了")
			}
			return gerror.Wrap(e, "记录帮砍失败")
		}

		// ④ 条件推进当前价与刀数（判行数——012/013/014 三轮评审确立）; 到底价置 2 + success_time
		//    （success_time 供砍价成交下单 FR-006 校验与进度展示）
		data := g.Map{
			rcols.CurrentPrice: money.ToYuanString(nextFen),
			rcols.CutCount:     gdb.Raw("cut_count + 1"),
		}
		if nextFen <= floorFen {
			data[rcols.Status] = 2 // 到底价（可下单）
			data[rcols.SuccessTime] = gtime.Now()
		}
		res, e := dao.BargainRecord.Ctx(ctx).
			Where(rcols.Id, recordId).
			Where(rcols.Status, 1).
			Data(data).Update()
		if e != nil {
			return gerror.Wrap(e, "推进砍价失败")
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return errcode.New(errcode.CodeActivityInvalid, "该砍价单已结束")
		}

		result = &model.BargainCutResult{
			CutAmount:    money.ToYuanString(cut),
			CurrentPrice: money.ToYuanString(nextFen),
			FloorReached: nextFen <= floorFen,
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
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
