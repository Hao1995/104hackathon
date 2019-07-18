CREATE TABLE `104hackathon-welfare`.`welfare_user_score` (
  `id` INT(11) UNSIGNED NOT NULL,
  `user_id` INT(11) UNSIGNED NOT NULL COMMENT 'Mapping to `users`.`id`',
  `welfare_no` INT(11) UNSIGNED NOT NULL COMMENT 'Mapping to `welfare`.`id`',
  `score` TINYINT COMMENT '-127~127. Score can be negative.',
  PRIMARY KEY (`id`),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (`welfare_no`) REFERENCES `welfares`(`id`) ON UPDATE CASCADE
)
ENGINE = InnoDB;
