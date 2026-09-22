# 端点到实现映射：售后域（013-after-sale）

**本次交付的 11 个端点**（当前全部为 `CodeNotImplemented` 桩；`make check-stub` 验收口径：shop 21→16、admin 64→58）。

## 一、C 端（shop，5 个）

| # | 端点 | controller 桩文件 | service 方法 | 鉴权 | 状态影响 |
|---|---|---|---|---|---|
| 1 | `POST /shop/after-sales` | `shop_v1_after_sale_create.go` | `NewAfterSaleLogic().Apply` | 会员（Bearer） | 新建 → 10 待审核 |
| 2 | `GET /shop/after-sales` | `shop_v1_after_sale_list.go` | `.List` | 会员 | 读（按状态筛选 + 分页） |
| 3 | `GET /shop/after-sales/{afterSaleNo}` | `shop_v1_after_sale_detail.go` | `.Detail` | 会员（仅本人） | 读 |
| 4 | `POST /shop/after-sales/{afterSaleNo}/cancel` | `shop_v1_after_sale_cancel.go` | `.Cancel` | 会员（仅本人） | {10,20,30} → 91 已撤销 |
| 5 | `POST /shop/after-sales/{afterSaleNo}/logistics` | `shop_v1_after_sale_logistics.go` | `.SubmitReturn` | 会员（仅本人） | 20 填寄回单号（状态不变） |

**注**：`shop_v1_refund_notify.go` 已于批次 06 连线（渠道退款回调），**不在本批**改动范围内；本批的 40→50 推进依赖它。

## 二、后台（admin，6 个）

| # | 端点 | controller 桩文件 | service 方法 | 权限点 | 状态影响 |
|---|---|---|---|---|---|
| 6 | `GET /admin/after-sales` | `admin_v1_admin_after_sale_list.go` | `.AdminList` | `aftersale:list` | 读（状态 + 类型筛选） |
| 7 | `GET /admin/after-sales/{afterSaleNo}` | `admin_v1_admin_after_sale_detail.go` | `.AdminDetail` | `aftersale:list` | 读（含 userId / 凭证） |
| 8 | `POST /admin/after-sales/{afterSaleNo}/approve` | `admin_v1_admin_after_sale_approve.go` | `.Approve` | `aftersale:audit` | 10→20（退货退款）或 10→30→40（仅退款，发起退款） |
| 9 | `POST /admin/after-sales/{afterSaleNo}/reject` | `admin_v1_admin_after_sale_reject.go` | `.Reject` | `aftersale:audit` | 10→90（原因必填） |
| 10 | `POST /admin/after-sales/{afterSaleNo}/confirm-receipt` | `admin_v1_admin_after_sale_confirm_receipt.go` | `.ConfirmReceipt` | `aftersale:audit` | 20→30→40（须已填寄回单号） |
| 11 | `POST /admin/after-sales/{afterSaleNo}/retry-refund` | `admin_v1_admin_after_sale_retry_refund.go` | `.RetryRefund` | `aftersale:refund` | 30→40（重试） |

权限点取自 `api/admin/v1/aftersale.go` 的注释（`aftersale:audit` / `aftersale:refund`；列表/详情沿用 `aftersale:list`，与既有后台列表惯例一致）。**注意**：`aftersale:*` 是否需要登记到 RBAC 权限种子表，取决于批次 01 的权限装载方式——实现时须核对（若权限表按前缀校验，缺失会直接 10005）。

## 三、需求 → 端点覆盖

| FR | 覆盖端点 |
|---|---|
| FR-001/002/003/004 | 1 |
| FR-005 | 2 / 3 |
| FR-006 | 6 / 7 |
| FR-007 | 8 |
| FR-008 | 9 |
| FR-009 | 5 |
| FR-010 | 10 |
| FR-011 | 8 / 10 / 11（发起）+ 既有 `refund_notify`（回调） |
| FR-012 | 11（+ `fail_reason` 留痕） |
| FR-013 | 4 |
| FR-014/015/016 | 8 / 10（完成后副作用：库存回补 / refund_status / 佣金冲销事件） |
| FR-017 | 全部写端点（条件更新 + affected 判定） |

## 四、分层与形态约束（`docs/layer-contracts.md`）

- controller 只做"绑定 → 调 service → 映射响应"，不含业务判断；售后单号不解析为数值（字符串直传，不需要 `parseID`）。
- 实现形态跟随 shop 域既有：**struct 方法**（`AfterSaleLogicImpl`，与 `OrderLogicImpl`/`CartLogicImpl` 同形），`NewAfterSaleLogic()` 构造。
- 跨域：本批读 `trade_order_item` / 写 `inventory` 均在 shop 域内；佣金冲销只投事件（不 import 其他域）。
- 生成物：`after_sale_order` 的 entity/do 由 `gf gen dao` 重新生成（本批唯一生成物变更）。
