// misc.go 商城周边域——表: logistics_company(V15) / store(V31) / risk_rule+risk_record(V21) / review 评价动作(V12)。
// 门店: B2C 线下载体（自提/核销/附近检索, 区划码+经纬度）; 风控: 规则与事件分离, 申诉流转。
package shop

import "context"

// ILogisticsLogic 物流公司字典。
type ILogisticsLogic interface {
	List(ctx context.Context, status int, page PageQuery) (*PageResult[LogisticsCompany], error)
	Create(ctx context.Context, in LogisticsCompanyInput) (int64, error)
	Update(ctx context.Context, id int64, in LogisticsCompanyInput) error
	Delete(ctx context.Context, id int64) error
}

type LogisticsCompany struct {
	Id           int64  `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	TrackingRule string `json:"trackingRule"`
	Status       int    `json:"status"`
}

type LogisticsCompanyInput struct {
	Code         string
	Name         string
	TrackingRule string
	Status       int
}

// IStoreLogic 门店。
type IStoreLogic interface {
	// PublicList 游客门店（区县筛选或经纬度附近+距离排序; 仅营业）。
	PublicList(ctx context.Context, q StoreQuery) (*PageResult[StoreItem], error)
	PublicDetail(ctx context.Context, storeId int64) (*StoreItem, error)
	AdminList(ctx context.Context, status int, keyword string, page PageQuery) (*PageResult[StoreItem], error)
	AdminCreate(ctx context.Context, in StoreInput) (int64, error)
	AdminUpdate(ctx context.Context, id int64, in StoreInput) error
	AdminDelete(ctx context.Context, id int64) error
}

type StoreQuery struct {
	DistrictCode string
	Longitude    float64
	Latitude     float64
	RadiusKm     int
	PageQuery
}

type StoreItem struct {
	Id            int64   `json:"id"`
	StoreNo       string  `json:"storeNo"`
	Name          string  `json:"name"`
	DistrictCode  string  `json:"districtCode"`
	DetailAddress string  `json:"detailAddress"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
	BusinessHours string  `json:"businessHours"`
	ContactPhone  string  `json:"contactPhone"`
	PickupEnabled bool    `json:"pickupEnabled"`
	Status        int     `json:"status"`
	DistanceM     int64   `json:"distanceM" dc:"附近检索返回"`
}

type StoreInput struct {
	Name          string
	ProvinceCode  string
	CityCode      string
	DistrictCode  string
	DetailAddress string
	Longitude     float64
	Latitude      float64
	BusinessHours string
	ContactPhone  string
	PickupEnabled bool
	Status        int
}

// IRiskLogic 风控。
type IRiskLogic interface {
	// Rules 规则列表。
	Rules(ctx context.Context, status int, page PageQuery) (*PageResult[RiskRuleItem], error)
	RuleCreate(ctx context.Context, in RiskRuleInput) (int64, error)
	RuleUpdate(ctx context.Context, id int64, in RiskRuleInput) error
	RuleDelete(ctx context.Context, id int64) error
	// Records 事件查询。
	Records(ctx context.Context, userId int64, appealStatus int, page PageQuery) (*PageResult[RiskRecordItem], error)
	// Hit 记录命中（拦截/标记; 下单/领券/帮砍等入口调用, 返回是否拦截）。
	Hit(ctx context.Context, userId int64, ruleId int64, objectType int, objectNo string) (blocked bool, err error)
	Appeal(ctx context.Context, recordId int64, pass bool, remark, operator string) error
}

type RiskRuleItem struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	RuleType      int    `json:"ruleType" dc:"1黑名单 2高频下单 3异常领券 4佣金套利 5休眠分级"`
	ConditionExpr string `json:"conditionExpr"`
	Action        int    `json:"action" dc:"1拦截 2标记"`
	Status        int    `json:"status"`
}

type RiskRuleInput struct {
	Name          string
	RuleType      int
	ConditionExpr string
	Action        int
	Status        int
}

type RiskRecordItem struct {
	Id           int64  `json:"id"`
	UserId       int64  `json:"userId"`
	RuleName     string `json:"ruleName"`
	ObjectType   int    `json:"objectType" dc:"1订单 2券 3提现 4售后"`
	ObjectNo     string `json:"objectNo"`
	Action       int    `json:"action"`
	AppealStatus int    `json:"appealStatus" dc:"0无 1申诉中 2通过 3驳回"`
	CreatedAt    string `json:"createdAt"`
}
