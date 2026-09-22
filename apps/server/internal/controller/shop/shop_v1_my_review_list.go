package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// MyReviewList 我的评价（含审核状态/追评/商家回复）
func (c *ControllerV1) MyReviewList(ctx context.Context, req *v1.MyReviewListReq) (res *v1.MyReviewListRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	out, err := shop.NewReviewLogic().MyList(ctx, userId, model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.MyReviewListRes{List: make([]v1.MyReviewItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.MyReviewItem{
			ReviewId:    fmtID(it.ReviewId),
			SpuName:     it.SpuName,
			Score:       it.Score,
			Content:     it.Content,
			AuditStatus: it.AuditStatus,
			Extra:       it.Extra,
			Reply:       it.Reply,
			CreatedAt:   it.CreatedAt,
		})
	}
	return res, nil
}
