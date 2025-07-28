# E-Commerce Project Study Case: Golang Backend Web Development

This project aims to develop a RESTful API for an Online Food Ordering website. It serves as a study case to test abilities in building a backend application using Golang, GORM, and JWT, ensuring that CRUD concepts, authentication, and database management are well-implemented.

## Project Overview
You are a mid-to-senior backend engineer tasked with this individual project. All guidelines provided must be followed to demonstrate proficiency in modern Go backend development practices.

## Objectives
1. Design and perform database queries for data management
2. Integrate ORM (GORM) to simplify database manipulation and schema migration
3. Understand RESTful API for database interaction
4. Connect backend API to perform CRUD operations
5. Create endpoints for creating, reading, updating, and deleting data
6. Structure routes and serve API using the Gin framework in Go
7. Design and implement user authentication using JSON Web Token (JWT)
8. Implement file upload/download for product images
9. Structure the project following best practices for backend applications using Golang

## Key Technologies
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: MySQL
- **Authentication**: JWT (JSON Web Token)
- **Language**: Go (Golang)

## Project Structure
The project follows a clean, layered structure, separating concerns (controller, service, repository, model):

```
root/
├── cmd/
│   └── server/
├── contract/
├── database/
├── docs/
├── internal/
│   ├── api/
│   │   ├── dto/
│   │   ├── handler/
│   │   └── router/
│   ├── config/
│   ├── contract/
│   ├── database/
│   │   └── seeder/
│   ├── entities/
│   ├── exception/
│   ├── logger/
│   ├── middleware/
│   ├── repository/
│   ├── service/
│   └── utils/
├── mcp/
└── pkg/
    ├── gin_helper/
    ├── jwt/
    ├── role_validator/
    └── web_response/
```

## API Endpoints

### User Authentication
- `POST /register` - User registration
- `POST /login` - User login

### Menu Management (Admin Only)
- `GET /menus` - View menu list
- `POST /menus` - Add new menu
- `PUT /menus/:id` - Update menu
- `DELETE /menus/:id` - Delete menu

### Cart Management
- `GET /cart` - View cart contents
- `POST /cart` - Add item to cart
- `PUT /cart/:id` - Change item quantity
- `DELETE /cart/:id` - Remove item from cart

### Order Management
- `POST /checkout` - Checkout cart
- `GET /orders` - View order history

### Reports (Admin Only)
- `GET /reports/sales` - Sales report

## Validation Rules
- Validate user input such as unique email, minimum 8-character password, and non-negative stock
- Use middleware for authentication and authorization
- Implement comprehensive input validation for all endpoints
- Ensure data integrity through proper validation layers

## Main Features

### 1. User Management
- User registration and login
- Authentication using JWT
- Role-based authorization for Admin and Customer
- Secure password hashing and validation

### 2. Menu Management (Admin)
- Add, update, delete, and view menu list
- Each menu has attributes: id, name, description, price, category, stock
- Stock validation and management
- Optional file upload for product images

### 3. Cart Management (Customer)
- Add items to cart
- Delete or change item quantity in cart
- View cart contents with total calculations
- Cart persistence per user

### 4. Order Checkout (Customer)
- Confirm order from cart items
- Validate food stock availability
- Save order to database with transaction integrity
- Update food stock after successful checkout
- Handle concurrent stock updates

### 5. Order History (Customer)
- View list of previous orders with details
- Order status tracking
- Order item breakdown

### 6. Sales Reports (Admin)
- View total sales per day/month
- View list of best-selling foods
- Sales analytics and insights

## Bonus Features (Optional)
- Menu search by name/category
- Pagination for menu list
- Unit testing using testify, mockgen, mockery
- Advanced filtering and sorting
- Performance optimization with caching
- Rate limiting and security enhancements

## Implementation Guidance
- Allocate time well, focusing on core features first
- Ensure code is neat and use relevant comments
- Use database with dummy data for testing
- Follow Go conventions and best practices
- Implement proper error handling and logging
- Create comprehensive API documentation

## API Documentation
Create API documentation using tools like Swagger or Postman, including:
- Request/response examples
- Error codes and messages
- Authentication requirements
- Rate limiting information

## Getting Started
1. Clone the repository
2. Set up MySQL database
3. Configure environment variables
4. Install dependencies (`go mod tidy`)
5. Run database migrations
6. Seed dummy data
7. Start the server
8. Access API documentation

## Development Standards
- Follow clean architecture principles
- Implement proper separation of concerns
- Use dependency injection where appropriate
- Write testable, maintainable code
- Implement comprehensive error handling
- Use structured logging throughout the application