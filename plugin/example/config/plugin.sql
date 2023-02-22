CREATE TABLE IF NOT EXISTS `{{user}}`  (
  `uuid` varchar(8) COLLATE latin1_bin NOT NULL DEFAULT '' COMMENT 'uuid',
  `name` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '名字',
  PRIMARY KEY (`uuid`)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;