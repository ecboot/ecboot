// pay_impl.go 支付实现（012-trade-wiring; 接口契约见 pay.go IPayLogic）。
// 幂等四层防线（research D2）:
//
//	① 条件更新（pay_order 10→20 且 affected=0 视为已处理 → 幂等应答）
//	② 留档（pay_callback_log 只追加, 成败均落 —— 对账与重放识别）
//	③ 金额校验（回调金额 vs 支付单应付; 不符 → 留档 + 拒绝）
//	④ 原文留档（rawBody 原样入库）
//
// 同事务推进: 支付单 20 → 订单 20（条件 status=10）→ 库存核销（total-n, locked-n）。
// 说明: 余额消费完成依赖账户域（user_account, 批次 11 实现）→ 本批落地库存核销与订单推进,
// 余额结算随批次 11 一并接（spec FR-006 范围修正, PROGRESS 已记账）。
package shop

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/library/money"
	"ecboot/internal/library/paychannel"
	"ecboot/internal/model"
)

// PayLogicImpl IPayLogic 实现。
type PayLogicImpl struct{}

func NewPayLogic() *PayLogicImpl { return &PayLogicImpl{} }

const (
	payStatusPending = 10
	payStatusSuccess = 20
	payStatusClosed  = 30

	orderStatusPending  = 10
	orderStatusPaid     = 20
	orderStatusCanceled = 90
)

// payChannelInstance 渠道实例（mock; 真实渠道按入参选择——另立特性）。
func payChannelInstance(_ int) paychannel.Channel { return paychannel.NewMock() }

// yuanToFen 元 → 分（复用项目 money 库: 内部一律 int64 分运算）。
func yuanToFen(yuan string) (int64, error) {
	fen, err := money.FromYuanString(yuan)
	if err != nil {
		return 0, errcode.New(errcode.CodeInvalidParam, "金额格式非法")
	}
	return fen, nil
}

// payNoOf 支付号（时间戳 + 支付单自增 ID 兜底唯一; 与既有 nextOrderNo 同风格）。
func payNoOf() (string, error) {
	return nextOrderNo() // 复用既有编号生成（前缀 ORD → 此处替换前缀）
}

// Create 发起支付（FR-005）: 校验订单态与归属 → 生成支付单 → 返回渠道唤起参数。
func (i *PayLogicImpl) Create(ctx context.Context, userId int64, orderNo string, channel int) (*model.PayCreated, error) {
	ocols := dao.TradeOrder.Columns()
	rec, err := dao.TradeOrder.Ctx(ctx).
		Where(ocols.OrderNo, orderNo).
		Where(ocols.UserId, userId).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单失败")
	}
	if rec.IsEmpty() {
		return nil, errcode.New(errcode.CodeOrderNotFound, "订单不存在")
	}
	if rec[ocols.Status].Int() != orderStatusPending {
		return nil, errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可支付")
	}
	payFen, err := yuanToFen(rec[ocols.PayAmount].String())
	if err != nil {
		return nil, err
	}
	if payFen <= 0 {
		return nil, errcode.New(errcode.CodePayCreateFailed, "应付金额必须大于 0")
	}

	raw, err := payNoOf()
	if err != nil {
		return nil, gerror.Wrap(err, "生成支付号失败")
	}
	payNo := "PAY" + raw[3:] // 复用编号主体, 换支付前缀
	pcols := dao.PayOrder.Columns()
	res, err := dao.PayOrder.Ctx(ctx).Data(g.Map{
		pcols.PayNo:      payNo,
		pcols.OrderId:    rec[ocols.Id].Int64(),
		pcols.OrderNo:    orderNo,
		pcols.UserId:     userId,
		pcols.PayChannel: channel,
		pcols.Amount:     rec[ocols.PayAmount].String(),
		pcols.Currency:   "CNY",
		pcols.Status:     payStatusPending,
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "创建支付单失败")
	}
	_ = res

	params, err := payChannelInstance(channel).CreateParams(payNo, payFen, "订单"+orderNo)
	if err != nil {
		return nil, gerror.Wrap(err, "生成渠道参数失败")
	}
	return &model.PayCreated{PayNo: payNo, ChannelParams: params}, nil
}

// Status 支付状态查询（归属校验）。
func (i *PayLogicImpl) Status(ctx context.Context, userId int64, payNo string) (*model.PayStatus, error) {
	pcols := dao.PayOrder.Columns()
	rec, err := dao.PayOrder.Ctx(ctx).
		Where(pcols.PayNo, payNo).
		Where(pcols.UserId, userId).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询支付单失败")
	}
	if rec.IsEmpty() {
		return nil, errcode.New(errcode.CodeNotFound, "支付单不存在")
	}
	return &model.PayStatus{
		Status: rec[pcols.Status].Int(),
		PaidAt: rec[pcols.SuccessTime].String(),
	}, nil
}

// channelCode 渠道字符串 → 数值码（pay_* 表 pay_channel 为 TINYINT: 1微信 2支付宝; mock 归 1）。
func channelCode(channel string) int {
	switch channel {
	case "wechat", "wx":
		return 1
	case "alipay":
		return 2
	default:
		return 1 // mock 等开发态渠道
	}
}

// writeCallbackLog 回调原文留档（只追加; 成败均落——对账依据）。
// 列语义（000007 注释）: notify_type 1支付结果/2退款结果; verify_status 0失败/1通过;
// process_status 0未处理或重复忽略/1已处理。
func writeCallbackLog(ctx context.Context, channel, payNo string, notifyType int, rawBody string, processed bool) {
	ps := 0
	if processed {
		ps = 1
	}
	_, _ = dao.PayCallbackLog.Ctx(ctx).Data(g.Map{
		"pay_no":         payNo,
		"pay_channel":    channelCode(channel),
		"notify_type":    notifyType,
		"verify_status":  1, // mock 渠道验签占位通过
		"process_status": ps,
		"raw_body":       rawBody,
	}).Insert()
}

// HandlePayNotify 支付回调（FR-006/007）: 幂等四层防线 + 同事务推进。
func (i *PayLogicImpl) HandlePayNotify(ctx context.Context, channel string, rawBody []byte) error {
	payload, err := paychannel.NewMock().ParseNotify(rawBody)
	if err != nil {
		return errcode.New(errcode.CodeInvalidParam, "回调报文非法")
	}
	raw := string(rawBody)

	// ⑧ 金额校验（先查支付单）
	pcols := dao.PayOrder.Columns()
	rec, err := dao.PayOrder.Ctx(ctx).Where(pcols.PayNo, payload.PayNo).One()
	if err != nil {
		return gerror.Wrap(err, "查询支付单失败")
	}
	if rec.IsEmpty() {
		writeCallbackLog(ctx, channel, payload.PayNo, 1, raw, false)
		return errcode.New(errcode.CodeNotFound, "支付单不存在")
	}
	expectFen, err := yuanToFen(rec[pcols.Amount].String())
	if err != nil {
		return err
	}
	if payload.AmountFen != expectFen {
		writeCallbackLog(ctx, channel, payload.PayNo, 1, raw, false)
		return errcode.New(errcode.CodeInvalidParam, "回调金额与应付不符")
	}

	// ①③④ 条件更新 + 留档 + 同事务推进
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, e := dao.PayOrder.Ctx(ctx).
			Where(pcols.PayNo, payload.PayNo).
			Where(pcols.Status, payStatusPending).
			Data(g.Map{
				pcols.Status:         payStatusSuccess,
				pcols.ChannelTradeNo: payload.ChannelTradeNo,
				pcols.SuccessTime:    gtime.Now(),
			}).Update()
		if e != nil {
			return gerror.Wrap(e, "更新支付单失败")
		}
		if n, _ := res.RowsAffected(); n == 0 {
			// 已处理（或已关闭）→ 幂等应答（不二次推进）
			writeCallbackLog(ctx, channel, payload.PayNo, 1, raw, false)
			return nil
		}
		// 订单推进（条件 status=10; 已取消 → 不推进, 留档）
		ocols := dao.TradeOrder.Columns()
		ores, e := dao.TradeOrder.Ctx(ctx).
			Where(ocols.OrderNo, rec[pcols.OrderNo].String()).
			Where(ocols.Status, orderStatusPending).
			Data(g.Map{ocols.Status: orderStatusPaid, ocols.PayTime: gtime.Now()}).
			Update()
		if e != nil {
			return gerror.Wrap(e, "推进订单失败")
		}
		if n, _ := ores.RowsAffected(); n == 0 {
			writeCallbackLog(ctx, channel, payload.PayNo, 1, raw, false)
			return nil // 支付成功但订单已取消等脏态: 不推进（资金由对账/退款处理）
		}
		// 库存核销（按订单行: total-n, locked-n）
		items, e := dao.TradeOrderItem.Ctx(ctx).
			Where(dao.TradeOrderItem.Columns().OrderNo, rec[pcols.OrderNo].String()).All()
		if e != nil {
			return gerror.Wrap(e, "查询订单项失败")
		}
		for _, it := range items {
			skuId := it[dao.TradeOrderItem.Columns().SkuId].Int64()
			qty := it[dao.TradeOrderItem.Columns().Quantity].Int()
			if _, e = tx.Model(dao.Inventory.Table()).
				Where(dao.Inventory.Columns().SkuId, skuId).
				Where(dao.Inventory.Columns().Locked+" >= ?", qty).
				Data(g.Map{
					dao.Inventory.Columns().Total:  gdb.Raw("total - " + itoa(qty)),
					dao.Inventory.Columns().Locked: gdb.Raw("locked - " + itoa(qty)),
				}).Update(); e != nil {
				return gerror.Wrap(e, "库存核销失败")
			}
		}
		// 状态流水
		_, _ = dao.TradeOrderLog.Ctx(ctx).Data(g.Map{
			"order_no":      rec[pcols.OrderNo].String(),
			"order_id":      rec[pcols.OrderId].Int64(),
			"from_status":   orderStatusPending,
			"to_status":     orderStatusPaid,
			"remark":        "支付成功",
			"operator_type": 1, // 1=系统（渠道回调）
			"operator_id":   "system:pay",
		}).Insert()
		writeCallbackLog(ctx, channel, payload.PayNo, 1, raw, true)
		return nil
	})
	return err
}

// HandleRefundNotify 退款回调（FR-008）: 幂等推进售后单状态（mock 渠道）。
func (i *PayLogicImpl) HandleRefundNotify(ctx context.Context, channel string, rawBody []byte) error {
	var p struct {
		PayNo       string `json:"payNo"`
		OutRefundNo string `json:"outRefundNo"`
		Success     bool   `json:"success"`
	}
	if err := json.Unmarshal(rawBody, &p); err != nil || !p.Success {
		return errcode.New(errcode.CodeInvalidParam, "退款回调报文非法")
	}
	// 售后单联动（条件更新: 退款中 → 已退款; 幂等——affected=0 视为已处理）
	acols := dao.AfterSaleOrder.Columns()
	res, err := dao.AfterSaleOrder.Ctx(ctx).
		Where(acols.RefundNo, p.OutRefundNo).
		Where(acols.Status, 30). // 30=退款中（V8 状态机）
		Data(g.Map{acols.Status: 40, acols.RefundTime: gtime.Now()}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "推进售后单失败")
	}
	n, _ := res.RowsAffected()
	writeCallbackLog(ctx, channel, p.PayNo, 2, string(rawBody), true)
	_ = n
	return nil
}

// CloseExpired 关单扫描（内部方法; 过期未支付的支付单关闭）。
func (i *PayLogicImpl) CloseExpired(ctx context.Context) (int64, error) {
	pcols := dao.PayOrder.Columns()
	res, err := dao.PayOrder.Ctx(ctx).
		Where(pcols.Status, payStatusPending).
		Where("expire_time IS NOT NULL AND expire_time < NOW()").
		Data(g.Map{pcols.Status: payStatusClosed, pcols.ClosedTime: gtime.Now()}).
		Update()
	if err != nil {
		return 0, gerror.Wrap(err, "关闭支付单失败")
	}
	n, _ := res.RowsAffected()
	return n, nil
}

// itoa 小整数转字符串（Raw 片段拼接; 值域为数量, 无注入面）。
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
