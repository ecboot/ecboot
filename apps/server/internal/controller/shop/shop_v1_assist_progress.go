package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// AssistProgress 助力进度（公开; 含助力人列表, 昵称脱敏）
func (c *ControllerV1) AssistProgress(ctx context.Context, req *v1.AssistProgressReq) (res *v1.AssistProgressRes, err error) {
	recordId, err := parseID(req.RecordId)
	if err != nil {
		return nil, err
	}
	p, err := shop.NewAssistLogic().Progress(ctx, recordId)
	if err != nil {
		return nil, err
	}
	res = &v1.AssistProgressRes{
		RecordId: fmtID(p.RecordId), ActivityId: fmtID(p.ActivityId),
		HelperCount: p.HelperCount, RequiredCount: p.RequiredCount, Status: p.Status,
		FinishTime: p.FinishTime, Helpers: make([]v1.AssistHelper, 0, len(p.Helpers)),
	}
	for _, h := range p.Helpers {
		res.Helpers = append(res.Helpers, v1.AssistHelper{
			UserId: fmtID(h.UserId), Nickname: h.Nickname, CreatedAt: h.CreatedAt,
		})
	}
	return res, nil
}
