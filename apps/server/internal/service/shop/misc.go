// misc.go 商城周边域——表: logistics_company(V15) / store(V31) / risk_rule+risk_record(V21) / review 评价动作(V12)。
// 门店: B2C 线下载体（自提/核销/附近检索, 区划码+经纬度）; 风控: 规则与事件分离, 申诉流转。
package shop

import (
	"context"

	"ecboot/internal/model"
)

// ILogisticsLogic 物流公司字典。
type ILogisticsLogic interface {
	List(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.LogisticsCompany], error)
	Create(ctx context.Context, in model.LogisticsCompanyInput) (int64, error)
	Update(ctx context.Context, id int64, in model.LogisticsCompanyInput) error
	Delete(ctx context.Context, id int64) error
	// Detail 详情（009 接口微扩: api 有管理详情端点而接口缺定义, 同 D6/D7 模式）。
	Detail(ctx context.Context, id int64) (*model.LogisticsCompany, error)
}

// IOperationLogic 运营装修（banner/楼层; 009-logistics-ops research D1 新建）。
// 落 shop 域: 装修是商城内容（C 端首页消费）, 管理端 admin 渠道消费同域（渠道与域解耦）。
type IOperationLogic interface {
	BannerList(ctx context.Context, position int, page model.PageReq) (*model.PageResult[model.OperBannerItem], error)
	BannerCreate(ctx context.Context, in model.OperBannerInput) (int64, error)
	BannerUpdate(ctx context.Context, id int64, in model.OperBannerInput) error
	BannerDelete(ctx context.Context, id int64) error
	FloorList(ctx context.Context, page model.PageReq) (*model.PageResult[model.OperFloorItem], error)
	FloorCreate(ctx context.Context, in model.OperFloorInput) (int64, error)
	FloorUpdate(ctx context.Context, id int64, in model.OperFloorInput) error
	FloorDelete(ctx context.Context, id int64) error
	// PublicBanners 在投轮播/弹窗（启用 + 投放时段内, 按位置; research D4）。
	PublicBanners(ctx context.Context, position int) ([]model.PublicBannerItem, error)
	// PublicFloors C 端楼层（启用按序; 商品楼层装配摘要, 失效剔除; research D2/D3）。
	PublicFloors(ctx context.Context) ([]model.PublicFloorItem, error)
}

// IStoreLogic 门店。
type IStoreLogic interface {
	// PublicList 游客门店（区县筛选或经纬度附近+距离排序; 仅营业）。
	PublicList(ctx context.Context, q model.StoreQuery) (*model.PageResult[model.StoreItem], error)
	PublicDetail(ctx context.Context, storeId int64) (*model.StoreItem, error)
	AdminList(ctx context.Context, status int, keyword string, page model.PageReq) (*model.PageResult[model.StoreItem], error)
	AdminCreate(ctx context.Context, in model.StoreInput) (int64, error)
	AdminUpdate(ctx context.Context, id int64, in model.StoreInput) error
	AdminDelete(ctx context.Context, id int64) error
	// AdminDetail 门店详情（008-store 接口微扩: api 有管理详情端点而接口缺失, 同 D6 模式）。
	AdminDetail(ctx context.Context, storeId int64) (*model.StoreItem, error)
}

// IRiskLogic 风控。
type IRiskLogic interface {
	// Rules 规则列表。
	Rules(ctx context.Context, status int, page model.PageReq) (*model.PageResult[model.RiskRuleItem], error)
	RuleCreate(ctx context.Context, in model.RiskRuleInput) (int64, error)
	RuleUpdate(ctx context.Context, id int64, in model.RiskRuleInput) error
	RuleDelete(ctx context.Context, id int64) error
	// Records 事件查询。
	Records(ctx context.Context, userId int64, appealStatus int, page model.PageReq) (*model.PageResult[model.RiskRecordItem], error)
	// Hit 记录命中（拦截/标记; 下单/领券/帮砍等入口调用, 返回是否拦截）。
	Hit(ctx context.Context, userId int64, ruleId int64, objectType int, objectNo string) (blocked bool, err error)
	Appeal(ctx context.Context, recordId int64, pass bool, remark, operator string) error
}
