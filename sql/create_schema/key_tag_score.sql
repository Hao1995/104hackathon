CREATE TABLE `key_tag_score` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `key_id` int(11) NOT NULL,
  `tag_id` int(11) NOT NULL,
  `score` int(11) NOT NULL COMMENT '此關鍵字的平均分數',
  PRIMARY KEY (`id`),
  KEY `FK_QueryKey_key` (`key_id`),
  KEY `FK_Tag_tag` (`tag_id`),
  CONSTRAINT `FK_QueryKey_key` FOREIGN KEY (`key_id`) REFERENCES `query_key` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_Tag_tag` FOREIGN KEY (`tag_id`) REFERENCES `tag` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci