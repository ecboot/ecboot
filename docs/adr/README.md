# Architecture Decision Records

架构决策记录索引（按时间序）。ADR 用于记录"难逆转、事后看意外、真实取舍"的决策；修订直接改原文并注明。

| # | 决策 | 状态 |
|---|---|---|
| [0001](0001-inventory-lock-model.md) | 库存采用"下单锁定"三段式，Redis 仅作前置挡板，DB 为唯一事实源 | 有效 |
| [0002](0002-no-physical-fk-immutable-documents.md) | 不建物理外键；交易单据与审计日志永不软删 | 有效（注销范围注记见文内） |
| [0003](0003-distribution-two-level-cap.md) | 分销关系链两级封顶（法规红线，结构强制） | 有效 |
