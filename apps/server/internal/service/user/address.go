// address.go 收货地址域——表: user_address（V1/V23 区划码）。
// 规则: 区划码为运费匹配口径（名称仅展示快照）; 默认地址每用户至多一个（应用层保证: 设默认时先清旧）。
package user

import "context"

// IAddressLogic 收货地址。
type IAddressLogic interface {
	List(ctx context.Context, userId int64) ([]AddressItem, error)
	Create(ctx context.Context, userId int64, in AddressInput) (int64, error)
	Update(ctx context.Context, userId, addressId int64, in AddressInput) error
	Delete(ctx context.Context, userId, addressId int64) error
	// SetDefault 设默认（同事务清同用户其他默认标记）。
	SetDefault(ctx context.Context, userId, addressId int64) error
	// GetForOrder 下单取址（校验归属, 返回含区划码——运费计算依据）。
	GetForOrder(ctx context.Context, userId, addressId int64) (*AddressItem, error)
}

type AddressInput struct {
	ReceiverName  string
	ReceiverPhone string
	Province      string
	City          string
	District      string
	DetailAddress string
	ProvinceCode  string
	CityCode      string
	DistrictCode  string
	IsDefault     bool
}

type AddressItem struct {
	Id            int64  `json:"id"`
	ReceiverName  string `json:"receiverName"`
	ReceiverPhone string `json:"receiverPhone" dc:"脱敏"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	DetailAddress string `json:"detailAddress"`
	ProvinceCode  string `json:"provinceCode"`
	CityCode      string `json:"cityCode"`
	DistrictCode  string `json:"districtCode"`
	IsDefault     bool   `json:"isDefault"`
}
