# Phase 0 研究与决策：评价域（014-review）

**Date**: 2026-09-22 | **Spec**: [spec.md](./spec.md)

本批性质：**接口已有三方法（Create/Extra/MyList）、一个端点缺方法、4 个 controller 桩连线**；预期零迁移。

## D1 审核口径：V1 自动通过（用户裁定 B）

- **Decision**: 新评价写入时 `audit_status = 1（通过）`；审核字段与 `IReviewLogic.Audit` 能力保留。商品评价列表只出"通过"的口径**不变**。
- **Rationale**: 本批只有 C 端 4 个端点，`api/` 里**没有后台审核契约**（13 批规划内无后台评价管理批）。若严格按既有注释"默认待审 + 仅通过可见"，则**无人能审核** → 商品评价列表恒为空、评价功能事实上不可用。
- **Alternatives considered**: ① 保持待审（忠实契约但功能死锁）；② 本批扩容补后台审核面（超批次范围、要改 218 端点账目）。用户选 ①→ 改为 V1 自动通过，并把"UGC 未经审核即公开"记为**合规债务**（公开上线前必修）。
- **代价与记账**: 已记 `PROGRESS §五`（批次 08 行）。

## D2 展示评价人：shop 域直接读 `user` 表（跟随既有先例）+ 域内自带脱敏

- **Decision**: 商品评价列表/我的评价展示评价人时，从 `user` 表读昵称，并在 **shop 域内**做脱敏（匿名一律显示固定占位，非匿名显示脱敏昵称）。
- **Rationale**: ① **既有先例**：`service/shop/order_mgmt_impl.go` 已直接读 `dao.User`（后台按手机号/ID 搜用户）——"兄弟域经事件协作"约束的是 **Go 包级 import**，不是表级读取；② `maskNickname` 位于 `service/user/invite_impl.go`，**不可**被 shop import（兄弟域），故在 shop 内置等价脱敏函数并注明原因。
- **Alternatives considered**: ① 通过 ports 新增"用户昵称"端口（更"干净"，但为一个展示字段新增端口+装配层改动，成本大于收益，宪法 V）；② 在 `product_review` 落昵称快照（改表=新增迁移，且昵称变更后快照会陈旧）。

## D3 一项一评的并发防线：唯一索引 + 冲突转业务码 + 条件写入

- **Decision**: `Create` 先做"归属/订单已完成/未评过"校验，写入依赖表上既有的 **`uk_order_item` 唯一索引**兜底；捕获唯一键冲突（MySQL 1062）后转成业务错误"已评价"。**不新增迁移**。
- **Rationale**: 表已有唯一索引（勘察确认），这是最强的防线；批次 07 的 C1 教训是"读-改窗口 + 无唯一约束 = 可重复落单"，本批天然免疫，但**必须把 1062 映射成业务码**（否则用户看到系统错误，且重试无意义）。
- **Alternatives considered**: 行锁事务（如批次 07 的 Apply）——在唯一索引已存在时属冗余；仍保留"先查再插"的顺序以获得友好错误信息，唯一索引负责竞态兜底。

## D4 追评限一次 + 90 天：条件更新 + 时间比较

- **Decision**: 追评以**条件更新**落库（`WHERE id=? AND user_id=? AND extra_content=''` 且 `created_at >= NOW()-90 天`），`affected=0` 时再回查区分"已追评/超期/不存在"以给准确错误；不新增列。
- **Rationale**: `extra_content` / `extra_time` 列已存在（000012），"未追评"即 `extra_content=''`——用它做条件更新的判据，比"读-判断-写"更抗并发（同批次 07 的铁律）。
- **Alternatives considered**: 新增 `extra_count` 列（无必要）。

## D5 商品评价汇总口径：只统计审核通过

- **Decision**: 汇总（平均分/分布/总数）与列表**同源同筛**（`audit_status=1 AND deleted=0 [AND score=?]`）；平均分保留 1 位小数。
- **Rationale**: 列表与汇总口径不一致是评价功能的经典缺陷（列表只显示 3 条而汇总说 100 条）；`model.ReviewSummary`（Avg/Distribution/Total）已存在，直接消费。
- **Alternatives considered**: 汇总统计全部审核状态（会泄露待审/驳回的存在性，且与列表矛盾）。

## D6 缺方法的契约微扩：`ProductList`

- **Decision**: 在 `IReviewLogic` 补 `ProductList(ctx, spuId int64, score int, page model.PageReq) (*model.PageResult[model.ReviewCard], error)`，并在 `ReviewCard` 之外提供汇总（见 D5 的返回形态）。
- **Rationale**: 第 4 个端点在 api 契约里存在（`api/shop/v1/product.go` 的 `ProductReviewListReq/Res`），而 service 接口漏了方法——属"接口漏定义"，同批次 01/02/04 的 `AdminDetail`/`AdminSpuListReq.Status` 先例。
- **形态**: 返回 `(summary, list)`——汇总与列表一起算，避免调用方两次查询口径漂移。

## D7 死契约记账：`Reply`/`Audit`/`PendingAudit` 不动

- **Decision**: 这三个后台方法**本批不实现**（无 API 端点，13 批规划内无对应批）。
- **Rationale**: 实现无出口的能力属死代码；但移除它们会影响未来后台评价面——故**保留并记账**（`PROGRESS §五`）。
- **注**: 本批不写回复；但列表要**透出已存在的** `reply_content`（后台手工/未来批次写入的）。

## 未知项收敛

| 未知 | 结论 |
|---|---|
| 是否新增迁移 | **否**（D1 改口径而非改结构；D3 依赖既有唯一索引；D4 用既有列） |
| 评价人展示 | shop 域读 `user` 表 + 域内脱敏（D2） |
| 缺方法 | 补 `ProductList`（D6） |
| 审核口径 | V1 自动通过（D1，用户裁定） |

**结论**：无遗留 NEEDS CLARIFICATION，可进入 Phase 1。
