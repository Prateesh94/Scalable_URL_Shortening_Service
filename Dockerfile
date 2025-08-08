# Build Stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN go build -o url-shortener .

# Final Stage
FROM alpine:latest

WORKDIR /app

# Copy the built binary
COPY --from=builder /app/url-shortener .

EXPOSE 8080

# Run the app (no config file needed, everything is env vars)
CMD ["./url-shortener"]