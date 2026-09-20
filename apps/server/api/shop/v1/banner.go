package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 轮播/弹窗（投放时段内, 按位置）
	BannerListReq struct {
		g.Meta   `path:"/banners" method:"GET" summary:"运营位列表"`
		Position int `json:"position" v:"required|in:1,2" dc:"位置:1首页轮播 2首页弹窗"`
	}
	BannerItem struct {
		Id       string `json:"id"`
		ImageUrl string `json:"imageUrl" dc:"图片"`
		LinkUrl  string `json:"linkUrl" dc:"跳转链接"`
	}
	BannerListRes struct {
		List []BannerItem `json:"list"`
	}

	// 楼层内容（含商品摘要装配）
	FloorListReq struct {
		g.Meta `path:"/floors" method:"GET" summary:"楼层内容"`
	}
	FloorProduct struct {
		SpuId string `json:"spuId"`
		Name  string `json:"name"`
		Image string `json:"image"`
		Price string `json:"price" dc:"价格(元)"`
	}
	FloorItem struct {
		FloorId  string         `json:"floorId"`
		FloorType int            `json:"floorType" dc:"类型:1金刚区 2商品楼层 3专题"`
		Title    string         `json:"title"`
		Products []FloorProduct `json:"products" dc:"商品楼层装配"`
		Config   map[string]any `json:"config" dc:"楼层配置(金刚区入口等)"`
	}
	FloorListRes struct {
		List []FloorItem `json:"list"`
	}
)
