CREATE TABLE `104hackathon-welfare`.`job_welfares` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `jobno` int(11) unsigned NOT NULL COMMENT 'Mapping to `jobs`.`jobno`',
  `welfare_no` int(11) unsigned NOT NULL COMMENT 'Mapping to `welfare`.`id`',
  PRIMARY KEY (`id`),
  KEY `jobno` (`jobno`),
  KEY `welfare_no` (`welfare_no`),
  CONSTRAINT `job_welfares_ibfk_1` FOREIGN KEY (`welfare_no`) REFERENCES `welfares` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci