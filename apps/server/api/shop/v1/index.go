package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	IndexReq struct {
		g.Meta `path:"/index" tags:"Shop" method:"GET" summary:"首页接口"`
	}

	// IndexRes 首页聚合（015 批次 09 用户裁定: 轮播 + 楼层 + 五类活动入口 + 可领券）。
	// 原为脚手架占位（仅 Msg 字段）; 本次按裁定扩为聚合结构（**复用各列表端点的既有类型**, 不新造分块类型）。
	// **各分块为空时返回空数组**, 任一活动类目无人配置也不影响首页可用。
	IndexRes struct {
		Banners        []BannerItem          `json:"banners" dc:"在投轮播(位置=1)"`
		Floors         []FloorItem           `json:"floors" dc:"启用楼层(含商品装配)"`
		FlashSales     []FlashSaleItem       `json:"flashSales" dc:"秒杀入口(进行中+预告, 至多N条)"`
		GroupBuys      []GroupBuyItem        `json:"groupBuys" dc:"拼团入口"`
		Bargains       []BargainActivityItem `json:"bargains" dc:"砍价入口"`
		Assists        []AssistActivityItem  `json:"assists" dc:"助力入口"`
		FullReductions []FullReductionItem   `json:"fullReductions" dc:"满减入口"`
		Coupons        []IndexCoupon         `json:"coupons" dc:"可领券(模板级口径: 启用未删且未领完)"`
	}

	// IndexCoupon 首页可领券条目。
	IndexCoupon struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		Threshold string `json:"threshold"`
		Discount  string `json:"discount"`
		ValidDesc string `json:"validDesc"`
	}
)
