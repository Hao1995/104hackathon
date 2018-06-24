CREATE DATABASE IF NOT EXISTS `docker-example`;
USE `docker-example`;

CREATE TABLE IF NOT EXISTS `users` ( 
    `id` INT NOT NULL AUTO_INCREMENT, 
    name VARCHAR(20), 
    message VARCHAR(200), 
    PRIMARY KEY (`id`)
);