package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	MessageItem struct {
		Id        string `json:"id"`
		Title     string `json:"title"`
		Content   string `json:"content" dc:"内容摘要"`
		BizType   int    `json:"bizType" dc:"1订单 2营销 3售后"`
		BizNo     string `json:"bizNo" dc:"关联单号"`
		IsRead    bool   `json:"isRead"`
		CreatedAt string `json:"createdAt"`
	}

	MessageListReq struct {
		g.Meta  `path:"/messages" method:"GET" summary:"站内信列表"`
		IsRead  int `json:"isRead" dc:"已读筛选:0未读 1已读(空=全部)"`
		PageReq
	}
	MessageListRes struct {
		PageRes
		List        []MessageItem `json:"list"`
		UnreadCount int           `json:"unreadCount" dc:"未读数"`
	}

	MessageReadReq struct {
		g.Meta `path:"/messages/{id}/read" method:"PUT" summary:"标记已读"`
		Id     string `json:"id" v:"required" dc:"消息ID"`
	}
	MessageReadRes struct {
		Success bool `json:"success"`
	}

	MessageReadAllReq struct {
		g.Meta `path:"/messages/read-all" method:"PUT" summary:"全部已读"`
	}
	MessageReadAllRes struct {
		Success bool `json:"success"`
	}

	NotifyPreference struct {
		Channel int  `json:"channel" dc:"1小程序订阅 2短信"`
		Enabled bool `json:"enabled" dc:"是否接收"`
	}

	NotifyPreferenceGetReq struct {
		g.Meta `path:"/notify-preferences" method:"GET" summary:"通知偏好"`
	}
	NotifyPreferenceGetRes struct {
		List []NotifyPreference `json:"list"`
	}

	NotifyPreferenceSetReq struct {
		g.Meta `path:"/notify-preferences" method:"PUT" summary:"设置通知偏好"`
		List   []NotifyPreference `json:"list" v:"required" dc:"偏好集合"`
	}
	NotifyPreferenceSetRes struct {
		Success bool `json:"success"`
	}
)
