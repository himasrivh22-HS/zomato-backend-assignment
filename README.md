# TaskFlow Backend API

A minimal yet production-style backend system for managing users, projects, and tasks with authentication.

---

## Overview

TaskFlow is a RESTful backend API that allows users to:

* Register and authenticate using JWT
* Create and manage projects
* Add and manage tasks within projects
* Assign tasks and track their status

Built with **Go**, **PostgreSQL**, and **Chi router**, focusing on clean architecture and real-world backend practices.

---

## Tech Stack

* **Language:** Go (Golang)
* **Router:** Chi
* **Database:** PostgreSQL
* **Authentication:** JWT + bcrypt
* **Environment Management:** godotenv

---

## Architecture Decisions

* Used **Chi router** for lightweight and flexible routing
* Implemented **JWT authentication** for stateless session handling
* Used **PostgreSQL** for strong relational data consistency
* Structured project into:

  * `handler` → request handling
  * `model` → data structures
  * `middleware` → authentication logic
  * `config` → database setup
* Avoided ORM to maintain full control over SQL queries

---

## Running Locally

### 1. Clone the repository

```bash
git clone https://github.com/<your-username>/taskflow-backend
cd taskflow-backend
```

---

### 2. Setup environment variables

Create `.env` file:

```env
DB_URL=postgres://postgres:postgres@localhost:5432/taskflow?sslmode=disable
JWT_SECRET=your-secret-key
PORT=8080
```

---

### 3. Install dependencies

```bash
go mod tidy
```

---

### 4. Run the server

```bash
go run cmd/server/main.go
```

---

### 5. Server runs at:

```plaintext
http://localhost:8080
```

---

##  Authentication

### Register

```
POST /auth/register
```

### Login

```
POST /auth/login
```

Returns JWT token:

```json
{
  "token": "<jwt-token>",
  "user": {
    "id": "...",
    "name": "...",
    "email": "..."
  }
}
```

---

##  API Endpoints

### Projects

* `GET /projects` → list projects
* `POST /projects` → create project
* `GET /projects/{id}` → get project details
* `PATCH /projects/{id}` → update project
* `DELETE /projects/{id}` → delete project

---

### Tasks

* `GET /projects/{id}/tasks` → list tasks
  Supports filters:

  ```
  ?status=todo
  ?assignee=<user_id>
  ```

* `POST /projects/{id}/tasks` → create task

* `PATCH /tasks/{id}` → update task

* `DELETE /tasks/{id}` → delete task

---

##  Error Handling

* `400` → validation errors
* `401` → unauthorized
* `403` → forbidden
* `404` → not found

Example:

```json
{
  "error": "validation failed",
  "fields": {
    "email": "is required"
  }
}
```

---

## Test Credentials


Email: test@example.com
Password: password123
```

---

## Features

* Secure password hashing (bcrypt)
* JWT-based authentication
* RESTful API design
* Filtering support for tasks
* Clean and modular code structure

---

## What I’d Improve With More Time

* Add pagination for large datasets
* Implement integration tests
* Add Docker support for full environment setup
* Add logging (zap/logrus)
* Role-based access control (RBAC)

---

## Submission Notes

This project focuses on:

* clean architecture
* correct API design
* real-world backend patterns

All core requirements are implemented and tested locally.

---
