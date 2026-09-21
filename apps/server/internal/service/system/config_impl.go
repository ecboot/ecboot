// config_impl.go 系统配置管理面（接口契约见 ops.go IConfigLogic）。
// 语义: 配置是覆盖层——停用/缺失回退代码默认, 不阻断业务（V30 迁移注释为准绳）;
// 本批只交付管理面, 业务读取组件推迟到首个消费批次（research D5）。
package system

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/model"
	"ecboot/internal/model/entity"
)

// ConfigList 配置全量列表（FR-019）。
func ConfigList(ctx context.Context) ([]model.ConfigItem, error) {
	recs, err := dao.SystemConfig.Ctx(ctx).
		OrderAsc(dao.SystemConfig.Columns().Id).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询系统配置失败")
	}
	var cfgs []entity.SystemConfig
	if err = recs.Structs(&cfgs); err != nil {
		return nil, gerror.Wrap(err, "解析系统配置失败")
	}
	list := make([]model.ConfigItem, 0, len(cfgs))
	for _, c := range cfgs {
		list = append(list, model.ConfigItem{
			Code:        c.Code,
			Value:       c.Value,
			ValueType:   c.ValueType,
			Name:        c.Name,
			Description: c.Description,
			Status:      c.Status,
		})
	}
	return list, nil
}

// ConfigUpdate 修改配置（FR-020/021）: 按类型校验值 + 启停切换。
func ConfigUpdate(ctx context.Context, code, value string, status int) error {
	rec, err := dao.SystemConfig.Ctx(ctx).
		Where(dao.SystemConfig.Columns().Code, code).
		One()
	if err != nil {
		return gerror.Wrap(err, "查询系统配置失败")
	}
	if rec.IsEmpty() {
		return errcode.New(errcode.CodeNotFound, "配置不存在")
	}
	var cfg entity.SystemConfig
	if err = rec.Struct(&cfg); err != nil {
		return gerror.Wrap(err, "解析系统配置失败")
	}
	if err = validateConfigValue(cfg.ValueType, value); err != nil {
		return err
	}
	data := map[string]any{"status": status}
	if value != "" {
		data["value"] = value
	}
	if _, err = dao.SystemConfig.Ctx(ctx).
		Where(dao.SystemConfig.Columns().Id, cfg.Id).
		Data(data).Update(); err != nil {
		return gerror.Wrap(err, "修改系统配置失败")
	}
	return nil
}

// validateConfigValue 值按 value_type 校验（FR-020: 1整数 2小数 3字符串 4布尔 5JSON）。
func validateConfigValue(valueType int, value string) error {
	switch valueType {
	case 1:
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return errcode.New(errcode.CodeConfigValueInvalid, "配置值须为整数")
		}
	case 2:
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return errcode.New(errcode.CodeConfigValueInvalid, "配置值须为小数")
		}
	case 4:
		if _, err := strconv.ParseBool(value); err != nil {
			return errcode.New(errcode.CodeConfigValueInvalid, "配置值须为布尔")
		}
	case 5:
		if !json.Valid([]byte(value)) {
			return errcode.New(errcode.CodeConfigValueInvalid, "配置值须为合法JSON")
		}
	}
	return nil // 3字符串: 任意值
}
