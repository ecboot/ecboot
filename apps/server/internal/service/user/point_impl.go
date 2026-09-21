// point_impl.go 积分账户与流水（011-member-center; 接口契约见 point.go IPointLogic）。
// 语义: 余额**可为负**（欠款语义, 如实展示不隐藏）; 每次变动同事务写流水（含变动后余额快照）;
// 获得类变动刷新 last_earned_at（滚动有效期口径: 自最后获得日起 12 个月）。
package user

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/model"
	"ecboot/internal/model/do"
)

// bizTypeConsume 下单消耗 / bizTypeRefund 退款回退 / bizTypeExpire 过期扣减。
const (
	bizTypeConsume = 3
	bizTypeRefund  = 4
	bizTypeExpire  = 9
)

// PointAccount 积分账户（FR-012）: 无账户行时返回 0 余额（不报错）。
func PointAccount(ctx context.Context, userId int64) (*model.PointAccountView, error) {
	cols := dao.PointAccount.Columns()
	rec, err := dao.PointAccount.Ctx(ctx).Where(cols.UserId, userId).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询积分账户失败")
	}
	if rec.IsEmpty() {
		return &model.PointAccountView{Balance: 0}, nil
	}
	return &model.PointAccountView{
		Balance:      rec[cols.Balance].Int(),
		LastEarnedAt: rec[cols.LastEarnedAt].String(),
	}, nil
}

// PointLogs 积分流水（FR-012）: bizType>0 筛选; 时间倒序 + 分页。
func PointLogs(ctx context.Context, userId int64, bizType int, page model.PageReq) (*model.PageResult[model.PointLogItem], error) {
	page = page.Normalized()
	cols := dao.PointLog.Columns()
	m := dao.PointLog.Ctx(ctx).Where(cols.UserId, userId)
	if bizType > 0 {
		m = m.Where(cols.BizType, bizType)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计积分流水失败")
	}
	recs, err := m.OrderDesc(cols.Id).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询积分流水失败")
	}
	list := make([]model.PointLogItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, model.PointLogItem{
			BizType:      r[cols.BizType].Int(),
			Points:       r[cols.Points].Int(),
			BalanceAfter: r[cols.BalanceAfter].Int(),
			OrderNo:      r[cols.OrderNo].String(),
			CreatedAt:    r[cols.CreatedAt].String(),
		})
	}
	return &model.PageResult[model.PointLogItem]{List: list, Total: int64(total)}, nil
}

// pointChange 积分变动核心（同事务: 确保账户行 → 变更余额 → 读回快照 → 写流水）。
func pointChange(ctx context.Context, userId int64, bizType, delta int, orderNo string, refreshEarned bool) error {
	acols := dao.PointAccount.Columns()
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cnt, err := dao.PointAccount.Ctx(ctx).Where(acols.UserId, userId).Count()
		if err != nil {
			return gerror.Wrap(err, "查询积分账户失败")
		}
		if cnt == 0 {
			if _, err = dao.PointAccount.Ctx(ctx).
				Data(do.PointAccount{UserId: userId, Balance: 0}).Insert(); err != nil {
				return gerror.Wrap(err, "初始化积分账户失败")
			}
		}
		data := do.PointAccount{Balance: gdb.Raw(fmt.Sprintf("balance + (%d)", delta))}
		if refreshEarned {
			data.LastEarnedAt = gtime.Now()
		}
		if _, err = dao.PointAccount.Ctx(ctx).
			Where(acols.UserId, userId).
			Data(data).
			Update(); err != nil {
			return gerror.Wrap(err, "更新积分余额失败")
		}
		rec, err := dao.PointAccount.Ctx(ctx).Fields(acols.Balance).
			Where(acols.UserId, userId).One()
		if err != nil {
			return gerror.Wrap(err, "读取积分余额失败")
		}
		if _, err = dao.PointLog.Ctx(ctx).Data(do.PointLog{
			UserId:       userId,
			BizType:      bizType,
			Points:       delta,
			BalanceAfter: rec[acols.Balance].Int(),
			OrderNo:      orderNo,
		}).Insert(); err != nil {
			return gerror.Wrap(err, "写积分流水失败")
		}
		return nil
	})
}

// PointEarn 获得积分（内部方法）: 刷新滚动有效期锚点。
func PointEarn(ctx context.Context, userId int64, bizType, points int, orderNo string) error {
	if points <= 0 {
		return errcode.New(errcode.CodeInvalidParam, "获得积分须为正数")
	}
	return pointChange(ctx, userId, bizType, points, orderNo, true)
}

// PointConsume 下单消耗（内部方法; 可致负——欠款语义）。
func PointConsume(ctx context.Context, userId int64, points int, orderNo string) error {
	if points <= 0 {
		return errcode.New(errcode.CodeInvalidParam, "消耗积分须为正数")
	}
	return pointChange(ctx, userId, bizTypeConsume, -points, orderNo, false)
}

// PointRefund 退款回退（内部方法）: **加回**积分（011 评审 I5 定档）——
// point_log.points 列注释的"消耗/回退负"为 V19 旧口径; 退款应退还用户已消耗的积分,
// 故本实现为正向变动（与 Consume 对称）。接口注释已同步。
func PointRefund(ctx context.Context, userId int64, points int, orderNo string) error {
	if points <= 0 {
		return errcode.New(errcode.CodeInvalidParam, "回退积分须为正数")
	}
	return pointChange(ctx, userId, bizTypeRefund, points, orderNo, false)
}

// PointExpireDormant 滚动过期清零（内部方法; 定时任务: last_earned_at 超 12 个月且余额>0）。
func PointExpireDormant(ctx context.Context) (int64, error) {
	acols := dao.PointAccount.Columns()
	recs, err := dao.PointAccount.Ctx(ctx).
		Where("last_earned_at IS NOT NULL").
		Where("last_earned_at < DATE_SUB(NOW(), INTERVAL 12 MONTH)").
		Where(acols.Balance + " > 0").
		All()
	if err != nil {
		return 0, gerror.Wrap(err, "查询过期积分账户失败")
	}
	var n int64
	for _, r := range recs {
		uid := r[acols.UserId].Int64()
		bal := r[acols.Balance].Int()
		if err = pointChange(ctx, uid, bizTypeExpire, -bal, "", false); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}
