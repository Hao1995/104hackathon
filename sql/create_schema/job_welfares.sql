CREATE TABLE `job_welfares` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `jobno` int(11) unsigned NOT NULL COMMENT 'Mapping to `jobs`.`jobno`',
  `welfare_no` int(11) unsigned NOT NULL COMMENT 'Mapping to `welfare`.`id`',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_job_welfare` (`jobno`,`welfare_no`),
  KEY `job_welfares_ibfk_2` (`welfare_no`),
  CONSTRAINT `job_welfares_ibfk_1` FOREIGN KEY (`jobno`) REFERENCES `jobs` (`jobno`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `job_welfares_ibfk_2` FOREIGN KEY (`welfare_no`) REFERENCES `welfares` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci