# 端点到实现映射：营销 C 端（015-marketing-c）

**本批 12 个端点**（全为 `CodeNotImplemented` 桩；验收口径：`check-stub` 的 **shop 12 → 0**）。

## 一、新建接口（C 端营销契约原本缺失，D7）

```go
// IMarketingLogic 营销公开面（列表与首页聚合）。
type IMarketingLogic interface {
    PublicGroupBuys(ctx context.Context, page model.PageReq) (*model.PageResult[model.PublicGroupBuyItem], error)
    PublicFlashSales(ctx context.Context, page model.PageReq) (*model.PageResult[model.PublicFlashSaleItem], error)
    PublicBargains(ctx context.Context, page model.PageReq) (*model.PageResult[model.PublicBargainItem], error)
    PublicAssists(ctx context.Context, page model.PageReq) (*model.PageResult[model.PublicAssistItem], error)
    PublicFullReductions(ctx context.Context, page model.PageReq) (*model.PageResult[model.PublicFullReductionItem], error)
    Index(ctx context.Context) (*model.IndexAggregate, error)
}

// IBargainLogic 砍价动作。
type IBargainLogic interface {
    Launch(ctx context.Context, userId, bargainItemId int64) (*model.BargainLaunchResult, error)
    Progress(ctx context.Context, recordId int64) (*model.BargainProgressView, error)
    Cut(ctx context.Context, userId, recordId int64) (*model.BargainCutResult, error)
}

// IAssistLogic 助力动作。
type IAssistLogic interface {
    Launch(ctx context.Context, userId, activityId int64) (*model.AssistLaunchResult, error)
    Progress(ctx context.Context, recordId int64) (*model.AssistProgressView, error)
    Help(ctx context.Context, userId, recordId int64) (*model.AssistHelpResult, error)
}
```

## 二、端点映射

| # | 端点 | controller 桩文件 | service 方法 | 鉴权 | 关键行为 |
|---|---|---|---|---|---|
| 1 | `GET /shop/index` | `shop_v1_index.go` | `NewMarketingLogic().Index` | 公开（`/shop/index` 需入白名单? 见注） | 轮播+楼层+五入口+券（D2） |
| 2 | `GET /shop/activities/group-buys` | `shop_v1_group_buy_list.go` | `.PublicGroupBuys` | 公开（`/shop/activities/` 已在白名单） | 只出进行中，按结束时间升序 |
| 3 | `GET /shop/activities/flash-sales` | `shop_v1_flash_sale_list.go` | `.PublicFlashSales` | 公开 | 进行中 + 预告 |
| 4 | `GET /shop/activities/bargains` | `shop_v1_bargain_activity_list.go` | `.PublicBargains` | 公开 | 进行中 |
| 5 | `GET /shop/activities/assists` | `shop_v1_assist_list.go` | `.PublicAssists` | 公开 | 进行中 |
| 6 | `GET /shop/full-reductions` | `shop_v1_full_reduction_list.go` | `.PublicFullReductions` | 公开（已在白名单） | 进行中 |
| 7 | `POST /shop/bargains` | `shop_v1_bargain_launch.go` | `NewBargainLogic().Launch` | 会员 | 建单 + 首刀 |
| 8 | `GET /shop/bargains/{recordId}` | `shop_v1_bargain_progress.go` | `.Progress` | 公开（`/shop/bargains/` 已在白名单） | 进度 + 帮砍列表（脱敏） |
| 9 | `POST /shop/bargains/{recordId}/cut` | `shop_v1_bargain_cut.go` | `.Cut` | 会员 | 一人一刀 + 不砍穿 + 风控端口 |
| 10 | `POST /shop/assists` | `shop_v1_assist_launch.go` | `NewAssistLogic().Launch` | 会员 | 发起次数上限 |
| 11 | `GET /shop/assists/{recordId}` | `shop_v1_assist_progress.go` | `.Progress` | 公开（`/shop/assists/` 已在白名单） | 进度 + 助力人列表（脱敏） |
| 12 | `POST /shop/assists/{recordId}/helpers` | `shop_v1_assist_help.go` | `.Help` | 会员 | 一人一助力 + 达标发奖意图 + 风控端口 |

**注（白名单）**: `/shop/activities/`、`/shop/full-reductions`、`/shop/bargains/`、`/shop/assists/` 均**已在** `publicPrefixes`（批次 04/09 之前就存在）→ 公开列表与进度端点天然可达；**会员动作端点（7/9/10/12）必须落在白名单之外**才能拿到 userId。`/shop/index` **不在**白名单 → 需确认是否加入（首页应对游客可见）——**实现时按 HTTP 探活决定，并补端点级可达性测试**（014 评审 C1 的教训）。

## 三、跨批改动（依协议记账）

| 文件 | 改动 | 理由 |
|---|---|---|
| `internal/service/shop/order_impl.go` | `collectLines` 取秒杀价快照；步骤 3 增加活动库存条件更新与 SKU 校验；**删除批次 07 的秒杀拦截闸**；`cancelBy` 回补 `sold_count` | 批次 07 ledger 明确移交本批清偿（D1） |
| `internal/service/shop/ports.go` | 新增 `IRiskHit`（风控）与 `IAssistReward`（发奖意图）端口 | D3/D6；未装配降级 |
| `internal/controller/shop/shop_v1_index.go` 等 12 桩 | 连线 | 本批端点 |
| `api/shop/v1/index.go` | `IndexRes` 从占位 `Msg` 扩为聚合结构（**api 契约变更**） | `/index` 原为脚手架占位（D2 用户裁定） |

## 四、需求 → 端点覆盖

| FR | 端点 |
|---|---|
| FR-001~005（含秒杀欠账） | 3 + 跨批 order 改动 |
| FR-006~010 | 4、7、8、9 |
| FR-011~014 | 5、10、11、12 |
| FR-015（风控） | 9、12 |
| FR-016 | 6 |
| FR-017/018 | 1（+ 各列表） |
| FR-019 | 全部（桩清零） |

## 五、分层与形态约束

- controller 只做"绑定 → 调 service → 映射响应"；**会员动作端点一律先过 `requireMember`**（批次 06 I8 防线）；公开端点不取 userId。
- 实现形态：struct 方法 + `NewXxxLogic()`（与 `CartLogicImpl`/`AfterSaleLogicImpl`/`ReviewLogicImpl` 同形）。
- 跨域读 `user.nickname`（帮砍/助力人脱敏）**跟随既有先例**（`order_mgmt_impl.go`、`review_impl.go`）；脱敏 helper 在 shop 域内自带（兄弟域不可 import）。
- 生成物：**无变更**（零迁移）。
