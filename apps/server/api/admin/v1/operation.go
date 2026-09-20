package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 轮播管理
	AdminBannerListReq struct {
		g.Meta `path:"/banners" method:"GET" summary:"轮播列表"`
		Position int `json:"position" dc:"位置筛选"`
		PageReq
	}
	AdminBannerItem struct {
		Id        string `json:"id"`
		Position  int    `json:"position" dc:"1轮播 2弹窗"`
		ImageUrl  string `json:"imageUrl"`
		LinkUrl   string `json:"linkUrl"`
		Sort      int    `json:"sort"`
		StartTime string `json:"startTime" dc:"投放起(空=立即)"`
		EndTime   string `json:"endTime" dc:"投放止(空=长期)"`
		Status    int    `json:"status"`
	}
	AdminBannerListRes struct {
		PageRes
		List []AdminBannerItem `json:"list"`
	}

	AdminBannerCreateReq struct {
		g.Meta   `path:"/banners" method:"POST" summary:"新增轮播"`
		Position int    `json:"position" v:"required|in:1,2" dc:"位置"`
		ImageUrl string `json:"imageUrl" v:"required" dc:"图片"`
		LinkUrl  string `json:"linkUrl" dc:"跳转"`
		Sort     int    `json:"sort" dc:"排序"`
		StartTime string `json:"startTime" dc:"投放起"`
		EndTime  string `json:"endTime" dc:"投放止"`
	}
	AdminBannerCreateRes struct {
		Id string `json:"id"`
	}

	AdminBannerUpdateReq struct {
		g.Meta   `path:"/banners/{id}" method:"PUT" summary:"修改轮播"`
		Id       string `json:"id" v:"required" dc:"ID"`
		ImageUrl string `json:"imageUrl" dc:"图片"`
		LinkUrl  string `json:"linkUrl" dc:"跳转"`
		Sort     int    `json:"sort" dc:"排序"`
		StartTime string `json:"startTime" dc:"投放起"`
		EndTime  string `json:"endTime" dc:"投放止"`
		Status   int    `json:"status" dc:"状态"`
	}
	AdminBannerUpdateRes struct {
		Success bool `json:"success"`
	}

	AdminBannerDeleteReq struct {
		g.Meta `path:"/banners/{id}" method:"DELETE" summary:"删除轮播(软删)"`
		Id     string `json:"id" v:"required" dc:"ID"`
	}
	AdminBannerDeleteRes struct {
		Success bool `json:"success"`
	}

	// 楼层管理
	AdminFloorListReq struct {
		g.Meta `path:"/floors" method:"GET" summary:"楼层列表"`
		PageReq
	}
	AdminFloorItem struct {
		Id        string         `json:"id"`
		FloorType int            `json:"floorType" dc:"1金刚区 2商品楼层 3专题"`
		Title     string         `json:"title"`
		Config    map[string]any `json:"config" dc:"配置(入口数组/商品ID列表)"`
		Sort      int            `json:"sort"`
		Status    int            `json:"status"`
	}
	AdminFloorListRes struct {
		PageRes
		List []AdminFloorItem `json:"list"`
	}

	AdminFloorCreateReq struct {
		g.Meta   `path:"/floors" method:"POST" summary:"新增楼层"`
		FloorType int            `json:"floorType" v:"required|in:1,2,3" dc:"类型"`
		Title    string         `json:"title" dc:"标题"`
		Config   map[string]any `json:"config" dc:"配置"`
		Sort     int            `json:"sort" dc:"排序"`
	}
	AdminFloorCreateRes struct {
		Id string `json:"id"`
	}

	AdminFloorUpdateReq struct {
		g.Meta   `path:"/floors/{id}" method:"PUT" summary:"修改楼层"`
		Id       string         `json:"id" v:"required" dc:"ID"`
		Title    string         `json:"title" dc:"标题"`
		Config   map[string]any `json:"config" dc:"配置"`
		Sort     int            `json:"sort" dc:"排序"`
		Status   int            `json:"status" dc:"状态"`
	}
	AdminFloorUpdateRes struct {
		Success bool `json:"success"`
	}

	AdminFloorDeleteReq struct {
		g.Meta `path:"/floors/{id}" method:"DELETE" summary:"删除楼层(软删)"`
		Id     string `json:"id" v:"required" dc:"ID"`
	}
	AdminFloorDeleteRes struct {
		Success bool `json:"success"`
	}
)
