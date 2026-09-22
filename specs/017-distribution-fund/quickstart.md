# Quickstart: 分销与资金（017-distribution-fund）冒烟验证

**Prerequisites**: 应用库 `ecboot` 迁移版本 43（000043 佣金记录唯一键——评审修复 C3/C4, 详见 PROGRESS §五）。命令在 `apps/server/` 下。

## 冒烟（TDD 全绿后手工抽查 5 组）

1. **关系与资质**: A 分享链接 → B 注册绑定 A → B 申请推广员 → 后台审核通过 → B 状态=通过 → B 关系页见上级 A（脱敏）→ 后台冻结 B → 状态=冻结。
2. **佣金全链**: B（上级 A，A 上级 C）购买命中规则商品并确认收货 → 佣金记录: A 待结算 L1、C 待结算 L2（基数=行实付）→ 结算任务 → A/C 余额入账+流水 → 该项退款完成 → A/C 负额冲销、余额扣回（可负）。
3. **提现状态机**: A 余额 100 申请提 60 → 余额 40/冻结 60 → 后台通过 → 打款登记成功（渠道单号 X）→ 40/冻结清零 → 用同渠道单号再登记 → 拒绝。
4. **提现拒绝回退**: 另一单申请 → 审核拒绝 → 冻结回退余额、流水 biz_type=4。
5. **激励与推广码**: 查推广码稳定 → 新用户经码注册 → 邀请激励记录一条（注册时机）→ 重复触发不重复。

## 自动化验收

```bash
cd apps/server
go test ./... -count=1          # 两连跑全绿
golangci-lint run               # 0 issues
for ch in admin common shop user; do
  stub=$(grep -l CodeNotImplemented internal/controller/$ch/*.go | wc -l)
  echo "$ch stub=$stub"         # admin 16 / common 1 / user 0 / shop 0
done
```
