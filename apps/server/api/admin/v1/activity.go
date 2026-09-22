package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/model"
)

// 拼团/秒杀/砍价/助力四类活动的管理契约。
// 四类活动同构（活动 → 场次商品 → 状态），结构刻意对齐，降低实现与认知成本。

type (
	// ---------- 拼团 ----------
	AdminGroupBuyListReq struct {
		g.Meta `path:"/group-buys" method:"GET" summary:"拼团活动列表"`
		Status int `json:"status" dc:"状态筛选"`
		model.PageReq
	}
	AdminGroupBuyItem struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		SpuId     string `json:"spuId"`
		GroupSize int    `json:"groupSize" dc:"成团人数"`
		PerLimit  int    `json:"perLimit" dc:"限购"`
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
		Status    int    `json:"status"`
	}
	AdminGroupBuyListRes struct {
		model.PageRes
		List []AdminGroupBuyItem `json:"list"`
	}

	// 权限: promotion:groupbuy:manage
	AdminGroupBuyCreateReq struct {
		g.Meta    `path:"/group-buys" method:"POST" summary:"创建拼团活动"`
		Name      string `json:"name" v:"required" dc:"名称"`
		SpuId     string `json:"spuId" v:"required" dc:"SPU"`
		GroupSize int    `json:"groupSize" v:"required|min:2" dc:"成团人数"`
		PerLimit  int    `json:"perLimit" dc:"限购" d:"1"`
		StartTime string `json:"startTime" v:"required" dc:"开始"`
		EndTime   string `json:"endTime" v:"required" dc:"结束"`
	}
	AdminGroupBuyCreateRes struct {
		Id string `json:"id"`
	}

	// 权限: promotion:groupbuy:manage
	AdminGroupBuyUpdateReq struct {
		g.Meta    `path:"/group-buys/{id}" method:"PUT" summary:"修改拼团活动"`
		Id        string `json:"id" v:"required" dc:"活动ID"`
		Name      string `json:"name" dc:"名称"`
		GroupSize int    `json:"groupSize" dc:"成团人数"`
		PerLimit  int    `json:"perLimit" dc:"限购"`
		StartTime string `json:"startTime" dc:"开始"`
		EndTime   string `json:"endTime" dc:"结束"`
		Status    *int   `json:"status" dc:"启停(nil=不修改; 0停用 1启用; 016 评审 I6 三态化)"`
	}
	AdminGroupBuyUpdateRes struct {
		Success bool `json:"success"`
	}

	// 权限: promotion:groupbuy:manage
	AdminGroupBuyDeleteReq struct {
		g.Meta `path:"/group-buys/{id}" method:"DELETE" summary:"删除拼团活动(软删)"`
		Id     string `json:"id" v:"required" dc:"活动ID"`
	}
	AdminGroupBuyDeleteRes struct {
		Success bool `json:"success"`
	}

	// 拼团场次商品（SKU 级成团价; 全量替换）
	// 权限: promotion:groupbuy:manage
	AdminGroupBuyItemsReq struct {
		g.Meta `path:"/group-buys/{id}/items" method:"PUT" summary:"拼团场次商品设置"`
		Id     string             `json:"id" v:"required" dc:"活动ID"`
		Items  []AdminActivitySku `json:"items" dc:"场次商品列表"`
	}
	AdminActivitySku struct {
		SkuId         string `json:"skuId" v:"required" dc:"SKU ID"`
		GroupPrice    string `json:"groupPrice" dc:"拼团:成团价"`
		FlashPrice    string `json:"flashPrice" dc:"秒杀:秒杀价"`
		StockCount    int    `json:"stockCount" dc:"秒杀:活动限量"`
		PerLimit      int    `json:"perLimit" dc:"秒杀:每人限购"`
		OriginalPrice string `json:"originalPrice" dc:"砍价:起始价"`
		FloorPrice    string `json:"floorPrice" dc:"砍价:底价"`
		MaxCutCount   int    `json:"maxCutCount" dc:"砍价:最大刀数,0不限"`
	}
	AdminGroupBuyItemsRes struct {
		Success bool `json:"success"`
	}

	// ---------- 秒杀 ----------
	AdminFlashSaleListReq struct {
		g.Meta `path:"/flash-sales" method:"GET" summary:"秒杀活动列表"`
		Status int `json:"status" dc:"状态筛选"`
		model.PageReq
	}
	AdminFlashSaleItem struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
		Status    int    `json:"status"`
	}
	AdminFlashSaleListRes struct {
		model.PageRes
		List []AdminFlashSaleItem `json:"list"`
	}

	// 权限: promotion:flashsale:manage
	AdminFlashSaleCreateReq struct {
		g.Meta    `path:"/flash-sales" method:"POST" summary:"创建秒杀活动"`
		Name      string `json:"name" v:"required" dc:"名称"`
		StartTime string `json:"startTime" v:"required" dc:"开始"`
		EndTime   string `json:"endTime" v:"required" dc:"结束"`
	}
	AdminFlashSaleCreateRes struct {
		Id string `json:"id"`
	}

	// 权限: promotion:flashsale:manage
	AdminFlashSaleUpdateReq struct {
		g.Meta    `path:"/flash-sales/{id}" method:"PUT" summary:"修改秒杀活动"`
		Id        string `json:"id" v:"required" dc:"活动ID"`
		Name      string `json:"name" dc:"名称"`
		StartTime string `json:"startTime" dc:"开始"`
		EndTime   string `json:"endTime" dc:"结束"`
		Status    *int   `json:"status" dc:"启停(nil=不修改; 016 评审 I6 三态化)"`
	}
	AdminFlashSaleUpdateRes struct {
		Success bool `json:"success"`
	}

	// 权限: promotion:flashsale:manage
	AdminFlashSaleDeleteReq struct {
		g.Meta `path:"/flash-sales/{id}" method:"DELETE" summary:"删除秒杀活动(软删)"`
		Id     string `json:"id" v:"required" dc:"活动ID"`
	}
	AdminFlashSaleDeleteRes struct {
		Success bool `json:"success"`
	}

	// 权限: promotion:flashsale:manage
	AdminFlashSaleItemsReq struct {
		g.Meta `path:"/flash-sales/{id}/items" method:"PUT" summary:"秒杀场次商品设置"`
		Id     string             `json:"id" v:"required" dc:"活动ID"`
		Items  []AdminActivitySku `json:"items" dc:"场次商品(秒杀价/限量/限购)"`
	}
	AdminFlashSaleItemsRes struct {
		Success bool `json:"success"`
	}

	// ---------- 砍价 ----------
	AdminBargainListReq struct {
		g.Meta `path:"/bargains" method:"GET" summary:"砍价活动列表"`
		Status int `json:"status" dc:"状态筛选"`
		model.PageReq
	}
	AdminBargainItem struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		SpuId     string `json:"spuId"`
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
		Status    int    `json:"status"`
	}
	AdminBargainListRes struct {
		model.PageRes
		List []AdminBargainItem `json:"list"`
	}

	// 权限: promotion:bargain:manage
	AdminBargainCreateReq struct {
		g.Meta    `path:"/bargains" method:"POST" summary:"创建砍价活动"`
		Name      string `json:"name" v:"required" dc:"名称"`
		SpuId     string `json:"spuId" v:"required" dc:"SPU"`
		StartTime string `json:"startTime" v:"required" dc:"开始"`
		EndTime   string `json:"endTime" v:"required" dc:"结束"`
	}
	AdminBargainCreateRes struct {
		Id string `json:"id"`
	}

	// 权限: promotion:bargain:manage
	AdminBargainUpdateReq struct {
		g.Meta    `path:"/bargains/{id}" method:"PUT" summary:"修改砍价活动"`
		Id        string `json:"id" v:"required" dc:"活动ID"`
		Name      string `json:"name" dc:"名称"`
		StartTime string `json:"startTime" dc:"开始"`
		EndTime   string `json:"endTime" dc:"结束"`
		Status    *int   `json:"status" dc:"启停(nil=不修改; 016 评审 I6 三态化)"`
	}
	AdminBargainUpdateRes struct {
		Success bool `json:"success"`
	}

	// 权限: promotion:bargain:manage
	AdminBargainDeleteReq struct {
		g.Meta `path:"/bargains/{id}" method:"DELETE" summary:"删除砍价活动(软删)"`
		Id     string `json:"id" v:"required" dc:"活动ID"`
	}
	AdminBargainDeleteRes struct {
		Success bool `json:"success"`
	}

	// 权限: promotion:bargain:manage
	AdminBargainItemsReq struct {
		g.Meta `path:"/bargains/{id}/items" method:"PUT" summary:"砍价场次商品设置"`
		Id     string             `json:"id" v:"required" dc:"活动ID"`
		Items  []AdminActivitySku `json:"items" dc:"场次商品(起始价/底价/刀数)"`
	}
	AdminBargainItemsRes struct {
		Success bool `json:"success"`
	}

	// ---------- 助力 ----------
	AdminAssistListReq struct {
		g.Meta `path:"/assists" method:"GET" summary:"助力活动列表"`
		Status int `json:"status" dc:"状态筛选"`
		model.PageReq
	}
	AdminAssistItem struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		RewardType    int    `json:"rewardType" dc:"1券 2积分"`
		RewardDesc    string `json:"rewardDesc" dc:"奖励说明"`
		RequiredCount int    `json:"requiredCount" dc:"所需人数"`
		PerLimit      int    `json:"perLimit" dc:"每人可发起"`
		StartTime     string `json:"startTime"`
		EndTime       string `json:"endTime"`
		Status        int    `json:"status"`
	}
	AdminAssistListRes struct {
		model.PageRes
		List []AdminAssistItem `json:"list"`
	}

	// 权限: promotion:assist:manage
	AdminAssistCreateReq struct {
		g.Meta        `path:"/assists" method:"POST" summary:"创建助力活动"`
		Name          string `json:"name" v:"required" dc:"名称"`
		RewardType    int    `json:"rewardType" v:"required|in:1,2" dc:"奖励类型"`
		RewardRef     string `json:"rewardRef" dc:"奖励载体ID(券ID;积分为空)"`
		PointAmount   int    `json:"pointAmount" dc:"积分奖励数(rewardType=2)"`
		RequiredCount int    `json:"requiredCount" v:"required|min:1" dc:"所需人数"`
		PerLimit      int    `json:"perLimit" dc:"每人可发起" d:"1"`
		StartTime     string `json:"startTime" v:"required" dc:"开始"`
		EndTime       string `json:"endTime" v:"required" dc:"结束"`
	}
	AdminAssistCreateRes struct {
		Id string `json:"id"`
	}

	// 权限: promotion:assist:manage
	AdminAssistUpdateReq struct {
		g.Meta        `path:"/assists/{id}" method:"PUT" summary:"修改助力活动"`
		Id            string `json:"id" v:"required" dc:"活动ID"`
		Name          string `json:"name" dc:"名称"`
		RequiredCount int    `json:"requiredCount" dc:"所需人数"`
		PerLimit      int    `json:"perLimit" dc:"每人可发起"`
		StartTime     string `json:"startTime" dc:"开始"`
		EndTime       string `json:"endTime" dc:"结束"`
		Status        *int   `json:"status" dc:"启停(nil=不修改; 016 评审 I6 三态化)"`
	}
	AdminAssistUpdateRes struct {
		Success bool `json:"success"`
	}

	// 权限: promotion:assist:manage
	AdminAssistDeleteReq struct {
		g.Meta `path:"/assists/{id}" method:"DELETE" summary:"删除助力活动(软删)"`
		Id     string `json:"id" v:"required" dc:"活动ID"`
	}
	AdminAssistDeleteRes struct {
		Success bool `json:"success"`
	}
)
