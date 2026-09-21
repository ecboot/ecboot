// invite_impl.go 邀请激励记录（011-member-center; 接口契约见 distribution.go IDistributionLogic.InviteRecords）。
// 语义: 仅返回本人（inviter_id）发起的记录; 被邀请人昵称**脱敏**展示。
package user

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	"ecboot/internal/dao"
	"ecboot/internal/model"
)

// maskNickname 昵称脱敏（首字符 + 掩码; 空昵称给占位）。
func maskNickname(n string) string {
	n = strings.TrimSpace(n)
	if n == "" {
		return "用户****"
	}
	r := []rune(n)
	if len(r) <= 1 {
		return string(r) + "*"
	}
	return string(r[0]) + strings.Repeat("*", len(r)-1)
}

// inviteRewardDesc 奖励说明（reward_type: 1注册即发 2首单后发——V16 reward_trigger 口径）。
func inviteRewardDesc(rewardType int) string {
	switch rewardType {
	case 1:
		return "注册奖励"
	case 2:
		return "首单奖励"
	default:
		return "邀请奖励"
	}
}

// InviteRecords 邀请激励记录（FR-014）: 仅本人 + 时间倒序 + 分页; 被邀请人脱敏。
func InviteRecords(ctx context.Context, userId int64, page model.PageReq) (*model.PageResult[model.InviteRecordItem], error) {
	page = page.Normalized()
	icols, ucols := dao.InviteRecord.Columns(), dao.User.Columns()
	m := dao.InviteRecord.Ctx(ctx).As("ir").
		LeftJoin(dao.User.Table()+" u", "u.id=ir.new_user_id").
		Where("ir."+icols.InviterId, userId)
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计邀请记录失败")
	}
	recs, err := m.Fields("ir."+icols.NewUserId+", ir."+icols.RewardType+", ir."+icols.Status+
		", ir."+icols.CreatedAt+", u."+ucols.Nickname).
		OrderDesc("ir."+icols.Id).
		Page(page.Page, page.PageSize).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询邀请记录失败")
	}
	list := make([]model.InviteRecordItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, model.InviteRecordItem{
			NewUser:    maskNickname(r[ucols.Nickname].String()),
			RewardDesc: inviteRewardDesc(r[icols.RewardType].Int()),
			Status:     r[icols.Status].Int(),
			CreatedAt:  r[icols.CreatedAt].String(),
		})
	}
	return &model.PageResult[model.InviteRecordItem]{List: list, Total: int64(total)}, nil
}
