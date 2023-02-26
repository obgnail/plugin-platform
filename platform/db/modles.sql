-- ----------------------------
-- Table structure for plugin_package
-- ----------------------------
-- 插件文件的信息
DROP TABLE IF EXISTS `plugin_package`;
CREATE TABLE `plugin_package`
(
    `id`          bigint(20)                         NOT NULL AUTO_INCREMENT,
    `app_uuid`    varchar(255) COLLATE latin1_bin    NOT NULL COMMENT '插件应用 id',
    `name`        varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '用户名称',
    `size`        bigint(20)                         NULL DEFAULT NULL COMMENT '大小(字节)',
    `version`     varchar(255) COLLATE latin1_bin    NOT NULL COMMENT '版本',
    `create_time` bigint(20)                         NULL DEFAULT NULL COMMENT '创建时间',
    `update_time` bigint(20)                         NULL DEFAULT NULL COMMENT '更新时间',
    `deleted`     tinyint(1)                         NULL DEFAULT NULL COMMENT '删除状态',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `app_uuid` (`app_uuid`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1
  COLLATE = latin1_bin;


-- ----------------------------
-- Table structure for plugin_config
-- ----------------------------
-- 插件配置 yaml下的config
DROP TABLE IF EXISTS `plugin_config`;
CREATE TABLE `plugin_config`
(
    `id`            bigint(20)                         NOT NULL AUTO_INCREMENT,
    `app_uuid`      varchar(255) COLLATE latin1_bin    NOT NULL COMMENT '插件应用 id',
    `instance_uuid` varchar(255) COLLATE latin1_bin    NOT NULL COMMENT '插件实例 id',
    `label`         varchar(255) CHARACTER SET utf8mb4      DEFAULT NULL COMMENT '配置标签',
    `arg_key`       varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '参数名称',
    `arg_value`     text CHARACTER SET utf8mb4         NOT NULL COMMENT '参数值',
    `type`          varchar(36)                        NOT NULL COMMENT '参数类型',
    `required`      tinyint(1)                         NULL DEFAULT NULL COMMENT '是否必填',
    `create_time`   bigint(20)                         NULL DEFAULT NULL COMMENT '创建时间',
    `update_time`   bigint(20)                         NULL DEFAULT NULL COMMENT '更新时间',
    `deleted`       tinyint(1)                         NULL DEFAULT NULL COMMENT '删除状态',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `app_uuid` (`app_uuid`) USING BTREE,
    KEY `instance_uuid` (`instance_uuid`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1
  COLLATE = latin1_bin;


-- ----------------------------
-- Table structure for plugin_instance
-- ----------------------------
-- 插件的信息表
DROP TABLE IF EXISTS `plugin_instance`;
CREATE TABLE `plugin_instance`
(
    `id`            bigint(20)                         NOT NULL AUTO_INCREMENT,
    `app_uuid`      varchar(255) COLLATE latin1_bin    NOT NULL COMMENT '插件应用 id',
    `instance_uuid` varchar(255) COLLATE latin1_bin    NOT NULL COMMENT '插件实例 id',
    `name`          varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '插件名称',
    `version`       varchar(255) COLLATE latin1_bin    NOT NULL COMMENT '插件版本',
    `description`   text CHARACTER SET utf8mb4         NOT NULL COMMENT '插件描述',
    `contact`       varchar(255) CHARACTER SET utf8mb4      DEFAULT NULL COMMENT '插件联系方式',
    `type`          int(11)                            NULL DEFAULT NULL COMMENT '插件类型',
    `status`        tinyint(1)                         NULL DEFAULT NULL COMMENT '插件状态 1 Preparing 2 Running 3 Invalid',
    `apis`          text CHARACTER SET utf8mb4         NOT NULL COMMENT '插件 apis',
    `create_time`   bigint(20)                         NULL DEFAULT NULL COMMENT '创建时间',
    `update_time`   bigint(20)                         NULL DEFAULT NULL COMMENT '更新时间',
    `deleted`       tinyint(1)                         NULL DEFAULT NULL COMMENT '删除状态',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `app_uuid` (`app_uuid`) USING BTREE,
    KEY `instance_uuid` (`instance_uuid`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1
  COLLATE = latin1_bin;


-- ----------------------------
-- Table structure for plugin_user
-- ----------------------------
-- 每个插件建一个user，绕过标品的鉴权
DROP TABLE IF EXISTS `plugin_user`;
CREATE TABLE `plugin_user`
(
    `id`            bigint(20)                         NOT NULL AUTO_INCREMENT,
    `user_uuid`     varchar(255) COLLATE latin1_bin    NOT NULL COMMENT '用户 id',
    `app_uuid`      varchar(255) COLLATE latin1_bin    NOT NULL COMMENT '插件应用 id',
    `instance_uuid` varchar(255) COLLATE latin1_bin    NOT NULL COMMENT '插件实例 id',
    `name`          varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '用户名称',
    `email`         varchar(255) COLLATE latin1_bin         DEFAULT NULL COMMENT '用户邮箱',
    `create_time`   bigint(20)                         NULL DEFAULT NULL COMMENT '创建时间',
    `update_time`   bigint(20)                         NULL DEFAULT NULL COMMENT '更新时间',
    `deleted`       tinyint(1)                         NULL DEFAULT NULL COMMENT '删除状态',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `app_uuid` (`app_uuid`) USING BTREE,
    KEY `instance_uuid` (`instance_uuid`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1
  COLLATE = latin1_bin;


-- ----------------------------------
-- Table structure for plugin_message
-- ----------------------------------
-- notify能力用的业务表(resource.GetAbility.GetNotify用到的表)
DROP TABLE IF EXISTS `plugin_message`;
CREATE TABLE `plugin_message`
(
    `id`           bigint(20)                     NOT NULL AUTO_INCREMENT,
    `message_body` text CHARACTER SET utf8mb4     NOT NULL,
    `message_id`   bigint(11)                     NOT NULL DEFAULT '0' COMMENT '消息id',
    `to_uuids`     varchar(255) COLLATE latin1_bin         DEFAULT NULL,
    `source`       varchar(16) COLLATE latin1_bin NOT NULL COMMENT '消息来源',
    `distinct`     varchar(16) COLLATE latin1_bin NOT NULL COMMENT '消息去向',
    `status`       tinyint(4)                     NOT NULL DEFAULT '0' COMMENT '状态(1:成功 2:失败)',
    `create_time`  bigint(20)                     NULL     DEFAULT NULL COMMENT '创建时间',
    `update_time`  bigint(20)                     NULL     DEFAULT NULL COMMENT '更新时间',
    `deleted`      tinyint(1)                     NULL     DEFAULT NULL COMMENT '删除状态',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1
  COLLATE = latin1_bin;
