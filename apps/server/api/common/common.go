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
	GetSmsCode(ctx context.Context, req *v1.GetSmsCodeReq) (res *v1.GetSmsCodeRes, err error)
}
