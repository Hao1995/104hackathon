CREATE TABLE `job_user_score` (
  `jobno` int(11) unsigned NOT NULL,
  `user_id` int(11) unsigned NOT NULL COMMENT 'Mapping to `users`.`id`',
  `good_score` int(11) NOT NULL,
  `bad_score` int(11) NOT NULL,
  PRIMARY KEY (`jobno`,`user_id`),
  KEY `userID_good_score` (`user_id`,`good_score`),
  KEY `userID_bad_score` (`user_id`,`bad_score`),
  CONSTRAINT `job_user_score_ibfk_1` FOREIGN KEY (`jobno`) REFERENCES `jobs` (`jobno`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `job_user_score_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci