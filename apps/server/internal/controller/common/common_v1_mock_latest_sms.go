package common

import (
	"context"

	"ecboot/api/common/v1"
	"ecboot/internal/service/system"
)

// MockLatestSms 联调取码（fail-closed: 仅 ECBOOT_MOCK=true 开放, 生产拒绝, research D4）
func (c *ControllerV1) MockLatestSms(ctx context.Context, req *v1.MockLatestSmsReq) (res *v1.MockLatestSmsRes, err error) {
	code, err := system.MockLatestSms(ctx, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return &v1.MockLatestSmsRes{SmsCode: code}, nil
}
