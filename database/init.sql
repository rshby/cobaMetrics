CREATE DATABASE coba_metrics;

USE coba_metrics;

CREATE TABLE accounts (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(255) NOT NULL UNIQUE ,
    username VARCHAR(255) NOT NULL UNIQUE ,
    password VARCHAR(5000) NOT NULL ,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)engine = InnoDB;