# KongDB Dockerfile
# Multi-stage build for optimized production image

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build KongDB binaries
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kongdb ./cmd/kongdb
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kongdb-node ./cmd/kongdb-node

# Run tests
RUN go test -v ./internal/...

# Production stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S kongdb && \
    adduser -u 1001 -S kongdb -G kongdb

# Set working directory
WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder /app/kongdb /app/kongdb
COPY --from=builder /app/kongdb-node /app/kongdb-node

# Create data directory
RUN mkdir -p /app/data && \
    chown -R kongdb:kongdb /app

# Switch to non-root user
USER kongdb

# Expose default port
EXPOSE 8080

# Set environment variables
ENV KONGDB_DATA_DIR=/app/data
ENV KONGDB_LOG_LEVEL=info

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Default command
CMD ["./kongdb"] 