-- 000039_after_sale_order_item_index: after_sale_order 补 order_item_id 索引（013 复审 I-2 / 批次07）
-- 背景: 批次 07 的 C1 修复把 `Apply` 的"读累计额度"改成**锁定读**（SELECT ... FOR UPDATE）——
--   那是抗竞态的承重墙（见 aftersale_impl.go 的注释）。但 `after_sale_order` **没有 order_item_id 索引**,
--   该锁定读因此是全表扫描 + 逐行加 X 锁: 实测"申请 A 行"会阻塞"申请 B 行"（1205 Lock wait timeout）,
--   即每个申请事务在其存续期阻塞全部售后写入, 表越大越久。补索引后锁定读收窄到该行的索引范围。
-- 注: 仅加索引, 不动列; 不影响生成物（dao/model 只由列生成）。
ALTER TABLE `after_sale_order`
  ADD KEY `idx_order_item` (`order_item_id`);
