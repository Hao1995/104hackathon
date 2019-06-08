CREATE TABLE `104hackathon`.`companies` (
  `custno` VARCHAR(40) NOT NULL COMMENT '公司代碼',
  `invoice` INT UNSIGNED NULL COMMENT '公司統編',
  `name` VARCHAR(200) NULL COMMENT '公司名稱',
  `profile` TEXT NULL COMMENT '公司簡介',
  `management` TEXT NULL COMMENT '經營理念',
  `welfare` TEXT NULL COMMENT '公司福利',
  `product` TEXT NULL COMMENT '公司業務介紹')
ENGINE = InnoDB;
