# Quickstart: 营销后台（016-marketing-admin）冒烟验证

**Prerequisites**: 应用库 `ecboot`（`127.0.0.1:3306`）迁移版本 42（本批零迁移）。命令在 `apps/server/` 下执行（本机无 make）。

## 冒烟（TDD 全绿后的手工抽查，5 组）

1. **券全链路**: 后台建券（满 100 减 20、总量 100）→ 列表可见（receivedCount=0）→ 详情一致 → C 端 `GET /shop/index` 可领券含该券 → 后台停发 → C 端消失 → 记录页空。
2. **满减往返**: 建带 2 档 + 指定商品范围的活动 → 详情回读一致（档位升序）→ 重复门槛提交 → 50008 → 修改为单档全场 → 详情一致 → C 端 `GET /shop/full-reductions` scopeDesc="全场"。
3. **秒杀配置联动**: 建秒杀活动 → items 设置 (SKU, 5.00, 10, 2) → C 端 `GET /shop/activities/flash-sales` 出现该场次（剩余 10）→ 该 SKU 秒杀下单成功 → 后台尝试移除该场次商品 → 已售保护拒绝。
4. **砍价配置**: 建砍价活动（SPU）→ items 设置 (SKU, 100.00, 60.00, 9) → C 端砍价列表可见 maxCutCount=9 → 发起帮砍链路正常（批次 09 行为不回归）。
5. **助力配置**: 建 rewardType=2 活动（pointAmount=100, requiredCount=3）→ C 端助力列表 rewardDesc 含"100" → C 端发起/助力链路正常。

## 自动化验收

```bash
cd apps/server
go test ./... -count=1          # 两连跑全绿
golangci-lint run               # 0 issues
for ch in admin common shop user; do
  stub=$(grep -l CodeNotImplemented internal/controller/$ch/*.go | wc -l)
  echo "$ch stub=$stub"         # admin 28 / common 1 / user 10 / shop 0
done
```
