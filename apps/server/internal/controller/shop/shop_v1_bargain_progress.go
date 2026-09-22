package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// BargainProgress 砍价进度（公开; 含帮砍列表, 昵称脱敏）
func (c *ControllerV1) BargainProgress(ctx context.Context, req *v1.BargainProgressReq) (res *v1.BargainProgressRes, err error) {
	recordId, err := parseID(req.RecordId)
	if err != nil {
		return nil, err
	}
	p, err := shop.NewBargainLogic().Progress(ctx, recordId)
	if err != nil {
		return nil, err
	}
	res = &v1.BargainProgressRes{
		RecordId: fmtID(p.RecordId), SkuId: fmtID(p.SkuId),
		OriginalPrice: p.OriginalPrice, CurrentPrice: p.CurrentPrice, FloorPrice: p.FloorPrice,
		CutCount: p.CutCount, Status: p.Status, ExpireTime: p.ExpireTime, OrderNo: p.OrderNo,
		Helpers: make([]v1.BargainHelper, 0, len(p.Helpers)),
	}
	for _, h := range p.Helpers {
		res.Helpers = append(res.Helpers, v1.BargainHelper{
			UserId: fmtID(h.UserId), Nickname: h.Nickname, CutAmount: h.Amount, CreatedAt: h.CreatedAt,
		})
	}
	return res, nil
}
