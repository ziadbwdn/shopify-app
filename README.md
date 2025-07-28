# E-Commerce Go Backend

This repository contains the backend for a comprehensive E-Commerce platform, built with Go. It follows a clean, layered architecture designed for maintainability, testability, and scalability.

## Features

-   **User Management**: Secure user registration and login with JWT-based authentication.
-   **Role-Based Access Control (RBAC)**: Distinction between `admin` and `customer` roles.
-   **Menu Management**: Admins can create, update, and delete menu items.
-   **Shopping Cart**: Users can add, update, remove, and clear items in their cart.
-   **Order Processing**: Users can checkout their cart to create an order. Admins can manage order statuses.
-   **Reporting**: Admins can generate sales and analytics reports.

## Architecture

The project adheres to a strict, unidirectional, layered architecture to ensure a clean separation of concerns.

1.  **Skeletal Layer**: Core data structures (`/internal/entities`) and database migrations.
2.  **Contract Layer**: Go interfaces (`/internal/contract`) defining the behavior for services and repositories.
3.  **Data Access Layer (Repository)**: Concrete implementation (`/internal/repository`) of data access logic using GORM. This is the only layer that interacts directly with the database.
4.  **Business Logic Layer (Service)**: Concrete implementation (`/internal/service`) of business logic, orchestrating data from repositories.
5.  **API/Presentation Layer**: The client-facing layer, consisting of:
    -   **DTOs** (`/internal/api/dto`): Data Transfer Objects for API requests/responses.
    -   **Handlers** (`/internal/api/handler`): Thin layer for handling HTTP requests, validating input, and calling services.
    -   **Router** (`/internal/api/router`): Defines API endpoints and connects them to handlers.
6.  **Helper/Utility Packages**: Reusable packages for JWT, password validation, etc. (`/pkg`).

## Getting Started

### Prerequisites

-   Go (version 1.18 or higher)
-   Docker and Docker Compose
-   A running MySQL instance (if not using Docker)

### 1. Environment Configuration

Create a `.env` file in the root of the project and populate it with the following variables:

```env
# Application Port
PORT=8080

# Database Credentials
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_database_password
DB_NAME=shopify_app

# JWT Secret
JWT_SECRET=your_super_secret_jwt_key
```

### 2. Running with Docker (Recommended)

The simplest way to get the application and a database running is with Docker Compose.

1.  **Build and Start Containers:**
    ```bash
    docker-compose up --build
    ```

2.  The application will be available at `http://localhost:8080`.

### 3. Running Manually

1.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

2.  **Run the Server:**
    ```bash
    go run cmd/server/main.go
    ```

The server will start on the port specified in your `.env` file.
