CREATE TABLE `104hackathon`.`train_action` (
  `action` VARCHAR(20) NOT NULL COMMENT 'viewJob:瀏覽職務 applyJob:應徵職務 saveJob:儲存職務 (註1)',
  `jobno` VARCHAR(10) NULL COMMENT '被點擊的工作',
  `date` VARCHAR(20) NULL COMMENT 'unixtime(millisecond)',
  `source` VARCHAR(10) NULL COMMENT 'app / web / mobileWeb',
  `device` VARCHAR(10) NULL COMMENT 'ios / android，只有source是app才有')
ENGINE = InnoDB;