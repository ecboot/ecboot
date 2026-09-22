---
description: "任务清单：分销与资金（批次 11）"
---

# Tasks: 分销与资金（017-distribution-fund）

**Feature**: `specs/017-distribution-fund` | **Plan**: [plan.md](./plan.md) | **Spec**: [spec.md](./spec.md) | **Contract**: [contracts/distribution-endpoint-mapping.md](./contracts/distribution-endpoint-mapping.md)

**Input**: 23 端点（user 11 + admin 12）、`IDistributionLogic` 13 方法 + admin 微扩 ~11 方法、跨域钩子 2 个（正向结算端口新增/反向冲销消费落地）、零权限种子 + 迁移 000043（评审修复轮破除"零迁移"原假设, 记账）。

**Tests**: TDD 红→绿；**红线场景 SC-3 全测试化**（两级封顶/自邀/幂等/状态机/并发提现）。

## Phase 1: Foundational

- [x] T001 接口微扩（D3 记账）: `IDistributionLogic` + admin 11 方法 + user 侧缺的 Records/RuleQuery/ShareCode; DTO 勘察补缺
- [x] T002 [P] 测试基座 `distribution_impl_test.go`: 用户/关系/资质/规则 fixture 与自清（精确名, 子表先删）

## Phase 2: US1 关系链与资质（P1）🎯 MVP

- [x] T003 [P] [US1] 测试（红）: 申请防重 60001/审核状态机（重放拒绝）/冻结解冻/**绑定自邀拒绝+一人一链唯一键兜底**/关系页脱敏/两级解析到顶即止
- [x] T004 [US1] 实现 `distribution_impl.go` 关系与资质方法; 连线 user 3 + admin 3 桩

## Phase 3: US2 佣金计提/结算/冲销（P1）

- [x] T005 [P] [US2] 测试（红）: 归因优先级（分享窗口>关系链>不计）/规则命中（商品>分类/未命中不计）/**计提幂等（事件重放）**/保护期结算入账（余额+流水+状态）/**退款冲销**（未结算失效/已结算负额冲销+可负扣回）/冻结推广员不计提
- [x] T006 [US2] 实现 SettleOrder/ConfirmSettle/ReverseOnRefund + 规则命中; shop 域 ports.go 新增正向投递端口 + bootstrap 装配（跨批改动记账）; 消费 ICommissionReverse
- [x] T007 [US2] 连线 admin 佣金记录列表桩

## Phase 4: US3 账户与提现（P1）

- [x] T008 [P] [US3] 测试（红）: 提现申请（余额条件冻结/**并发双申请只成一笔**/金额校验/负余额拒）/审核通过 20 拒绝 50 回退/打款成功 40+渠道单号幂等/失败 60 回退/**状态机不可逆**（10 直打款拒、40 再审拒）/流水双快照
- [x] T009 [US3] 实现 Account/AccountLogs/WithdrawApply/WithdrawList + admin 审核/打款; 连线 user 4 + admin 2 桩

## Phase 5: US4 激励/推广码/查询面（P2）

- [x] T010 [P] [US4] 测试（红）: 推广码稳定/邀请激励一人一次/规则查询命中优先级/各列表分页与筛选
- [x] T011 [US4] 实现 ShareCode/InviteRecords/RuleQuery; 连线 user 3 + admin 1 桩

## Phase 6: Polish & 收尾

- [x] T012 端点级测试: admin 权限探活（**非超管 10005 对照**——批次 10 I4 内化）+ user 会员端点凭证可达/未登录 10003
- [x] T013 `go test ./...` 两连跑全绿; `golangci-lint run` 0 issues
- [x] T014 `check-stub` 对账: user 10→0、admin 28→16（四渠道 39→16）; quickstart 冒烟
- [x] T015 更新 `specs/PROGRESS.md` 批次 11 状态 ✅ 与完成 commit（同 commit）并提交

## Notes（本批硬约束）

- **宪法红线 #1（刑事级）**: 关系链两级封顶——任何路径不得第三级; 重点测试
- **资金铁律**: 账户变更=条件 UPDATE+判行数+同事务双快照流水; 状态机不可逆; 打款渠道单号唯一幂等
- 变更先记账; 事件钩子未装配降级告警（批次 07 同型）; 禁止手改生成物
- 提现打款 V1 人工登记; 结算保护期 V1 常量配置; 推广员等级无表结构 → 恒 0 占位记账

## 完成记录（2026-09-22）

- **23 端点全清**: `check-stub` 对账 **user 10→0、admin 28→16**, 四渠道合计 **39→16**
- **实现**: 新写 `distribution_impl.go`（IDistributionLogic 16 方法 + IDistributionAdminLogic 12 方法, **编译期接口断言**）; 连线 23 桩; 契约微扩 D3（admin 接口 + user 侧 Records/RuleQuery/ShareCode, 记账）
- **跨域钩子**: shop ports 新增 `ICommissionSettle`（Confirm 确认收货投递, 未装配告警降级——计提幂等可补偿重放）; bootstrap 装配正向 SettleOrder + 反向 ICommissionReverse 消费适配器（after_sale_order.order_item_id 按项精确冲销（评审 I3 修正））
- **红线场景 SC-3 全测试化**: 两级封顶（两次单列查询到顶即止）/自邀拒绝+一人一链唯一键兜底/**计提幂等**（事件重放）/保护期结算入账（余额+双快照流水）/**退款冲销**（负额+可负扣回+冲销幂等）/**提现并发双申请只成一笔**（余额条件冻结）/**渠道单号幂等防重复打款**/**状态机不可逆**
- **实现期新发现（测试当场抓到）**: ①user fixture 的 phone_hash 撞 uk_phone_hash（每用户唯一值修复）; ②share_code 列宽 varchar(16) UNIQUE → sonyflake base36 编码（≤14 字符全局唯一）
- **验证**: `go test ./...` 两连跑全绿; `golangci-lint` **0 issues**; 迁移 000043（评审修复轮新增）+零权限种子（000032 七权限点就位）; wiring 测试含**非超管 10005 对照**（批次 10 I4 内化）
