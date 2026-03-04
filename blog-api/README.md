# Blog Posts API (Go + MySQL)

Simple CRUD REST API for `posts` using Go `net/http` + MySQL.

## Requirements
- Go 1.22+ (uses ServeMux patterns like `GET /posts/{id}`)
- MySQL 8+
- Driver: `github.com/go-sql-driver/mysql`

## Files
- `main.go` (server)
- `config.json` (MySQL DSN)
- `schema.sql` (table schema)

## Setup

1) Create `schema.sql` and run it in MySQL:

-- schema.sql
CREATE DATABASE IF NOT EXISTS blog
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

USE blog;

CREATE TABLE IF NOT EXISTS posts (
  id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  category VARCHAR(100) NOT NULL,
  tags JSON NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

2) Create `config.json` in the project root:

{
  "mysql_dsn": "root:YOUR_PASSWORD@tcp(127.0.0.1:3306)/blog?parseTime=true&charset=utf8mb4&collation=utf8mb4_0900_ai_ci&loc=Asia%2FSeoul"
}

3) Install deps / run:

go mod tidy
go run .

Server: http://localhost:8080

## Endpoints
- POST   /posts
- GET    /posts
- GET    /posts/{id}
- PUT    /posts/{id}
- DELETE /posts/{id}

## Data Model
PostInput (request):
- title (string)
- content (string)
- category (string)
- tags (string[])

Post (response):
- id (int64)
- title, content, category (string)
- tags (string[])
- createdAt, updatedAt (time)
