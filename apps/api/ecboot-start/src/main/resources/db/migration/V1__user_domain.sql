-- ============================================================
-- V1 用户域 (user domain)
-- 表: user(用户), user_address(收货地址)
-- 登录双通道: 手机号 + 微信小程序 (password_hash 可空)
-- ============================================================

CREATE TABLE `user` (
  `id`               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `nickname`         VARCHAR(64)     NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar`           VARCHAR(512)    NOT NULL DEFAULT '' COMMENT '头像URL(OSS/CDN,只存链接)',
  `phone`            VARCHAR(20)     NOT NULL COMMENT '手机号(登录凭证,全局唯一)',
  `password_hash`    VARCHAR(100)    NULL DEFAULT NULL COMMENT '密码哈希(bcrypt);微信首次登录未设密码时为NULL',
  `wx_openid`        VARCHAR(64)     NULL DEFAULT NULL COMMENT '微信openid(小程序登录凭证,唯一)',
  `wx_unionid`       VARCHAR(64)     NULL DEFAULT NULL COMMENT '微信unionid(开放平台多应用打通预留)',
  `gender`           TINYINT         NOT NULL DEFAULT 0 COMMENT '性别:0未知 1男 2女',
  `status`           TINYINT         NOT NULL DEFAULT 1 COMMENT '账号状态:1正常 2禁用',
  `register_channel` TINYINT         NOT NULL DEFAULT 1 COMMENT '注册渠道:1微信小程序 2H5 3后台创建',
  `deleted`          TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是(删除后手机号仍占用,复用走后台改绑)',
  `created_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_phone` (`phone`),
  UNIQUE KEY `uk_wx_openid` (`wx_openid`),
  KEY `idx_wx_unionid` (`wx_unionid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';

CREATE TABLE `user_address` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '地址ID',
  `user_id`        BIGINT UNSIGNED NOT NULL COMMENT '所属用户ID',
  `receiver_name`  VARCHAR(64)     NOT NULL COMMENT '收货人姓名',
  `receiver_phone` VARCHAR(20)     NOT NULL COMMENT '收货人手机号',
  `province`       VARCHAR(32)     NOT NULL COMMENT '省',
  `city`           VARCHAR(32)     NOT NULL COMMENT '市',
  `district`       VARCHAR(32)     NOT NULL DEFAULT '' COMMENT '区/县(直筒子市可为空)',
  `detail_address` VARCHAR(255)    NOT NULL COMMENT '详细地址(街道门牌)',
  `is_default`     TINYINT         NOT NULL DEFAULT 0 COMMENT '默认地址:0否 1是(每用户至多一个,应用层保证)',
  `deleted`        TINYINT         NOT NULL DEFAULT 0 COMMENT '软删除:0否 1是',
  `created_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户收货地址表';
