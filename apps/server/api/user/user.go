// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package user

import (
	"context"

	"ecboot/api/user/v1"
)

type IUserV1 interface {
	AddressList(ctx context.Context, req *v1.AddressListReq) (res *v1.AddressListRes, err error)
	AddressCreate(ctx context.Context, req *v1.AddressCreateReq) (res *v1.AddressCreateRes, err error)
	AddressUpdate(ctx context.Context, req *v1.AddressUpdateReq) (res *v1.AddressUpdateRes, err error)
	AddressDelete(ctx context.Context, req *v1.AddressDeleteReq) (res *v1.AddressDeleteRes, err error)
	AddressSetDefault(ctx context.Context, req *v1.AddressSetDefaultReq) (res *v1.AddressSetDefaultRes, err error)
	SmsLogin(ctx context.Context, req *v1.SmsLoginReq) (res *v1.SmsLoginRes, err error)
	WxLogin(ctx context.Context, req *v1.WxLoginReq) (res *v1.WxLoginRes, err error)
	TokenRefresh(ctx context.Context, req *v1.TokenRefreshReq) (res *v1.TokenRefreshRes, err error)
	Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error)
	Me(ctx context.Context, req *v1.MeReq) (res *v1.MeRes, err error)
	CouponAvailableList(ctx context.Context, req *v1.CouponAvailableListReq) (res *v1.CouponAvailableListRes, err error)
	CouponReceive(ctx context.Context, req *v1.CouponReceiveReq) (res *v1.CouponReceiveRes, err error)
	MyCouponList(ctx context.Context, req *v1.MyCouponListReq) (res *v1.MyCouponListRes, err error)
	DistApply(ctx context.Context, req *v1.DistApplyReq) (res *v1.DistApplyRes, err error)
	DistStatus(ctx context.Context, req *v1.DistStatusReq) (res *v1.DistStatusRes, err error)
	DistRelation(ctx context.Context, req *v1.DistRelationReq) (res *v1.DistRelationRes, err error)
	DistRuleQuery(ctx context.Context, req *v1.DistRuleQueryReq) (res *v1.DistRuleQueryRes, err error)
	DistRecordList(ctx context.Context, req *v1.DistRecordListReq) (res *v1.DistRecordListRes, err error)
	DistAccount(ctx context.Context, req *v1.DistAccountReq) (res *v1.DistAccountRes, err error)
	DistAccountLogList(ctx context.Context, req *v1.DistAccountLogListReq) (res *v1.DistAccountLogListRes, err error)
	WithdrawCreate(ctx context.Context, req *v1.WithdrawCreateReq) (res *v1.WithdrawCreateRes, err error)
	WithdrawList(ctx context.Context, req *v1.WithdrawListReq) (res *v1.WithdrawListRes, err error)
	InviteRecordList(ctx context.Context, req *v1.InviteRecordListReq) (res *v1.InviteRecordListRes, err error)
	ShareCode(ctx context.Context, req *v1.ShareCodeReq) (res *v1.ShareCodeRes, err error)
	FavoriteList(ctx context.Context, req *v1.FavoriteListReq) (res *v1.FavoriteListRes, err error)
	FavoriteAdd(ctx context.Context, req *v1.FavoriteAddReq) (res *v1.FavoriteAddRes, err error)
	FavoriteRemove(ctx context.Context, req *v1.FavoriteRemoveReq) (res *v1.FavoriteRemoveRes, err error)
	FootprintList(ctx context.Context, req *v1.FootprintListReq) (res *v1.FootprintListRes, err error)
	FootprintClear(ctx context.Context, req *v1.FootprintClearReq) (res *v1.FootprintClearRes, err error)
	MessageList(ctx context.Context, req *v1.MessageListReq) (res *v1.MessageListRes, err error)
	MessageRead(ctx context.Context, req *v1.MessageReadReq) (res *v1.MessageReadRes, err error)
	MessageReadAll(ctx context.Context, req *v1.MessageReadAllReq) (res *v1.MessageReadAllRes, err error)
	NotifyPreferenceGet(ctx context.Context, req *v1.NotifyPreferenceGetReq) (res *v1.NotifyPreferenceGetRes, err error)
	NotifyPreferenceSet(ctx context.Context, req *v1.NotifyPreferenceSetReq) (res *v1.NotifyPreferenceSetRes, err error)
	ProfileDetail(ctx context.Context, req *v1.ProfileDetailReq) (res *v1.ProfileDetailRes, err error)
	ProfileUpdate(ctx context.Context, req *v1.ProfileUpdateReq) (res *v1.ProfileUpdateRes, err error)
	PointAccount(ctx context.Context, req *v1.PointAccountReq) (res *v1.PointAccountRes, err error)
	PointLogList(ctx context.Context, req *v1.PointLogListReq) (res *v1.PointLogListRes, err error)
}
