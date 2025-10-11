# Stage 1: Build the Go binary with CGO
FROM golang:1.24-alpine AS builder

# Install gcc and dependencies for CGO/SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary with CGO enabled
ENV CGO_ENABLED=1
RUN go build -o library-management .

# Stage 2: Run the application
FROM alpine:latest
WORKDIR /root

# Install runtime dependencies for CGO/SQLite
RUN apk add --no-cache libc6-compat sqlite

COPY --from=builder /app/library-management .

EXPOSE 8080
CMD ["./library-management"]
