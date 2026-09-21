-- 000035_admin_seed: 种子超管账号（纯数据，000032 先例）
-- 背景: 000010/000032 均无 admin_user 数据，后台无可登录入口（specs/007-admin-base research D3）。
-- 初始密码: Ecboot@Admin2026（bcrypt cost=10）——部署后首次登录必须修改！
INSERT INTO `admin_user` (`username`, `password_hash`, `real_name`, `is_super`, `status`)
VALUES ('admin', '$2a$10$AX.WGxKFiEoaIxKDzcbdGOIwtl555zjc1zLPFLDoz8mHu80DW1xg2', '平台超管', 1, 1);
