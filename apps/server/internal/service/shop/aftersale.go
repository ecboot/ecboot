// aftersale.go 售后域——表: after_sale_order（V8/V23 数量/币种）。
// 规则: 按订单项粒度（quantity ≤ 行数量）; 状态机:
// 10待审核→(同意)仅退款30/退货退款20→确认收货30→渠道退款40→成功50; 10→90拒绝; 未终态→91撤销;
// 退款渠道幂等: out_refund_no = after_sale_no; 完成回补库存(按 quantity)+更新订单 refund_status+佣金冲销事件。
package shop

import "context"

// IAfterSaleLogic 售后。
type IAfterSaleLogic interface {
	// Apply 申请（校验订单已完成/项未超额; 创建 10）。
	Apply(ctx context.Context, userId int64, in AfterSaleApplyInput) (string, error)
	List(ctx context.Context, userId int64, status int, page PageQuery) (*PageResult[AfterSaleSummary], error)
	Detail(ctx context.Context, userId int64, afterSaleNo string) (*AfterSaleDetail, error)
	Cancel(ctx context.Context, userId int64, afterSaleNo string) error
	// SubmitReturn 填写寄回单号（type=2 且状态 20）。
	SubmitReturn(ctx context.Context, userId int64, afterSaleNo, logisticsNo string) error

	// ---- 后台 ----
	AdminList(ctx context.Context, status, saleType int, page PageQuery) (*PageResult[AfterSaleSummary], error)
	AdminDetail(ctx context.Context, afterSaleNo string) (*AfterSaleDetail, error)
	Approve(ctx context.Context, afterSaleNo, operator string) error
	Reject(ctx context.Context, afterSaleNo, reason, operator string) error
	// ConfirmReceipt 退货确认收货（→30 待退款）。
	ConfirmReceipt(ctx context.Context, afterSaleNo, operator string) error
	// RetryRefund 退款重试（渠道失败后）。
	RetryRefund(ctx context.Context, afterSaleNo, operator string) error
}

type AfterSaleApplyInput struct {
	OrderItemId   int64
	Type          int // 1仅退款 2退货退款
	Quantity      int
	Reason        string
	Description   string
	VoucherImages []string
}

type AfterSaleSummary struct {
	AfterSaleNo  string `json:"afterSaleNo"`
	OrderNo      string `json:"orderNo"`
	Type         int    `json:"type"`
	Quantity     int    `json:"quantity"`
	RefundAmount string `json:"refundAmount"`
	Status       int    `json:"status"`
	CreatedAt    string `json:"createdAt"`
}

type AfterSaleDetail struct {
	AfterSaleNo       string   `json:"afterSaleNo"`
	OrderNo           string   `json:"orderNo"`
	OrderItemId       int64    `json:"orderItemId"`
	Type              int      `json:"type"`
	Quantity          int      `json:"quantity"`
	Reason            string   `json:"reason"`
	Description       string   `json:"description"`
	VoucherImages     []string `json:"voucherImages"`
	RefundAmount      string   `json:"refundAmount"`
	ReturnLogisticsNo string   `json:"returnLogisticsNo"`
	RefundNo          string   `json:"refundNo" dc:"渠道退款单号"`
	RejectReason      string   `json:"rejectReason"`
	AuditTime         string   `json:"auditTime"`
	RefundTime        string   `json:"refundTime"`
	Status            int      `json:"status"`
}
