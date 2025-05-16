Event Management API (Go + MySQL)
A simple RESTful API built with Go's net/http and database/sql packages for managing events.
 Tech Stack
Language: Go (Golang)
Database: MySQL
Packages Used:
github.com/go-sql-driver/mysql
encoding/json
net/http
database/sql

 Database Setup
Create a MySQL database called event_managment (or modify the DSN in main.go).
Run the following SQL to create the events table:
CREATE TABLE events (
  id INT PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  location VARCHAR(255),
  start_time DATETIME,
  end_time DATETIME,
  created_by VARCHAR(255),
  created_at DATETIME,
  updated_at DATETIME
);


Getting Started
Clone the repo and open the project directory.
Update your MySQL connection in main.go:
dsn := "root:root@tcp(127.0.0.1:3306)/event_management?parseTime=true"
Run the application:
go run main.go
Server will run at:
http://localhost:8080
