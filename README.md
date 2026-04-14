# TaskFlow Backend API

A minimal, production-style backend system for managing users, projects, and tasks with authentication.

---

## Overview

TaskFlow is a RESTful API that allows users to:

* Register and authenticate using JWT
* Create and manage projects
* Add and manage tasks within projects
* Assign tasks and track their status

The project is built with a focus on clean architecture, proper API design, and real-world backend practices.

---

## Tech Stack

* Language: Go (Golang)
* Router: Chi
* Database: PostgreSQL
* Authentication: JWT + bcrypt
* Environment: godotenv
* Containerization: Docker + Docker Compose

---

## Architecture Decisions

* Used Chi router for lightweight and flexible routing
* Implemented JWT authentication for stateless session handling
* Used PostgreSQL for relational data consistency
* Avoided ORM to maintain full control over SQL queries

### Project Structure

```
internal/
  ├── handler      # HTTP handlers
  ├── middleware   # Authentication middleware
  ├── model        # Data models
  ├── config       # Database configuration
cmd/
  └── server       # Application entry point
migrations/        # SQL migrations and seed data
```

### Tradeoffs

* Skipped ORM for simplicity and transparency
* Used basic validation instead of a validation library
* Focused on core requirements rather than advanced optimizations

---

## Running Locally (Docker Recommended)

### 1. Clone repository

```
git clone https://github.com/<your-username>/taskflow-backend
cd taskflow-backend
```

---

### 2. Setup environment variables

Create a `.env` file:

```
DB_URL=postgres://postgres:postgres@db:5432/taskflow?sslmode=disable
JWT_SECRET=your-secret-key
PORT=8080
```

---

### 3. Run the application

```
docker compose up --build
```

---

### 4. Access API

```
http://localhost:8080
```

---

## Migrations

Migrations are located in:

```
/migrations
```

They include:

* Schema creation
* Seed data (user, project, tasks)

If migrations are not automatically applied:

```
migrate -path ./migrations -database "$DB_URL" up
```

---

## Authentication

### Register

```
POST /auth/register
```

### Login

```
POST /auth/login
```

### Response

```
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

## API Endpoints

### Projects

* GET /projects → list projects
* POST /projects → create project
* GET /projects/{id} → project details with tasks
* PATCH /projects/{id} → update project
* DELETE /projects/{id} → delete project

---

### Tasks

* GET /projects/{id}/tasks → list tasks

Supports filtering:

```
?status=todo
?assignee=<user_id>
```

* POST /projects/{id}/tasks → create task
* PATCH /tasks/{id} → update task
* DELETE /tasks/{id} → delete task

---

## Error Handling

| Status | Meaning          |
| ------ | ---------------- |
| 400    | Validation error |
| 401    | Unauthorized     |
| 403    | Forbidden        |
| 404    | Not found        |

### Example

```
{
  "error": "validation failed",
  "fields": {
    "title": "is required"
  }
}
```

---

## Test Credentials

```
Email: test@example.com  
Password: password123
```

---

## Features

* Secure password hashing using bcrypt
* JWT-based authentication
* Role-based access control (owner and assignee)
* Task filtering support
* Structured error responses
* Clean and modular architecture

---

## What I'd Improve With More Time

* Add pagination (page, limit)
* Add integration tests
* Add structured logging (zap or logrus)
* Implement project statistics endpoint
* Improve validation using a dedicated library

---

## Submission Notes

This project focuses on:

* Clean architecture
* RESTful API design
* Proper authentication and authorization
* Production-style error handling
* Dockerized setup for easy evaluation

All core requirements are implemented and tested locally.
