# Build stage
FROM golang:1.24.1 AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/api/main.go

# Run stage
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/main .

# Install CA certificates for HTTPS connections
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o main ./cmd/api/main.go

EXPOSE 8080
CMD ["./main"]