// review.go 评价动作域——表: product_review（V12）。
// 规则: 一项一评（UNIQUE 兜底）; 仅已完成订单项; 追评一次限 90 天; 商家回复一次; 审核后展示。
package shop

import (
	"context"

	"ecboot/internal/model"
)

// IReviewLogic 评价动作。
type IReviewLogic interface {
	// Create 提交评价（校验订单归属+已完成+未评过; 内容审核状态默认待审）。
	Create(ctx context.Context, userId int64, in model.ReviewCreateInput) (int64, error)
	// Extra 追评（一次; 主评 90 天内）。
	Extra(ctx context.Context, userId int64, reviewId int64, content string, images []string) error
	// Reply 商家回复（一次; 后台）。
	Reply(ctx context.Context, reviewId int64, content, operator string) error
	// Audit 审核（通过/驳回）。
	Audit(ctx context.Context, reviewId int64, pass bool, operator string) error
	// MyList 我的评价。
	MyList(ctx context.Context, userId int64, page model.PageReq) (*model.PageResult[model.MyReviewItem], error)
	// PendingAudit 待审核列表（后台）。
	PendingAudit(ctx context.Context, page model.PageReq) (*model.PageResult[model.AdminReviewItem], error)
}
