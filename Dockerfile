# Build stage
FROM golang:1.23.4-alpine AS builder

# Install git, air and swag
RUN apk add --no-cache git && \
    go install github.com/air-verse/air@latest

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Development stage
FROM builder AS dev

# Set working directory
WORKDIR /app

# Copy air config
COPY .air.toml .

# Run air
CMD ["air", "-c", ".air.toml"]

# Production stage
FROM alpine:latest AS prod

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/app .

# Run the application
CMD ["./app"] 