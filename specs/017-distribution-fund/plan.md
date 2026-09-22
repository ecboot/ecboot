# Implementation Plan: 分销与资金（批次11）

**Branch**: `017-distribution-fund` | **Spec**: [spec.md](./spec.md) | **Created**: 2026-09-22

## Technical Context

- Go + GoFrame v2.10.3；user 侧接口 `internal/service/user/distribution.go`（13 方法，含表/红线/归因详注）；admin 侧方法**契约缺口**需微扩（D3）。
- 表 9 张全既有（000016/000028）；权限点 7 个全就位（000032），零权限种子。评审修复 C3/C4 增**迁移 000043**（佣金记录唯一键, 原假设"零迁移"破除并记账）。
- 跨域：shop（订单/售后）→ user（分销结算/冲销）经 ports 投递——既有 `ICommissionReverse`（shop 投递侧已接）补 user 消费实现；新增正向 `ICommissionSettle` 投递端口（确认收货→SettleOrder）。

## Constitution Check

| 门 | 判定 |
|---|---|
| 红线 #1 分销 ≤2 级 | 结构强制（user_relation 无祖父列）+ 应用层二级解析两次单列查询到顶即止；重点测试 |
| 模块化单体/分层 | shop→user 只经 ports 事件投递，无包互引 |
| 可验证交付 | SC-1 桩数 39→16；红线场景 SC-3 全测试化 |
| 简单优先 | 提现打款 V1 人工登记（渠道对接后置）；cron 触发形态契约期裁定 |

## Design Decisions

- **D1 资金防线**（012~015 六轮铁律 + 资金域强化）: 账户变更=单行条件 UPDATE（`balance >= x` 判行数）+ 同事务 account_log 双快照流水（只追加）；提现状态机全部条件更新（WHERE status=期望态）判行数；打款幂等 `uk_channel_order` 1062 → 业务码。
- **D2 计提幂等**: commission_record 以（order_item_id, beneficiary_user_id, level）业务键防重——计提前查存在性（同事务内），事件重放不重复落账。
- **D3 契约微扩（记账）**: `IDistributionLogic` 扩 admin 管理方法（DistributorList/Audit/Freeze、RuleList/Create/Update/Delete、AdminRecordList、AdminWithdrawList/Audit/Pay、AdminInviteRecords）——或按域内惯例落为同接口追加；`distribution_user` 无 level 列 → api Level 恒 0 占位（记账降级）。
- **D4 归因**: share_record 窗口内最近一次分享该 SPU 的用户 > user_relation 两级 > 不计佣；窗口读 config（000028 设计），V1 常量配置化读取。
- **D5 佣金算式**: base=订单项 pay_amount（实付分）；amount = base × rate% 四舍五入到分（`money` 库运算，禁 float；口径注释固化）。
- **D6 冲销**: 退款完成事件 → 未结算（status=1）置 3 失效；已结算（status=2）插入负额记录（reversal_of_id 回指）+ 账户条件扣回（可负）+ 流水 biz_type=5；幂等=原记录已有回指或已存在冲销记录。
- **D7 结算任务**: ConfirmSettle 批量迁移 `status=1 AND 保护期满`（保护期天数读配置，V1 常量）→ 逐条事务入账；返回迁移条数（测试断言用）。
- **D8 端点级测试**: admin 权限探活（非超管 10005 对照——批次 10 I4 教训直接内化）+ user 会员端点凭证可达。

## Risks

- 资金状态机/幂等是最高风险面——SC-3 红线场景全覆盖 + 评审必做变异实验。
- 跨域事件丢失（确认收货时端口未装配）→ 降级告警 + 记账可补偿（结算任务不依赖事件补账 V1 不做，挂账）。
