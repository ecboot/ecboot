// dto_user.go user 域 DTO（自 internal/service/user 迁入 2026-09-21, 契约出入参统一归属 model）。
package model

type AddressInput struct {
	ReceiverName  string
	ReceiverPhone string
	Province      string
	City          string
	District      string
	DetailAddress string
	ProvinceCode  string
	CityCode      string
	DistrictCode  string
	IsDefault     bool
}

type AddressItem struct {
	Id            int64  `json:"id"`
	ReceiverName  string `json:"receiverName"`
	ReceiverPhone string `json:"receiverPhone" dc:"脱敏"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	DetailAddress string `json:"detailAddress"`
	ProvinceCode  string `json:"provinceCode"`
	CityCode      string `json:"cityCode"`
	DistrictCode  string `json:"districtCode"`
	IsDefault     bool   `json:"isDefault"`
}

type LoginOutcome struct {
	Token        string
	RefreshToken string
	UserId       int64
	IsNew        bool
}

type AvailableCoupon struct {
	CouponId   int64  `json:"couponId"`
	Name       string `json:"name"`
	Type       int    `json:"type" dc:"1满减 2无门槛"`
	Threshold  string `json:"threshold"`
	Discount   string `json:"discount"`
	ValidDesc  string `json:"validDesc"`
	CanReceive bool   `json:"canReceive"`
}

type MyCouponItem struct {
	UserCouponId int64  `json:"userCouponId"`
	Name         string `json:"name"`
	Threshold    string `json:"threshold"`
	Discount     string `json:"discount"`
	ExpireTime   string `json:"expireTime"`
	Status       int    `json:"status"`
}

type UsableCouponItem struct {
	UserCouponId int64  `json:"userCouponId"`
	Name         string `json:"name"`
	Discount     string `json:"discount"`
}

type DistStatus struct {
	Status int `json:"status" dc:"0未申请 1待审核 2通过 3冻结"`
	Level  int `json:"level" dc:"推广员等级(V29 预留)"`
}

type DistRelationResult struct {
	Inviter  map[string]string `json:"inviter" dc:"直接上级(脱敏)"`
	Invitees []DistInvitee     `json:"invitees"`
	Total    int64             `json:"total"`
}

type DistInvitee struct {
	UserId   int64  `json:"userId"`
	Nickname string `json:"nickname" dc:"脱敏"`
	BindTime string `json:"bindTime"`
}

type DistAccount struct {
	Balance string `json:"balance" dc:"可负"`
	Frozen  string `json:"frozen"`
}

type AccountLogItem struct {
	BizType      int    `json:"bizType" dc:"1入账 2提现冻结 3完成 4回退 5冲销 6消费冻结 7消费完成 8消费退回"`
	Amount       string `json:"amount" dc:"有符号"`
	BalanceAfter string `json:"balanceAfter"`
	BizNo        string `json:"bizNo"`
	CreatedAt    string `json:"createdAt"`
}

type WithdrawItem struct {
	WithdrawNo string `json:"withdrawNo"`
	Amount     string `json:"amount"`
	Status     int    `json:"status" dc:"10待审 20过审 30打款中 40成功 50拒绝 60失败回退"`
	CreatedAt  string `json:"createdAt"`
}

type FavoriteItem struct {
	SpuId    int64  `json:"spuId"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Price    string `json:"price" dc:"现价(元)"`
	Sellable bool   `json:"sellable"`
	Invalid  bool   `json:"invalid" dc:"已下架标注"`
}

type FootprintItem struct {
	SpuId      int64  `json:"spuId"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	Price      string `json:"price"`
	LastViewAt string `json:"lastViewAt"`
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

type PointAccountView struct {
	Balance      int    `json:"balance" dc:"可负"`
	LastEarnedAt string `json:"lastEarnedAt" dc:"滚动过期口径"`
}

type PointLogItem struct {
	BizType      int    `json:"bizType"`
	Points       int    `json:"points"`
	BalanceAfter int    `json:"balanceAfter"`
	OrderNo      string `json:"orderNo"`
	CreatedAt    string `json:"createdAt"`
}

type ProfileDetail struct {
	UserId      string `json:"userId"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Gender      int    `json:"gender"`
	PhoneMasked string `json:"phone" dc:"138****1234"`
	Level       int64  `json:"level" dc:"等级ID"`
	LevelName   string `json:"levelName"`
	GrowthValue int    `json:"growthValue"`
}

type ProfileUpdateInput struct {
	Nickname string
	Avatar   string
	Gender   int
}

type UserLoginLogItem struct {
	Channel   int    `json:"channel"`
	Status    int    `json:"status"`
	Ip        string `json:"ip" dc:"脱敏"`
	UserAgent string `json:"userAgent"`
	CreatedAt string `json:"createdAt"`
}

type WxIdentity struct {
	Openid  string
	Unionid string
	Phone   string // getPhoneNumber 一键取号结果（mock 模式由命令显式传入）
}

// InviteRecordItem 邀请激励记录（011-member-center; 对齐 api 契约）。
type InviteRecordItem struct {
	NewUser    string `json:"newUser" dc:"新用户(脱敏昵称)"`
	RewardDesc string `json:"rewardDesc" dc:"奖励说明"`
	Status     int    `json:"status" dc:"1已发放"`
	CreatedAt  string `json:"createdAt"`
}
