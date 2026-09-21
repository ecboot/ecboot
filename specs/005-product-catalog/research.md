# Research: 商品目录域实现（Phase 0）

2026-09-21。6 项实现决策（服务接口已在 internal/service/shop 定义，本文件定实现策略）。

## D1 价格区间与排序：SPU 冗余列（price_min/price_max）+ 变更时刷新

- **Decision**: 迁移 000034 为 `product_spu` 增加 `price_min/price_max DECIMAL(10,2) NULL`；SKU 创建/编辑/启停/删除时**同事务刷新**冗余值；列表查询与价格排序直接走冗余列（可索引）。
- **Rationale**: 列表/搜索的"价格区间筛选+价格排序"是高频路径，JOIN 子查询聚合每页都算；冗余列以"SKU 变更时一次写"换"每次读零聚合"——商品读远多于写，收益明确。
- **Alternatives**: 查询时 JOIN 聚合（被拒：每页聚合 + 排序无法走索引）；ES（演进路径，不改契约）。

## D2 dao 扩表生成

- **Decision**: `gf gen dao -l ... -t "product_category,product_brand,product_spu,product_sku,inventory,inventory_log,product_review"` 一次生成七表 dao/entity/do。
- **Rationale**: 强类型数据访问（契约 §二 禁 map 传参）依赖 do/entity 生成物。
- **注**: product_review 生成仅为评价汇总/列表读取（写属评价特性）。

## D3 列表可售与价格口径

- **Decision**: 列表"价格区间"= 启用 SKU 的 min~max；全部 SKU 禁用时区间为空、卡片标不可下单；"可售"判定 = SPU 上架 AND SKU 启用 AND 可售>0（详情 SKU 级同口径）。
- **Rationale**: spec FR-003/FR-004 与边界场景一致；口径集中在一个转换函数（entity→DTO）。

## D4 管理端权限校验落地（004 遗留补齐）

- **Decision**: 管理组 Auth 中间件在本特性实现**简单版权限校验**：Bearer 后台会话（后台登录属后续特性——V1 先以 `is_super` 种子账号 + 直通形态运行），权限点校验实现为 `consts` 常量查 `admin_role_permission` 关联（is_super 直通）；管理写接口挂接已定义权限点。后台登录/角色管理本身属治理特性。
- **Rationale**: 004 FR-014 的权限点挂接不能无限期占位；登录缺失用种子超管账号过渡（应用初始化写入，密码 bcrypt）。
- **Alternatives**: 等后台登录特性后再开管理接口（被拒：商品管理是商品域验收的一部分）。

## D5 搜索实现与转义

- **Decision**: `name LIKE '%kw%'`（ESCAPE 转义 %/_/'），分类/品牌/价格筛选同列表；排序复用 D1。
- **Rationale**: LIKE 在万级数据 + idx 前缀下满足契约 1 秒内；转义防通配符注入。
- **Alternatives**: 全文索引/ES（演进）。

## D6 测试数据策略

- **Decision**: service 层测试**自建数据+前置清理**（按唯一键 hash/no 删除），确定性配置注入（gdb/gredis SetConfig，与 003 同款 init）；断言到行为（错误码/状态/区间）。
- **Rationale**: 003 已验证该模式；不依赖迁移种子，测试可重复。

## D7 足迹 UPSERT 归属

- **Decision**: 详情接口软触足迹（viewer>0 时）——UPSERT 语句内联在商品实现（单行 ON DUPLICATE KEY UPDATE last_view_at/view_count），不抽独立组件（V1 最简；消息/清理任务属用户域特性）。
- **Rationale**: spec FR-006; 一条 SQL 的量级不值得抽象。

## 关键事实核对

- service/shop/product.go 与 inventory.go 接口已定义（方法签名冻结）✓
- errcode 3xxxx 段已登记（30001~30007）✓
- 权限点常量已就绪（consts/permission.go）✓
- specs_hash 为 DB 生成列（V23），插入自动计算，并发冲突→30006 映射 ✓
