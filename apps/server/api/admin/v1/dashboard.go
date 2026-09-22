package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// 交易看板
	AdminDashboardTradeReq struct {
		g.Meta    `path:"/dashboard/trade" method:"GET" summary:"交易看板"`
		StartTime string `json:"startTime" dc:"起 RFC3339"`
		EndTime   string `json:"endTime" dc:"止"`
	}
	AdminDashboardTradeRes struct {
		OrderCount     int64  `json:"orderCount" dc:"订单数"`
		SalesAmount    string `json:"salesAmount" dc:"销售额(元)"`
		RefundAmount   string `json:"refundAmount" dc:"退款额(元)"`
		PendingDeliver int64  `json:"pendingDeliver" dc:"待发货单数"`
	}

	// 会员看板
	AdminDashboardMemberReq struct {
		g.Meta    `path:"/dashboard/member" method:"GET" summary:"会员看板"`
		StartTime string `json:"startTime" dc:"起"`
		EndTime   string `json:"endTime" dc:"止"`
	}
	AdminDashboardMemberRes struct {
		NewCount     int64 `json:"newCount" dc:"新增会员"`
		ActiveCount  int64 `json:"activeCount" dc:"活跃会员(口径:周期内 last_active_at 有活动; 无窗口取近 30 天)"`
		DormantCount int64 `json:"dormantCount" dc:"休眠会员(阈值读 system_config dormant.tier1.days, 缺省90天)"`
	}

	// 商品看板
	AdminDashboardProductReq struct {
		g.Meta `path:"/dashboard/product" method:"GET" summary:"商品看板"`
	}
	AdminDashboardProductRes struct {
		OnSaleCount   int64 `json:"onSaleCount" dc:"在售商品"`
		LowStockCount int64 `json:"lowStockCount" dc:"低库存预警数"`
		PendingReview int64 `json:"pendingReview" dc:"待审核评价数"`
	}
)
