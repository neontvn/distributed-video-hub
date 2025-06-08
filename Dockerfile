FROM golang:1.24-alpine

# Install required dependencies
RUN apk add --no-cache \
    ca-certificates \
    sqlite-libs \
    build-base \
    gcc \
    musl-dev \
    sqlite-dev

WORKDIR /app

# Copy all code
COPY . .

# Let Go figure out the correct dependencies
RUN go mod tidy

# Create storage directories
RUN mkdir -p /app/storage/8090 \
    /app/storage/8091 \
    /app/storage/8092 \
    /app/data

# Build the applications
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/web-server ./cmd/web && \
    CGO_ENABLED=1 GOOS=linux go build -o /app/storage-server ./cmd/storage

# Default command (will be overridden by docker-compose)
CMD ["/app/web-server"] 