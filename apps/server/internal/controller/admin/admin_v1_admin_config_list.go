package admin

import (
	"context"

	"ecboot/api/admin/v1"
	"ecboot/internal/service/system"
)

// AdminConfigList 系统配置列表
func (c *ControllerV1) AdminConfigList(ctx context.Context, req *v1.AdminConfigListReq) (res *v1.AdminConfigListRes, err error) {
	items, err := system.ConfigList(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.AdminConfigListRes{List: make([]v1.AdminConfigItem, 0, len(items))}
	for _, it := range items {
		res.List = append(res.List, v1.AdminConfigItem{
			Code:        it.Code,
			Value:       it.Value,
			ValueType:   it.ValueType,
			Name:        it.Name,
			Description: it.Description,
			Status:      it.Status,
		})
	}
	return res, nil
}
