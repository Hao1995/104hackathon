ALTER TABLE `jobs`
ADD FOREIGN KEY (`custno`) REFERENCES `companies`(`custno`) 
	ON UPDATE CASCADE;

-- // Lack data ... Do not make the foreign key.
-- ALTER TABLE `jobs`
-- ADD FOREIGN KEY (`jobcat1`) REFERENCES `job_categories`(`id`) 
-- 	ON UPDATE CASCADE;
-- ALTER TABLE `jobs`
-- ADD FOREIGN KEY (`jobcat2`) REFERENCES `job_categories`(`id`) 
-- 	ON UPDATE CASCADE;
-- ALTER TABLE `jobs`
-- ADD FOREIGN KEY (`jobcat3`) REFERENCES `job_categories`(`id`) 
-- 	ON UPDATE CASCADE;

-- ALTER TABLE `jobs`
-- ADD FOREIGN KEY (`exp_jobcat1`) REFERENCES `job_categories`(`id`) 
-- 	ON UPDATE CASCADE;
-- ALTER TABLE `jobs`
-- ADD FOREIGN KEY (`exp_jobcat2`) REFERENCES `job_categories`(`id`) 
-- 	ON UPDATE CASCADE;
-- ALTER TABLE `jobs`
-- ADD FOREIGN KEY (`exp_jobcat3`) REFERENCES `job_categories`(`id`) 
-- 	ON UPDATE CASCADE;
    
ALTER TABLE `jobs`
ADD FOREIGN KEY (`major_cat1`) REFERENCES `departments`(`id`) 
	ON UPDATE CASCADE;
ALTER TABLE `jobs`
ADD FOREIGN KEY (`major_cat2`) REFERENCES `departments`(`id`) 
	ON UPDATE CASCADE;
ALTER TABLE `jobs`
ADD FOREIGN KEY (`major_cat3`) REFERENCES `departments`(`id`) 
	ON UPDATE CASCADE;

ALTER TABLE `jobs`
ADD FOREIGN KEY (`industry`) REFERENCES `industries`(`id`) 
	ON UPDATE CASCADE;

INSERT INTO `districts`(`id`,`name`,`desc`, `hide`) VALUES(6005003013, '美國塞班島', '', '否');
ALTER TABLE `jobs`
ADD FOREIGN KEY (`addr_no`) REFERENCES `districts`(`id`) 
	ON UPDATE CASCADE;

ALTER TABLE `train_action`
ADD FOREIGN KEY (`jobno`) REFERENCES `jobs`(`jobno`) 
	ON UPDATE CASCADE;
ALTER TABLE `train_click`
ADD FOREIGN KEY (`jobno`) REFERENCES `jobs`(`jobno`) 
	ON UPDATE CASCADE;