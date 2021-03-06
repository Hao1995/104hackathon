CREATE TABLE `departments` (
  `id` int(11) UNSIGNED NOT NULL COMMENT '類目代碼',
  `name` varchar(50) DEFAULT NULL COMMENT '類目名稱',
  `desc` text COMMENT '說明',
  `hide` varchar(4) COMMENT '是否隱藏',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `districts` (
  `id` bigint(11) UNSIGNED NOT NULL COMMENT '類目代碼',
  `name` varchar(50) DEFAULT NULL COMMENT '類目名稱',
  `desc` text COMMENT '說明',
  `hide` varchar(4) COMMENT '是否隱藏',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `industries` (
  `id` int(11) UNSIGNED NOT NULL COMMENT '類目代碼',
  `name` varchar(50) DEFAULT NULL COMMENT '類目名稱',
  `desc` text COMMENT '說明',
  `hide` varchar(4) COMMENT '是否隱藏',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `job_categories` (
  `id` int(11) UNSIGNED NOT NULL COMMENT '類目代碼',
  `name` varchar(50) DEFAULT NULL COMMENT '類目名稱',
  `desc` text COMMENT '說明',
  `hide` varchar(4) COMMENT '是否隱藏',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;