# syntax=docker/dockerfile:1

# Builder stage
FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
# Download dependencies first to leverage Docker cache
RUN go mod download
COPY . .
# Build the application
# Ensure cmd/main.go is the correct path to your main package
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /ginDYMall cmd/main.go

# Final stage
FROM alpine:3.16
# Add ca-certificates for HTTPS calls if needed by the application
RUN apk add --no-cache ca-certificates
WORKDIR /app
# Copy the binary from the builder stage
COPY --from=builder /ginDYMall .
# Copy the entire config directory
# The application will pick the correct config file based on APP_ENV at runtime
COPY config ./config/
# Expose the port the application runs on (should match config.yaml's HttpPort)
EXPOSE 5001
# Set default APP_ENV if needed, or expect it to be passed at runtime
# For example, to run in production by default if not specified:
# ENV APP_ENV=prod
# If APP_ENV is expected to be passed during `docker run -e APP_ENV=prod`, then no need to set ENV here.
# Or, provide a default for development:
ENV APP_ENV=dev
ENTRYPOINT ["./ginDYMall"]
