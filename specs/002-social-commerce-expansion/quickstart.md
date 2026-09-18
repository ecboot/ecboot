# Quickstart: 社交电商能力扩展 Schema 验证指南

五个可执行场景，覆盖 spec 全部 7 条成功指标（SC-001~SC-007）。命令均可复制执行；MySQL 经 compose 提供（服务名 `mysql`，库 `mydatabase`，用户 `myuser`/`secret`）。

## 前提准备

```bash
cd apps/api && docker compose up -d mysql
# 执行迁移的两种方式（任选）：
# A. 启动应用，Flyway 自动执行（需本机 JDK25；首次会连 Redis/ES，建议 Docker Desktop 已启动）
cd apps/api && ./mvnw spring-boot:run -pl ecboot-start -am
# B. 仅验证 SQL：用 mysql 客户端按版本序手工执行 V11~V21（可重复用于 CI 冒烟）
docker compose exec mysql mysql -umyuser -psecret mydatabase < 某个V文件
```

## 场景一：全量迁移执行（SC-001）

**步骤**：空库（或已应用 V1~V10 的库）执行全部迁移。
**预期**：Flyway 输出无报错；校验版本与表数：

```bash
docker compose exec mysql mysql -umyuser -psecret mydatabase -e "
SELECT version, description, success FROM flyway_schema_history ORDER BY installed_rank;"
# 预期：V1..V21 共 21 行 success=1
docker compose exec mysql mysql -umyuser -psecret mydatabase -e "
SELECT COUNT(*) AS tables_count FROM information_schema.tables WHERE table_schema='mydatabase' AND table_name NOT LIKE 'flyway%';"
# 预期：54（既有 26 + 新增 28）
```

## 场景二：约束注入——违反唯一约束必须被数据库拒绝（SC-003）

```bash
# ① 重复收藏（FR-009）
docker compose exec mysql mysql -umyuser -psecret mydatabase -e "
INSERT INTO user_favorite(user_id, spu_id) VALUES (1, 100);
INSERT INTO user_favorite(user_id, spu_id) VALUES (1, 100);"
# 预期：第 2 条报 ERROR 1062 Duplicate entry '1-100' for key 'uk_user_spu'

# ② 重复主评价（FR-006）
#    INSERT product_review 同一 order_item_id 两次 → ERROR 1062 for key 'uk_order_item'

# ③ 同团同用户重复参团（FR-015）
#    INSERT group_buy_team_member 同 (team_id, user_id) 两次 → ERROR 1062 for key 'uk_team_user'

# ④ 同活动同门槛（FR-020）
#    INSERT promotion_activity_ladder 同 (activity_id, threshold_amount) 两次 → ERROR 1062

# ⑤ 分销两级封顶结构审查（SC-002 / 契约 R1）
docker compose exec mysql mysql -umyuser -psecret mydatabase -e "
SHOW COLUMNS FROM user_relation;"
# 预期：仅 user_id / inviter_id / bind_channel / bind_time / created_at / updated_at
#      ——无祖父列、无路径列、无层级列（三级关系结构上无处落库）
```

## 场景三：优惠三构成恒等式（SC-006）

样例数据（手工 INSERT 一笔三优惠订单，尾差记末行）：

```sql
-- trade_order: total=210.00, promotion_amount=35.00, freight=0, pay=175.00
--   coupon_amount=10.00, full_reduction_amount=20.00, point_amount=5.00(500分)
-- trade_order_item 两行：分摊合计须等于头明细
SELECT order_no,
       coupon_amount + full_reduction_amount + point_amount AS detail_sum,
       promotion_amount
FROM trade_order
WHERE detail_sum <> promotion_amount;   -- 预期：空结果集（恒等成立）
SELECT order_id,
       SUM(coupon_amount) c, SUM(full_reduction_amount) f, SUM(point_amount) p
FROM trade_order_item GROUP BY order_id
HAVING c <> (SELECT coupon_amount FROM trade_order WHERE id = order_id)
    OR f <> (SELECT full_reduction_amount FROM trade_order WHERE id = order_id)
    OR p <> (SELECT point_amount FROM trade_order WHERE id = order_id);
-- 预期：空结果集
```

## 场景四：手机号加密改造断言（SC-004 / 契约 R3）

```bash
docker compose exec mysql mysql -umyuser -psecret mydatabase -e "
SHOW INDEX FROM user WHERE Key_name LIKE 'uk_phone%';
SHOW COLUMNS FROM user LIKE 'phone%';"
# 预期：仅 uk_phone_hash 存在、uk_phone 已删除；
#      phone 列注释为"手机号密文"、phone_hash CHAR(64) NOT NULL、无尾号/明文片段列
```

## 场景五：余额可负与活动库存分账断言（SC-005 / 契约 R2/R4）

```bash
docker compose exec mysql mysql -umyuser -psecret mydatabase -e "
SHOW COLUMNS FROM user_account LIKE 'balance';
SHOW COLUMNS FROM point_account LIKE 'balance';"
# 预期：两列类型为 decimal(10,2)（有符号，无 unsigned）

docker compose exec mysql mysql -umyuser -psecret mydatabase -e "
SHOW COLUMNS FROM flash_sale_item;"
# 预期：stock_count/sold_count 在本表；inventory 表结构与此零耦合（分账成立）
```

## 附注

- 回滚边界：按域一文件（契约 §1）；Flyway 不做 down，回滚=恢复备份后重放至目标版本。
- 本指南不包含实现代码；迁移脚本明细属 tasks/实现阶段产出。
- 文档同步三件（schema-design.md / CONTEXT.md / ADR-0003）随实现任务一并交付，验收对照 contracts §5。
