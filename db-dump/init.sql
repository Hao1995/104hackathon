SELECT DATABASE();
USE users;
CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(20), 
    message VARCHAR(200), 
    PRIMARY KEY (id)
);