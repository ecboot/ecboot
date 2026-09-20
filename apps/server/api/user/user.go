// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package user

import (
	"context"

	"ecboot/api/user/v1"
)

type IUserV1 interface {
	SmsLogin(ctx context.Context, req *v1.SmsLoginReq) (res *v1.SmsLoginRes, err error)
}
