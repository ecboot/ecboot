# Verification Record: 特性 002 社交电商扩展 Schema（宪法 IV 留证）

执行环境：Windows 11 / Docker 29.6.2（mysql:8.4 容器，宿主端口 13306）/ Java 25 / Spring Boot 4.1.1。
验证方式：`cd apps/api && ./mvnw test -pl ecboot-start -am`（上下文启动即 Flyway 按版本序应用）+ `docker compose exec mysql mysql` 断言查询。

## 基线（T001）

- 修复前置缺陷后基线可验证：`apps/api/pom.xml` enforcer `<executions>` 截断语法错误（自 git 历史 57aab52 完整恢复）；`application.yaml` 补数据源；宿主 3306 被本机 MySQL 服务占用 → 映射改 13306。
- 结果：`BUILD SUCCESS`；`flyway_schema_history` **10 行全 success=1**；业务表 **26** 张。

## 逐域断言（V11~V21，每域一次 BUILD SUCCESS + flyway success=1）

| 版本 | 断言 | 结果 |
|---|---|---|
| V11 | `SHOW INDEX user`：仅 `uk_phone_hash`（uk_phone 已删）；`phone_hash CHAR(64) NOT NULL UNI`；`phone VARCHAR(256)` 密文语义、无明文片段列 | ✅（quickstart 场景四） |
| V12 | 同 `order_item_id` 二次插入 | ✅ `ERROR 1062 (23000)` |
| V13 | 同 `(user_id, spu_id)` 二次收藏 | ✅ `ERROR 1062` |
| V14 | `notify_task.status` 默认 10、状态机注释与 `idx(status,next_retry_time)` 齐备 | ✅ |
| V15 | 字典插入启用/停用两行成功（SF/YT 样例） | ✅ |
| V16 | `SHOW COLUMNS user_relation`：仅 user_id/inviter_id/bind_channel/bind_time——**无祖父列/路径列/层级列**（SC-002/契约 R1）；`user_account.balance` 为有符号 decimal（R2）；分销 8 表齐建 | ✅ |
| V17 | `trade_order.group_buy_team_id` 可空（普通订单零影响）；团员表 `UNIQUE(team_id,user_id)` 重复插入 | ✅ `ERROR 1062` |
| V18 | `flash_sale_item.sold_count` 独立存在，与 `inventory` 零耦合（R4） | ✅ |
| V19 | `point_account.balance` 有符号 int；`user.level` 可空 | ✅ |
| V20 | 档位 `UNIQUE(activity_id,threshold_amount)` 重复插入 ✅ `ERROR 1062`；恒等式（见下） | ✅ |
| V21 | `risk_record.appeal_status` 四态枚举、双索引 | ✅ |

## 恒等式校验（quickstart 场景三，SC-006）

样例订单 QT1（total 210.00 / promotion 35.00 = coupon 10 + 满减 20 + 积分 5 / pay 175.00，订单项分摊与头一致）：

```text
head_violations=0  item_violations=0  pay_violations=0
（SELECT 计数违反三条恒等式的行数，全部为 0）
```

## 空库全量重放（T028，SC-001）

`DROP/CREATE DATABASE ecboot_fresh` 后按**版本序**（`sort -V`，修正了首轮 glob 字典序误序）重放 V1~V21：

- 失败文件数：**0**
- `fresh_tables = 54` ✅
- `trade_order.group_buy_team_id`（可空）/ `coupon_amount`（默认 0.00）增列在位 ✅

> 注：Flyway 版本序执行的完整性由主库 `flyway_schema_history` V1~V21 全部 success=1 逐版证实；空库重放补充证明 DDL 全量顺序可执行。

## 首轮重放的教训（留档）

shell glob `V*.sql` 为字典序（V1, V10, …, V19, V2, V20…），V17/V19/V20 曾因先于 V6 执行而 ALTER 失败——**人工批量执行迁移必须 `sort -V`**；Flyway 自身无此问题。
