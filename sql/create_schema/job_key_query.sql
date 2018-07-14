CREATE TABLE `job_key_score` (
`jobno` varchar(20) NOT NULL COMMENT '職位',
  `job` varchar(500) NOT NULL COMMENT '職位',
  `key` varchar(200) NOT NULL,
  `good_score` int(11) NOT NULL,
  `bad_score` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci