# Stage 1: Builder
FROM golang:1.24-alpine AS builder

# ติดตั้ง dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Stage 2: Runner
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary
COPY --from=builder /app/main .

# Copy proto files (if any)
COPY --from=builder /app/proto ./proto

# Copy .env file
COPY .env .

# Expose ports
EXPOSE 8081 50051

# Run
CMD ["./main"]