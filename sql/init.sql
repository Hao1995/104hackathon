CREATE DATABASE IF NOT EXISTS `104hackathon`;
USE `104hackathon`;

CREATE TABLE IF NOT EXISTS `users` ( 
    `id` INT NOT NULL AUTO_INCREMENT, 
    name VARCHAR(20), 
    message VARCHAR(200), 
    PRIMARY KEY (`id`)
);