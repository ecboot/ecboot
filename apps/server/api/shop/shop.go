// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package shop

import (
	"context"

	"ecboot/api/shop/v1"
)

type IShopV1 interface {
	GroupBuyList(ctx context.Context, req *v1.GroupBuyListReq) (res *v1.GroupBuyListRes, err error)
	FlashSaleList(ctx context.Context, req *v1.FlashSaleListReq) (res *v1.FlashSaleListRes, err error)
	BargainActivityList(ctx context.Context, req *v1.BargainActivityListReq) (res *v1.BargainActivityListRes, err error)
	AssistList(ctx context.Context, req *v1.AssistListReq) (res *v1.AssistListRes, err error)
	FullReductionList(ctx context.Context, req *v1.FullReductionListReq) (res *v1.FullReductionListRes, err error)
	AfterSaleCreate(ctx context.Context, req *v1.AfterSaleCreateReq) (res *v1.AfterSaleCreateRes, err error)
	AfterSaleList(ctx context.Context, req *v1.AfterSaleListReq) (res *v1.AfterSaleListRes, err error)
	AfterSaleDetail(ctx context.Context, req *v1.AfterSaleDetailReq) (res *v1.AfterSaleDetailRes, err error)
	AfterSaleCancel(ctx context.Context, req *v1.AfterSaleCancelReq) (res *v1.AfterSaleCancelRes, err error)
	AfterSaleLogistics(ctx context.Context, req *v1.AfterSaleLogisticsReq) (res *v1.AfterSaleLogisticsRes, err error)
	BannerList(ctx context.Context, req *v1.BannerListReq) (res *v1.BannerListRes, err error)
	FloorList(ctx context.Context, req *v1.FloorListReq) (res *v1.FloorListRes, err error)
	CartDetail(ctx context.Context, req *v1.CartDetailReq) (res *v1.CartDetailRes, err error)
	CartAddItem(ctx context.Context, req *v1.CartAddItemReq) (res *v1.CartAddItemRes, err error)
	CartUpdateItem(ctx context.Context, req *v1.CartUpdateItemReq) (res *v1.CartUpdateItemRes, err error)
	CartRemoveItem(ctx context.Context, req *v1.CartRemoveItemReq) (res *v1.CartRemoveItemRes, err error)
	CartCheckout(ctx context.Context, req *v1.CartCheckoutReq) (res *v1.CartCheckoutRes, err error)
	Index(ctx context.Context, req *v1.IndexReq) (res *v1.IndexRes, err error)
	OrderCreate(ctx context.Context, req *v1.OrderCreateReq) (res *v1.OrderCreateRes, err error)
	OrderList(ctx context.Context, req *v1.OrderListReq) (res *v1.OrderListRes, err error)
	OrderDetail(ctx context.Context, req *v1.OrderDetailReq) (res *v1.OrderDetailRes, err error)
	OrderCancel(ctx context.Context, req *v1.OrderCancelReq) (res *v1.OrderCancelRes, err error)
	OrderConfirm(ctx context.Context, req *v1.OrderConfirmReq) (res *v1.OrderConfirmRes, err error)
	PayCreate(ctx context.Context, req *v1.PayCreateReq) (res *v1.PayCreateRes, err error)
	PayStatus(ctx context.Context, req *v1.PayStatusReq) (res *v1.PayStatusRes, err error)
	PayNotify(ctx context.Context, req *v1.PayNotifyReq) (res *v1.PayNotifyRes, err error)
	RefundNotify(ctx context.Context, req *v1.RefundNotifyReq) (res *v1.RefundNotifyRes, err error)
	CategoryTree(ctx context.Context, req *v1.CategoryTreeReq) (res *v1.CategoryTreeRes, err error)
	BrandList(ctx context.Context, req *v1.BrandListReq) (res *v1.BrandListRes, err error)
	ProductList(ctx context.Context, req *v1.ProductListReq) (res *v1.ProductListRes, err error)
	ProductDetail(ctx context.Context, req *v1.ProductDetailReq) (res *v1.ProductDetailRes, err error)
	ProductSearch(ctx context.Context, req *v1.ProductSearchReq) (res *v1.ProductSearchRes, err error)
	ProductReviewList(ctx context.Context, req *v1.ProductReviewListReq) (res *v1.ProductReviewListRes, err error)
	ReviewCreate(ctx context.Context, req *v1.ReviewCreateReq) (res *v1.ReviewCreateRes, err error)
	ReviewExtra(ctx context.Context, req *v1.ReviewExtraReq) (res *v1.ReviewExtraRes, err error)
	MyReviewList(ctx context.Context, req *v1.MyReviewListReq) (res *v1.MyReviewListRes, err error)
}
