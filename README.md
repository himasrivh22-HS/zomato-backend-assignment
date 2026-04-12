# Zomato Backend Assignment

## Overview
This project is a backend system built using Go and PostgreSQL. It supports user registration and login with secure authentication using JWT.

## Features
- User Registration API
- User Login API
- Password hashing using bcrypt
- JWT-based authentication
- PostgreSQL database integration

## Tech Stack
- Go (Golang)
- PostgreSQL
- Docker
- JWT

## Setup Instructions

1. Start database:
   docker-compose up -d

2. Run server:
   go run cmd/server/main.go

## API Endpoints

### Register
POST /auth/register

### Login
POST /auth/login

## Notes
- Passwords are securely hashed using bcrypt
- JWT tokens are generated for authentication