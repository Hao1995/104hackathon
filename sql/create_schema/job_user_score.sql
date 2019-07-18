CREATE TABLE `104hackathon-welfare`.`job_user_score` (
  `jobno` INT(11) UNSIGNED NOT NULL,
  `user_id` INT(11) UNSIGNED NOT NULL COMMENT 'Mapping to `users`.`id`',
  `score` INT(11) NOT NULL,
  PRIMARY KEY (`jobno`),
  FOREIGN KEY (`jobno`) REFERENCES `jobs`(`jobno`) ON UPDATE CASCADE,
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
)
ENGINE = InnoDB;
