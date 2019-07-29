CREATE TABLE `train_click` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `action` varchar(10) DEFAULT NULL COMMENT 'Max size is 10(clickApply).',
  `jobno` int(11) unsigned DEFAULT NULL,
  `date` timestamp NULL DEFAULT NULL,
  `joblist` text,
  `querystring` text,
  `source` varchar(10) DEFAULT NULL,
  `key` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `key` (`key`),
  KEY `train_click_ibfk_1` (`jobno`),
  CONSTRAINT `train_click_ibfk_1` FOREIGN KEY (`jobno`) REFERENCES `jobs` (`jobno`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=748936 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

--- VARCHAR(255) About 255/4 = 63 words