### **RULE.md**

# Golden Rules for E-Commerce Golang Backend Development

This document defines the non-negotiable rules and architectural principles for this project. All development must adhere strictly to these guidelines to ensure a clean, maintainable, and testable codebase. Deviation is not permitted.

## Rule #1: The Unidirectional Development Flow

All feature development **must** follow a strict bottom-up, layer-by-layer sequence. Do not write any API handler code until the underlying service, repository, and entity layers are complete and functional.

The mandatory order of implementation is:

1.  **Skeletal Layer (Entities, Database, any other core layer):**
    *   Define the data model in `internal/entities`. Mainly this includes any kind of GORM tags
    *   Define and set up error handler package in `internal/exception/`.
    *   Define and set up utility function in `internal/utils/`.
    *   Set up the corresponding database migration in `internal/database/`.
    *   This is the absolute first step. The data structure dictates everything else.

2.  **Contract Layer (Interfaces):**
    *   Define the necessary `Repository` and `Service` interfaces in `internal/contract`.
    *   These contracts define *what* each layer does, not *how*. This enables dependency injection and parallel development if needed.

3.  **Data Access Layer (Repository):**
    *   Implement the repository interfaces defined in the contract layer within `internal/repository`.
    *   **Responsibility:** This layer's **only** function is to perform direct database operations (CRUD). It contains GORM queries.
    *   **Strictly Prohibited:** No business logic, no data manipulation (other than what's required by the DB), no calling other repositories or services.

4.  **Business Logic Layer (Service):**
    *   Implement the service interfaces defined in the contract layer within `internal/service`.
    *   **Responsibility:** This layer contains all business logic. It orchestrates calls to one or more repositories, validates data according to business rules (e.g., checking stock, calculating totals), and handles complex operations.
    *   It depends on repository *interfaces*, not concrete implementations, which are injected into it.

5.  **API/Presentation Layer (DTO, Handler, Router):**
    *   This is the **final** step.
    *   **DTOs (`internal/api/dto`):** Define Data Transfer Objects for API request bodies and responses. Use validation tags here for incoming request data. Never expose database entities directly to the client.
    *   **Handler (`internal/api/handler`):** The handler's role is to be a thin layer that:
        1.  Parses and validates the incoming HTTP request into a DTO.
        2.  Calls a single method on the appropriate service.
        3.  Translates the service's response into a response DTO and sends it back to the client using the `pkg/web_response` utility.
    *   **Router (`internal/api/router`):** Define the API endpoints, connect them to their respective handlers, and apply necessary middleware (auth, logging).

6.  **Helper/Utility Packages (`pkg/`):**
    *   Develop reusable, generic packages like `jwt`, `role_validator`, etc., as needed. These packages must have no dependencies on the `internal` directory.

## Rule #2: Strict Separation of Concerns (SoC)

Each layer has a single, well-defined responsibility. Do not mix responsibilities.

-   **Handler:** Manages HTTP only.
-   **Service:** Manages business logic only.
-   **Repository:** Manages database interaction only.
-   **Entity:** Represents a database table schema only.
-   **DTO:** Represents API request/response data only.

## Rule #3: Dependency Injection is Mandatory

-   Dependencies must flow in one direction: `Handler` → `Service` → `Repository`.
-   Always depend on abstractions (interfaces from `internal/contract`), not on concrete implementations. This is crucial for testability and decoupling.

## Rule #4: Refactoring Guidelines

-   **No Circular Refactoring:** Do not refactor code in a way that violates the unidirectional flow. A change in the `handler` should never necessitate a change in the `repository` interface.
-   **Refactor After Testing:** Only refactor code that is complete and has passing tests. Verify with tests again after refactoring to ensure no regressions were introduced.
-   **Refactoring is for Improvement, Not Style:** Refactor only to improve clarity, performance, or maintainability, or to remove technical debt.

## Rule #5: Error Handling and Logging

-   Use the centralized logger from `internal/logger`.
-   Use the centralized middleware logic from `internal/middleware`.
-   Use the centralized error handler package from `internal/exception`.
-   Errors must be handled gracefully at each layer. A repository error should be wrapped and returned to the service, which may in turn wrap and return it to the handler.
-   The handler is responsible for converting final errors into the appropriate HTTP status code and standardized error response using `pkg/web_response`.

Adherence to these rules is paramount for the success and scalability of the project.