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
	"time"

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

// pay_order.status 契约（000007 列注释）: 10待支付 20支付成功 30支付失败 90已关闭
// 评审 I3: 原 `payStatusClosed = 30` 误把"关闭"写成"失败"（且 data-model 图纸同错）→ 已修正。
// 注: 30（失败）本批无写入路径——渠道回调 success=false 按"拒绝并等渠道重试"处理, 不推进状态机
// （若置 30, 后续真实成功回调的条件更新 `status=10` 将不命中, 反而制造资金脏态）。
const (
	payStatusPending = 10
	payStatusSuccess = 20
	payStatusClosed  = 90

	orderStatusPending  = 10
	orderStatusPaid     = 20
	orderStatusCanceled = 90
)

// payChannelInstance 渠道实例（mock; 真实渠道按入参选择——另立特性）。
func payChannelInstance(_ int) paychannel.Channel { return paychannel.NewMock() }

// errPayOrderNotPayable 回调命中"不可支付的支付单"（已关闭/已失败）——仅作事务内信号,
// 由事务外转为告警 + 契约错误码（评审 Critical: 不得当幂等吞掉）。
var errPayOrderNotPayable = gerror.New("支付单已关闭或已失败")

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

	// 评审 C3: 同订单同时至多一张待支付单（000007 注释为服务层义务）——先失效旧的待支付单
	pcols0 := dao.PayOrder.Columns()
	if _, err = dao.PayOrder.Ctx(ctx).
		Where(pcols0.OrderNo, orderNo).
		Where(pcols0.Status, payStatusPending).
		Data(g.Map{pcols0.Status: payStatusClosed, pcols0.ClosedTime: gtime.Now()}).
		Update(); err != nil {
		return nil, gerror.Wrap(err, "失效历史支付单失败")
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
		pcols.ExpireTime: gtime.Now().Add(30 * time.Minute), // 与订单超时口径同源（评审 I3）
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
// 评审 I2: ① 写入失败不再静默（留档是对账唯一依据, 丢失必须可见）;
// ② raw_body 为 JSON NOT NULL——非 JSON 原文由 jsonColumn 兜底编码后再入库。
func writeCallbackLog(ctx context.Context, channel, payNo string, notifyType int, rawBody string, processed bool) {
	ps := 0
	if processed {
		ps = 1
	}
	if _, err := dao.PayCallbackLog.Ctx(ctx).Data(g.Map{
		"pay_no":         payNo,
		"pay_channel":    channelCode(channel),
		"notify_type":    notifyType,
		"verify_status":  1, // mock 渠道验签占位通过
		"process_status": ps,
		"raw_body":       jsonColumn(rawBody),
	}).Insert(); err != nil {
		g.Log().Errorf(ctx, "回调留档写入失败 pay_no=%q notify_type=%d err=%v", payNo, notifyType, err)
	}
}

// jsonColumn 保证写入 JSON 列的字符串是合法 JSON。
// 背景（评审 I2）: pay_callback_log.raw_body 为 `JSON NOT NULL`, 直接写非 JSON 原文会报
// 3140 Invalid JSON text 并被原实现的 `_, _ =` 吞掉——**恰恰是"报文非法"这一最该留档的场景留不下任何记录**。
// 退化策略: 非 JSON 原文编码为 JSON 字符串字面量（内容完整保留, 仅多了转义, 仍满足"原文留档可重放"）。
func jsonColumn(raw string) string {
	if json.Valid([]byte(raw)) {
		return raw
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return `""`
	}
	return string(b)
}

// HandlePayNotify 支付回调（FR-006/007）: 幂等四层防线 + 同事务推进。
func (i *PayLogicImpl) HandlePayNotify(ctx context.Context, channel string, rawBody []byte) error {
	payload, err := paychannel.NewMock().ParseNotify(rawBody)
	if err != nil {
		// 评审 I2: 报文非法是最该留档的场景（原实现直接 return 无留档）
		writeCallbackLog(ctx, channel, "", 1, string(rawBody), false)
		return errcode.New(errcode.CodeInvalidParam, "回调报文非法")
	}
	raw := string(rawBody)

	// 评审 C2: success 标志必须校验——否则"支付失败"通知会走完成功路径（刷单/资损）
	if !payload.Success {
		writeCallbackLog(ctx, channel, payload.PayNo, 1, raw, false)
		return errcode.New(errcode.CodeInvalidParam, "渠道标记支付失败")
	}

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

	// ①③④ 条件更新 + 同事务推进
	// 评审 I2: **留档一律移出事务**——写在事务内则回滚即丢（原实现两处在事务内、一处在事务外, 语义不一）。
	processed := false // 是否真正推进（process_status=1）
	dirtyOrderNo := "" // 资金已入账但订单不可推进的脏态订单号（评审 C3b）
	var dirtyOrderId int64
	closedPayNo := "" // 命中已关闭/已失败支付单的支付号（评审 Critical）
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
			// 评审 Critical: affected=0 有**两种成因**, 此前被混为一谈以致静默资损。
			// ① status=20 → 真的已处理, 幂等应答 ✓
			// ② status∈{90已关闭, 30支付失败} → 渠道是就**旧单**扣的款（用户在浏览器/微信里扫的还是旧码）,
			//    钱真的到了。此时若照旧 `return nil`, 渠道收到 SUCCESS、支付单停在 90、订单停在 10、
			//    库存不核销、且 C3b 的口径（pay=20 且 order=90）也捞不到 → "钱付了、单没了、无人知道"。
			//    C3a 的关单逻辑正是主动制造 ② 的路径, 故必须在此处与幂等区分开。
			cur, ce := dao.PayOrder.Ctx(ctx).Fields(pcols.Status).Where(pcols.PayNo, payload.PayNo).Value()
			if ce == nil && cur.Int() == payStatusSuccess {
				return nil // ① 真幂等
			}
			closedPayNo = payload.PayNo // ② 交事务外告警与留档（不得应答 SUCCESS）
			return errPayOrderNotPayable
		}
		// 订单推进（条件 status=10; 已取消 → 不推进, 记脏态）
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
			// 支付成功但订单已取消等脏态: 不推进（资金由对账/退款处理）; 标记留待事务外留档
			dirtyOrderNo = rec[pcols.OrderNo].String()
			dirtyOrderId = rec[pcols.OrderId].Int64()
			return nil
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
			ires, ie := tx.Model(dao.Inventory.Table()).
				Where(dao.Inventory.Columns().SkuId, skuId).
				Where(dao.Inventory.Columns().Locked+" >= ?", qty).
				Data(g.Map{
					dao.Inventory.Columns().Total:  gdb.Raw("total - " + itoa(qty)),
					dao.Inventory.Columns().Locked: gdb.Raw("locked - " + itoa(qty)),
				}).Update()
			if ie != nil {
				return gerror.Wrap(ie, "库存核销失败")
			}
			// 评审 I1: 必须判行数——affected=0 意味着库存行缺失或锁定不足（如秒杀单不锁 inventory），
			// 静默跳过会造成"订单已推进但库存未扣"的账实不符。
			if n, _ := ires.RowsAffected(); n == 0 {
				return gerror.Newf("库存核销未命中(sku=%d qty=%d): 锁定不足或库存行缺失", skuId, qty)
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
		processed = true
		return nil
	})

	// 留档（事务外, 评审 I2）: 事务回滚也留——"回调来过但没处理成"正是对账最需要的一行。
	writeCallbackLog(ctx, channel, payload.PayNo, 1, raw, processed)
	if err != nil {
		// 评审 Critical: 命中已关闭/已失败支付单——钱可能真的到了, 不得应答 SUCCESS（渠道须收到 FAIL 并留痕）。
		if closedPayNo != "" {
			g.Log().Errorf(ctx,
				"[资金异常] 支付成功回调命中已关闭支付单, 须人工核对并退款: pay_no=%s order_no=%s",
				closedPayNo, rec[pcols.OrderNo].String())
			return errcode.New(errcode.CodeStatusNotAllowed, "支付单已关闭, 需人工核对退款")
		}
		return err
	}

	// 评审 C3b: 脏态资金对账标记——钱已入账（支付单保持 20, 不得篡改事实）而订单不可推进。
	// 对账口径（收敛为一条, 覆盖 ② 与 ③ 两种形态）:
	//   **存在 notify_type=1 且 process_status=0 的留档** 即"回调来过但没处理成":
	//   ② 已关闭支付单命中（pay=90/30, 见上）  ③ 订单已取消/不可推进（pay=20 而 order≠20）。
	// 此处再打一条固定前缀的 Error 级告警（可捞）; 不写 trade_order_log 标记行——该表经 OrderDetail
	// 直接暴露给 C 端时间线, 会泄露内部资金文案。
	if dirtyOrderNo != "" {
		g.Log().Errorf(ctx, "[资金异常] 支付成功但订单不可推进, 待退款对账: pay_no=%s order_no=%s order_id=%d",
			payload.PayNo, dirtyOrderNo, dirtyOrderId)
	}
	return nil
}

// HandleRefundNotify 退款回调（FR-008）: 幂等推进售后单状态（mock 渠道）。
func (i *PayLogicImpl) HandleRefundNotify(ctx context.Context, channel string, rawBody []byte) error {
	var p struct {
		PayNo       string `json:"payNo"`
		OutRefundNo string `json:"outRefundNo"`
		Success     bool   `json:"success"`
	}
	if err := json.Unmarshal(rawBody, &p); err != nil {
		writeCallbackLog(ctx, channel, "", 2, string(rawBody), false)
		return errcode.New(errcode.CodeInvalidParam, "退款回调报文非法")
	}
	if !p.Success {
		writeCallbackLog(ctx, channel, p.PayNo, 2, string(rawBody), false)
		return errcode.New(errcode.CodeInvalidParam, "渠道标记退款失败")
	}
	// 评审 C1 修正: ① 匹配键用 after_sale_no（000008 契约: out_refund_no = after_sale_no;
	// 原用 refund_no 是渠道回填列, 无人写入 → WHERE 永不命中）；② 状态 40退款中 → 50已完成
	// （原 30→40 实为"待退款→退款中", 语义错位）；③ affected=0 不得记"已处理"。
	acols := dao.AfterSaleOrder.Columns()
	advanced := false
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, e := dao.AfterSaleOrder.Ctx(ctx).
			Where(acols.AfterSaleNo, p.OutRefundNo).
			Where(acols.Status, 40). // 40=退款中
			Data(g.Map{acols.Status: 50, acols.RefundTime: gtime.Now()}).
			Update()
		if e != nil {
			return gerror.Wrap(e, "推进售后单失败")
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return nil // 未匹配/状态不符 → 事务外判定（不记"已处理"）
		}
		advanced = true
		// 013（批次07）: 完成副作用与状态推进**同事务**——仅退货退款回补库存 + 重算订单退款状态 +
		// 投递佣金冲销事件。仅在本次条件更新真正命中时执行一次, 故重复回调不会二次回补/二次投递。
		return afterSaleFinishedSideEffects(ctx, p.OutRefundNo, "system:refund")
	})
	if err != nil {
		// 评审 I2: 事务失败（含完成副作用失败）也必须留档 + 告警——"回调来过但没处理成"正对对账口径,
		// 且此刻钱已出（状态 40 意味着渠道退款成功）, 不告警则长期无人知道。
		writeCallbackLog(ctx, channel, p.PayNo, 2, string(rawBody), false)
		g.Log().Errorf(ctx, "[资金异常] 退款回调处理失败, 售后单未推进, 需人工核对: out_refund_no=%s err=%v",
			p.OutRefundNo, err)
		return err
	}
	if !advanced {
		// 评审 M2: 区分"真幂等"与"不匹配"——状态已是 50（已完成）说明这是重复回调, 应静默应答成功,
		// 否则渠道会为一个本该忽略的回调无限重试; 其余情况（未匹配/已撤销/状态不符）仍返错并要求留档。
		cur, ce := dao.AfterSaleOrder.Ctx(ctx).Fields(acols.Status).
			Where(acols.AfterSaleNo, p.OutRefundNo).Value()
		if ce == nil && cur.Int() == 50 {
			writeCallbackLog(ctx, channel, p.PayNo, 2, string(rawBody), true)
			return nil
		}
		writeCallbackLog(ctx, channel, p.PayNo, 2, string(rawBody), false)
		return errcode.New(errcode.CodeNotFound, "退款单未匹配或状态不符")
	}
	writeCallbackLog(ctx, channel, p.PayNo, 2, string(rawBody), true)
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
