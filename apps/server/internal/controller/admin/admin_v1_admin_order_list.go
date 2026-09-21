package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/model"
	"ecboot/internal/service/shop"
)

// AdminOrderList 后台订单列表（状态/订单号/用户/时间筛选 + 分页）
func (c *ControllerV1) AdminOrderList(ctx context.Context, req *v1.AdminOrderListReq) (res *v1.AdminOrderListRes, err error) {
	out, err := shop.NewOrderLogic().AdminList(ctx, model.AdminOrderQuery{
		Status: req.Status, OrderNo: req.OrderNo, UserKeyword: req.UserKeyword,
		StartTime: req.StartTime, EndTime: req.EndTime,
		PageReq: model.PageReq{Page: req.Page, PageSize: req.PageSize},
	})
	if err != nil {
		return nil, err
	}
	res = &v1.AdminOrderListRes{List: make([]v1.AdminOrderItem, 0, len(out.List))}
	res.Total = out.Total
	for _, it := range out.List {
		row := v1.AdminOrderItem{
			OrderNo: it.OrderNo, UserId: it.UserId, Status: it.Status,
			PayAmount: it.PayAmount, CreatedAt: it.CreatedAt,
			Items: make([]v1.AdminOrderItemBrief, 0, len(it.Items)),
		}
		for _, li := range it.Items {
			row.Items = append(row.Items, v1.AdminOrderItemBrief{
				SpuName: li.SpuName, SkuSpecs: li.SkuSpecs, Image: li.Image,
				Quantity: li.Quantity, Price: li.Price,
			})
		}
		res.List = append(res.List, row)
	}
	return res, nil
}
