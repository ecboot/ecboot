package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 分享行为上报（归因来源记录; 游客可报, 携带会员凭证则归因到人）
	ShareReportReq struct {
		g.Meta     `path:"/shares" method:"POST" summary:"分享行为上报"` // 需登录（017 修复: 游客无分享者, 上报拒绝 10003）
		SpuId      string `json:"spuId" dc:"分享商品SPU(整页分享为空)"`
		Channel    int    `json:"channel" v:"required|in:1,2,3,4" dc:"渠道:1小程序卡片 2海报 3口令 4朋友圈社群"`
		SceneValue string `json:"sceneValue" dc:"小程序场景值"`
	}
	ShareReportRes struct {
		Success bool `json:"success" dc:"固定true"`
	}
)
