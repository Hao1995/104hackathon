CREATE TABLE `query_key` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(10) NOT NULL COMMENT '搜索關鍵字的名稱',
  `score` int(11) NOT NULL DEFAULT '0' COMMENT '此關鍵字的平均分數',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=104501 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci