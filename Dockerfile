# Use Go 1.24 with Alpine
FROM golang:1.24-alpine

# Install compiler and libc for CGO + SQLite
RUN apk add --no-cache gcc g++ musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Enable CGO (required for go-sqlite3)
ENV CGO_ENABLED=1

RUN go build -o library-app .

EXPOSE 8080

CMD ["./library-app"]

