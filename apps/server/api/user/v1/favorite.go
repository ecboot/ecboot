package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/model"
)

type (
	FavoriteItem struct {
		SpuId    string `json:"spuId"`
		SpuName  string `json:"spuName" dc:"商品名"`
		Image    string `json:"image" dc:"主图"`
		Price    string `json:"price" dc:"现价(元)"`
		Sellable bool   `json:"sellable" dc:"可售"`
		Invalid  bool   `json:"invalid" dc:"已失效(下架)"`
	}

	FavoriteListReq struct {
		g.Meta `path:"/favorites" method:"GET" summary:"收藏列表"`
		model.PageReq
	}
	FavoriteListRes struct {
		model.PageRes
		List []FavoriteItem `json:"list"`
	}

	// 收藏（已软删则复活, 重复收藏幂等）
	FavoriteAddReq struct {
		g.Meta `path:"/favorites/{spuId}" method:"PUT" summary:"收藏商品"`
		SpuId  string `json:"spuId" v:"required" dc:"SPU ID"`
	}
	FavoriteAddRes struct {
		Success bool `json:"success"`
	}

	FavoriteRemoveReq struct {
		g.Meta `path:"/favorites/{spuId}" method:"DELETE" summary:"取消收藏"`
		SpuId  string `json:"spuId" v:"required" dc:"SPU ID"`
	}
	FavoriteRemoveRes struct {
		Success bool `json:"success"`
	}
)
