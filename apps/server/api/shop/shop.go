// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package shop

import (
	"context"

	"ecboot/api/shop/v1"
)

type IShopV1 interface {
	Index(ctx context.Context, req *v1.IndexReq) (res *v1.IndexRes, err error)
}
