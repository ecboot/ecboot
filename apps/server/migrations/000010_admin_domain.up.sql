-- ============================================================
-- V10 后台域 (admin domain) — RBAC + 审计
-- 表: admin_user, admin_role, admin_user_role,
--     admin_permission, admin_role_permission,
--     admin_login_log(登录审计), admin_operation_log(操作审计)
-- RBAC: 账号-角色 M:N, 角色-权限 M:N; is_super 跳过权限校验
-- 权限树: 菜单(1)/按钮操作(2)/接口(3)统一在 admin_permission
-- 审计: 登录日志(成功/失败) + 操作日志(AOP切面异步写)
-- 初始超级管理员由应用启动初始化(密码哈希需应用生成)
-- ============================================================

CREATE TABLE `admin_user` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '后台账号ID',
  `username`        VARCHAR(32)     NOT NULL COMMENT '登录名(唯一)',
  `password_hash`   VARCHAR(100)    NOT NULL COMMENT '密码哈希(bcrypt)',
  `real_name`       VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '姓名',
  `is_super`        TINYINT         NOT NULL DEFAULT 0 COMMENT '超级管理员:1是(跳过权限校验) 0否',
  `status`          TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1正常 2禁用',
  `last_login_time` DATETIME        NULL DEFAULT NULL COMMENT '最后登录时间',
  `deleted`         TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='后台账号表(与C端user完全隔离)';

CREATE TABLE `admin_role` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name`        VARCHAR(32)     NOT NULL COMMENT '角色名称(如:运营/客服/财务)',
  `code`        VARCHAR(32)     NOT NULL COMMENT '角色编码(唯一,如 ops/service/finance)',
  `description` VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '角色描述',
  `status`      TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0停用',
  `deleted`     TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='后台角色表(权限集合)';

CREATE TABLE `admin_user_role` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `admin_id`   BIGINT UNSIGNED NOT NULL COMMENT '后台账号ID',
  `role_id`    BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_admin_role` (`admin_id`, `role_id`),
  KEY `idx_role` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='账号-角色关联表(M:N)';

CREATE TABLE `admin_permission` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '权限ID',
  `parent_id`  BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父权限ID,0为根(菜单树)',
  `name`       VARCHAR(64)     NOT NULL COMMENT '权限名称(如:商品管理/SPU上架)',
  `code`       VARCHAR(64)     NOT NULL COMMENT '权限编码(如 product:spu:create;接口权限与API路径对应)',
  `type`       TINYINT         NOT NULL COMMENT '类型:1菜单 2按钮/操作 3接口',
  `path`       VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '前端路由地址(菜单)或API路径(接口)',
  `sort`       INT             NOT NULL DEFAULT 0 COMMENT '同级排序',
  `status`     TINYINT         NOT NULL DEFAULT 1 COMMENT '状态:1启用 0禁用',
  `deleted`    TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='后台权限表(菜单/按钮/接口统一树)';

CREATE TABLE `admin_role_permission` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `role_id`       BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
  `permission_id` BIGINT UNSIGNED NOT NULL COMMENT '权限ID',
  `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_perm` (`role_id`, `permission_id`),
  KEY `idx_perm` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色-权限关联表(M:N)';

CREATE TABLE `admin_login_log` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `username`    VARCHAR(32)     NOT NULL COMMENT '尝试登录的用户名(含失败,账号可能不存在)',
  `admin_id`    BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '成功时回填后台账号ID',
  `login_status` TINYINT        NOT NULL COMMENT '结果:1成功 2失败(密码错误) 3失败(账号禁用/不存在)',
  `ip`          VARCHAR(45)     NOT NULL DEFAULT '' COMMENT '来源IP(IPv6最长45字符)',
  `user_agent`  VARCHAR(512)    NOT NULL DEFAULT '' COMMENT '浏览器User-Agent',
  `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  PRIMARY KEY (`id`),
  KEY `idx_username_time` (`username`, `created_at`),
  KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='后台登录审计日志(只追加,含失败尝试,防暴力破解分析依据)';

CREATE TABLE `admin_operation_log` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `admin_id`      BIGINT UNSIGNED NOT NULL COMMENT '操作账号ID',
  `username`      VARCHAR(32)     NOT NULL COMMENT '操作人用户名(冗余快照,账号删除后仍可读)',
  `module`        VARCHAR(64)     NOT NULL COMMENT '业务模块(如 商品/订单/售后/权限)',
  `operation`     VARCHAR(64)     NOT NULL COMMENT '操作(如 上下架/改价/退款审核/角色授权)',
  `method`        VARCHAR(8)      NOT NULL DEFAULT '' COMMENT 'HTTP方法(GET/POST/PUT/DELETE)',
  `request_uri`   VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '请求路径',
  `request_params` JSON           NULL COMMENT '请求参数(敏感字段脱敏后留档)',
  `result_status` TINYINT         NOT NULL DEFAULT 1 COMMENT '结果:1成功 0失败',
  `error_msg`     VARCHAR(512)    NOT NULL DEFAULT '' COMMENT '失败原因',
  `ip`            VARCHAR(45)     NOT NULL DEFAULT '' COMMENT '来源IP',
  `cost_ms`       INT UNSIGNED    NOT NULL DEFAULT 0 COMMENT '耗时(毫秒)',
  `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_admin_time` (`admin_id`, `created_at`),
  KEY `idx_module_time` (`module`, `created_at`),
  KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='后台操作审计日志(只追加,AOP切面异步写,写失败不阻断业务)';
