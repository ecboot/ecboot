package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 提交评价（一项一评; 004 契约, 复用 V12 表）
	ReviewCreateReq struct {
		g.Meta      `path:"/reviews" method:"POST" summary:"提交评价"`
		OrderItemId string   `json:"orderItemId" v:"required" dc:"订单项ID"`
		Score       int      `json:"score" v:"required|between:1,5" dc:"评分1-5"`
		Content     string   `json:"content" dc:"评价内容"`
		Images      []string `json:"images" dc:"评价图片"`
		IsAnonymous bool     `json:"isAnonymous" dc:"匿名"`
	}
	ReviewCreateRes struct {
		ReviewId string `json:"reviewId" dc:"评价ID"`
	}

	// 追评（一次, 90 天内）
	ReviewExtraReq struct {
		g.Meta   `path:"/reviews/{reviewId}/extra" method:"POST" summary:"追加评价"`
		ReviewId string   `json:"reviewId" v:"required" dc:"评价ID"`
		Content  string   `json:"content" v:"required" dc:"追评内容"`
		Images   []string `json:"images" dc:"追评图片"`
	}
	ReviewExtraRes struct {
		Success bool `json:"success"`
	}

	// 我的评价
	MyReviewListReq struct {
		g.Meta `path:"/reviews/mine" method:"GET" summary:"我的评价"`
		PageReq
	}
	MyReviewItem struct {
		ReviewId    string `json:"reviewId"`
		SpuName     string `json:"spuName"`
		Score       int    `json:"score"`
		Content     string `json:"content"`
		AuditStatus int    `json:"auditStatus" dc:"0待审 1通过 2驳回"`
		Extra       string `json:"extra" dc:"追评"`
		Reply       string `json:"reply" dc:"商家回复"`
		CreatedAt   string `json:"createdAt"`
	}
	MyReviewListRes struct {
		PageRes
		List []MyReviewItem `json:"list"`
	}
)
