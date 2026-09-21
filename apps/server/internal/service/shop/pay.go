// pay.go 支付域——表: pay_order / pay_callback_log（V7）。
// 规则: 一单多尝试（同单同时至多一个待支付——服务层保证）; 回调幂等四层:
// 条件UPDATE抢占(status=10) + UNIQUE(渠道,渠道单号)防重放 + 金额/币种校验 + 原文留档;
// 支付成功事务: 核销库存 + 优惠核销 + 余额消费完成 + 订单20 + 分销/积分事件（FR 契约 §支付）。
package shop

import "context"

// IPayLogic 支付。
type IPayLogic interface {
	// Create 发起支付（校验应付>0; 返回渠道唤起参数——渠道适配器另立特性, 当前 mock 参数）。
	Create(ctx context.Context, userId int64, orderNo string, channel int) (*PayCreated, error)
	// Status 状态查询。
	Status(ctx context.Context, userId int64, payNo string) (*PayStatus, error)
	// HandlePayNotify 支付回调（验签语义+幂等+同事务: 支付单20/订单20/库存核销/余额消费完成/通知事件; 原文留档）。
	HandlePayNotify(ctx context.Context, channel string, rawBody []byte) error
	// HandleRefundNotify 退款回调（售后域联动: 退款单状态推进）。
	HandleRefundNotify(ctx context.Context, channel string, rawBody []byte) error
	// CloseExpired 关单扫描（订单取消/超时的支付单关闭, 定时任务）。
	CloseExpired(ctx context.Context) (int64, error)
}

type PayCreated struct {
	PayNo         string         `json:"payNo"`
	ChannelParams map[string]any `json:"channelParams" dc:"渠道唤起参数(mock 为占位)"`
}

type PayStatus struct {
	Status int    `json:"status" dc:"10待支付 20成功 30失败 90关闭"`
	PaidAt string `json:"paidAt"`
}
