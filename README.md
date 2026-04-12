# TaskFlow Backend (Go + PostgreSQL)

## Overview

This project is a backend service for a task management system built using Go (Golang) and PostgreSQL.
It allows users to register, log in, create projects, and manage tasks within those projects.

---

## Tech Stack

* Go (Golang)
* PostgreSQL
* Docker
* JWT Authentication
* Chi Router

---

## Features

* User Registration & Login
* Password hashing using bcrypt
* JWT-based authentication
* Protected APIs using middleware
* Project CRUD operations
* Task CRUD operations
* Tasks linked to projects
* RESTful API design

---

## Running Locally

```bash
git clone https://github.com/himasrivh22-HS/zomato-backend-assignment
cd zomato-backend-assignment
go run cmd/server/main.go
```

Server runs at:
http://localhost:8080

---

## Authentication

All endpoints (except register/login) require:

Authorization: Bearer <your_token>

---

## API Endpoints

### Auth

* POST /auth/register
* POST /auth/login

### Projects

* GET /projects
* POST /projects
* GET /projects/{id}
* PATCH /projects/{id}
* DELETE /projects/{id}

### Tasks

* GET /projects/{id}/tasks
* POST /projects/{id}/tasks
* PATCH /tasks/{id}
* DELETE /tasks/{id}

---

## Example Request

### Create Project

```json
{
  "name": "Test Project",
  "description": "Demo project",
  "owner_id": "user-id"
}
```

---

## What I Would Improve

* Add pagination for large data
* Add filtering (status, assignee)
* Add input validation
* Add unit and integration tests
* Improve error handling

---

## Repository

https://github.com/himasrivh22-HS/zomato-backend-assignment
