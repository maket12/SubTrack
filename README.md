# 📦 Subtrack --- Online Subscription Tracking Service

**Subtrack** is a REST service designed to aggregate and manage user
online subscriptions.\
It provides full CRUDL operations and calculates the total cost of
subscriptions for any selected period.\
This project was implemented as part of the Junior Golang Developer test
assignment.

## ✨ Features

-   Create, read, update and delete subscription records
-   List subscriptions with filtering options
-   Calculate total subscription cost over a specified period
-   Filtering by `user_id` and `service_name`
-   Support for subscriptions with no end date
-   Input validation (UUID, dates, price)
-   Swagger API documentation
-   Structured request & error logging
-   Configuration via `.env`
-   Fully containerized using **Docker Compose**

## 🛠️ Tech Stack

-   **Go** (Gin, sqlx, uuid, slog)
-   **PostgreSQL**
-   **Swaggo** (Swagger documentation)
-   **Docker / Docker Compose**

## 📁 Project Structure

    /cmd/subtrack/main.go
    /docs/                  — Swagger spec
    /internal/
        adapter/
            in/
                http/       — HTTP handlers
            out/
                db/         — database layer
        app/                — business logic
        config/
        domain/             — domain entities/ports
    migrations/             — SQL migrations
    docker-compose.yml
    Dockerfile

## 🚀 Run the Service
``` bash
docker-compose up --build
```

Swagger UI available at:

    http://localhost:8080/swagger/index.html

## 📘 API Overview

-   Full CRUDL for subscription entities
-   Endpoint for calculating total subscription cost for a specified
    period
-   Filtering by user and service name
-   Monthly-based cost calculation logic

## 📄 Test Assignment --- Effective Mobile

This project implements a REST microservice following the requirements
of the
**Effective Mobile --- Junior Golang Developer test assignment**.
