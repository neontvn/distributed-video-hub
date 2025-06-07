#!/bin/bash

# Create storage directories if they don't exist
mkdir -p storage/8090
mkdir -p storage/8091
mkdir -p storage/8092

# Start the storage servers in the background
go run ./cmd/storage -host localhost -port 8090 "./storage/8090" &
STORAGE1_PID=$!

go run ./cmd/storage -host localhost -port 8091 "./storage/8091" &
STORAGE2_PID=$!

go run ./cmd/storage -host localhost -port 8092 "./storage/8092" &
STORAGE3_PID=$!

# Give storage servers a moment to start up
sleep 2

# Start the web server
go run ./cmd/web \
    sqlite "./metadata.db" \
    nw "localhost:8081,localhost:8090,localhost:8091,localhost:8092" &
WEB_PID=$!


# Wait for all background processes
wait 