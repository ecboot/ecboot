package shop

import (
	"context"

	"ecboot/api/shop/v1"
	"ecboot/internal/service/shop"
)

// AssistHelp 助力（会员; 一人一助力）
func (c *ControllerV1) AssistHelp(ctx context.Context, req *v1.AssistHelpReq) (res *v1.AssistHelpRes, err error) {
	userId, err := requireMember(ctx)
	if err != nil {
		return nil, err
	}
	recordId, err := parseID(req.RecordId)
	if err != nil {
		return nil, err
	}
	out, err := shop.NewAssistLogic().Help(ctx, userId, recordId)
	if err != nil {
		return nil, err
	}
	return &v1.AssistHelpRes{Success: true, Done: out.Done}, nil
}
