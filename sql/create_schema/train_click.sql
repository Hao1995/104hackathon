CREATE TABLE `train_click` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `action` varchar(10) DEFAULT NULL COMMENT 'Max size is 10(clickApply).',
  `jobno` int(11) UNSIGNED DEFAULT NULL,
  `date` timestamp NULL DEFAULT NULL,
  `joblist` text,
  `querystring` text,
  `source` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci