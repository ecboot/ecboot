package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	FootprintItem struct {
		SpuId      string `json:"spuId"`
		SpuName    string `json:"spuName"`
		Image      string `json:"image"`
		Price      string `json:"price" dc:"现价(元)"`
		LastViewAt string `json:"lastViewAt" dc:"最近浏览时间"`
	}

	FootprintListReq struct {
		g.Meta `path:"/footprints" method:"GET" summary:"浏览足迹"`
		PageReq
	}
	FootprintListRes struct {
		PageRes
		List []FootprintItem `json:"list"`
	}

	FootprintClearReq struct {
		g.Meta `path:"/footprints" method:"DELETE" summary:"清空足迹"`
	}
	FootprintClearRes struct {
		Success bool `json:"success"`
	}
)
