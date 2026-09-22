// review_impl.go 评价域实现（014-review / 批次 08）。
// 契约: 接口签名见 review.go; 口径与规则见 specs/014-review/{spec,data-model}.md。
//
// 三条硬约束（沿用 012/013 修复轮确立的惯例）:
//  1. "一项一评"依赖表上既有 **uk_order_item 唯一索引**兜底, 且唯一键冲突（1062）**必须转业务码**,
//     不得裸抛系统错误（友好前置校验负责给出准确提示, 唯一索引负责竞态兜底）。
//  2. 追评用**条件更新**（extra_content=” 且 90 天内）, 不写"读-判断-写"的竞态窗口。
//  3. 展示侧只出 audit_status=1; 列表与汇总**同源同筛**（避免"列表 3 条而汇总 100 条"）。
package shop

import (
	"context"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/model"
	"ecboot/internal/model/do"
)

// 审核状态（product_review.audit_status, 000012 列注释）。
const (
	reviewAuditPending = 0 // 待审
	reviewAuditPassed  = 1 // 通过
	reviewAuditReject  = 2 // 驳回
)

// 追评时效（天）——既有契约注释"追评一次限 90 天"。
const reviewExtraDays = 90

// ReviewLogicImpl IReviewLogic 实现。
type ReviewLogicImpl struct{}

func NewReviewLogic() *ReviewLogicImpl { return &ReviewLogicImpl{} }

// isDupKey 唯一键冲突（MySQL 1062）。
func isDupKey(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "1062") || strings.Contains(s, "Duplicate entry")
}

// Create 提交评价（FR-001~004）: 校验归属/已完成/未评过 → 落快照 → 初始审核状态=通过（V1 自动通过）。
func (i *ReviewLogicImpl) Create(
	ctx context.Context, userId int64, in model.ReviewCreateInput,
) (int64, error) {
	if in.Score < 1 || in.Score > 5 {
		return 0, errcode.New(errcode.CodeInvalidParam, "评分须为 1~5")
	}
	item, order, err := reviewItemAndOrder(ctx, in.OrderItemId, userId)
	if err != nil {
		return 0, err
	}
	ocols := dao.TradeOrder.Columns()
	if order[ocols.Status].Int() != orderStatusCompleted {
		// 订单未完成/已取消 → 状态不允许（40006）
		return 0, errcode.New(errcode.CodeStatusNotAllowed, "仅已完成的订单可评价")
	}

	cols := dao.ProductReview.Columns()
	exist, err := dao.ProductReview.Ctx(ctx).
		Where(cols.OrderItemId, in.OrderItemId).
		Where(cols.Deleted, 0).Count()
	if err != nil {
		return 0, gerror.Wrap(err, "查询评价失败")
	}
	if exist > 0 {
		return 0, errcode.New(errcode.CodeAlreadyReviewed, "该订单项已评价")
	}

	icols := dao.TradeOrderItem.Columns()
	res, err := dao.ProductReview.Ctx(ctx).Data(do.ProductReview{
		OrderItemId: in.OrderItemId,
		OrderNo:     order[ocols.OrderNo].String(),
		UserId:      userId,
		SpuId:       item[icols.SpuId].Int64(),
		SkuId:       item[icols.SkuId].Int64(),
		SpuName:     item[icols.SpuName].String(),  // 下单时快照
		SkuSpecs:    item[icols.SkuSpecs].String(), // JSON 快照
		Score:       in.Score,
		Content:     in.Content,
		Images:      mustJSON(reviewImages(in.Images)),
		IsAnonymous: boolFlag(in.IsAnonymous),
		AuditStatus: reviewAuditPassed, // V1 自动通过（用户裁定 B; 见 review.go 文件头口径）
	}).Insert()
	if err != nil {
		if isDupKey(err) { // 并发下唯一索引兜底: 转业务码, 不裸抛
			return 0, errcode.New(errcode.CodeAlreadyReviewed, "该订单项已评价")
		}
		return 0, gerror.Wrap(err, "提交评价失败")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, gerror.Wrap(err, "读取评价ID失败")
	}
	return id, nil
}

// reviewItemAndOrder 取订单项与其订单并校验归属（他人资源一律按"不存在"处理）。
func reviewItemAndOrder(ctx context.Context, orderItemId, userId int64) (gdb.Record, gdb.Record, error) {
	icols, ocols := dao.TradeOrderItem.Columns(), dao.TradeOrder.Columns()
	item, err := dao.TradeOrderItem.Ctx(ctx).Where(icols.Id, orderItemId).One()
	if err != nil {
		return nil, nil, gerror.Wrap(err, "查询订单项失败")
	}
	if item.IsEmpty() {
		return nil, nil, errcode.New(errcode.CodeNotFound, "订单项不存在")
	}
	order, err := dao.TradeOrder.Ctx(ctx).Where(ocols.Id, item[icols.OrderId].Int64()).One()
	if err != nil {
		return nil, nil, gerror.Wrap(err, "查询订单失败")
	}
	if order.IsEmpty() || order[ocols.UserId].Int64() != userId {
		return nil, nil, errcode.New(errcode.CodeNotFound, "订单项不存在")
	}
	return item, order, nil
}

// reviewImages 图片数组：空时给空数组（JSON 列存 [] 而非 NULL）。
func reviewImages(images []string) []string {
	if images == nil {
		return []string{}
	}
	return images
}

// boolFlag 布尔 → 0/1。
func boolFlag(b bool) int {
	if b {
		return 1
	}
	return 0
}

// maskReviewUser 评价人展示名：匿名给固定占位, 非匿名给脱敏昵称。
// 说明（research D2）: `maskNickname` 位于兄弟域 `service/user`（`invite_impl.go`）, 按分层契约
// **不可 import**, 故在 shop 域内保留等价实现（口径与之一致: 首字符 + 掩码, 空昵称占位）。
func maskReviewUser(nickname string, anonymous int) string {
	if anonymous == 1 {
		return "匿名用户"
	}
	n := strings.TrimSpace(nickname)
	if n == "" {
		return "用户****"
	}
	r := []rune(n)
	if len(r) <= 1 {
		return string(r) + "*"
	}
	return string(r[0]) + strings.Repeat("*", len(r)-1)
}

// parseSpecs sku_specs JSON → 规格表（展示用; 失败给空表）。
func parseSpecs(raw string) map[string]string {
	out := map[string]string{}
	if strings.TrimSpace(raw) == "" {
		return out
	}
	_ = gconv.Struct(raw, &out)
	return out
}

// parseImages images JSON → 图片数组（空给空数组）。
func parseImages(raw string) []string {
	out := []string{}
	if strings.TrimSpace(raw) == "" {
		return out
	}
	_ = gconv.Struct(raw, &out)
	return out
}

// Extra 追评（FR-006）: 本人主评 + 未追评 + 主评 90 天内, 一次为限。
// 实现用**条件更新**（`extra_content=”` 且 `created_at >= NOW()-90d`）而不是"读-判断-写"——
// 并发两次追评只有一次能把空串写成内容（affected=0 时回查区分"已追评/超期"以给准确提示）。
//
// 注（契约缺口, 已记账）: api 的 `ReviewExtraReq.Images` **既无存储列（表无 extra_images）也无出参字段**
// （`ReviewItem`/`MyReviewItem` 都没有 extraImages）——属契约自相矛盾（疑似从提交评价的请求拷贝而来）。
// 本批**不为其新增列**（那会写入永远无法展示的数据）, 该字段当前被忽略; 若要支持追评图片,
// 需同时补「表列 + 两个出参字段」, 属 API 契约变更, 待裁定（PROGRESS §五 已记账）。
func (i *ReviewLogicImpl) Extra(ctx context.Context, userId int64, reviewId int64, content string, images []string) error {
	_ = images // 见上方契约缺口说明: 暂不持久化
	if strings.TrimSpace(content) == "" {
		return errcode.New(errcode.CodeInvalidParam, "追评内容必填")
	}
	cols := dao.ProductReview.Columns()
	rec, err := dao.ProductReview.Ctx(ctx).
		Where(cols.Id, reviewId).Where(cols.Deleted, 0).One()
	if err != nil {
		return gerror.Wrap(err, "查询评价失败")
	}
	if rec.IsEmpty() || rec[cols.UserId].Int64() != userId {
		return errcode.New(errcode.CodeNotFound, "评价不存在") // 他人资源按不存在处理
	}

	res, err := dao.ProductReview.Ctx(ctx).
		Where(cols.Id, reviewId).
		Where(cols.UserId, userId).
		Where(cols.ExtraContent, ""). // 未追评
		Where("created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", reviewExtraDays).
		Data(do.ProductReview{ExtraContent: content, ExtraTime: gtime.Now()}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "追评失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		// 条件未命中: 回查区分"已追评"与"超期"（给用户可行动的提示）
		cur, ce := dao.ProductReview.Ctx(ctx).Where(cols.Id, reviewId).One()
		if ce != nil || cur.IsEmpty() {
			return errcode.New(errcode.CodeExtraReviewDenied, "不可追评")
		}
		if strings.TrimSpace(cur[cols.ExtraContent].String()) != "" {
			return errcode.New(errcode.CodeExtraReviewDenied, "该评价已追评过")
		}
		return errcode.New(errcode.CodeExtraReviewDenied, "已超过 90 天追评期限")
	}
	return nil
}

// MyList 我的评价（FR-005/010）: 仅本人; **不筛审核状态**（待审/驳回也要让本人看到状态）。
func (i *ReviewLogicImpl) MyList(
	ctx context.Context, userId int64, page model.PageReq,
) (*model.PageResult[model.MyReviewItem], error) {
	page = page.Normalized()
	cols := dao.ProductReview.Columns()
	base := func() *gdb.Model {
		return dao.ProductReview.Ctx(ctx).Where(cols.UserId, userId).Where(cols.Deleted, 0)
	}
	total, err := base().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计我的评价失败")
	}
	recs, err := base().OrderDesc(cols.Id).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询我的评价失败")
	}
	list := make([]model.MyReviewItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, model.MyReviewItem{
			ReviewId:    r[cols.Id].Int64(),
			SpuName:     r[cols.SpuName].String(),
			Score:       r[cols.Score].Int(),
			Content:     r[cols.Content].String(),
			Extra:       r[cols.ExtraContent].String(),
			Reply:       r[cols.ReplyContent].String(),
			AuditStatus: r[cols.AuditStatus].Int(),
			CreatedAt:   r[cols.CreatedAt].String(),
		})
	}
	return &model.PageResult[model.MyReviewItem]{List: list, Total: int64(total)}, nil
}

// ProductList 商品评价列表与汇总（FR-007/008/009）: 只出审核通过; score>0 时按星级筛选。
// 汇总与列表**同源同筛**（同一个构造器 base()），避免"列表 3 条而汇总 100 条"的口径漂移。
func (i *ReviewLogicImpl) ProductList(
	ctx context.Context, spuId int64, score int, page model.PageReq,
) (*model.ReviewSummary, *model.PageResult[model.ReviewCard], error) {
	page = page.Normalized()
	cols := dao.ProductReview.Columns()

	base := func() *gdb.Model {
		m := dao.ProductReview.Ctx(ctx).
			Where(cols.SpuId, spuId).
			Where(cols.AuditStatus, reviewAuditPassed). // 只出审核通过（V1 自动通过后即刻可见）
			Where(cols.Deleted, 0)
		if score > 0 {
			m = m.Where(cols.Score, score)
		}
		return m
	}

	total, err := base().Count()
	if err != nil {
		return nil, nil, gerror.Wrap(err, "统计商品评价失败")
	}
	// 汇总: 总数为 0 时不给平均分（避免 0/0）
	avg := "0.0"
	dist := map[string]int{}
	if total > 0 {
		sum, se := base().Sum(cols.Score)
		if se != nil {
			return nil, nil, gerror.Wrap(se, "统计评分失败")
		}
		avg = strconv.FormatFloat(sum/float64(total), 'f', 1, 64)
		distRecs, de := base().Fields(cols.Score + ", COUNT(*) AS c").Group(cols.Score).All()
		if de != nil {
			return nil, nil, gerror.Wrap(de, "统计评分分布失败")
		}
		for _, r := range distRecs {
			dist[strconv.Itoa(r[cols.Score].Int())] = r["c"].Int()
		}
	}

	recs, err := base().OrderDesc(cols.Id).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, nil, gerror.Wrap(err, "查询商品评价失败")
	}

	// 评价人昵称批量取（一次 IN 查询, 避免 N+1）; 读 user 表**跟随既有先例**（order_mgmt_impl 已如此）
	uids := make([]int64, 0, len(recs))
	for _, r := range recs {
		uids = append(uids, r[cols.UserId].Int64())
	}
	nick := map[int64]string{}
	if len(uids) > 0 {
		ucols := dao.User.Columns()
		urecs, ue := dao.User.Ctx(ctx).Fields(ucols.Id, ucols.Nickname).WhereIn(ucols.Id, uids).All()
		if ue != nil {
			return nil, nil, gerror.Wrap(ue, "查询评价人失败")
		}
		for _, u := range urecs {
			nick[u[ucols.Id].Int64()] = u[ucols.Nickname].String()
		}
	}

	list := make([]model.ReviewCard, 0, len(recs))
	for _, r := range recs {
		uid := r[cols.UserId].Int64()
		anon := r[cols.IsAnonymous].Int()
		list = append(list, model.ReviewCard{
			ReviewId:  r[cols.Id].Int64(),
			User:      maskReviewUser(nick[uid], anon),
			Score:     r[cols.Score].Int(),
			Content:   r[cols.Content].String(),
			Images:    parseImages(r[cols.Images].String()),
			Specs:     parseSpecs(r[cols.SkuSpecs].String()),
			Reply:     r[cols.ReplyContent].String(),
			Extra:     r[cols.ExtraContent].String(),
			CreatedAt: r[cols.CreatedAt].String(),
		})
	}
	summary := &model.ReviewSummary{Avg: avg, Distribution: dist, Total: total}
	return summary, &model.PageResult[model.ReviewCard]{List: list, Total: int64(total)}, nil
}
