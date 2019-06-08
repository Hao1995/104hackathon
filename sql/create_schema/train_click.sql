CREATE TABLE `104hackathon`.`train_click` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `action` VARCHAR(20) NOT NULL COMMENT '瀏覽104時的行為記錄',
  `jobno` VARCHAR(10) NULL COMMENT '職務代碼',
  `date` VARCHAR(20) NULL COMMENT '此筆log的時間戳記',
  `joblist` TEXT NULL COMMENT '工作列表',
  `querystring` TEXT NULL COMMENT '網址參數',
  `source` VARCHAR(10) NULL COMMENT '產品名稱',
  `key` VARCHAR(200) NULL COMMENT '整理完的關鍵字',
  PRIMARY KEY(`id`))
ENGINE = InnoDB;
