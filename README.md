# Zomato Backend Assignment

##  Overview
This project is a backend system built using Go (Golang) and PostgreSQL.  
It supports user authentication, project management, and task management with secure JWT-based authorization.

---

##  Features
- User Registration & Login
- Password hashing using bcrypt
- JWT-based authentication
- Protected APIs using middleware
- Project management APIs
- Task management APIs
- PostgreSQL database integration
- Docker support

---

## Tech Stack
- Go (Golang)
- PostgreSQL
- Docker
- JWT

---

## Setup Instructions

### 1. Start Database

### 2. Run Server

---

## Authentication
All endpoints (except register/login) require:

---

## API Endpoints

### Auth
- POST /auth/register  
- POST /auth/login  

### Projects
- POST /projects  

### Tasks
- POST /tasks  
- GET /tasks/list  
- POST /tasks/update?id=  
- POST /tasks/delete?id=  

---

##  Testing
You can test APIs using:
- Postman
- PowerShell (Invoke-RestMethod)

---

## Notes
- Passwords are securely hashed using bcrypt  
- JWT tokens are used for authentication  
- Tasks are linked to projects using project_id  

---

##  Repository
https://github.com/himasrivh22-HS/zomato-backend-assignment
