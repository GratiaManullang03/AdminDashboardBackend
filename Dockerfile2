# Build stage
FROM golang:1.24.1 AS builder
WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build statically linked binary in the BUILD stage
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o main ./cmd/api/main.go

# Run stage (Alpine is smaller and works well with statically linked Go binaries)
FROM alpine:latest
WORKDIR /app

# Just copy the compiled binary from the build stage
COPY --from=builder /app/main .

# Install CA certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

EXPOSE 3000
CMD ["./main"]