-To run this project you need: Go, Docker and Docker Compose, and Web Browser.



-Copy the example config in bash:
   cp config.yaml.example config.yaml

-You can change the server port or database DSN in config.yaml file. By default it is set as 8080.

-This app is using MySQL and to start MySQL Container, from the project root folder, enter this command:

docker compose up -d

-Next, we have to create database and table inside MySQL Container:

CREATE DATABASE myapp;
USE myapp;

CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    age INT NOT NULL
);

-From the root folder, run the Go application:

go run ./cmd/main.go

-Open the application in browser:

http://localhost:8080/users


If everything is OK, you should see the main UI of application.



