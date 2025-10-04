# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod first for caching
COPY go.mod ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 go build -o /library-management main.go

# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /

# Copy only the compiled binary (no source files!)
COPY --from=builder /library-management /library-management

EXPOSE 8080

CMD ["/library-management"]

