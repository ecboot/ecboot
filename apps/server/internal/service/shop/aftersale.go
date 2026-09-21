// aftersale.go 售后域——表: after_sale_order（V8/V23 数量/币种）。
// 规则: 按订单项粒度（quantity ≤ 行数量）; 状态机:
// 10待审核→(同意)仅退款30/退货退款20→确认收货30→渠道退款40→成功50; 10→90拒绝; 未终态→91撤销;
// 退款渠道幂等: out_refund_no = after_sale_no; 完成回补库存(按 quantity)+更新订单 refund_status+佣金冲销事件。
package shop

import (
	"context"

	"ecboot/internal/model"
)

// IAfterSaleLogic 售后。
type IAfterSaleLogic interface {
	// Apply 申请（校验订单已完成/项未超额; 创建 10）。
	Apply(ctx context.Context, userId int64, in model.AfterSaleApplyInput) (string, error)
	List(ctx context.Context, userId int64, status int, page model.PageReq) (*model.PageResult[model.AfterSaleSummary], error)
	Detail(ctx context.Context, userId int64, afterSaleNo string) (*model.AfterSaleDetail, error)
	Cancel(ctx context.Context, userId int64, afterSaleNo string) error
	// SubmitReturn 填写寄回单号（type=2 且状态 20）。
	SubmitReturn(ctx context.Context, userId int64, afterSaleNo, logisticsNo string) error

	// ---- 后台 ----
	AdminList(ctx context.Context, status, saleType int, page model.PageReq) (*model.PageResult[model.AfterSaleSummary], error)
	AdminDetail(ctx context.Context, afterSaleNo string) (*model.AfterSaleDetail, error)
	Approve(ctx context.Context, afterSaleNo, operator string) error
	Reject(ctx context.Context, afterSaleNo, reason, operator string) error
	// ConfirmReceipt 退货确认收货（→30 待退款）。
	ConfirmReceipt(ctx context.Context, afterSaleNo, operator string) error
	// RetryRefund 退款重试（渠道失败后）。
	RetryRefund(ctx context.Context, afterSaleNo, operator string) error
}
