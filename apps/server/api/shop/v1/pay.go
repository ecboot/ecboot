package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 发起支付（创建支付单; 返回渠道唤起参数）
	PayCreateReq struct {
		g.Meta     `path:"/pay" method:"POST" summary:"发起支付"`
		OrderNo    string `json:"orderNo" v:"required" dc:"订单号"`
		PayChannel int    `json:"payChannel" dc:"支付渠道:1微信支付" d:"1"`
	}
	PayCreateRes struct {
		PayNo         string         `json:"payNo" dc:"支付单号"`
		ChannelParams map[string]any `json:"channelParams" dc:"渠道唤起参数"`
	}

	// 支付状态查询
	PayStatusReq struct {
		g.Meta `path:"/pay/status" method:"GET" summary:"支付状态查询"`
		PayNo  string `json:"payNo" v:"required" dc:"支付单号"`
	}
	PayStatusRes struct {
		Status int    `json:"status" dc:"10待支付 20成功 30失败 90已关闭"`
		PaidAt string `json:"paidAt" dc:"支付时间"`
	}

	// 支付结果回调（开放路径; 渠道→平台, 需签名验证）
	PayNotifyReq struct {
		g.Meta `path:"/pay/notify" method:"POST" summary:"支付结果回调(渠道)"`
	}
	PayNotifyRes struct {
		Code    string `json:"code" dc:"渠道应答码"`
		Message string `json:"message" dc:"渠道应答消息"`
	}

	// 退款结果回调（开放路径; 需签名验证）
	RefundNotifyReq struct {
		g.Meta `path:"/refund/notify" method:"POST" summary:"退款结果回调(渠道)"`
	}
	RefundNotifyRes struct {
		Code    string `json:"code" dc:"渠道应答码"`
		Message string `json:"message" dc:"渠道应答消息"`
	}
)
