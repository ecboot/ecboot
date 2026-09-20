package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 发起砍价（会员; 选择场次商品开一刀）
	BargainLaunchReq struct {
		g.Meta        `path:"/bargains" method:"POST" summary:"发起砍价"`
		BargainItemId string `json:"bargainItemId" v:"required" dc:"砍价场次商品ID"`
	}
	BargainLaunchRes struct {
		RecordId     string `json:"recordId" dc:"砍价单ID"`
		CurrentPrice string `json:"currentPrice" dc:"当前价(发起即首刀,元)"`
	}

	// 砍价进度（公开; 含帮砍列表）
	BargainProgressReq struct {
		g.Meta   `path:"/bargains/{recordId}" method:"GET" summary:"砍价进度"`
		RecordId string `json:"recordId" v:"required" dc:"砍价单ID"`
	}
	BargainHelper struct {
		UserId    string `json:"userId"`
		Nickname  string `json:"nickname" dc:"帮砍人(脱敏)"`
		CutAmount string `json:"cutAmount" dc:"本刀金额"`
		CreatedAt string `json:"createdAt"`
	}
	BargainProgressRes struct {
		RecordId      string          `json:"recordId"`
		SkuId         string          `json:"skuId"`
		OriginalPrice string          `json:"originalPrice" dc:"起始价"`
		CurrentPrice  string          `json:"currentPrice" dc:"当前价"`
		FloorPrice    string          `json:"floorPrice" dc:"底价"`
		CutCount      int             `json:"cutCount" dc:"已砍刀数"`
		Status        int             `json:"status" dc:"1砍价中 2到底价 3已下单 4超时 5取消"`
		ExpireTime    string          `json:"expireTime" dc:"截止"`
		OrderNo       string          `json:"orderNo" dc:"成交订单号"`
		Helpers       []BargainHelper `json:"helpers" dc:"帮砍列表"`
	}

	// 帮砍（会员; 一人一刀; 风控挂载位: risk_rule 2/3/4）
	BargainCutReq struct {
		g.Meta   `path:"/bargains/{recordId}/cut" method:"POST" summary:"帮砍一刀"`
		RecordId string `json:"recordId" v:"required" dc:"砍价单ID"`
	}
	BargainCutRes struct {
		CutAmount    string `json:"cutAmount" dc:"本刀砍掉金额"`
		CurrentPrice string `json:"currentPrice" dc:"砍后当前价"`
		FloorReached bool   `json:"floorReached" dc:"是否已到底价(可下单)"`
	}
)
