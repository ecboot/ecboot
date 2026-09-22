package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/model"
)

type (
	// 拼团活动列表（进行中）
	GroupBuyListReq struct {
		g.Meta `path:"/activities/group-buys" method:"GET" summary:"拼团活动列表"`
		model.PageReq
	}
	GroupBuySku struct {
		SkuId      string `json:"skuId"`
		GroupPrice string `json:"groupPrice" dc:"成团价(元)"`
	}
	GroupBuyItem struct {
		ActivityId string        `json:"activityId"`
		Name       string        `json:"name" dc:"活动名称"`
		SpuId      string        `json:"spuId" dc:"商品SPU"`
		SpuName    string        `json:"spuName" dc:"商品名称"`
		Image      string        `json:"image" dc:"商品主图"`
		Items      []GroupBuySku `json:"items" dc:"SKU成团价"`
		GroupSize  int           `json:"groupSize" dc:"成团人数"`
		EndTime    string        `json:"endTime" dc:"活动截止"`
	}
	GroupBuyListRes struct {
		model.PageRes
		List []GroupBuyItem `json:"list"`
	}

	// 秒杀列表（进行中与预告）
	FlashSaleListReq struct {
		g.Meta `path:"/activities/flash-sales" method:"GET" summary:"秒杀活动列表"`
		model.PageReq
	}
	FlashSaleSku struct {
		SkuId       string `json:"skuId"`
		FlashPrice  string `json:"flashPrice" dc:"秒杀价(元)"`
		StockRemain int    `json:"stockRemain" dc:"活动剩余量"`
		PerLimit    int    `json:"perLimit" dc:"每人限购"`
	}
	FlashSaleItem struct {
		ActivityId string         `json:"activityId"`
		Name       string         `json:"name"`
		StartTime  string         `json:"startTime"`
		EndTime    string         `json:"endTime"`
		Upcoming   bool           `json:"upcoming" dc:"是否预告(未开始; 015 修复轮微扩)"`
		Items      []FlashSaleSku `json:"items" dc:"场次商品"`
	}
	FlashSaleListRes struct {
		model.PageRes
		List []FlashSaleItem `json:"list"`
	}

	// 砍价活动列表
	BargainActivityListReq struct {
		g.Meta `path:"/activities/bargains" method:"GET" summary:"砍价活动列表"`
		model.PageReq
	}
	BargainSku struct {
		ItemId        string `json:"itemId" dc:"场次商品ID(发起砍价必传; 015 微扩)"`
		SkuId         string `json:"skuId"`
		OriginalPrice string `json:"originalPrice" dc:"起始价"`
		FloorPrice    string `json:"floorPrice" dc:"底价"`
		MaxCutCount   int    `json:"maxCutCount" dc:"最大刀数(015 微扩; 展示用)"`
	}
	BargainActivityItem struct {
		ActivityId string       `json:"activityId"`
		Name       string       `json:"name"`
		SpuId      string       `json:"spuId"`
		SpuName    string       `json:"spuName"`
		Image      string       `json:"image"`
		Items      []BargainSku `json:"items" dc:"SKU价格区间"`
		EndTime    string       `json:"endTime"`
	}
	BargainActivityListRes struct {
		model.PageRes
		List []BargainActivityItem `json:"list"`
	}

	// 助力活动列表
	AssistListReq struct {
		g.Meta `path:"/activities/assists" method:"GET" summary:"助力活动列表"`
		model.PageReq
	}
	AssistActivityItem struct {
		ActivityId    string `json:"activityId"`
		Name          string `json:"name"`
		RewardDesc    string `json:"rewardDesc" dc:"奖励说明"`
		RequiredCount int    `json:"requiredCount" dc:"所需人数"`
		EndTime       string `json:"endTime"`
	}
	AssistListRes struct {
		model.PageRes
		List []AssistActivityItem `json:"list"`
	}

	// 当前满减活动（可按商品过滤范围命中）
	FullReductionListReq struct {
		g.Meta `path:"/full-reductions" method:"GET" summary:"满减活动列表"`
		SpuId  string `json:"spuId" dc:"按商品过滤范围命中"`
		model.PageReq
	}
	FullReductionLadder struct {
		Threshold string `json:"threshold" dc:"满X元"`
		Discount  string `json:"discount" dc:"减Y元"`
	}
	FullReductionItem struct {
		ActivityId string                `json:"activityId"`
		Name       string                `json:"name"`
		Ladders    []FullReductionLadder `json:"ladders" dc:"档位"`
		ScopeDesc  string                `json:"scopeDesc" dc:"适用范围摘要(全场/分类/指定商品; 015 修复轮微扩)"`
		EndTime    string                `json:"endTime" dc:"015 修复轮微扩"`
	}
	FullReductionListRes struct {
		model.PageRes
		List []FullReductionItem `json:"list"`
	}
)
