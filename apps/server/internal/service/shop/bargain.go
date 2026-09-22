// bargain.go 砍价玩法动作（015-marketing-c 批次 09 新建）。
// 机制（research D4/D5）:
//   - 发起即砍第一刀; 每刀金额 = (起始价 − 底价) ÷ 最大刀数 **向下取整到分**;
//     当本刀会砍穿底价时**恰好补齐到**底价（当前价永不低于底价, 多刀之和恒等于"起始价 − 底价"）。
//   - **一人一刀**: 依赖 `bargain_helper` 的 `uk_record_helper(record_id, helper_user_id)` 唯一键兜底,
//     冲突(1062)转业务码; 计数与当前价用**条件更新 + 判 RowsAffected** 推进。
//   - 发起者本人不可帮砍自己的单; 超时单不可再砍; 帮砍入口**须挂风控**（迁移 000029 明确要求, 走端口降级）。
package shop

import (
	"context"

	"ecboot/internal/model"
)

// IBargainLogic 砍价动作。
type IBargainLogic interface {
	// Launch 发起砍价（发起即首刀; 受活动"每人可发起次数"限制）。
	Launch(ctx context.Context, userId, bargainItemId int64) (*model.BargainLaunchResult, error)
	// Progress 砍价进度（公开; 含帮砍列表, 昵称脱敏）。
	Progress(ctx context.Context, recordId int64) (*model.BargainProgressView, error)
	// Cut 帮砍一刀（一人一刀; 本人不可; 不砍穿）。
	Cut(ctx context.Context, userId, recordId int64) (*model.BargainCutResult, error)
}
