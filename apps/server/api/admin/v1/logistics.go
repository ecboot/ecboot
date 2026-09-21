package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"ecboot/internal/model"
)

type (
	// 物流公司字典列表
	AdminLogisticsListReq struct {
		g.Meta `path:"/logistics-companies" method:"GET" summary:"物流公司列表"`
		Status int `json:"status" dc:"状态筛选"`
		model.PageReq
	}
	AdminLogisticsItem struct {
		Id           string `json:"id"`
		Code         string `json:"code" dc:"编码(唯一)"`
		Name         string `json:"name"`
		TrackingRule string `json:"trackingRule" dc:"单号校验规则"`
		Status       int    `json:"status"`
	}
	AdminLogisticsListRes struct {
		model.PageRes
		List []AdminLogisticsItem `json:"list"`
	}

	// 权限: logistics:company:manage
	AdminLogisticsCreateReq struct {
		g.Meta       `path:"/logistics-companies" method:"POST" summary:"新增物流公司"`
		Code         string `json:"code" v:"required" dc:"编码"`
		Name         string `json:"name" v:"required" dc:"名称"`
		TrackingRule string `json:"trackingRule" dc:"单号规则"`
	}
	AdminLogisticsCreateRes struct {
		Id string `json:"id"`
	}

	// 权限: logistics:company:manage
	AdminLogisticsUpdateReq struct {
		g.Meta       `path:"/logistics-companies/{id}" method:"PUT" summary:"修改物流公司"`
		Id           string `json:"id" v:"required" dc:"ID"`
		Name         string `json:"name" dc:"名称"`
		TrackingRule string `json:"trackingRule" dc:"单号规则"`
		Status       int    `json:"status" dc:"1启用 0停用"`
	}
	AdminLogisticsUpdateRes struct {
		Success bool `json:"success"`
	}

	// 权限: logistics:company:manage
	AdminLogisticsDeleteReq struct {
		g.Meta `path:"/logistics-companies/{id}" method:"DELETE" summary:"删除物流公司(软删)"`
		Id     string `json:"id" v:"required" dc:"ID"`
	}
	AdminLogisticsDeleteRes struct {
		Success bool `json:"success"`
	}

	// 物流公司详情
	AdminLogisticsDetailReq struct {
		g.Meta `path:"/logistics-companies/{id}" method:"GET" summary:"物流公司详情"`
		Id     string `json:"id" v:"required" dc:"ID"`
	}
	AdminLogisticsDetailRes struct {
		Id           string `json:"id"`
		Code         string `json:"code"`
		Name         string `json:"name"`
		TrackingRule string `json:"trackingRule"`
		Status       int    `json:"status"`
	}
)
