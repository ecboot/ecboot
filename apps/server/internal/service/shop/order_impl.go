// order_impl.go IOrderLogic 实现——下单九步事务编排（data-model §一）与订单管理。
// 全部状态迁移为条件 UPDATE（affected=0 → 40006）；
// 事务传递：g.DB().Begin → tx.Model(...)/dao.Xxx.TX(tx) 链式操作，统一 Commit/Rollback。
package shop

import (
	"context"

	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"strings"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/library/money"
	"ecboot/internal/model"
	"ecboot/internal/model/do"
)

// OrderLogicImpl IOrderLogic 实现。
type OrderLogicImpl struct{}

func NewOrderLogic() *OrderLogicImpl { return &OrderLogicImpl{} }

func nextOrderNo() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "ORD" + time.Now().Format("060102150405") + hex.EncodeToString(b)[:4], nil
}

// ---------- C 端 ----------

// orderLine 下单行（聚合 SKU/SPU/库存信息）。
type orderLine struct {
	SkuId          int64
	SkuNo          string // trade_order_item.sku_no 非空列（快照）
	SpuId          int64
	SpuName        string
	SkuName        string
	SkuImage       string
	Specs          map[string]string
	PriceFen       int64
	LineFen        int64
	Quantity       int
	SaleRestricted bool
}

// Create 下单九步事务编排（幂等/四玩法/优惠分摊/库存锁定）。
func (i *OrderLogicImpl) Create(ctx context.Context, userId int64, in model.OrderCreateInput) (*model.OrderCreated, error) {
	if in.RequestToken == "" {
		return nil, errcode.New(errcode.CodeInvalidParam, "缺少幂等凭证")
	}

	// 评审 Important: 秒杀玩法（批次 09/10 营销域）尚未接线——原分支只锁活动库存、**不锁 inventory**,
	// 且 collectLines 恒用 product_sku.price（从不读 flash_sale_item.flash_price, 而 000018 的契约是
	// "秒杀价经 trade_order_item.price 快照承载"）。这样的单 ① 按原价计费 ② 支付回调的库存核销必然
	// 未命中（I1 判定为账实不符）→ 整单回滚, 即**永远无法支付**。与其产出无法履约的单, 不如明确拒绝。
	// 待批次 09 接通「秒杀价快照 + inventory 锁 + 取消时回补 sold_count」后删除本闸（PROGRESS §五 已记账）。
	if in.FlashSaleItemId > 0 {
		return nil, errcode.New(errcode.CodeActivityInvalid, "秒杀玩法未上线")
	}

	// ---- 事务外预备: 收货地址快照 ----
	addr, err := dao.UserAddress.Ctx(ctx).
		Where(dao.UserAddress.Columns().Id, in.AddressId).
		Where(dao.UserAddress.Columns().UserId, userId).One()
	if err != nil {
		return nil, err
	}
	if addr.IsEmpty() {
		return nil, errcode.New(errcode.CodeInvalidParam, "收货地址不存在")
	}
	provinceCode := addr["province_code"].String()

	// ---- 事务 ----
	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, err
	}
	// 评审 Important（修复轮自查漏项）: 用**提交标志**兜底回滚, 不再依赖"每条错误路径都记得给 err 赋值"
	// 这种易漏的不变量——原写法在 `len(lines)==0`（任何客户端 POST 空购买项即可触发）与"收货地限售"
	// 两条早退路径上漏赋值, 于是 BEGIN 之后既不提交也不回滚, 事务与连接一直悬到请求 ctx 被取消归还;
	// 长生命周期 ctx 的调用方（定时任务/后台 worker）会真实泄漏, 且一旦有人把判空挪到取锁之后即变持锁泄漏。
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	orderNo, err := nextOrderNo()
	if err != nil {
		return nil, err
	}

	// 步骤1: 行聚合
	lines, totalFen, err := i.collectLines(ctx, tx, userId, in)
	if err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, errcode.New(errcode.CodeInvalidParam, "无有效购买项")
	}

	// 步骤2: 限售校验
	for _, ln := range lines {
		if ln.SaleRestricted && provinceCode != "" && provinceRestricted(ctx, ln.SpuId, provinceCode) {
			return nil, errcode.New(errcode.CodeInvalidParam, "商品在收货地区限售")
		}
	}

	// 步骤3: 库存锁定
	// 012 修复轮 I11: 条件更新**必须判 RowsAffected**——未命中即库存不足。
	// 原实现只看 error（条件不命中时 error 为 nil, 行数为 0）, 于是库存 1 也能下 2 件的单（静默超卖）。
	// 回滚由函数头的 `committed` 标志兜底（见事务开头的注释）, 此处无需关心 err 赋值约定。
	for _, ln := range lines {
		// 普通/拼团/砍价: inventory 锁定（条件更新防超卖）
		var invRes sql.Result
		if invRes, err = tx.Model("inventory").Ctx(ctx).
			Where("sku_id", ln.SkuId).
			Where("total - locked >= ?", ln.Quantity).
			Data(g.Map{"locked": gdb.Raw("locked + " + fmt.Sprint(ln.Quantity))}).
			Update(); err != nil {
			return nil, err
		}
		if n, _ := invRes.RowsAffected(); n == 0 {
			err = errcode.New(errcode.CodeStockInsufficient, "库存不足")
			return nil, err
		}
	}

	// 步骤4: 优惠计算（先满减后券; 积分; 余额冻结最后）
	fullReductionFen := calcFullReductionFen(ctx, provinceCode, totalFen)
	couponFen := int64(0)
	if in.UserCouponId > 0 {
		couponFen = calcCouponDiscountFen(ctx, userId, in.UserCouponId, totalFen)
	}
	pointFen := int64(0)
	if in.UsePoint {
		pointFen = calcPointDeductFen(ctx, userId, true, totalFen-couponFen-fullReductionFen)
	}
	promotionFen := couponFen + fullReductionFen + pointFen
	freightFen := int64(0) // V1 包邮; 运费模板计算随交易完善
	payFen := totalFen - promotionFen + freightFen

	// ---- 步骤5: 订单快照落库 ----
	orderData := g.Map{
		"order_no":               orderNo,
		"user_id":                userId,
		"order_channel":          in.Channel,
		"status":                 10,
		"currency":               "CNY",
		"total_amount":           money.ToYuanString(totalFen),
		"promotion_amount":       money.ToYuanString(promotionFen),
		"coupon_amount":          money.ToYuanString(couponFen),
		"full_reduction_amount":  money.ToYuanString(fullReductionFen),
		"point_amount":           money.ToYuanString(pointFen),
		"point_used":             0,
		"account_amount":         money.ToYuanString(0),
		"freight_amount":         money.ToYuanString(freightFen),
		"pay_amount":             money.ToYuanString(payFen),
		"receiver_name":          addr["receiver_name"].String(),
		"receiver_phone":         addr["receiver_phone"].String(),
		"receiver_province":      addr["province"].String(),
		"receiver_city":          addr["city"].String(),
		"receiver_district":      addr["district"].String(),
		"receiver_detail":        addr["detail_address"].String(),
		"receiver_province_code": addr["province_code"].String(),
		"receiver_city_code":     addr["city_code"].String(),
		"receiver_district_code": addr["district_code"].String(),
		"user_remark":            in.UserRemark,
		"request_token":          in.RequestToken,
	}
	if in.UserCouponId > 0 {
		orderData["user_coupon_id"] = in.UserCouponId
	}
	if in.GroupBuyTeamId > 0 {
		orderData["group_buy_team_id"] = in.GroupBuyTeamId
	}
	if in.BargainRecordId > 0 {
		orderData["bargain_record_id"] = in.BargainRecordId
	}
	if _, err = tx.Model("trade_order").Ctx(ctx).Data(orderData).Insert(); err != nil {
		return nil, err
	}
	orderId, err := tx.Model("trade_order").Ctx(ctx).
		Where("order_no", orderNo).Value("id")
	if err != nil {
		return nil, err
	}

	// ---- 步骤6: 订单项快照（行分摊, 尾差记末行） ----
	// 012 修复轮 C4a: sku_no/spu_name/sku_name/original_price 均为**非空无默认值列**（000006），
	// 原实现漏写 → 每次下单报 1364 回滚（该端点此前零测试 + 冒烟未覆盖下单, 故长期未暴露）。
	// 012 修复轮（评审 Important）: 三个**行分摊列**按各自构成分别分摊——表契约（000019/000020）
	// 明写"合计=订单头对应明细, 尾差记末行", 而原实现只写聚合列, 三列恒 0.00:
	// 订单头显示"券抵 5 元"而每行 coupon_amount=0, 售后按行取数（trade_order_item 是售后金额依据）即对不上。
	lineFens := lineFensOf(lines)
	couponAlloc := money.AllocateProRata(couponFen, lineFens)
	frAlloc := money.AllocateProRata(fullReductionFen, lineFens)
	pointAlloc := money.AllocateProRata(pointFen, lineFens)
	for idx, ln := range lines {
		// 逐构成分摊后求和: 各构成的分摊和恒等于订单头对应值（AllocateProRata 保和）, 故行 promo 和 == promotionFen
		linePromo := couponAlloc[idx] + frAlloc[idx] + pointAlloc[idx]
		linePay := ln.LineFen - linePromo
		if _, err = tx.Model("trade_order_item").Ctx(ctx).Data(g.Map{
			"order_no":              orderNo,
			"order_id":              orderId,
			"spu_id":                ln.SpuId,
			"sku_id":                ln.SkuId,
			"sku_no":                ln.SkuNo,
			"spu_name":              ln.SpuName,
			"sku_name":              ln.SkuName,
			"sku_image":             ln.SkuImage,
			"quantity":              ln.Quantity,
			"sku_specs":             mustJSON(ln.Specs),
			"original_price":        money.ToYuanString(ln.PriceFen),
			"price":                 money.ToYuanString(ln.PriceFen),
			"coupon_amount":         money.ToYuanString(couponAlloc[idx]),
			"full_reduction_amount": money.ToYuanString(frAlloc[idx]),
			"point_amount":          money.ToYuanString(pointAlloc[idx]),
			"promotion_amount":      money.ToYuanString(linePromo),
			"pay_amount":            money.ToYuanString(linePay),
		}).Insert(); err != nil {
			return nil, err
		}
	}

	// ---- 步骤7: 状态流水 ----
	// 012 修复轮 C4b: trade_order_log **无 operator 列**（只有 operator_type/operator_id）——
	// 与评审 I6（Cancel 同缺陷）同根因; 下单方为用户 → 2=用户（表注释: 1系统 2用户 3管理员）。
	if _, err = tx.Model("trade_order_log").Ctx(ctx).Data(do.TradeOrderLog{
		OrderNo:      orderNo,
		OrderId:      orderId,
		ToStatus:     10,
		Remark:       "订单创建",
		OperatorType: 2,
		OperatorId:   fmt.Sprintf("user:%d", userId),
	}).Insert(); err != nil {
		return nil, err
	}

	// ---- 提交事务 ----
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	committed = true

	return &model.OrderCreated{
		OrderNo:   orderNo,
		PayAmount: money.ToYuanString(payFen),
	}, nil
}

// lineFensOf 提取行金额列表。
func lineFensOf(lines []orderLine) []int64 {
	out := make([]int64, len(lines))
	for i, l := range lines {
		out[i] = l.LineFen
	}
	return out
}

// List 我的订单列表。
func (i *OrderLogicImpl) List(ctx context.Context, userId int64, status int, page model.PageReq) (*model.PageResult[model.OrderSummary], error) {
	page = page.Normalized()
	m := dao.TradeOrder.Ctx(ctx).
		Where(dao.TradeOrder.Columns().UserId, userId)
	if status > 0 {
		m = m.Where(dao.TradeOrder.Columns().Status, status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	orders, err := m.Page(page.Page, page.PageSize).Order("id desc").All()
	if err != nil {
		return nil, err
	}
	list := []model.OrderSummary{}
	for _, r := range orders {
		ono := r["order_no"].String()
		items, _ := dao.TradeOrderItem.Ctx(ctx).
			Where(dao.TradeOrderItem.Columns().OrderNo, ono).All()
		oi := []model.OrderItemBrief{}
		for _, it := range items {
			oi = append(oi, model.OrderItemBrief{
				SpuName:  it["spu_name"].String(),
				SkuSpecs: specsMap(it["sku_specs"].String()),
				Image:    it["sku_image"].String(),
				Quantity: it["quantity"].Int(),
				Price:    it["price"].String(),
			})
		}
		list = append(list, model.OrderSummary{
			OrderNo:   ono,
			Status:    r["status"].Int(),
			Amount:    orderAmountBook(r),
			Items:     oi,
			CreatedAt: r["created_at"].String(),
		})
	}
	return &model.PageResult[model.OrderSummary]{List: list, Total: int64(total)}, nil
}

// orderAmountBook 从订单记录构造金额账本。
func orderAmountBook(r gdb.Record) model.AmountBook {
	return model.AmountBook{
		TotalAmount:         r["total_amount"].String(),
		CouponAmount:        r["coupon_amount"].String(),
		FullReductionAmount: r["full_reduction_amount"].String(),
		PointAmount:         r["point_amount"].String(),
		AccountAmount:       r["account_amount"].String(),
		FreightAmount:       r["freight_amount"].String(),
		PayAmount:           r["pay_amount"].String(),
	}
}

// OrderDetail 订单详情（含状态时间线; 他人资源按不存在处理）。
func (i *OrderLogicImpl) OrderDetail(ctx context.Context, userId int64, orderNo string) (*model.OrderDetailView, error) {
	m := dao.TradeOrder.Ctx(ctx).Where(dao.TradeOrder.Columns().OrderNo, orderNo)
	if userId > 0 { // userId=0 为后台视角（不限制归属, 012 管理面复用）
		m = m.Where(dao.TradeOrder.Columns().UserId, userId)
	}
	rec, err := m.One()
	if err != nil {
		return nil, err
	}
	if rec.IsEmpty() {
		return nil, errcode.New(errcode.CodeOrderNotFound, "订单不存在")
	}
	items, _ := dao.TradeOrderItem.Ctx(ctx).
		Where(dao.TradeOrderItem.Columns().OrderNo, orderNo).All()
	oi := []model.OrderItemBrief{}
	for _, it := range items {
		oi = append(oi, model.OrderItemBrief{
			SpuName:  it["spu_name"].String(),
			SkuSpecs: specsMap(it["sku_specs"].String()),
			Image:    it["sku_image"].String(),
			Quantity: it["quantity"].Int(),
			Price:    it["price"].String(),
		})
	}
	logs, _ := dao.TradeOrderLog.Ctx(ctx).
		Where(dao.TradeOrderLog.Columns().OrderNo, orderNo).
		Order("id asc").All()
	var statusLogs []model.OrderStatusLog
	for _, lg := range logs {
		statusLogs = append(statusLogs, model.OrderStatusLog{
			FromStatus: lg["from_status"].Int(),
			ToStatus:   lg["to_status"].Int(),
			Remark:     lg["remark"].String(),
			CreatedAt:  lg["created_at"].String(),
		})
	}
	return &model.OrderDetailView{
		OrderNo:      orderNo,
		Status:       rec["status"].Int(),
		RefundStatus: rec["refund_status"].Int(),
		Amount:       orderAmountBook(rec),
		Items:        oi,
		Receiver: map[string]string{
			"name":  rec["receiver_name"].String(),
			"phone": rec["receiver_phone"].String(),
			"address": rec["receiver_province"].String() + rec["receiver_city"].String() +
				rec["receiver_district"].String() + rec["receiver_detail"].String(),
		},
		UserRemark:   rec["user_remark"].String(),
		SellerRemark: rec["seller_remark"].String(),
		StatusLogs:   statusLogs,
		CreatedAt:    rec["created_at"].String(),
	}, nil
}

// 取消方口径（012 修复轮 I6 补漏）——**两个列的枚举不同, 禁止混用**:
//
//	trade_order.cancel_type        1用户 2系统超时 3管理员（000006 列注释）
//	trade_order_log.operator_type  1系统 2用户 3管理员（000006 列注释）
//
// 修复前: cancel_type 恒写 1（后台/超时取消都被记成"用户取消"）, operator_type 用的却是
// cancel_type 的枚举（用户取消写成 1=系统）, 且 userId=0 把"系统超时"与"管理员"混为一谈。
const (
	cancelTypeUser   = 1
	cancelTypeSystem = 2
	cancelTypeAdmin  = 3

	opTypeSystem = 1
	opTypeUser   = 2
	opTypeAdmin  = 3
)

// Cancel 用户取消（C 端; 仅待付款）——ownerUserId 限定归属, 他人订单按不存在处理。
// 后台取消/超时取消走 cancelBy（取消方不同 → 审计两列不同）。
func (i *OrderLogicImpl) Cancel(ctx context.Context, userId int64, orderNo, reason string) error {
	return i.cancelBy(ctx, userId, orderNo, reason, cancelTypeUser, opTypeUser, fmt.Sprintf("user:%d", userId))
}

// cancelBy 取消实现（ownerUserId>0 时校验归属; 后台/系统视角传 0 不限归属）。
func (i *OrderLogicImpl) cancelBy(
	ctx context.Context, ownerUserId int64, orderNo, reason string, cancelType, opType int, opId string,
) error {
	m := dao.TradeOrder.Ctx(ctx).Where(dao.TradeOrder.Columns().OrderNo, orderNo)
	if ownerUserId > 0 {
		m = m.Where(dao.TradeOrder.Columns().UserId, ownerUserId)
	}
	rec, err := m.One()
	if err != nil {
		return err
	}
	if rec.IsEmpty() {
		return errcode.New(errcode.CodeOrderNotFound, "订单不存在")
	}
	if rec["status"].Int() != 10 {
		return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可取消")
	}
	res, err := dao.TradeOrder.Ctx(ctx).
		Where(dao.TradeOrder.Columns().OrderNo, orderNo).
		Where(dao.TradeOrder.Columns().Status, 10).
		Data(g.Map{
			dao.TradeOrder.Columns().Status:       90,
			dao.TradeOrder.Columns().CancelType:   cancelType,
			dao.TradeOrder.Columns().CancelReason: reason,
			dao.TradeOrder.Columns().CancelTime:   time.Now(),
		}).Update()
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可取消")
	}

	// 副作用: 释放库存（按行）; 状态流水
	items, _ := dao.TradeOrderItem.Ctx(ctx).
		Where(dao.TradeOrderItem.Columns().OrderNo, orderNo).All()
	for _, it := range items {
		skuId := it["sku_id"].Int64()
		qty := it["quantity"].Int()
		_, _ = g.DB().Exec(ctx,
			"UPDATE inventory SET locked=locked-? WHERE sku_id=? AND locked>=?", qty, skuId, qty)
	}
	// 评审 I6 修正: trade_order_log **无 operator 列**（只有 operator_type/operator_id）——
	// 原写法报 1054 且被吞, 致每次取消都没有状态流水。
	if _, e := dao.TradeOrderLog.Ctx(ctx).Data(do.TradeOrderLog{
		OrderNo:      orderNo,
		OrderId:      rec["id"].Int64(),
		FromStatus:   10,
		ToStatus:     90,
		Remark:       reason,
		OperatorType: opType,
		OperatorId:   opId,
	}).Insert(); e != nil {
		// 不再吞错（审计流水缺失必须可见）
		g.Log().Errorf(ctx, "订单取消流水写入失败 order_no=%s err=%v", orderNo, e)
	}
	return nil
}

// collectLines 行聚合：购物车项（普通）或直购 SKU → 下单行列表。
func (i *OrderLogicImpl) collectLines(ctx context.Context, tx gdb.TX, userId int64, in model.OrderCreateInput) ([]orderLine, int64, error) {
	var lines []orderLine
	var totalFen int64

	addSku := func(skuId int64, qty int) error {
		sku, err := dao.ProductSku.Ctx(ctx).
			Where(dao.ProductSku.Columns().Id, skuId).
			Where(dao.ProductSku.Columns().Deleted, 0).
			Where(dao.ProductSku.Columns().Status, 1).One()
		if err != nil {
			return err
		}
		if sku.IsEmpty() {
			return errcode.New(errcode.CodeProductNotFound, "商品不存在或已下架")
		}
		spu, err := dao.ProductSpu.Ctx(ctx).
			Where(dao.ProductSpu.Columns().Id, sku["spu_id"].Int64()).One()
		if err != nil {
			return err
		}
		priceFen, err := money.FromYuanString(sku["price"].String())
		if err != nil {
			return err
		}
		lines = append(lines, orderLine{
			SkuId: skuId, SkuNo: sku["sku_no"].String(), SpuId: sku["spu_id"].Int64(),
			SpuName: spu["name"].String(), SkuName: sku["name"].String(),
			SkuImage: sku["image"].String(),
			Specs:    specsMap(sku["specs"].String()),
			PriceFen: priceFen, LineFen: priceFen * int64(qty), Quantity: qty,
		})
		return nil
	}

	for _, itemId := range in.CartItemIds {
		item, err := dao.CartItem.Ctx(ctx).
			Where(dao.CartItem.Columns().Id, itemId).
			Where(dao.CartItem.Columns().UserId, userId).One()
		if err != nil {
			return nil, 0, err
		}
		if item.IsEmpty() {
			return nil, 0, errcode.New(errcode.CodeInvalidParam, "购物车项不存在")
		}
		if err := addSku(item["sku_id"].Int64(), item["quantity"].Int()); err != nil {
			return nil, 0, err
		}
	}
	if in.SkuId > 0 && in.Quantity > 0 {
		if err := addSku(in.SkuId, in.Quantity); err != nil {
			return nil, 0, err
		}
	}
	for _, l := range lines {
		totalFen += l.LineFen
	}
	return lines, totalFen, nil
}

// provinceRestricted SPU 是否在指定省份限售。
func provinceRestricted(ctx context.Context, spuId int64, provinceCode string) bool {
	if provinceCode == "" {
		return false
	}
	v, err := dao.ProductSpu.Ctx(ctx).
		Where(dao.ProductSpu.Columns().Id, spuId).
		Value(dao.ProductSpu.Columns().SaleRestrictCodes)
	if err != nil || v == nil || v.IsNil() {
		return false
	}
	return strings.Contains(v.String(), provinceCode)
}

// Confirm 确认收货（012 FR-004）: 待收货(30) → 已完成(40) + 状态流水; 非 30 → 40006; 他人订单 → 40005。
// 佣金计提事件位: 完成态由分销结算（批次 11 / 定时消费）按订单完成推进, 本方法只落状态与流水。
func (i *OrderLogicImpl) Confirm(ctx context.Context, userId int64, orderNo string) error {
	cols := dao.TradeOrder.Columns()
	rec, err := dao.TradeOrder.Ctx(ctx).
		Where(cols.OrderNo, orderNo).
		Where(cols.UserId, userId).One()
	if err != nil {
		return err
	}
	if rec.IsEmpty() {
		return errcode.New(errcode.CodeOrderNotFound, "订单不存在")
	}
	if rec[cols.Status].Int() != 30 {
		return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可确认收货")
	}
	// 条件更新防并发重复确认
	res, err := dao.TradeOrder.Ctx(ctx).
		Where(cols.OrderNo, orderNo).
		Where(cols.Status, 30).
		Data(do.TradeOrder{Status: 40, FinishTime: gtime.Now()}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "确认收货失败")
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errcode.New(errcode.CodeStatusNotAllowed, "当前状态不可确认收货")
	}
	_, _ = dao.TradeOrderLog.Ctx(ctx).Data(do.TradeOrderLog{
		OrderNo:      orderNo,
		OrderId:      rec[cols.Id].Int64(),
		FromStatus:   30,
		ToStatus:     40,
		Remark:       "确认收货",
		OperatorType: 2, // 2=用户（表注释: 1系统 2用户 3管理员）
		OperatorId:   fmt.Sprintf("user:%d", userId),
	}).Insert()
	return nil
}
