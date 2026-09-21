// review.go 评价动作域——表: product_review（V12）。
// 规则: 一项一评（UNIQUE 兜底）; 仅已完成订单项; 追评一次限 90 天; 商家回复一次; 审核后展示。
package shop

import "context"

// IReviewLogic 评价动作。
type IReviewLogic interface {
	// Create 提交评价（校验订单归属+已完成+未评过; 内容审核状态默认待审）。
	Create(ctx context.Context, userId int64, in ReviewCreateInput) (int64, error)
	// Extra 追评（一次; 主评 90 天内）。
	Extra(ctx context.Context, userId int64, reviewId int64, content string, images []string) error
	// Reply 商家回复（一次; 后台）。
	Reply(ctx context.Context, reviewId int64, content, operator string) error
	// Audit 审核（通过/驳回）。
	Audit(ctx context.Context, reviewId int64, pass bool, operator string) error
	// MyList 我的评价。
	MyList(ctx context.Context, userId int64, page PageQuery) (*PageResult[MyReviewItem], error)
	// PendingAudit 待审核列表（后台）。
	PendingAudit(ctx context.Context, page PageQuery) (*PageResult[AdminReviewItem], error)
}

type ReviewCreateInput struct {
	OrderItemId int64
	Score       int
	Content     string
	Images      []string
	IsAnonymous bool
}

type MyReviewItem struct {
	ReviewId    int64  `json:"reviewId"`
	SpuName     string `json:"spuName"`
	Score       int    `json:"score"`
	Content     string `json:"content"`
	Extra       string `json:"extra"`
	Reply       string `json:"reply"`
	AuditStatus int    `json:"auditStatus" dc:"0待审 1通过 2驳回"`
	CreatedAt   string `json:"createdAt"`
}

type AdminReviewItem struct {
	ReviewId    int64  `json:"reviewId"`
	OrderNo     string `json:"orderNo"`
	UserMasked  string `json:"user" dc:"脱敏"`
	SpuName     string `json:"spuName" dc:"快照"`
	Score       int    `json:"score"`
	Content     string `json:"content"`
	AuditStatus int    `json:"auditStatus"`
	CreatedAt   string `json:"createdAt"`
}
