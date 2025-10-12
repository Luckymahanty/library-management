# Stage 1: Build the Go binary
FROM golang:1.24-bullseye AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the project
COPY . .

# Enable CGO for SQLite
ENV CGO_ENABLED=1

# Build the Go binary
RUN go build -o library-management .

# Stage 2: Run the application
FROM debian:bullseye-slim

WORKDIR /root

# Install SQLite (for runtime access)
RUN apt-get update && apt-get install -y sqlite3 && rm -rf /var/lib/apt/lists/*

# Copy compiled binary and frontend folder
COPY --from=builder /app/library-management .
COPY --from=builder /app/frontend ./frontend

EXPOSE 8080
CMD ["./library-management"]
