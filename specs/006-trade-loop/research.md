# Research: 交易链路实现（Phase 0）

2026-09-21。8 项实现决策。

## D1 事务编排：单一本地事务 + 跨域接口回调

- **Decision**: 下单在 `g.DB().Begin()` 单事务内完成：幂等抢占 → 限售校验 → 库存锁定 → 优惠计算 → 券/积分核销 → 余额冻结 → 快照落库 → 状态流水 → Commit。各域锁定/核销方法接受 `*gdb.TX`（005 inventory.Lock 的 `tx any` 已预留）。
- **Rationale**: 单体单库最优解——本地事务 ACID 免补偿；微服务化时该事务边界即拆分断点。
- **Alternatives**: Saga/消息最终一致（被拒：单库场景是过度设计）。

## D2 幂等：唯一约束 + 重复查询返回原单

- **Decision**: `request_token` 落 `uk_user_token` 唯一索引；INSERT 捕获 1062 → 按 token 反查原单返回（校验参数指纹一致，不一致返回幂等冲突）。
- **Rationale**: 数据库最终防线；参数指纹（token+userId+关键金额哈希）防参数篡改。

## D3 优惠计算顺序（既定规则落地）

- **Decision**: 限售校验 → 库存锁定 → 满减自动命中最优档 → 券核销 → 积分抵扣 → 余额冻结（最后）；分摊按行金额比例、尾差记末行（券/满减/积分/余额四套各自独立分摊）。
- **Rationale**: 先满减后券已定档；余额最后冻结（失败即整体回滚，无需补偿）；比例分摊+尾差末行是恒等式成立的唯一保证。

## D4 mock 支付渠道

- **Decision**: `paychannel.Mock` 生成占位唤起参数；测试/联调直接调用 `HandlePayNotify` 驱动成功（构造标准报文+金额）；真实渠道另立特性实现同接口。
- **Rationale**: 回调处理逻辑（幂等四层+联动）是真实价值所在；渠道差异隔离在适配器。

## D5 回调幂等四层（复用 004 契约）

- **Decision**: ① `UPDATE pay_order SET status=20 WHERE pay_no=? AND status=10` 抢占；② 渠道单号唯一；③ 金额校验（回调金额 ≡ 应付现金）；④ 原文留档 `pay_callback_log`。重复通知 → 直接应答成功（幂等应答）。

## D6 金额计算：内部分（int64）

- **Decision**: service 内部金额一律以「分」为 int64 运算（元 string ↔ 分 int64 转换 helper 放 `internal/library/money`）；恒等式在分域断言，DTO 出参再转元 string。
- **Rationale**: AGENTS 金额红线的 Go 落点；int64 分无精度问题且恒等式断言精确。

## D7 秒杀/拼团下单联动

- **Decision**: 秒杀扣 `flash_sale_item` 分账（`sold_count+n WHERE stock-sold>=n`）；拼团参团校验团状态+人数上限（成员行同事务插入）；砍价成交校验状态=待下单并置已下单。取消时逆向回补（秒杀回已售、拼团删成员行释放名额、砍价单回待下单）。

## D8 管理端权限

- **Decision**: 延续 005 形态（Auth 直通 + operator=admin:seed）；权限点真实校验随治理特性。

## 关键事实核对

- 005 已交付 inventory 锁定/核销/释放/回补四操作的接口签名（tx 感知）✓
- user 域接口（IUserCouponLogic.Consume/ReturnBack、IPointLogic.Consume/Refund、IUserAccountLogic 冻结/解冻）已在 004/005 设计，本特性按依赖倒置注册 ✓
- 恒等式三条 + 余额第四式（contracts）为验收断言 ✓
