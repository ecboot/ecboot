# Phase 1 数据模型：评价域（014-review）

**依据**: 迁移 `000012_product_review`（建表）+ `000022_review_fixes`（补 `order_no`/`sku_id` 索引）。**本批零迁移**（D1 改口径、D3/D4 用既有列与索引）。

## 一、product_review（既有表，本批只读写，不改结构）

| 列 | 类型 | 说明 / 本批用法 |
|---|---|---|
| id | BIGINT | 评价 ID（出参 `reviewId`） |
| **order_item_id** | BIGINT | **唯一键 `uk_order_item`** —— "一项一评"的**数据库兜底**（D3） |
| order_no / user_id / spu_id / sku_id | — | 归属与筛选键（`idx_spu_audit_time`、`idx_user`、`idx_order_no`、`idx_sku`） |
| spu_name / sku_specs | VARCHAR(128) / JSON | **下单时快照**（写入时从订单项快照取，展示"购买规格"） |
| score | TINYINT(1–5) | 评分 |
| content | VARCHAR(1024) | 评价内容（**可空**——允许只打分） |
| images | JSON | 评价图片（空时写 `[]` 而非 NULL） |
| is_anonymous | TINYINT | 匿名：仅影响**对外展示**（D2） |
| **audit_status** | TINYINT | 0 待审 / 1 通过 / 2 驳回。**本批写入即 1**（D1，用户裁定 V1 自动通过）；**列表只出 1** |
| extra_content / extra_time | VARCHAR(1024) / DATETIME | 追评（**一次**，`extra_content=''` 即未追评——D4 的条件更新判据） |
| reply_content / reply_time | VARCHAR(512) / DATETIME | 商家回复（**本批只读透出**，不写入——D7） |
| deleted / created_at / updated_at | — | 软删与时间（列表一律 `deleted=0`） |

## 二、状态与可见性

```
创建评价（唯一键兜底）
  └─ audit_status = 1 通过        ← V1 自动通过（D1）
         │
         ├─ 商品评价列表 / 汇总：**只出 通过**（且 deleted=0）
         └─ 我的评价：出全部状态（含 待审/驳回），本人可见自己的内容与状态

追评：extra_content 由 '' → 内容（条件更新，仅一次 + 主评 90 天内）
```

| 动作 | 前置条件 | 效果 | 失败码 |
|---|---|---|---|
| Create | 订单项属于本人 且 所属订单**已完成** 且 未评过 | 落评价（audit_status=1） | 未归属→不存在（40009）；**未完成/已取消→40006（状态不允许）**；已评过→40010（`CodeAlreadyReviewed`，含唯一键 1062 转换）；评分越界→10001 |
| Extra | 本人评价 且 `extra_content=''` 且 `created_at >= NOW()-90d` | 写 extra_content/extra_time | 已追评/超期→40011（`CodeExtraReviewDenied`）；他人→不存在 |

**错误码（既有，无需新增）**: `CodeAlreadyReviewed = 40010 已评价`、`CodeExtraReviewDenied = 40011 追评违规（已追评/超期）`。

## 三、查询口径（列表与汇总同源，D5）

```
商品评价列表: WHERE spu_id=? AND audit_status=1 AND deleted=0 [AND score=?] ORDER BY id DESC 分页
汇总(同源同筛，去掉分页): 总数 = COUNT(*)；平均分 = ROUND(AVG(score),1)；分布 = GROUP BY score
我的评价:     WHERE user_id=? AND deleted=0 ORDER BY id DESC 分页（不筛审核状态）
```

## 四、展示脱敏（D2）

| 情形 | 对外展示（商品评价列表） | 我的评价 |
|---|---|---|
| `is_anonymous=1` | 固定占位"匿名用户"（不读昵称） | 本人可见自己的内容（本条也无需展示昵称） |
| `is_anonymous=0` | 昵称**脱敏**（首字符 + 掩码；空昵称给占位）——shop 域内实现（`maskNickname` 在兄弟域 `service/user`，不可 import） | 同上 |

**昵称来源**: `user` 表按 `user_id` 读（**既有先例**: `order_mgmt_impl.go` 已直接读 `dao.User`；包级 import 隔离仍守）。

## 五、校验规则汇总（可测）

| 规则 | 出处 | 失败码 |
|---|---|---|
| 评分 1–5 | api `v:"between:1,5"` + 服务侧兜底 | 10001 |
| 订单项归属本人（他人按不存在） | FR-001 | 不存在（40009） |
| 所属订单已完成 | FR-001 | 40006（状态不允许） |
| 一项一评（含并发） | FR-002 | 40010（唯一键 1062 转业务码） |
| 追评：本人 + 未追评 + 90 天内 | FR-006 | 40011 / 不存在 |
| 追评内容必填 | FR-006 | 10001 |
| 列表只出审核通过 | FR-007/008 | —（口径） |
| 匿名脱敏 | FR-009 | —（口径） |
