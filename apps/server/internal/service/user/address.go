// address.go 收货地址域——表: user_address（V1/V23 区划码）。
// 规则: 区划码为运费匹配口径（名称仅展示快照）; 默认地址每用户至多一个（应用层保证: 设默认时先清旧）。
package user

import (
	"context"

	"ecboot/internal/model"
)

// IAddressLogic 收货地址。
type IAddressLogic interface {
	List(ctx context.Context, userId int64) ([]model.AddressItem, error)
	Create(ctx context.Context, userId int64, in model.AddressInput) (int64, error)
	Update(ctx context.Context, userId, addressId int64, in model.AddressInput) error
	Delete(ctx context.Context, userId, addressId int64) error
	// SetDefault 设默认（同事务清同用户其他默认标记）。
	SetDefault(ctx context.Context, userId, addressId int64) error
	// GetForOrder 下单取址（校验归属, 返回含区划码——运费计算依据）。
	GetForOrder(ctx context.Context, userId, addressId int64) (*model.AddressItem, error)
}
