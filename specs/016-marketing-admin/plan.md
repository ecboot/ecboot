# Implementation Plan: 营销后台（批次10）

**Branch**: `016-marketing-admin` | **Spec**: [spec.md](./spec.md) | **Created**: 2026-09-22

## Technical Context

- Go + GoFrame v2.10.3 模块化单体（`apps/server`）：`api/admin/v1` → `internal/controller/admin` → `internal/service/shop`（`ICouponLogic`/`IActivityLogic` 已有定义、待补实现）→ dao/model（生成物）。
- 表全部既有（券 `coupon`/`user_coupon`；满减 `promotion_activity`/`_ladder`/`_scope`；拼团/秒杀/砍价/助力活动+场次商品表 000017/000018/000020/000029）。**零迁移**。
- 权限点全部已在 000032 种子与 `consts/permission.go` 就位（勘察确认），**零权限迁移**；挂载沿用后台鉴权中间件的权限判定形态（批次 01~04 同口径）。
- 金额口径 `DECIMAL(10,2)` 元 / Go int64 分（`library/money`）；分布式单号类 ID `library/idgen`（本批无新增单号场景）。

## Constitution Check

| 门 | 判定 |
|---|---|
| 模块化单体/分层 | 控制器→service→dao；四渠道隔离；admin 渠道不动 user/common |
| 统一技术栈 | 无新依赖 |
| 中文优先 | 代码注释/错误文案中文 |
| 可验证交付 | SC-1 桩数 58→28 数字对账；TDD 红→绿 |
| 简单优先 | 接口已有（实现跟随既有契约）；死契约方法（XxxDetail×3）保留不实现、记账 |

## Design Decisions

- **D1 写入防线**（012~015 五轮评审铁律）：状态/配置变更一律条件 UPDATE 判 RowsAffected；唯一键兜底（满减档位 `uk_activity_threshold`→50008、场次商品 `uk_activity_sku`→业务码）并 1062 转业务码。
- **D2 全量替换语义**：满减档位/范围、三类场次商品为"全量替换"——事务内先删后插（软删活动行不动）；**已售保护**：秒杀场次商品被订单项引用（`trade_order_item.flash_sale_item_id`）且已被售（sold_count>0）时禁止移除（拒绝并说明，防活动库存账悬空）。
- **D3 契约微扩（记账）**：①`ICouponLogic` 补 `AdminDetail`（api 有 `GET /coupons/{id}` 而接口漏定义, 批次 01/08 同例）；②`model.CouponTemplate` +`ValidType` 字段（api 列表项需输出）。
- **D4 死契约**：`IActivityLogic` 的 `GroupBuyDetail`/`FlashSaleDetail`/`BargainDetail` 无对应端点（api 只有 items 写端点）——保留接口方法不实现（批次 08 Reply/Audit 同模式），实现文件不提供该方法会有编译错误吗？——接口由实现者全部实现，故**三个 Detail 方法一并实现**（复用同一活动装配函数，成本低且避免"文字承诺"缺口）；无端点消费、记账。
- **D5 助力奖励落库**：rewardType=1 → `reward_ref`=券 ID（校验券存在且未删）；rewardType=2 → `reward_ref`=0、`config` JSON 存 `{"pointAmount":N}`；`RewardDesc` 组装"邀 N 人得券/积分 X"。
- **D6 per_limit 语义**：接受 >1 配置、如实存储与回显，契约注释声明">1 因 `uk_activity_user` 不生效（每人一次）"——不静默改值（批次 09 I1 决策延续）。
- **D7 N+1 收敛**：列表内档位/范围/SPU 简介改为批量 IN 查询（同文件域内清偿批次 09 M8/复审-8；活动/场次商品行数有界）。
- **D8 C 端联动口径**：后台停发/软删即时影响 C 端——复用 C 端列表既有过滤（status=1 AND deleted=0），管理侧不另设口径。

## Risks

- 满减 scope 配置错误直接影响下单计价（计价侧认范围属 012 文件、已挂账）——本批以"详情回读==提交"的往返测试钉住配置面。
- 场次商品替换与进行中活动的并发下单竞态：替换走事务 + 已售保护校验在事务内以锁定读复核。
