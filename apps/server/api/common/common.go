// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package common

import (
	"context"

	"ecboot/api/common/v1"
)

type ICommonV1 interface {
	GetCaptcha(ctx context.Context, req *v1.GetCaptchaReq) (res *v1.GetCaptchaRes, err error)
	Ping(ctx context.Context, req *v1.PingReq) (res *v1.PingRes, err error)
	ShareReport(ctx context.Context, req *v1.ShareReportReq) (res *v1.ShareReportRes, err error)
	GetSmsCode(ctx context.Context, req *v1.GetSmsCodeReq) (res *v1.GetSmsCodeRes, err error)
	StoreList(ctx context.Context, req *v1.StoreListReq) (res *v1.StoreListRes, err error)
	StoreDetail(ctx context.Context, req *v1.StoreDetailReq) (res *v1.StoreDetailRes, err error)
}
