# Distributed Video Hub

A distributed video streaming and sharing platform built with Next.js frontend and Go backend.

## Features

- **Video Upload**: Drag-and-drop or click-to-browse file upload
- **Video Streaming**: DASH.js powered video player with adaptive streaming
- **Video Management**: Upload, view, and delete videos
- **Dark Mode Support**: Toggle between light and dark themes
- **Responsive Design**: Works on desktop, tablet, and mobile devices

## Prerequisites

- Node.js 18+
- Go 1.21+
- Docker and Docker Compose (optional, for easy setup)

## Quick Start (Using Docker)

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd distributed-video-hub
   ```

2. **Start all services with Docker Compose**:
   ```bash
   docker-compose up -d
   ```

3. **Open your browser** and go to [http://localhost:3000](http://localhost:3000)

## Manual Setup

### Backend Setup

1. **Start the Go backend services**:
   ```bash
   # Terminal 1: Start the web server
   go run cmd/web/main.go
   
   # Terminal 2: Start the storage service
   go run cmd/storage/main.go
   
   # Terminal 3: Start the admin service
   go run cmd/admin/main.go
   ```

### Frontend Setup

1. **Navigate to the frontend directory**:
   ```bash
   cd frontend
   ```

2. **Install dependencies**:
   ```bash
   npm install
   ```

3. **Start the development server**:
   ```bash
   npm run dev
   ```

4. **Open your browser** and go to [http://localhost:3000](http://localhost:3000)

## Usage

1. **Upload a video**: Use the upload form on the home page
2. **View videos**: Browse the video grid to see all uploaded videos
3. **Play videos**: Click on any video to start streaming
4. **Delete videos**: Use the delete button on video cards or in the video player

## Project Structure

```
├── frontend/                 # Next.js frontend application
├── cmd/                     # Go service entry points
├── internal/                # Go backend services
├── proto/                   # Protocol buffer definitions
└── docker-compose.yml       # Docker setup for all services
```

## Technologies Used

- **Frontend**: Next.js 14, TypeScript, Tailwind CSS, DASH.js
- **Backend**: Go, gRPC, Protocol Buffers
- **Storage**: File-based storage system
- **Deployment**: Docker and Docker Compose

## Troubleshooting

- **Frontend not connecting to backend**: Make sure all Go services are running on their default ports
- **Videos not loading**: Check that the storage service is running and accessible
- **Upload failures**: Verify file size is under 10MB and format is supported

## Development

To run in development mode:

```bash
# Backend (in separate terminals)
go run cmd/web/main.go
go run cmd/storage/main.go
go run cmd/admin/main.go

# Frontend
cd frontend
npm run dev
```
