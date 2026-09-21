package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/model"
)

type (
	// 售后列表（状态筛选）
	AdminAfterSaleListReq struct {
		g.Meta `path:"/after-sales" method:"GET" summary:"售后列表"`
		Status int `json:"status" dc:"状态筛选"`
		Type   int `json:"type" dc:"类型筛选"`
		model.PageReq
	}
	AdminAfterSaleItem struct {
		AfterSaleNo  string `json:"afterSaleNo"`
		OrderNo      string `json:"orderNo"`
		UserId       string `json:"userId"`
		Type         int    `json:"type"`
		Quantity     int    `json:"quantity"`
		RefundAmount string `json:"refundAmount"`
		Status       int    `json:"status"`
		CreatedAt    string `json:"createdAt"`
	}
	AdminAfterSaleListRes struct {
		model.PageRes
		List []AdminAfterSaleItem `json:"list"`
	}

	AdminAfterSaleDetailReq struct {
		g.Meta      `path:"/after-sales/{afterSaleNo}" method:"GET" summary:"售后详情"`
		AfterSaleNo string `json:"afterSaleNo" v:"required" dc:"售后单号"`
	}
	AdminAfterSaleDetailRes struct {
		AfterSaleNo       string   `json:"afterSaleNo"`
		OrderNo           string   `json:"orderNo"`
		UserId            string   `json:"userId"`
		Type              int      `json:"type"`
		Quantity          int      `json:"quantity"`
		Reason            string   `json:"reason"`
		Description       string   `json:"description"`
		VoucherImages     []string `json:"voucherImages"`
		RefundAmount      string   `json:"refundAmount"`
		ReturnLogisticsNo string   `json:"returnLogisticsNo"`
		RejectReason      string   `json:"rejectReason"`
		RefundNo          string   `json:"refundNo" dc:"渠道退款单号"`
		Status            int      `json:"status"`
	}

	// 同意（仅退款→待退款; 退货退款→待寄回）
	// 权限: aftersale:audit
	AdminAfterSaleApproveReq struct {
		g.Meta      `path:"/after-sales/{afterSaleNo}/approve" method:"POST" summary:"同意售后"`
		AfterSaleNo string `json:"afterSaleNo" v:"required" dc:"售后单号"`
	}
	AdminAfterSaleApproveRes struct {
		Success bool `json:"success"`
	}

	// 拒绝
	// 权限: aftersale:audit
	AdminAfterSaleRejectReq struct {
		g.Meta       `path:"/after-sales/{afterSaleNo}/reject" method:"POST" summary:"拒绝售后"`
		AfterSaleNo  string `json:"afterSaleNo" v:"required" dc:"售后单号"`
		RejectReason string `json:"rejectReason" v:"required" dc:"拒绝原因"`
	}
	AdminAfterSaleRejectRes struct {
		Success bool `json:"success"`
	}

	// 退货确认收货（→待退款）
	// 权限: aftersale:audit
	AdminAfterSaleConfirmReceiptReq struct {
		g.Meta      `path:"/after-sales/{afterSaleNo}/confirm-receipt" method:"POST" summary:"退货确认收货"`
		AfterSaleNo string `json:"afterSaleNo" v:"required" dc:"售后单号"`
	}
	AdminAfterSaleConfirmReceiptRes struct {
		Success bool `json:"success"`
	}

	// 退款重试
	// 权限: aftersale:refund
	AdminAfterSaleRetryRefundReq struct {
		g.Meta      `path:"/after-sales/{afterSaleNo}/retry-refund" method:"POST" summary:"退款重试"`
		AfterSaleNo string `json:"afterSaleNo" v:"required" dc:"售后单号"`
	}
	AdminAfterSaleRetryRefundRes struct {
		Success bool `json:"success"`
	}
)
