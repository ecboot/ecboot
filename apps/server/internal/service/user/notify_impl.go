// notify_impl.go 站内信与通知偏好（011-member-center; 接口契约见 message.go INotifyLogic）。
// 站内信: 列表返回**全量未读计数**（非当前页）; 标记已读做归属校验（他人→10006）。
// 通知偏好: 渠道维度（1小程序订阅 2短信）; **未建行=默认全开**; 设置幂等 UPSERT（uk user×channel）;
// 站内信不受偏好控制（表内固定渠道 3, 与本偏好表无关）。
// Enqueue/DispatchTask 为内部方法（交易事件入队与投递）, 本批实现基础形态供后续批次切换。
package user

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"ecboot/internal/dao"
	"ecboot/internal/errcode"
	"ecboot/internal/model"
	"ecboot/internal/model/do"
)

// mustJSONMap map → JSON 字符串（通知参数落库）。
func mustJSONMap(m map[string]any) string {
	if m == nil {
		return "{}"
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// notifyChannels 偏好渠道全集（站内信不在其中——不受偏好控制; 评审 Minor11: 定长数组防调用方污染）。
var notifyChannels = [2]int{1, 2}

// Messages 站内信列表（FR-009）: isRead=-1 全部 / 0 未读 / 1 已读; 未读计数为**全量**。
func Messages(ctx context.Context, userId int64, isRead int, page model.PageReq) (*model.MessageListResult, error) {
	page = page.Normalized()
	cols := dao.UserMessage.Columns()
	base := func() *gdb.Model {
		return dao.UserMessage.Ctx(ctx).Where(cols.UserId, userId)
	}
	m := base()
	if isRead == 0 || isRead == 1 {
		m = m.Where(cols.IsRead, isRead)
	}
	total, err := m.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计站内信失败")
	}
	// 未读计数: 全量（不受当前筛选影响）
	unread, err := base().Where(cols.IsRead, 0).Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计未读失败")
	}
	recs, err := m.OrderDesc(cols.Id).Page(page.Page, page.PageSize).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询站内信失败")
	}
	list := make([]model.MessageItem, 0, len(recs))
	for _, r := range recs {
		list = append(list, model.MessageItem{
			Id:        r[cols.Id].Int64(),
			Title:     r[cols.Title].String(),
			Content:   r[cols.Content].String(),
			BizType:   r[cols.BizType].Int(),
			BizNo:     r[cols.BizNo].String(),
			IsRead:    r[cols.IsRead].Int() == 1,
			CreatedAt: r[cols.CreatedAt].String(),
		})
	}
	return &model.MessageListResult{List: list, Total: int64(total), UnreadCount: int64(unread)}, nil
}

// MarkRead 标记单条已读（FR-009）: 归属校验（他人→10006）。
func MarkRead(ctx context.Context, userId, messageId int64) error {
	cols := dao.UserMessage.Columns()
	cnt, err := dao.UserMessage.Ctx(ctx).
		Where(cols.Id, messageId).Where(cols.UserId, userId).Count()
	if err != nil {
		return gerror.Wrap(err, "查询站内信失败")
	}
	if cnt == 0 {
		return errcode.New(errcode.CodeNotFound, "消息不存在")
	}
	if _, err = dao.UserMessage.Ctx(ctx).
		Where(cols.Id, messageId).Where(cols.UserId, userId).
		Data(do.UserMessage{IsRead: 1, ReadTime: gtime.Now()}).
		Fields(cols.IsRead, cols.ReadTime).
		Update(); err != nil {
		return gerror.Wrap(err, "标记已读失败")
	}
	return nil
}

// MarkAllRead 全部已读（FR-009）。
func MarkAllRead(ctx context.Context, userId int64) error {
	cols := dao.UserMessage.Columns()
	_, err := dao.UserMessage.Ctx(ctx).
		Where(cols.UserId, userId).
		Where(cols.IsRead, 0).
		Data(do.UserMessage{IsRead: 1, ReadTime: gtime.Now()}).
		Fields(cols.IsRead, cols.ReadTime).
		Update()
	if err != nil {
		return gerror.Wrap(err, "全部已读失败")
	}
	return nil
}

// Preferences 通知偏好（FR-010）: 渠道全集返回; **未建行=默认全开**。
func Preferences(ctx context.Context, userId int64) ([]model.NotifyPreference, error) {
	cols := dao.UserNotifyPreference.Columns()
	recs, err := dao.UserNotifyPreference.Ctx(ctx).
		Where(cols.UserId, userId).
		All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询偏好失败")
	}
	enabledByChannel := map[int]bool{}
	for _, r := range recs {
		enabledByChannel[r[cols.Channel].Int()] = r[cols.Enabled].Int() == 1
	}
	out := make([]model.NotifyPreference, 0, len(notifyChannels))
	for _, ch := range notifyChannels {
		enabled, ok := enabledByChannel[ch]
		if !ok {
			enabled = true // 未建行 = 默认全开
		}
		out = append(out, model.NotifyPreference{Channel: ch, Enabled: enabled})
	}
	return out, nil
}

// SetPreferences 设置偏好（FR-011）: 幂等 UPSERT（uk user×channel）; 仅接受的渠道入表。
func SetPreferences(ctx context.Context, userId int64, prefs []model.NotifyPreference) error {
	cols := dao.UserNotifyPreference.Columns()
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, p := range prefs {
			if p.Channel != 1 && p.Channel != 2 {
				return errcode.New(errcode.CodeInvalidParam, "偏好渠道须为1小程序订阅或2短信")
			}
			enabled := 0
			if p.Enabled {
				enabled = 1
			}
			// UPSERT: 存在则更新, 不存在则插入
			cnt, e := dao.UserNotifyPreference.Ctx(ctx).
				Where(cols.UserId, userId).Where(cols.Channel, p.Channel).Count()
			if e != nil {
				return gerror.Wrap(e, "查询偏好失败")
			}
			if cnt > 0 {
				_, e = dao.UserNotifyPreference.Ctx(ctx).
					Where(cols.UserId, userId).Where(cols.Channel, p.Channel).
					Data(do.UserNotifyPreference{Enabled: enabled}).
					Fields(cols.Enabled).
					Update()
			} else {
				_, e = dao.UserNotifyPreference.Ctx(ctx).
					Data(do.UserNotifyPreference{UserId: userId, Channel: p.Channel, Enabled: enabled}).
					Insert()
			}
			if e != nil {
				return gerror.Wrap(e, "设置偏好失败")
			}
		}
		return nil
	})
	return err
}

// Enqueue 创建通知任务（内部方法）: 按偏好拆渠道行——被关闭的渠道**不建行**; 站内信不受偏好控制。
func Enqueue(ctx context.Context, userId int64, bizType int, bizNo, templateCode string, params map[string]any) error {
	prefs, err := Preferences(ctx, userId)
	if err != nil {
		return err
	}
	enabled := map[int]bool{}
	for _, p := range prefs {
		enabled[p.Channel] = p.Enabled
	}
	channels := []int{3} // 站内信总是入队
	for _, ch := range notifyChannels {
		if enabled[ch] {
			channels = append(channels, ch)
		}
	}
	raw := mustJSONMap(params)
	for _, ch := range channels {
		if _, err = dao.NotifyTask.Ctx(ctx).Data(do.NotifyTask{
			UserId:       userId,
			Channel:      ch,
			BizType:      bizType,
			BizNo:        bizNo,
			TemplateCode: templateCode,
			Params:       raw,
			Status:       10, // 待发送
		}).Insert(); err != nil {
			return gerror.Wrap(err, "创建通知任务失败")
		}
	}
	return nil
}

// DispatchTask 投递任务（内部方法）。
// TODO(011 评审 I2): 当前为**基础形态**——站内信直接落 user_message 且 title=模板编码、content=参数 JSON,
// 尚未按接口契约"渲染模板(notify_template)→渠道发送→重试计数(retry_count/next_retry_time, 失败≤上限回待发送)"
// 实现; 站内信落库与状态更新亦未同事务。**接线前需补齐**（本批无调用点, 不阻塞）。
func DispatchTask(ctx context.Context, taskId int64) error {
	cols := dao.NotifyTask.Columns()
	rec, err := dao.NotifyTask.Ctx(ctx).Where(cols.Id, taskId).Where(cols.Status, 10).One()
	if err != nil {
		return gerror.Wrap(err, "查询通知任务失败")
	}
	if rec.IsEmpty() {
		return nil // 已投递或不存在（幂等）
	}
	if rec[cols.Channel].Int() == 3 {
		if _, err = dao.UserMessage.Ctx(ctx).Data(do.UserMessage{
			UserId:  rec[cols.UserId].Int64(),
			Title:   rec[cols.TemplateCode].String(),
			Content: rec[cols.Params].String(),
			BizType: rec[cols.BizType].Int(),
			BizNo:   rec[cols.BizNo].String(),
		}).Insert(); err != nil {
			return gerror.Wrap(err, "站内信落库失败")
		}
	}
	if _, err = dao.NotifyTask.Ctx(ctx).
		Where(cols.Id, taskId).
		Data(do.NotifyTask{Status: 20, SentTime: gtime.Now()}).
		Fields(cols.Status, cols.SentTime).
		Update(); err != nil {
		return gerror.Wrap(err, "标记发送失败")
	}
	return nil
}
