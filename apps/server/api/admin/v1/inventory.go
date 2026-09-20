package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 库存列表（可售/锁定）
	AdminInventoryListReq struct {
		g.Meta  `path:"/inventories" method:"GET" summary:"库存列表"`
		SkuId   string `json:"skuId" dc:"SKU筛选"`
		Keyword string `json:"keyword" dc:"商品名/编码"`
		PageReq
	}
	AdminInventoryItem struct {
		SkuId      string `json:"skuId"`
		SkuNo      string `json:"skuNo"`
		SkuName    string `json:"skuName"`
		Total      int    `json:"total" dc:"总库存"`
		Locked     int    `json:"locked" dc:"锁定"`
		Available  int    `json:"available" dc:"可售=total-locked"`
		WarnCount  int    `json:"warnCount" dc:"预警阈值"`
	}
	AdminInventoryListRes struct {
		PageRes
		List []AdminInventoryItem `json:"list"`
	}

	// 库存调整（留痕 inventory_log）
	AdminInventoryAdjustReq struct {
		g.Meta `path:"/inventories/{skuId}/adjust" method:"POST" summary:"库存调整"`
		SkuId  string `json:"skuId" v:"required" dc:"SKU ID"`
		Delta  int    `json:"delta" v:"required" dc:"调整量(±)"`
		Remark string `json:"remark" v:"required" dc:"调整原因"`
	}
	AdminInventoryAdjustRes struct {
		TotalAfter int `json:"totalAfter" dc:"调整后总库存"`
	}

	// 预警列表（可售≤阈值）
	AdminInventoryWarnReq struct {
		g.Meta `path:"/inventories/warnings" method:"GET" summary:"库存预警列表"`
		PageReq
	}
	AdminInventoryWarnRes struct {
		PageRes
		List []AdminInventoryItem `json:"list"`
	}
)
