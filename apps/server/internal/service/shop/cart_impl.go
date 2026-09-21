// cart_impl.go ICartLogic 实现（购物车与结算试算）。
// 试算口径: 实时价态; 失效行分离; 限售校验; 恒等式 payAmount ≡ total − promotion + freight。
package shop

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/library/money"
	"ecboot/internal/model"
)

// CartLogicImpl ICartLogic 实现。
type CartLogicImpl struct{}

func NewCartLogic() *CartLogicImpl { return &CartLogicImpl{} }

func (i *CartLogicImpl) Detail(ctx context.Context, userId int64) (*model.CartView, error) {
	all, err := dao.CartItem.Ctx(ctx).
		Where(dao.CartItem.Columns().UserId, userId).
		Order("id desc").All()
	if err != nil {
		return nil, err
	}
	view := &model.CartView{Items: []model.CartLine{}}
	for _, r := range all {
		line, err := i.hydrateLine(ctx, r["id"].Int64(), r["sku_id"].Int64(), r["quantity"].Int(), r["checked"].Int() == 1)
		if err != nil {
			return nil, err
		}
		view.Items = append(view.Items, *line)
	}
	return view, nil
}

func (i *CartLogicImpl) AddItem(ctx context.Context, userId, skuId int64, quantity int) (int64, error) {
	if quantity <= 0 || quantity > 99 {
		return 0, errcode.New(errcode.CodeInvalidParam, "数量需在1~99之间")
	}
	existing, err := dao.CartItem.Ctx(ctx).
		Where(dao.CartItem.Columns().UserId, userId).
		Where(dao.CartItem.Columns().SkuId, skuId).One()
	if err != nil {
		return 0, err
	}
	if !existing.IsEmpty() {
		newQty := existing["quantity"].Int() + quantity
		if newQty > 99 {
			newQty = 99
		}
		_, err = dao.CartItem.Ctx(ctx).Where(dao.CartItem.Columns().Id, existing["id"].Int64()).
			Data(g.Map{dao.CartItem.Columns().Quantity: newQty}).Update()
		return existing["id"].Int64(), err
	}
	res, err := dao.CartItem.Ctx(ctx).Data(g.Map{
		dao.CartItem.Columns().UserId: userId, dao.CartItem.Columns().SkuId: skuId,
		dao.CartItem.Columns().Quantity: quantity, dao.CartItem.Columns().Checked: 1,
	}).InsertAndGetId()
	return res, err
}

func (i *CartLogicImpl) UpdateItem(ctx context.Context, userId, itemId int64, quantity int, checked *bool) error {
	data := g.Map{}
	if quantity > 0 {
		if quantity > 99 {
			return errcode.New(errcode.CodeInvalidParam, "数量上限99")
		}
		data[dao.CartItem.Columns().Quantity] = quantity
	}
	if checked != nil {
		data[dao.CartItem.Columns().Checked] = *checked
	}
	if len(data) == 0 {
		return nil
	}
	_, err := dao.CartItem.Ctx(ctx).Where(dao.CartItem.Columns().Id, itemId).
		Where(dao.CartItem.Columns().UserId, userId).Data(data).Update()
	return err
}

func (i *CartLogicImpl) RemoveItem(ctx context.Context, userId, itemId int64) error {
	_, err := dao.CartItem.Ctx(ctx).
		Where(dao.CartItem.Columns().Id, itemId).
		Where(dao.CartItem.Columns().UserId, userId).Delete()
	return err
}

// hydrateLine 购物车行实时价态（SKU 联查）。
func (i *CartLogicImpl) hydrateLine(ctx context.Context, itemId, skuId int64, quantity int, checked bool) (*model.CartLine, error) {
	sku, err := dao.ProductSku.Ctx(ctx).
		Where(dao.ProductSku.Columns().Id, skuId).
		Where(dao.ProductSku.Columns().Deleted, 0).One()
	if err != nil {
		return nil, err
	}
	sellable := !sku.IsEmpty() && sku["status"].Int() == 1
	if sellable {
		inv, _ := dao.Inventory.Ctx(ctx).
			Where(dao.Inventory.Columns().SkuId, skuId).One()
		if inv.IsEmpty() || inv["total"].Int()-inv["locked"].Int() <= 0 {
			sellable = false
		}
	}
	line := &model.CartLine{
		ItemId:   itemId,
		SkuId:    skuId,
		Specs:    specsMap(sku["specs"].String()),
		Image:    sku["image"].String(),
		Price:    sku["price"].String(),
		Sellable: sellable,
		Quantity: quantity,
		Checked:  checked,
	}
	spu, _ := dao.ProductSpu.Ctx(ctx).
		Where(dao.ProductSpu.Columns().Id, sku["spu_id"].Int64()).One()
	if !spu.IsEmpty() {
		line.SpuName = spu["name"].String()
	}
	return line, nil
}

// Checkout 结算试算（FR-002/FR-003: AmountBook 全分项 + 勾稽 + 失效提示 + 限售校验）。
func (i *CartLogicImpl) Checkout(ctx context.Context, userId int64, q model.CheckoutQuery) (*model.CheckoutResult, error) {
	all, err := dao.CartItem.Ctx(ctx).
		Where(dao.CartItem.Columns().UserId, userId).
		Where(dao.CartItem.Columns().Checked, 1).All()
	if err != nil {
		return nil, err
	}
	result := &model.CheckoutResult{Items: []model.CartLine{}, Errors: []string{}}

	// 收货地址（限售区域口径）
	var provinceCode string
	if q.AddressId > 0 {
		addr, _ := dao.UserAddress.Ctx(ctx).
			Where(dao.UserAddress.Columns().Id, q.AddressId).
			Where(dao.UserAddress.Columns().UserId, userId).One()
		if !addr.IsEmpty() {
			provinceCode = addr["province_code"].String()
		}
	}

	var totalFen int64
	for _, r := range all {
		sku, err := dao.ProductSku.Ctx(ctx).
			Where(dao.ProductSku.Columns().Id, r["sku_id"].Int64()).One()
		if err != nil {
			return nil, err
		}
		if sku.IsEmpty() || sku["status"].Int() != 1 || sku["deleted"].Int() == 1 {
			result.Errors = append(result.Errors, "存在已失效商品,已自动置灰")
			continue
		}
		spu, err := dao.ProductSpu.Ctx(ctx).
			Where(dao.ProductSpu.Columns().Id, sku["spu_id"].Int64()).One()
		if err != nil {
			return nil, err
		}
		if spu.IsEmpty() || spu["status"].Int() != 1 {
			result.Errors = append(result.Errors, "存在已下架商品")
			continue
		}
		// 限售校验（收货省 ∈ SPU 禁售列表）
		if provinceCode != "" {
			var restrict []string
			_ = gconv.Struct(spu["sale_restrict_codes"].String(), &restrict)
			for _, code := range restrict {
				if code == provinceCode {
					result.Errors = append(result.Errors, "该商品在收货地区限售")
					break
				}
			}
		}
		priceFen, err := money.FromYuanString(sku["price"].String())
		if err != nil {
			return nil, err
		}
		totalFen += priceFen * r["quantity"].Int64()
		result.Items = append(result.Items, model.CartLine{
			ItemId: r["id"].Int64(), SkuId: r["sku_id"].Int64(),
			SpuName: spu["name"].String(),
			Specs:   specsMap(sku["specs"].String()),
			Price:   sku["price"].String(), Quantity: r["quantity"].Int(),
		})
	}

	// 优惠计算（满减命中最优档 → 券门槛校验 → 积分抵扣）
	fullReductionFen := calcFullReductionFen(ctx, provinceCode, totalFen)
	couponFen := calcCouponDiscountFen(ctx, userId, q.CouponId, totalFen)
	pointFen := calcPointDeductFen(ctx, userId, q.UsePoint, totalFen-couponFen-fullReductionFen)

	// 运费 V1: 全包邮（运费模板计算随交易完善）
	freightFen := int64(0)

	promotionFen := couponFen + fullReductionFen + pointFen
	payFen := totalFen - promotionFen + freightFen
	if payFen < 0 {
		payFen = 0
	}

	result.Amount = model.AmountBook{
		TotalAmount:         money.ToYuanString(totalFen),
		CouponAmount:        money.ToYuanString(couponFen),
		FullReductionAmount: money.ToYuanString(fullReductionFen),
		PointAmount:         money.ToYuanString(pointFen),
		FreightAmount:       money.ToYuanString(freightFen),
		PayAmount:           money.ToYuanString(payFen),
	}
	return result, nil
}
