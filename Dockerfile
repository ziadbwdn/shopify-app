# Stage 1: Build the application
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application
# -ldflags="-w -s" strips debug information, reducing binary size
# CGO_ENABLED=0 creates a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /server cmd/server/main.go

# Stage 2: Create the final, minimal image
FROM alpine:latest

# Set the working directory
WORKDIR /

# Copy the built binary from the builder stage
COPY --from=builder /server /server

# Copy the .env file. It's better to manage secrets with Docker secrets or environment variables in production,
# but for local development, this is convenient.
COPY .env .

# Expose the port the application runs on
EXPOSE 8080

# Command to run the application
CMD ["/server"]
