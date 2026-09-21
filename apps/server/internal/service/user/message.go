// message.go 消息中心——表: user_message（站内必达兜底）+ notify_task（异步投递）+ notify_template（V26）。
// 规则(research D4): 一渠道一行; 发送失败重试有上限(终态40); 渠道关闭→跳过(50);
// 站内信必达; 模板渲染参数契约 notify_template.params_schema。
package user

import "context"

// INotifyLogic 消息通知。
type INotifyLogic interface {
	// Messages 站内信列表（含未读计数）。
	Messages(ctx context.Context, userId int64, isRead int, page PageQuery) (*MessageListResult, error)
	// MarkRead 标记已读。
	MarkRead(ctx context.Context, userId, messageId int64) error
	// MarkAllRead 全部已读。
	MarkAllRead(ctx context.Context, userId int64) error
	// Preferences 订阅偏好查询。
	Preferences(ctx context.Context, userId int64) ([]NotifyPreference, error)
	// SetPreferences 设置偏好（跳过语义: notify_task 不建行, 站内信除外）。
	SetPreferences(ctx context.Context, userId int64, prefs []NotifyPreference) error
	// Enqueue 创建通知任务（订单状态迁移等业务事件调用; 按偏好拆渠道行）。
	Enqueue(ctx context.Context, userId int64, bizType int, bizNo, templateCode string, params map[string]any) error
	// DispatchTask 投递执行（定时/即时: 渲染模板→渠道发送→重试计数; 失败≤上限回待发送）。
	DispatchTask(ctx context.Context, taskId int64) error
}

type MessageListResult struct {
	List        []MessageItem `json:"list"`
	Total       int64         `json:"total"`
	UnreadCount int64         `json:"unreadCount"`
}

type MessageItem struct {
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	BizType   int    `json:"bizType"`
	BizNo     string `json:"bizNo"`
	IsRead    bool   `json:"isRead"`
	CreatedAt string `json:"createdAt"`
}

type NotifyPreference struct {
	Channel int  `json:"channel" dc:"1小程序订阅 2短信"`
	Enabled bool `json:"enabled"`
}
