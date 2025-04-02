# Build stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/api/main.go

# Run stage
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/main .

RUN go mod download && go build -o main ./cmd/api/main.go

# Install CA certificates for HTTPS connections
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

EXPOSE 8080
CMD ["./main"]