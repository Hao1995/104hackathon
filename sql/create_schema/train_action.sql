CREATE TABLE `train_action` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `jobno` int(11) UNSIGNED DEFAULT NULL COMMENT '被點擊的工作',
  `date` timestamp NULL DEFAULT NULL,
  `action` varchar(10) DEFAULT NULL COMMENT 'Max size is 10(applyJob). viewJob:瀏覽職務 applyJob:應徵職務 saveJob:儲存職務 (註1)',
  `source` varchar(10) DEFAULT NULL COMMENT 'app / web / mobileWeb',
  `device` VARCHAR(10) DEFAULT NULL COMMENT 'ios / android，只有source是app才有',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;