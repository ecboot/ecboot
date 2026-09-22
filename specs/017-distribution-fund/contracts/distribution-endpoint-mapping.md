# 端点到实现映射：分销与资金（017-distribution-fund）

**23 端点**（user 11 + admin 12；验收：`check-stub` user 10→0、admin 28→16，四渠道 39→16）。

## 一、接口

- user 侧 `IDistributionLogic` 13 方法已有（Apply/Status/Relations/BindRelation/InviteRecords/ShareReport/SettleOrder/ConfirmSettle/ReverseOnRefund/Account/AccountLogs/WithdrawApply/WithdrawList）。
- admin 侧微扩（D3, 记账）: `AdminDistributorList/Audit/Freeze`、`AdminRuleList/Create/Update/Delete`、`AdminRecordList`、`AdminWithdrawList/Audit/Pay`、`AdminInviteRecords`。

## 二、端点映射

### user 渠道（11, 全会员）

| # | 端点 | controller 桩 | service |
|---|---|---|---|
| 1 | `POST /user/distribution/apply` | user_v1_dist_apply.go | `.Apply` |
| 2 | `GET /user/distribution/status` | user_v1_dist_status.go | `.Status` |
| 3 | `GET /user/distribution/relations` | user_v1_dist_relations.go | `.Relations` |
| 4 | `GET /user/distribution/commission-rules` | user_v1_dist_rule_query.go | 规则命中查询（新方法 RuleQuery） |
| 5 | `GET /user/distribution/records` | user_v1_dist_record_list.go | `.Records`（新方法: 接口无、补） |
| 6 | `GET /user/distribution/account` | user_v1_dist_account.go | `.Account` |
| 7 | `GET /user/distribution/account/logs` | user_v1_dist_account_logs.go | `.AccountLogs` |
| 8 | `POST /user/distribution/withdraws` | user_v1_withdraw_create.go | `.WithdrawApply` |
| 9 | `GET /user/distribution/withdraws` | user_v1_withdraw_list.go | `.WithdrawList` |
| 10 | `GET /user/distribution/invite-records` | user_v1_invite_record_list.go | `.InviteRecords` |
| 11 | `GET /user/share-code` | user_v1_share_code.go | `.ShareCode`（新方法） |

### admin 渠道（12）

| # | 端点 | controller 桩 | service | 权限 |
|---|---|---|---|---|
| 12 | `GET /admin/distributors` | admin_v1_admin_distributor_list.go | `.AdminDistributorList` | distribution:read |
| 13 | `POST /admin/distributors/{id}/audit` | admin_v1_admin_distributor_audit.go | `.AdminDistributorAudit` | distribution:audit |
| 14 | `POST /admin/distributors/{id}/freeze` | admin_v1_admin_distributor_freeze.go | `.AdminDistributorFreeze` | distribution:audit |
| 15-18 | `GET/POST/PUT/DELETE /admin/commission-rules*` | admin_v1_admin_dist_rule_*.go | `.AdminRule{List,Create,Update,Delete}` | rule:read / rule:manage |
| 19 | `GET /admin/commission-records` | admin_v1_admin_dist_record_list.go | `.AdminRecordList` | distribution:read |
| 20 | `GET /admin/withdraws` | admin_v1_admin_withdraw_list.go | `.AdminWithdrawList` | withdraw:read |
| 21 | `POST /admin/withdraws/{no}/audit` | admin_v1_admin_withdraw_audit.go | `.AdminWithdrawAudit` | withdraw:audit |
| 22 | `POST /admin/withdraws/{no}/pay` | admin_v1_admin_withdraw_pay.go | `.AdminWithdrawPay` | withdraw:pay |
| 23 | `GET /admin/invite-records` | admin_v1_admin_invite_record_list.go | `.AdminInviteRecords` | distribution:read |

## 三、关键校验与防线

| 场景 | 防线 | 失败码 |
|---|---|---|
| 绑定自邀/一人一链 | 应用校验 inviter≠user + uk_user 1062→业务码 | 60001 系 |
| 计提幂等 | (order_item_id, beneficiary, level) 业务键同事务查重 | 静默跳过+日志 |
| 结算/状态迁移 | 条件 UPDATE WHERE status=期望态 判行数 | 条件不命中静默跳过/拒绝 |
| 提现冻结 | `balance >= amount` 条件扣减判行数 | 余额不足业务码 |
| 打款幂等 | uk_channel_order 1062 → 业务码 | 40006 系 |
| 未传 status 类更新 | D1 铁律: affected=0 区分"不存在"vs"同值"（批次 10 I3 教训内化） | — |
| 金额 | FromYuanString 拒 >2 位小数/非正; 佣金=基数×比例% 四舍五入到分（money 库） | 10001 |
