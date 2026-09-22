// review.go 评价动作域——表: product_review（V12）。
// 规则: 一项一评（UNIQUE 兜底）; 仅已完成订单项; 追评一次限 90 天; 商家回复一次;
// **审核口径（014 裁定）**: 本批只交付 C 端 4 端点、无后台审核端点 → V1 **新评价自动通过**
// （audit_status=1）使评价即刻可用; 展示侧仍严格"只出通过"。审核字段与 Audit 能力保留,
// 待后续补后台审核面时切回"先审后发"（合规债务见 PROGRESS §五 / research D1）。
package shop

import (
	"context"

	"ecboot/internal/model"
)

// IReviewLogic 评价动作。
type IReviewLogic interface {
	// Create 提交评价（校验订单归属+已完成+未评过; 写 upsert 快照; 初始审核状态=通过, 见文件头口径）。
	Create(ctx context.Context, userId int64, in model.ReviewCreateInput) (int64, error)
	// Extra 追评（一次; 主评 90 天内）。
	Extra(ctx context.Context, userId int64, reviewId int64, content string, images []string) error
	// Reply 商家回复（一次; 后台）。
	Reply(ctx context.Context, reviewId int64, content, operator string) error
	// Audit 审核（通过/驳回）。
	Audit(ctx context.Context, reviewId int64, pass bool, operator string) error
	// MyList 我的评价。
	MyList(ctx context.Context, userId int64, page model.PageReq) (*model.PageResult[model.MyReviewItem], error)
	// ProductList 商品评价列表与汇总（FR-007/008）: 只出审核通过; score>0 时按星级筛选。
	// 汇总与列表**同源同筛**（调用方只查一次, 避免口径漂移）。
	ProductList(ctx context.Context, spuId int64, score int, page model.PageReq) (*model.ReviewSummary, *model.PageResult[model.ReviewCard], error)
	// PendingAudit 待审核列表（后台）。
	PendingAudit(ctx context.Context, page model.PageReq) (*model.PageResult[model.AdminReviewItem], error)
}
