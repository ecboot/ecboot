-- 000035_admin_seed 回退：移除种子超管（纯前进式占位亦可，此处给出对称清理）
DELETE FROM `admin_user` WHERE `username` = 'admin' AND `is_super` = 1;
