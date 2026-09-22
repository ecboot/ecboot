// assist.go 助力玩法动作（015-marketing-c 批次 09 新建）。
// 机制（research D3/D5）:
//   - 发起生成参与记录, 受活动"每人可发起次数"与时间窗限制。
//   - **一人一助力**: 依赖 `assist_helper` 的 `uk_record_helper(record_id, helper_user_id)` 唯一键兜底;
//     发起者本人不可助力自己的记录; 助力计数用条件更新推进。
//   - 达标（helper_count >= required_count）→ 置"已完成发奖" + finish_time, 并经**端口投递发奖意图**
//     （实际发放跨 user 域, 属后续批次; 端口未装配则告警降级——同批次 07 的 ICommissionReverse 形态）。
//   - 并发达标只允许投递**一次**发奖意图; 助力入口**须挂风控**（迁移 000029 明确要求, 走端口降级）。
package shop

import (
	"context"

	"ecboot/internal/model"
)

// IAssistLogic 助力动作。
type IAssistLogic interface {
	// Launch 发起助力（受"每人可发起次数"限制）。
	Launch(ctx context.Context, userId, activityId int64) (*model.AssistLaunchResult, error)
	// Progress 助力进度（公开; 含助力人列表, 昵称脱敏）。
	Progress(ctx context.Context, recordId int64) (*model.AssistProgressView, error)
	// Help 助力（一人一助力; 达标发奖）。
	Help(ctx context.Context, userId, recordId int64) (*model.AssistHelpResult, error)
}
