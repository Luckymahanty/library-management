# Library Management System (Go + DevOps)
Name - LibraryHub

A simple library management backend built in **Go**, with CRUD operations for books.  
Future integration: frontend (HTML/CSS/JS), Docker, AWS (EC2,S3 , VPC), Terraform, and CI/CD pipelines.

## 🚀 Features
- Add books
- View books
- Delete books
- In-memory database
- Brow Book
- Return Book

## 🔧 Run locally
```bash
go run main.go
```
## Run In Docker
```
 docker compose down -v
 docker compose build --no-cache
 docker compose up
 ```
 ## AWS Resourses
 create an ec2 instance than copy the backend an whole project to it.
 And create and S3 whrere the frontend store .
 I dont use RDS because of my free tire limit reach you can use it .
