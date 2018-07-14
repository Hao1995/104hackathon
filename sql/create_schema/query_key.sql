CREATE TABLE `query_key` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(200) NOT NULL COMMENT '搜索關鍵字的名稱',
  `good_score` int(11) NOT NULL DEFAULT '0',
  `bad_score` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=229502 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci