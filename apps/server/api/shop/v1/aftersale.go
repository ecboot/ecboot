package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/model"
)

type (
	// 申请售后（按订单项; 仅退款/退货退款）
	AfterSaleCreateReq struct {
		g.Meta        `path:"/after-sales" method:"POST" summary:"申请售后"`
		OrderItemId   string   `json:"orderItemId" v:"required" dc:"订单项ID"`
		Type          int      `json:"type" v:"required|in:1,2" dc:"类型:1仅退款 2退货退款"`
		Quantity      int      `json:"quantity" v:"required|min:1" dc:"退货数量"`
		Reason        string   `json:"reason" v:"required" dc:"售后原因"`
		Description   string   `json:"description" dc:"问题描述"`
		VoucherImages []string `json:"voucherImages" dc:"凭证图片"`
	}
	AfterSaleCreateRes struct {
		AfterSaleNo string `json:"afterSaleNo" dc:"售后单号"`
	}

	AfterSaleListItem struct {
		AfterSaleNo  string `json:"afterSaleNo"`
		OrderNo      string `json:"orderNo"`
		Type         int    `json:"type" dc:"1仅退款 2退货退款"`
		RefundAmount string `json:"refundAmount" dc:"退款金额"`
		Status       int    `json:"status" dc:"10待审核 20待寄回 30待退款 40退款中 50已完成 90已拒绝 91已撤销"`
		CreatedAt    string `json:"createdAt"`
	}
	AfterSaleListReq struct {
		g.Meta `path:"/after-sales" method:"GET" summary:"售后列表"`
		Status int `json:"status" dc:"状态筛选"`
		model.PageReq
	}
	AfterSaleListRes struct {
		model.PageRes
		List []AfterSaleListItem `json:"list"`
	}

	AfterSaleDetailReq struct {
		g.Meta      `path:"/after-sales/{afterSaleNo}" method:"GET" summary:"售后详情"`
		AfterSaleNo string `json:"afterSaleNo" v:"required" dc:"售后单号"`
	}
	AfterSaleDetailRes struct {
		AfterSaleNo       string   `json:"afterSaleNo"`
		Status            int      `json:"status"`
		Type              int      `json:"type"`
		Quantity          int      `json:"quantity"`
		Reason            string   `json:"reason"`
		Description       string   `json:"description"`
		VoucherImages     []string `json:"voucherImages"`
		RefundAmount      string   `json:"refundAmount"`
		ReturnLogisticsNo string   `json:"returnLogisticsNo" dc:"寄回单号"`
		RejectReason      string   `json:"rejectReason" dc:"拒绝原因"`
		AuditTime         string   `json:"auditTime"`
		RefundTime        string   `json:"refundTime"`
	}

	AfterSaleCancelReq struct {
		g.Meta      `path:"/after-sales/{afterSaleNo}/cancel" method:"POST" summary:"撤销售后申请"`
		AfterSaleNo string `json:"afterSaleNo" v:"required" dc:"售后单号"`
	}
	AfterSaleCancelRes struct {
		Success bool `json:"success"`
	}

	AfterSaleLogisticsReq struct {
		g.Meta            `path:"/after-sales/{afterSaleNo}/logistics" method:"POST" summary:"填写寄回单号(退货退款)"`
		AfterSaleNo       string `json:"afterSaleNo" v:"required" dc:"售后单号"`
		ReturnLogisticsNo string `json:"returnLogisticsNo" v:"required" dc:"寄回物流单号"`
	}
	AfterSaleLogisticsRes struct {
		Success bool `json:"success"`
	}
)
