CREATE TABLE IF NOT EXISTS `{{upload_tips}}`  (
  `team_uuid` varchar(8) COLLATE latin1_bin NOT NULL DEFAULT '' COMMENT '团队uuid',
  `content` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '内容',
  `update_time` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`team_uuid`)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;