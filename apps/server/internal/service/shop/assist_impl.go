// assist_impl.go 助力动作实现（015-marketing-c 批次 09）。
// 规则见 assist.go 的接口注释; 达标投递发奖意图（端口, 实际发放跨 user 域, research D3）。
package shop

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/model"
	"ecboot/internal/model/do"
)

// AssistLogicImpl IAssistLogic 实现。
type AssistLogicImpl struct{}

func NewAssistLogic() *AssistLogicImpl { return &AssistLogicImpl{} }

// Launch 发起助力（FR-012）: 校验活动时间窗与"每人可发起次数" → 建参与记录。
// I1（评审修复）: per_limit 的语义已被 schema 钉死为"每人一次"——assist_record 上有
// UNIQUE KEY uk_activity_user(activity_id,user_id)（000029）, per_limit>1 结构性不可达
// （计数放行后第二次 INSERT 仍会撞键）。此处保留计数检查只为 per_limit=1 的友好前置拦截,
// 真正的并发防线是唯一键兜底 + 1062 → 50006（原来裸抛会被归一成 10002 系统错误）。
// 批次 10 营销后台配置 per_limit 时必须按此约束（>1 不生效）。
func (i *AssistLogicImpl) Launch(
	ctx context.Context, userId, activityId int64,
) (*model.AssistLaunchResult, error) {
	act, err := loadAssistActivityForPlay(ctx, activityId)
	if err != nil {
		return nil, err
	}
	acols, rcols := dao.AssistActivity.Columns(), dao.AssistRecord.Columns()
	// 每人可发起次数上限（assist_activity.per_limit; 见函数头 I1 说明——schema 只支持 1）
	cnt, err := dao.AssistRecord.Ctx(ctx).
		Where(rcols.ActivityId, activityId).
		Where(rcols.UserId, userId).Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计发起次数失败")
	}
	if perLimit := act[acols.PerLimit].Int(); perLimit > 0 && cnt >= perLimit {
		return nil, errcode.New(errcode.CodeAssistUsedUp, "该活动的发起次数已用完")
	}

	res, err := dao.AssistRecord.Ctx(ctx).Data(do.AssistRecord{
		ActivityId:  activityId,
		UserId:      userId,
		HelperCount: 0,
		Status:      1, // 进行中
	}).Insert()
	if err != nil {
		if isDupKey(err) { // 并发双击: 唯一键兜底 → 业务码（不再裸抛 10002）
			return nil, errcode.New(errcode.CodeAssistUsedUp, "该活动的发起次数已用完")
		}
		return nil, gerror.Wrap(err, "创建助力记录失败")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "读取记录ID失败")
	}
	return &model.AssistLaunchResult{RecordId: id}, nil
}

// Progress 助力进度（FR-013）: 公开; 助力人列表昵称脱敏; 活动结束未达标 → 惰性判定为"已过期"。
func (i *AssistLogicImpl) Progress(ctx context.Context, recordId int64) (*model.AssistProgressView, error) {
	rec, act, err := loadAssistRecord(ctx, recordId)
	if err != nil {
		return nil, err
	}
	rcols, acols := dao.AssistRecord.Columns(), dao.AssistActivity.Columns()
	helpers, err := assistHelpers(ctx, recordId)
	if err != nil {
		return nil, err
	}
	status := rec[rcols.Status].Int()
	if status == 1 && act[acols.EndTime].GTime() != nil && !act[acols.EndTime].GTime().After(gtime.Now()) {
		status = 3 // 已过期（惰性判定, 不改表）
	}
	return &model.AssistProgressView{
		RecordId:      recordId,
		ActivityId:    rec[rcols.ActivityId].Int64(),
		HelperCount:   rec[rcols.HelperCount].Int(),
		RequiredCount: act[acols.RequiredCount].Int(),
		Status:        status,
		FinishTime:    rec[rcols.FinishTime].String(),
		Helpers:       helpers,
	}, nil
}

// Help 助力（FR-013/014）: 一人一助力 / 本人不可 / 活动结束不可 / 风控;
// 达标时**仅实际完成状态迁移的那一次**投递发奖意图（并发下天然只一次）。
func (i *AssistLogicImpl) Help(
	ctx context.Context, userId, recordId int64,
) (*model.AssistHelpResult, error) {
	if blocked, err := riskHitForPlay(ctx, userId, 2, "assist_help"); err != nil {
		return nil, err
	} else if blocked {
		return nil, errcode.New(errcode.CodeForbidden, "操作过于频繁, 请稍后再试")
	}

	rec, act, err := loadAssistRecord(ctx, recordId)
	if err != nil {
		return nil, err
	}
	rcols, acols := dao.AssistRecord.Columns(), dao.AssistActivity.Columns()
	if rec[rcols.UserId].Int64() == userId {
		return nil, errcode.New(errcode.CodeActivityInvalid, "不能助力自己的记录")
	}
	if rec[rcols.Status].Int() != 1 {
		return nil, errcode.New(errcode.CodeActivityInvalid, "该助力记录已结束")
	}
	if act[acols.EndTime].GTime() != nil && !act[acols.EndTime].GTime().After(gtime.Now()) {
		return nil, errcode.New(errcode.CodeActivityInvalid, "该助力活动已结束")
	}

	// M2（评审修复）: 助力流水、计数推进与达标置位收进**同一事务**同生共死——
	// 原先流水先写、计数推进失败（affected=0）不回滚, 该会员会被永久记为"已助力"（不能重试）。
	var done bool
	if err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 一人一助力: 先写流水（唯一键 uk_record_helper 兜底并发 → 1062 转 50007）
		if _, e := dao.AssistHelper.Ctx(ctx).Data(do.AssistHelper{
			RecordId:     recordId,
			HelperUserId: userId,
		}).Insert(); e != nil {
			if isDupKey(e) {
				return errcode.New(errcode.CodeAlreadyAssisted, "您已经助力过了")
			}
			return gerror.Wrap(e, "记录助力失败")
		}

		// 条件推进助力计数（判行数）
		res, e := dao.AssistRecord.Ctx(ctx).
			Where(rcols.Id, recordId).
			Where(rcols.Status, 1).
			Data(g.Map{rcols.HelperCount: gdb.Raw("helper_count + 1")}).
			Update()
		if e != nil {
			return gerror.Wrap(e, "推进助力计数失败")
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return errcode.New(errcode.CodeActivityInvalid, "该助力记录已结束")
		}

		// 达标判定: 必须用**自增后的最新计数**（重读）。
		// 教训: 若用自增前读到的 `rec.HelperCount + 1` 判定, 并发助力时每个请求都看到旧计数,
		// 谁都判不出"达标" → 记录永远停在"进行中"（并发用例实测 EXPECT 1 == 2 抓到）。
		required := act[acols.RequiredCount].Int()
		cur, e := dao.AssistRecord.Ctx(ctx).
			Where(rcols.Id, recordId).Value(rcols.HelperCount)
		if e != nil {
			return gerror.Wrap(e, "读取助力人数失败")
		}
		if cur.Int() < required {
			return nil // 未达标: 提交流水与计数, Done=false
		}
		// 达标: **条件置位**（status 1→2 且判行数）——只有真正完成迁移的那一次才投递发奖意图,
		// 并发下其余请求 affected=0 → 不投递（FR-014: 发奖只一次）
		fin, e := dao.AssistRecord.Ctx(ctx).
			Where(rcols.Id, recordId).
			Where(rcols.Status, 1).
			Data(g.Map{rcols.Status: 2, rcols.FinishTime: gtime.Now()}).
			Update()
		if e != nil {
			return gerror.Wrap(e, "完成助力记录失败")
		}
		if n, _ := fin.RowsAffected(); n == 0 {
			return nil // 已被并发者置位 → 不重复投递
		}
		done = true
		return nil
	}); err != nil {
		return nil, err
	}
	if !done {
		return &model.AssistHelpResult{Done: false}, nil
	}
	// 投递发奖意图（端口; 未装配告警降级——实际发放跨 user 域, 属后续批次）。
	// 置于事务提交之后: 端口是跨域副作用, 不应被卷入本事务的回滚语义。
	if AssistReward != nil {
		AssistReward.GrantForAssist(ctx, rec[rcols.UserId].Int64(), rec[rcols.ActivityId].Int64(),
			recordId, act[acols.RewardType].Int(), act[acols.RewardRef].Int64())
	} else {
		g.Log().Warningf(ctx, "发奖端口未装配(后续批次): 助力达标未投递 user_id=%d record_id=%d",
			rec[rcols.UserId].Int64(), recordId)
	}
	return &model.AssistHelpResult{Done: true}, nil
}

// loadAssistActivityForPlay 取助力活动并校验（启用未删 + 时间窗内）。
func loadAssistActivityForPlay(ctx context.Context, activityId int64) (gdb.Record, error) {
	cols := dao.AssistActivity.Columns()
	act, err := dao.AssistActivity.Ctx(ctx).
		Where(cols.Id, activityId).
		Where(cols.Status, 1).Where(cols.Deleted, 0).
		Where(cols.StartTime + " <= NOW()").Where(cols.EndTime + " >= NOW()").One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询助力活动失败")
	}
	if act.IsEmpty() {
		return nil, errcode.New(errcode.CodeActivityInvalid, "助力活动不在进行中")
	}
	return act, nil
}

// loadAssistRecord 取参与记录与其活动。
func loadAssistRecord(ctx context.Context, recordId int64) (gdb.Record, gdb.Record, error) {
	rcols := dao.AssistRecord.Columns()
	rec, err := dao.AssistRecord.Ctx(ctx).Where(rcols.Id, recordId).One()
	if err != nil {
		return nil, nil, gerror.Wrap(err, "查询助力记录失败")
	}
	if rec.IsEmpty() {
		return nil, nil, errcode.New(errcode.CodeActivityNotFound, "助力记录不存在")
	}
	acols := dao.AssistActivity.Columns()
	act, err := dao.AssistActivity.Ctx(ctx).Where(acols.Id, rec[rcols.ActivityId].Int64()).One()
	if err != nil {
		return nil, nil, gerror.Wrap(err, "查询助力活动失败")
	}
	if act.IsEmpty() {
		return nil, nil, errcode.New(errcode.CodeActivityNotFound, "助力活动不存在")
	}
	return rec, act, nil
}

// assistHelpers 助力人列表（昵称脱敏）。
func assistHelpers(ctx context.Context, recordId int64) ([]model.PlayHelperItem, error) {
	cols := dao.AssistHelper.Columns()
	recs, err := dao.AssistHelper.Ctx(ctx).Where(cols.RecordId, recordId).OrderAsc(cols.Id).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询助力记录失败")
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
			Nickname:  maskReviewUser(nicks[uid], 0),
			CreatedAt: r[cols.CreatedAt].String(),
		})
	}
	return out, nil
}

// collectHelperUids 取 helper 行的用户 id（去重）。
func collectHelperUids(recs gdb.Result, col string) []int64 {
	seen := map[int64]bool{}
	out := []int64{}
	for _, r := range recs {
		id := r[col].Int64()
		if id > 0 && !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	return out
}

// helperNicknames 批量取昵称（一次 IN 查询; 软删用户不展示——个保法最小化, 与评价域同口径）。
func helperNicknames(ctx context.Context, uids []int64) (map[int64]string, error) {
	out := map[int64]string{}
	if len(uids) == 0 {
		return out, nil
	}
	ucols := dao.User.Columns()
	recs, err := dao.User.Ctx(ctx).Fields(ucols.Id, ucols.Nickname).
		WhereIn(ucols.Id, uids).Where(ucols.Deleted, 0).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询用户昵称失败")
	}
	for _, r := range recs {
		out[r[ucols.Id].Int64()] = r[ucols.Nickname].String()
	}
	return out, nil
}
