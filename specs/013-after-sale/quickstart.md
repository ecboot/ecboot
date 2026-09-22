# 验证指南：售后域（013-after-sale）

**Prerequisites**: 应用库 `ecboot`（`127.0.0.1:3306`）已迁移到版本 **38**；测试与运行同源（见 `AGENTS.md`）。命令在 `apps/server/` 下执行（本机无 `make`，直接用等价命令）。

## 一、验收命令（DoD 数字）

```bash
go build ./...                                   # 编译
go test ./... -count=1                           # 全量回归（含既有交易链路零退化）
golangci-lint run                                # 期望 0 issues（本批须保持）
migrate -path migrations -database "mysql://root:root@tcp(127.0.0.1:3306)/ecboot" up   # 应用 000038
for ch in admin shop; do echo -n "$ch: "; grep -l CodeNotImplemented internal/controller/$ch/*.go | wc -l; done
# 期望: admin 58, shop 16（本批 11 端点清零）
```

## 二、端到端场景（服务层测试覆盖，按顺序即为一轮完整演示）

场景脚本落 `internal/service/shop/aftersale_*_test.go`（TDD，先红后绿），均打真实库并自清。

1. **仅退款全链路（P1）**
   造已完成订单（含订单项）→ `Apply`（type=1, q=1）→ 断言状态 10、refund_amount 按行实付折算
   → `Approve` → 断言状态 **40 退款中**（已发起 mock 退款）、audit_time/operator_id 落库
   → 投递成功退款回调 → 断言 **50 已完成** + refund_time 落库 + 订单 `refund_status=1 部分退款`。
2. **退货退款全链路（P1）**
   `Apply`（type=2）→ `Approve` → 断言 **20 待买家寄回**且**未**发起渠道退款
   → 未填寄回单号时 `ConfirmReceipt` → 拒绝
   → `SubmitReturn` 填单号 → `ConfirmReceipt` → 断言 40 → 回调 → 50，且**库存回补**了 quantity（仅退款不回补，须对照断言）。
3. **超额与重复申请（SC-003）**
   行数量 2 → 申请 2 → 再次申请 1 → 40008 且无新单；申请数量 3 → 40008；他人订单项 → 40009。
4. **撤销边界（用户裁定 D4）**
   10/20/30 均可撤 → 91，且可退数量恢复（可再次申请成功）；**40 撤销 → 40006**；50 撤销 → 40006。
5. **退款幂等与重试（SC-004/FR-012）**
   40 → 重复投递成功回调 → 状态不变、refund_time 不被改写；未匹配回调 → 拒绝且留档；
   渠道失败路径（构造 mock 失败）→ 停 30 + `fail_reason` → `RetryRefund` → 40 且 fail_reason 清空。
6. **0 元边界（research D6）**
   行实付为 0（全额优惠）→ 审核通过后**不调用渠道**，直接 50 已完成、refund_no 为空。
7. **首次退款不做二次退款（SC-004）**
   50 状态再 `RetryRefund` → 40006（拒绝），不产生第二次出款。

## 三、桩清零与既有链路回归

- `make check-stub` 等价命令输出：**shop 21→16、admin 64→58**（本批 11 个端点涉及的 controller 文件不再含 `CodeNotImplemented`）。
- 既有链路回归：批次 01~06 的全部测试必须保持绿（`go test ./...` 两连跑）——本批触及 `inventory`/`trade_order`（写）故尤其相关。

## 四、不做的事（边界）

- 不改 `refund_notify`（批次 06 已连线并修正）；不改支付/订单域既有实现。
- 不实现佣金**结算**（只投事件出口，结算属批次 11）。
- 不接真实微信/支付宝退款渠道（mock 渠道；真实渠道另立特性）。
