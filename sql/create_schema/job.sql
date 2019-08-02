CREATE TABLE `104hackathon-welfare-1st`.`job` (
  `custno` VARCHAR(40) NOT NULL,
  `jobno` INT UNSIGNED NULL,
  `job` TEXT NULL,
  `jobcat1` INT NULL,
  `jobcat2` INT NULL,
  `jobcat3` INT NULL,
  `edu` TINYINT UNSIGNED NULL,
  `salary_low` INT UNSIGNED NULL,
  `salary_high` INT UNSIGNED NULL,
  `role` TINYINT UNSIGNED NULL,
  `language1` MEDIUMINT UNSIGNED NULL,
  `language2` MEDIUMINT UNSIGNED NULL,
  `language3` MEDIUMINT UNSIGNED NULL,
  `period` TINYINT NULL,
  `major_cat` INT UNSIGNED NULL,
  `major_cat2` INT UNSIGNED NULL,
  `major_cat3` INT UNSIGNED NULL,
  `industry` INT UNSIGNED NULL,
  `worktime` VARCHAR(20) NULL,
  `role_status` TINYINT UNSIGNED NULL,
  `s2` TINYINT NULL,
  `s3` TINYINT NULL,
  `addr_no` INT UNSIGNED NULL,
  `s9` TINYINT NULL,
  `need_emp` INT UNSIGNED NULL,
  `need_emp1` INT UNSIGNED NULL,
  `startby` TINYINT UNSIGNED NULL,
  `exp_jobcat1` INT UNSIGNED NULL,
  `exp_jobcat2` INT UNSIGNED NULL,
  `exp_jobcat3` INT UNSIGNED NULL,
  `description` TEXT NULL,
  `others` TEXT NULL)
ENGINE = InnoDB;
