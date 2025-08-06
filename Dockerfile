# Build stage
FROM golang:1.24-alpine AS builder

# Install required build tools
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main -ldflags="-w -s" .

# Final stage
FROM alpine:3.19

# Add non root user
RUN adduser -D -g '' appuser

# Create process directory with correct permissions
RUN mkdir -p /app/process && chown -R appuser:appuser /app

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/main .

# Set ownership
RUN chown -R appuser:appuser /app

# Use non root user
USER appuser

# Expose port
EXPOSE 8080


# Run the application
CMD ["./main"]
