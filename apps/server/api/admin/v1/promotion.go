package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	// ---------- 优惠券模板 ----------
	AdminCouponListReq struct {
		g.Meta `path:"/coupons" method:"GET" summary:"券模板列表"`
		Status int `json:"status" dc:"状态筛选"`
		PageReq
	}
	AdminCouponItem struct {
		Id             string `json:"id"`
		Name           string `json:"name"`
		Type           int    `json:"type" dc:"1满减 2无门槛"`
		Threshold      string `json:"threshold"`
		Discount       string `json:"discount"`
		TotalCount     int    `json:"totalCount" dc:"发放总量,0不限"`
		ReceivedCount  int    `json:"receivedCount" dc:"已领"`
		PerLimit       int    `json:"perLimit" dc:"每人限领"`
		ValidType      int    `json:"validType" dc:"1固定区间 2领取后N天"`
		ValidDesc      string `json:"validDesc" dc:"有效期描述"`
		Status         int    `json:"status"`
	}
	AdminCouponListRes struct {
		PageRes
		List []AdminCouponItem `json:"list"`
	}

	AdminCouponCreateReq struct {
		g.Meta         `path:"/coupons" method:"POST" summary:"创建券模板"`
		Name           string `json:"name" v:"required" dc:"名称"`
		Type           int    `json:"type" v:"required|in:1,2" dc:"类型"`
		Threshold      string `json:"threshold" dc:"门槛"`
		Discount       string `json:"discount" v:"required" dc:"抵扣"`
		TotalCount     int    `json:"totalCount" dc:"总量,0不限"`
		PerLimit       int    `json:"perLimit" dc:"每人限领"`
		ValidType      int    `json:"validType" v:"required|in:1,2" dc:"有效期方式"`
		ValidStartAt   string `json:"validStartAt" dc:"固定区间开始"`
		ValidEndAt     string `json:"validEndAt" dc:"固定区间结束"`
		ValidDays      int    `json:"validDays" dc:"领取后N天"`
	}
	AdminCouponCreateRes struct {
		Id string `json:"id"`
	}

	AdminCouponUpdateReq struct {
		g.Meta     `path:"/coupons/{id}" method:"PUT" summary:"修改券模板"`
		Id         string `json:"id" v:"required" dc:"券ID"`
		Name       string `json:"name" dc:"名称"`
		Threshold  string `json:"threshold" dc:"门槛"`
		Discount   string `json:"discount" dc:"抵扣"`
		TotalCount int    `json:"totalCount" dc:"总量"`
		PerLimit   int    `json:"perLimit" dc:"限领"`
		Status     int    `json:"status" dc:"1启用 0停发"`
	}
	AdminCouponUpdateRes struct {
		Success bool `json:"success"`
	}

	AdminCouponDeleteReq struct {
		g.Meta `path:"/coupons/{id}" method:"DELETE" summary:"删除券模板(软删)"`
		Id     string `json:"id" v:"required" dc:"券ID"`
	}
	AdminCouponDeleteRes struct {
		Success bool `json:"success"`
	}

	// 券发放/使用记录
	AdminCouponRecordListReq struct {
		g.Meta `path:"/coupons/{id}/records" method:"GET" summary:"券记录"`
		Id     string `json:"id" v:"required" dc:"券ID"`
		PageReq
	}
	AdminCouponRecordItem struct {
		UserCouponId string `json:"userCouponId"`
		UserId       string `json:"userId"`
		Status       int    `json:"status" dc:"1未用 2已用 3过期 4退回"`
		OrderNo      string `json:"orderNo" dc:"核销单号"`
		CreatedAt    string `json:"createdAt" dc:"领取时间"`
	}
	AdminCouponRecordListRes struct {
		PageRes
		List []AdminCouponRecordItem `json:"list"`
	}

	// ---------- 满减活动（档位+范围嵌套提交） ----------
	AdminFullReductionLadder struct {
		Threshold string `json:"threshold" v:"required" dc:"满X元"`
		Discount  string `json:"discount" v:"required" dc:"减Y元"`
	}
	AdminFullReductionScope struct {
		ScopeType int    `json:"scopeType" v:"required|in:1,2,3" dc:"1全场 2分类 3商品"`
		TargetId  string `json:"targetId" dc:"分类/商品ID(全场为空)"`
	}
	AdminFullReductionListReq struct {
		g.Meta `path:"/full-reductions" method:"GET" summary:"满减活动列表"`
		Status int `json:"status" dc:"状态筛选"`
		PageReq
	}
	AdminFullReductionItem struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
		Status    int    `json:"status"`
	}
	AdminFullReductionListRes struct {
		PageRes
		List []AdminFullReductionItem `json:"list"`
	}

	AdminFullReductionCreateReq struct {
		g.Meta  `path:"/full-reductions" method:"POST" summary:"创建满减活动"`
		Name    string                      `json:"name" v:"required" dc:"名称"`
		StartTime string                    `json:"startTime" v:"required" dc:"开始"`
		EndTime string                      `json:"endTime" v:"required" dc:"结束"`
		Ladders []AdminFullReductionLadder  `json:"ladders" v:"required" dc:"档位"`
		Scopes  []AdminFullReductionScope   `json:"scopes" dc:"范围(空=全场)"`
	}
	AdminFullReductionCreateRes struct {
		Id string `json:"id"`
	}

	AdminFullReductionDetailReq struct {
		g.Meta `path:"/full-reductions/{id}" method:"GET" summary:"满减详情"`
		Id     string `json:"id" v:"required" dc:"活动ID"`
	}
	AdminFullReductionDetailRes struct {
		Id        string                     `json:"id"`
		Name      string                     `json:"name"`
		StartTime string                     `json:"startTime"`
		EndTime   string                     `json:"endTime"`
		Status    int                        `json:"status"`
		Ladders   []AdminFullReductionLadder `json:"ladders"`
		Scopes    []AdminFullReductionScope  `json:"scopes"`
	}

	AdminFullReductionUpdateReq struct {
		g.Meta  `path:"/full-reductions/{id}" method:"PUT" summary:"修改满减活动"`
		Id      string                      `json:"id" v:"required" dc:"活动ID"`
		Name    string                      `json:"name" dc:"名称"`
		StartTime string                    `json:"startTime" dc:"开始"`
		EndTime string                      `json:"endTime" dc:"结束"`
		Status  int                         `json:"status" dc:"状态"`
		Ladders []AdminFullReductionLadder  `json:"ladders" dc:"档位(全量替换)"`
		Scopes  []AdminFullReductionScope   `json:"scopes" dc:"范围(全量替换)"`
	}
	AdminFullReductionUpdateRes struct {
		Success bool `json:"success"`
	}

	AdminFullReductionDeleteReq struct {
		g.Meta `path:"/full-reductions/{id}" method:"DELETE" summary:"删除满减活动(软删)"`
		Id     string `json:"id" v:"required" dc:"活动ID"`
	}
	AdminFullReductionDeleteRes struct {
		Success bool `json:"success"`
	}
)
