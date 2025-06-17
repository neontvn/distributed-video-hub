[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/cdomgBmc)

# TritonTube Next.js Frontend

This is a Next.js frontend that recreates the exact same UI as the original TritonTube Go application. It communicates with the Go backend API to provide a modern, responsive video streaming interface.

## Features

- **Exact UI Recreation**: 100% pixel-perfect recreation of the original TritonTube interface
- **Dark Mode Support**: Toggle between light and dark themes
- **Video Upload**: Drag-and-drop or click-to-browse file upload with 10MB limit
- **Video Grid**: Responsive grid layout showing all uploaded videos
- **Video Player**: DASH.js powered video player with adaptive streaming
- **Delete Functionality**: Delete videos with confirmation dialogs
- **Responsive Design**: Works on desktop, tablet, and mobile devices

## Prerequisites

- Node.js 18+ 
- The Go backend server running on `localhost:8080`

## Installation

1. Install dependencies:
```bash
npm install
```

2. Make sure the Go backend is running on `localhost:8080`

3. Start the development server:
```bash
npm run dev
```

4. Open [http://localhost:3000](http://localhost:3000) in your browser

## API Endpoints

The Next.js app communicates with the Go backend through these API routes:

- `GET /api/videos` - List all videos
- `GET /api/videos/[id]` - Get specific video details
- `POST /api/upload` - Upload a new video
- `DELETE /api/delete/[id]` - Delete a video
- `GET /api/content/[videoId]/[filename]` - Serve video content files

## Project Structure

```
src/
├── app/
│   ├── api/                    # API routes that proxy to Go backend
│   ├── videos/[id]/           # Individual video page
│   ├── layout.tsx             # Root layout with dark mode setup
│   └── page.tsx               # Home page with upload and video grid
├── components/
│   ├── Navigation.tsx         # Navigation bar with dark mode toggle
│   ├── UploadForm.tsx         # File upload form
│   ├── VideoCard.tsx          # Individual video card component
│   └── VideoPlayer.tsx        # DASH.js video player
└── globals.css                # Global styles
```

## Technologies Used

- **Next.js 14** - React framework with App Router
- **TypeScript** - Type safety
- **Tailwind CSS** - Utility-first CSS framework
- **DASH.js** - Adaptive streaming video player
- **React Hooks** - State management

## UI Components

### Navigation
- Fixed top navigation bar
- TritonTube branding
- Dark mode toggle button
- Back button on video pages

### Upload Form
- Drag-and-drop file upload area
- File size validation (10MB limit)
- Upload progress indicator
- Error handling

### Video Grid
- Responsive grid layout (1-4 columns based on screen size)
- Video thumbnails with fallback icons
- Hover effects with delete buttons
- Empty state when no videos

### Video Player
- Full-width video player
- DASH.js adaptive streaming
- Video metadata display
- Delete functionality

## Dark Mode

The application supports system preference-based dark mode with manual toggle. The dark mode uses YouTube-inspired colors:

- Background: `#0F0F0F`
- Sidebar: `#212121`
- Hover: `#272727`

## Development

To run in development mode:

```bash
npm run dev
```

To build for production:

```bash
npm run build
npm start
```

## Backend Integration

This frontend is designed to work with the Go backend that provides:

- Video metadata storage
- Video content storage
- DASH manifest generation
- Thumbnail generation

Make sure the Go backend is running on `localhost:8080` before starting the Next.js application.
