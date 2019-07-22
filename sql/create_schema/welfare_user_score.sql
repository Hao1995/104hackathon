CREATE TABLE `welfare_user_score` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned NOT NULL COMMENT 'Mapping to `users`.`id`',
  `welfare_no` int(11) unsigned NOT NULL COMMENT 'Mapping to `welfare`.`id`',
  `score` tinyint(4) DEFAULT NULL COMMENT '-127~127. Score can be negative.',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_welfare_user` (`user_id`,`welfare_no`),
  KEY `welfare_no` (`welfare_no`),
  CONSTRAINT `welfare_user_score_ibfk_1` FOREIGN KEY (`welfare_no`) REFERENCES `welfares` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `welfare_user_score_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci