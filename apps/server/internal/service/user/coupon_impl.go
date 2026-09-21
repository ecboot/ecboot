// coupon_impl.go 会员优惠券（012-trade-wiring; 接口契约见 coupon.go IUserCouponLogic）。
// 语义:
//   - 领取**同事务**: 限领校验（个人已领 < per_limit）→ 防超发（条件更新 received_count+1 且 < total_count）
//     → 计算有效期落 user_coupon（valid_type 1 固定区间 / 2 领取后 N 天）。
//   - 过期**惰性判定**: 查询时按 expire_time 与当前时间比较（不改表、不依赖定时任务）;
//     status=1 且已过期者在"已过期"视图中展示, 不在"未使用"中出现。
//   - 核销绑单号（Consume）; 退回解绑且**有效期不变**（ReturnBack）。
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
)

// 券状态（user_coupon.status; 3 为惰性判定视图, 不落库）。
const (
	couponStatusUnused   = 1
	couponStatusUsed     = 2
	couponStatusExpired  = 3
	couponStatusReturned = 4
)

// 模板有效期方式。
const (
	validTypeFixedRange = 1
	validTypeAfterDays  = 2
)

// AvailableTemplates 可领模板（FR-013）: 过滤停发/领完; 每张标注个人是否可领（限领未满）。
func AvailableTemplates(ctx context.Context, userId int64, page model.PageReq) (*model.PageResult[model.AvailableCoupon], error) {
	page = page.Normalized()
	cols := dao.Coupon.Columns()
	m := dao.Coupon.Ctx(ctx).
		Where(cols.Status, 1).
		Where(cols.Deleted, 0)
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计券模板失败")
	}
	recs, err := m.OrderDesc(cols.Id).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询券模板失败")
	}
	list := make([]model.AvailableCoupon, 0, len(recs))
	for _, r := range recs {
		cid := r[cols.Id].Int64()
		// 个人已领数
		mine, err := dao.UserCoupon.Ctx(ctx).
			Where(dao.UserCoupon.Columns().UserId, userId).
			Where(dao.UserCoupon.Columns().CouponId, cid).
			Count()
		if err != nil {
			return nil, gerror.Wrap(err, "查询已领券失败")
		}
		perLimit := r[cols.PerLimit].Int()
		received := r[cols.ReceivedCount].Int()
		totalCount := r[cols.TotalCount].Int()
		can := mine < perLimit && (totalCount == 0 || received < totalCount) // total_count=0 视为不限量
		list = append(list, model.AvailableCoupon{
			CouponId:   cid,
			Name:       r[cols.Name].String(),
			Type:       r[cols.Type].Int(),
			Threshold:  r[cols.ThresholdAmount].String(),
			Discount:   r[cols.DiscountAmount].String(),
			ValidDesc:  couponValidDesc(r),
			CanReceive: can,
		})
	}
	return &model.PageResult[model.AvailableCoupon]{List: list, Total: int64(total)}, nil
}

// couponValidDesc 有效期描述（列表展示）。
func couponValidDesc(r gdb.Record) string {
	cols := dao.Coupon.Columns()
	if r[cols.ValidType].Int() == validTypeAfterDays {
		return fmt.Sprintf("领取后 %d 天有效", r[cols.ValidDays].Int())
	}
	return r[cols.ValidStartAt].String() + " 至 " + r[cols.ValidEndAt].String()
}

// Receive 领取（FR-013）: 同事务「限领校验 → 防超发条件更新 → 落券」; 超限/领完 → 50001。
func Receive(ctx context.Context, userId, couponId int64) (int64, error) {
	var id int64
	ccols, ucols := dao.Coupon.Columns(), dao.UserCoupon.Columns()
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 评审 I4: 先对模板行加**悲观锁**（LockUpdate）串行化同模板的并发领取;
		// 否则两个并发事务都读到 mine=0 都通过限领（尤其 total_count=0 不限量时整段防超发被跳过, 完全无锁）。
		rec, e := dao.Coupon.Ctx(ctx).
			Where(ccols.Id, couponId).
			Where(ccols.Status, 1).
			Where(ccols.Deleted, 0).
			LockUpdate().
			One()
		if e != nil {
			return gerror.Wrap(e, "查询券模板失败")
		}
		if rec.IsEmpty() {
			return errcode.New(errcode.CodeCouponSoldOut, "优惠券不存在或已停发")
		}
		// 限领校验: 锁定读（RR 快照下普通 SELECT 读旧值, 须 FOR UPDATE 才见最新已提交行）
		var mine int
		mv, e := dao.UserCoupon.Ctx(ctx).
			Fields("COUNT(*) AS c").
			Where(ucols.UserId, userId).
			Where(ucols.CouponId, couponId).
			LockUpdate().
			One()
		if e != nil {
			return gerror.Wrap(e, "查询已领券失败")
		}
		if mv != nil {
			mine = mv["c"].Int()
		}
		if mine >= rec[ccols.PerLimit].Int() {
			return errcode.New(errcode.CodeCouponSoldOut, "已达个人限领数量")
		}
		// 防超发（条件更新; total_count=0 视为不限量）
		totalCount := rec[ccols.TotalCount].Int()
		if totalCount > 0 {
			res, e := dao.Coupon.Ctx(ctx).
				Where(ccols.Id, couponId).
				Where(ccols.ReceivedCount+" < ", totalCount).
				Data(g.Map{ccols.ReceivedCount: gdb.Raw("received_count + 1")}).
				Update()
			if e != nil {
				return gerror.Wrap(e, "领取失败")
			}
			if n, _ := res.RowsAffected(); n == 0 {
				return errcode.New(errcode.CodeCouponSoldOut, "优惠券已领完")
			}
		}
		// 有效期计算
		expire := couponExpireTime(rec)
		res, e := dao.UserCoupon.Ctx(ctx).Data(g.Map{
			ucols.UserId:     userId,
			ucols.CouponId:   couponId,
			ucols.Status:     couponStatusUnused,
			ucols.ExpireTime: expire,
		}).Insert()
		if e != nil {
			return gerror.Wrap(e, "发券失败")
		}
		nid, e := res.LastInsertId()
		if e != nil {
			return gerror.Wrap(e, "读取券ID失败")
		}
		id = nid
		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

// couponExpireTime 按模板有效期方式计算到期时间。
func couponExpireTime(rec gdb.Record) *gtime.Time {
	cols := dao.Coupon.Columns()
	if rec[cols.ValidType].Int() == validTypeAfterDays {
		days := rec[cols.ValidDays].Int()
		if days <= 0 {
			days = 30
		}
		return gtime.Now().AddDate(0, 0, days)
	}
	// 固定区间: 取模板结束时间（模板开始时间不单独落券——券在区间内均可用, 由模板端控制发放）
	if t, err := gtime.StrToTime(rec[cols.ValidEndAt].String()); err == nil {
		return t
	}
	return gtime.Now().AddDate(0, 0, 30)
}

// Mine 我的券（FR-013）: status 1未使用/2已使用/3已过期/4已退回; 过期**惰性判定**。
func Mine(ctx context.Context, userId int64, status int, page model.PageReq) (*model.PageResult[model.MyCouponItem], error) {
	page = page.Normalized()
	ucols, ccols := dao.UserCoupon.Columns(), dao.Coupon.Columns()
	m := dao.UserCoupon.Ctx(ctx).As("uc").
		LeftJoin(dao.Coupon.Table()+" c", "c.id=uc.coupon_id").
		Where("uc."+ucols.UserId, userId)
	switch status {
	case couponStatusUnused:
		m = m.Where("uc."+ucols.Status, couponStatusUnused).Where("uc." + ucols.ExpireTime + " > NOW()")
	case couponStatusUsed:
		m = m.Where("uc."+ucols.Status, couponStatusUsed)
	case couponStatusExpired:
		// 惰性: status=1（未使用）但已过期 → 视为已过期（不改表）
		m = m.Where("uc."+ucols.Status, couponStatusUnused).Where("uc." + ucols.ExpireTime + " <= NOW()")
	case couponStatusReturned:
		m = m.Where("uc."+ucols.Status, couponStatusReturned)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计我的券失败")
	}
	recs, err := m.Fields("uc.id, uc.status, uc.expire_time, c."+ccols.Name+", c."+ccols.ThresholdAmount+", c."+ccols.DiscountAmount).
		OrderDesc("uc."+ucols.Id).
		Page(page.Page, page.PageSize).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询我的券失败")
	}
	list := make([]model.MyCouponItem, 0, len(recs))
	for _, r := range recs {
		st := r["status"].Int()
		if st == couponStatusUnused && r["expire_time"].GTime() != nil &&
			!r["expire_time"].GTime().After(gtime.Now()) {
			st = couponStatusExpired // 惰性判定展示态
		}
		list = append(list, model.MyCouponItem{
			UserCouponId: r["id"].Int64(),
			Name:         r["name"].String(),
			Threshold:    r["threshold_amount"].String(),
			Discount:     r["discount_amount"].String(),
			ExpireTime:   r["expire_time"].String(),
			Status:       st,
		})
	}
	return &model.PageResult[model.MyCouponItem]{List: list, Total: int64(total)}, nil
}

// UsableForOrder 结算可用券匹配（FR-014）: 未使用 + 未过期 + 门槛 <= 商品金额 → 按抵扣降序。
func UsableForOrder(ctx context.Context, userId int64, goodsAmount string) ([]model.UsableCouponItem, error) {
	ucols, ccols := dao.UserCoupon.Columns(), dao.Coupon.Columns()
	recs, err := dao.UserCoupon.Ctx(ctx).As("uc").
		LeftJoin(dao.Coupon.Table()+" c", "c.id=uc.coupon_id").
		Where("uc."+ucols.UserId, userId).
		Where("uc."+ucols.Status, couponStatusUnused).
		Where("uc."+ucols.ExpireTime+" > NOW()").
		Where("c."+ccols.ThresholdAmount+" <= ?", goodsAmount).
		Where("c."+ccols.Deleted, 0).
		Fields("uc.id, c." + ccols.Name + ", c." + ccols.DiscountAmount).
		OrderDesc("c." + ccols.DiscountAmount).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "匹配可用券失败")
	}
	out := make([]model.UsableCouponItem, 0, len(recs))
	for _, r := range recs {
		out = append(out, model.UsableCouponItem{
			UserCouponId: r["id"].Int64(),
			Name:         r["name"].String(),
			Discount:     r["discount_amount"].String(),
		})
	}
	return out, nil
}

// Consume 核销（FR-014）: 下单事务内调用——置已使用 + 绑单号（条件更新, 幂等）。
func Consume(ctx context.Context, userId, userCouponId int64, orderNo string) error {
	ucols := dao.UserCoupon.Columns()
	_, err := dao.UserCoupon.Ctx(ctx).
		Where(ucols.Id, userCouponId).
		Where(ucols.UserId, userId).
		Where(ucols.Status, couponStatusUnused).
		Data(g.Map{
			ucols.Status:   couponStatusUsed,
			ucols.OrderNo:  orderNo,
			ucols.UsedTime: gtime.Now(),
		}).Update()
	if err != nil {
		return gerror.Wrap(err, "核销优惠券失败")
	}
	return nil
}

// ReturnBack 退回（FR-014）: 订单取消——置未使用 + 解绑单号; **有效期不变**。
func ReturnBack(ctx context.Context, userCouponId int64) error {
	ucols := dao.UserCoupon.Columns()
	_, err := dao.UserCoupon.Ctx(ctx).
		Where(ucols.Id, userCouponId).
		Where(ucols.Status, couponStatusUsed).
		Data(g.Map{
			ucols.Status:   couponStatusUnused,
			ucols.OrderNo:  "",
			ucols.UsedTime: nil,
		}).Update()
	if err != nil {
		return gerror.Wrap(err, "退回优惠券失败")
	}
	return nil
}
