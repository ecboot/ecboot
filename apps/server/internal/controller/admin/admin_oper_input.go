package admin

import "ecboot/internal/model"

// bannerInputFromReq 轮播请求 → 服务入参（创建/修改共用; 全量覆盖语义）。
func bannerInputFromReq(position int, imageUrl, linkUrl string, sort int,
	startTime, endTime string, status int) model.OperBannerInput {
	return model.OperBannerInput{
		Position: position, ImageUrl: imageUrl, LinkUrl: linkUrl, Sort: sort,
		StartTime: startTime, EndTime: endTime, Status: status,
	}
}

// floorInputFromReq 楼层请求 → 服务入参（创建传 floorType; 修改该字段由 api 契约固定不传）。
func floorInputFromReq(floorType int, title string, config map[string]any, sort, status int) model.OperFloorInput {
	return model.OperFloorInput{FloorType: floorType, Title: title, Config: config, Sort: sort, Status: status}
}
