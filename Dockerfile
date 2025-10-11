# Stage 1: Build the Go binary with CGO enabled
FROM golang:1.24-alpine AS builder

# Install gcc and dependencies for CGO/SQLite
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Enable CGO when building
ENV CGO_ENABLED=1
RUN go build -o library-management .

# Stage 2: Run the application
FROM debian:bullseye-slim
WORKDIR /root/
COPY --from=builder /app/library-management .
EXPOSE 8080
CMD ["./library-management"]


