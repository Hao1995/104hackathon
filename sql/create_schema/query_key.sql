CREATE TABLE `query_key` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(10) NOT NULL COMMENT '搜索關鍵字的名稱',
  `score` INT NOT NULL DEFAULT 0 COMMENT '此關鍵字的平均分數',
  PRIMARY KEY (`id`)
)
ENGINE = InnoDB;
