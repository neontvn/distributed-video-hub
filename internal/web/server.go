// Lab 7: Implement a web server

package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type server struct {
	Addr string
	Port int

	metadataService VideoMetadataService
	contentService  VideoContentService

	mux *http.ServeMux
}

func NewServer(
	metadataService VideoMetadataService,
	contentService VideoContentService,
) *server {
	return &server{
		metadataService: metadataService,
		contentService:  contentService,
	}
}

func (s *server) Start(lis net.Listener) error {
	s.mux = http.NewServeMux()

	// API endpoints for Next.js
	s.mux.HandleFunc("/api/videos", s.handleAPIVideos)
	s.mux.HandleFunc("/api/videos/", s.handleAPIVideo)
	s.mux.HandleFunc("/api/upload", s.handleAPIUpload)
	s.mux.HandleFunc("/api/delete/", s.handleAPIDelete)
	s.mux.HandleFunc("/api/content/", s.handleAPIVideoContent)

	// CORS middleware for all API routes
	s.mux.HandleFunc("/", s.handleCORS)

	return http.Serve(lis, s.mux)
}

// API Response structures
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type VideoAPIResponse struct {
	Id         string `json:"id"`
	UploadedAt string `json:"uploadedAt"`
}

type UploadAPIResponse struct {
	VideoId string `json:"videoId"`
}

// Helper function to send JSON responses
func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Helper function to send error responses
func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	sendJSONResponse(w, statusCode, APIResponse{
		Success: false,
		Error:   message,
	})
}

// CORS handler for all routes
func (s *server) handleCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// If not OPTIONS, return 404 for unknown routes
	sendErrorResponse(w, http.StatusNotFound, "Endpoint not found")
}

// API endpoint: GET /api/videos - List all videos
func (s *server) handleAPIVideos(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	videos, err := s.metadataService.List()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Error fetching video metadata")
		log.Println("Metadata service error:", err)
		return
	}

	var videoResponses []VideoAPIResponse
	for _, video := range videos {
		videoResponses = append(videoResponses, VideoAPIResponse{
			Id:         video.Id,
			UploadedAt: video.UploadedAt.Format("2006-01-02 15:04:05"),
		})
	}

	sendJSONResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    videoResponses,
	})
}

// API endpoint: GET /api/videos/{videoId} - Get specific video
func (s *server) handleAPIVideo(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	videoId := r.URL.Path[len("/api/videos/"):]
	if videoId == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Video ID is required")
		return
	}

	video, err := s.metadataService.Read(videoId)
	if err != nil || video == nil {
		sendErrorResponse(w, http.StatusNotFound, "Video not found")
		return
	}

	videoResponse := VideoAPIResponse{
		Id:         video.Id,
		UploadedAt: video.UploadedAt.Format("2006-01-02 15:04:05"),
	}

	sendJSONResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    videoResponse,
	})
}

// API endpoint: POST /api/upload - Upload video
func (s *server) handleAPIUpload(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	err := r.ParseMultipartForm(0)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Error parsing form data")
		log.Println("Error parsing form data", err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Error retrieving file")
		log.Println("Error retrieving file", err)
		return
	}
	defer file.Close()

	videoID := generateVideoID(header.Filename)

	existing, err := s.metadataService.Read(videoID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Error checking for existing video")
		log.Println("Error checking for existing video:", err)
		return
	}
	if existing != nil {
		sendErrorResponse(w, http.StatusConflict, "Video ID already exists")
		return
	}

	filedata, err := io.ReadAll(file)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Error reading file content")
		log.Println("Error reading file content: ", err)
		return
	}

	err = s.metadataService.Create(videoID, time.Now())
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Error in saving metadata")
		log.Println("Error in saving metadata:", err)
		return
	}

	err = s.contentService.Write(videoID, header.Filename, filedata)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Error saving file to content service")
		log.Println("Content service write error:", err)
		return
	}

	sendJSONResponse(w, http.StatusCreated, APIResponse{
		Success: true,
		Data: UploadAPIResponse{
			VideoId: videoID,
		},
	})
}

// API endpoint: DELETE /api/delete/{videoId} - Delete video
func (s *server) handleAPIDelete(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodDelete {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	videoId := r.URL.Path[len("/api/delete/"):]
	if videoId == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Video ID is required")
		return
	}

	// First check if video exists
	video, err := s.metadataService.Read(videoId)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Error checking video existence")
		log.Println("Error checking video:", err)
		return
	}
	if video == nil {
		sendErrorResponse(w, http.StatusNotFound, "Video not found")
		return
	}

	err = s.metadataService.Delete(videoId)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Error deleting video metadata")
		log.Println("Error deleting metadata:", err)
		return
	}

	files, err := s.contentService.ListFiles(videoId)
	if err != nil {
		log.Printf("Warning: Could not list files for video %s: %v", videoId, err)
		files = []string{"manifest.mpd", "thumbnail.jpg"}
		files = append(files, "init-0.m4s", "init-1.m4s")
		for stream := 0; stream <= 1; stream++ {
			for i := 1; i <= 100; i++ {
				files = append(files, fmt.Sprintf("chunk-%d-%05d.m4s", stream, i))
			}
		}
	}

	// Delete each file
	var deleteErrors []error
	for _, filename := range files {
		err := s.contentService.Delete(videoId, filename)
		if err != nil {
			if !strings.Contains(err.Error(), "no such file") {
				log.Printf("Error deleting file %s: %v", filename, err)
				deleteErrors = append(deleteErrors, err)
			}
		}
	}

	if len(deleteErrors) > 0 {
		sendErrorResponse(w, http.StatusInternalServerError, "Error deleting some video files")
		log.Printf("Errors deleting files for video %s: %v", videoId, deleteErrors)
		return
	}

	sendJSONResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Video deleted successfully"},
	})
}

// API endpoint: GET /api/content/{videoId}/{filename} - Serve video content
func (s *server) handleAPIVideoContent(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// parse /api/content/<videoId>/<filename>
	contentPath := r.URL.Path[len("/api/content/"):]
	parts := strings.Split(contentPath, "/")
	if len(parts) != 2 {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid content path")
		return
	}
	videoId := parts[0]
	filename := parts[1]
	log.Println("Video ID:", videoId, "Filename:", filename)

	content, err := s.contentService.Read(videoId, filename)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Error reading video content")
		log.Println("Content service read error:", err)
		return
	}

	// Set appropriate Content-Type based on filename
	if strings.HasSuffix(filename, ".mpd") {
		w.Header().Set("Content-Type", "application/dash+xml")
	} else if strings.HasSuffix(filename, ".mp4") {
		w.Header().Set("Content-Type", "video/mp4")
	} else if strings.HasSuffix(filename, ".jpg") {
		w.Header().Set("Content-Type", "image/jpeg")
	} else if strings.HasSuffix(filename, ".m4s") {
		w.Header().Set("Content-Type", "video/mp4")
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

func generateVideoID(filename string) string {
	return fmt.Sprintf("%s", strings.TrimSuffix(filename, filepath.Ext(filename)))
}
