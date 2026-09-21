package user

import (
	"context"

	"ecboot/api/user/v1"
	"ecboot/internal/middleware"
	"ecboot/internal/model"
	"ecboot/internal/service/user"
)

// MessageList 站内信列表（含全量未读计数）
func (c *ControllerV1) MessageList(ctx context.Context, req *v1.MessageListReq) (res *v1.MessageListRes, err error) {
	out, err := user.Messages(ctx, middleware.CtxUserIdFrom(ctx), req.IsRead,
		model.PageReq{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	res = &v1.MessageListRes{
		List:        make([]v1.MessageItem, 0, len(out.List)),
		UnreadCount: int(out.UnreadCount),
	}
	res.Total = out.Total
	for _, it := range out.List {
		res.List = append(res.List, v1.MessageItem{
			Id: fmtID(it.Id), Title: it.Title, Content: it.Content,
			BizType: it.BizType, BizNo: it.BizNo, IsRead: it.IsRead, CreatedAt: it.CreatedAt,
		})
	}
	return res, nil
}
