FROM golang:1.24-alpine

# Install required dependencies including ffmpeg
RUN apk add --no-cache \
    ca-certificates \
    sqlite-libs \
    build-base \
    gcc \
    musl-dev \
    sqlite-dev \
    ffmpeg

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Download Go dependencies
RUN go mod tidy

# Create necessary directories
RUN mkdir -p /app/storage/8090 \
    /app/storage/8091 \
    /app/storage/8092 \
    /app/data

# Build the binaries
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/web-server ./cmd/web && \
    CGO_ENABLED=1 GOOS=linux go build -o /app/storage-server ./cmd/storage

# Default command (can be overridden by docker-compose)
CMD ["/app/web-server"]