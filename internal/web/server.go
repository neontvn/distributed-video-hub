// Lab 7: Implement a web server

package web

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
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
	s.mux.HandleFunc("/upload", s.handleUpload)
	s.mux.HandleFunc("/videos/", s.handleVideo)
	s.mux.HandleFunc("/content/", s.handleVideoContent)
	s.mux.HandleFunc("/delete/", s.handleDelete)
	s.mux.HandleFunc("/", s.handleIndex)

	return http.Serve(lis, s.mux)
}

func (s *server) handleIndex(w http.ResponseWriter, r *http.Request) {

	videos, err := s.metadataService.List()
	if err != nil {
		http.Error(w, "Error fetching video metadata", http.StatusInternalServerError)
		log.Println("Metadata service error:", err)
		return
	}

	type VideoData struct {
		Id         string
		EscapedId  string
		UploadTime string
	}

	var videoData []VideoData
	for _, video := range videos {
		videoData = append(videoData, VideoData{
			Id:         video.Id,
			EscapedId:  url.PathEscape(video.Id),
			UploadTime: video.UploadedAt.Format("2006-01-02 15:04:05"),
		})
	}

	tmpl, err := template.New("index").Parse(indexHTML)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		log.Println("Parse error:", err)
		return
	}

	err = tmpl.Execute(w, videoData)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Println("Execute error:", err)
		return
	}
}

func generateVideoID(filename string) string {
	return fmt.Sprintf("%s", strings.TrimSuffix(filename, filepath.Ext(filename)))
}

func (s *server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(0)
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		log.Println("Error parsing form data", err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		log.Println("Error retrieving file", err)
		return
	}

	defer file.Close()

	videoID := generateVideoID(header.Filename)

	existing, err := s.metadataService.Read(videoID)
	if err != nil {
		http.Error(w, "Error checking for existing video", http.StatusInternalServerError)
		log.Println("Error checking for existing video:", err)
		return
	}
	if existing != nil {
		http.Error(w, "Video ID already exists", http.StatusConflict)
		return
	}

	filedata, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file content", http.StatusInternalServerError)
		log.Println("Error reading file content: ", err)
		return
	}

	err = s.metadataService.Create(videoID, time.Now())
	if err != nil {
		http.Error(w, "Error in saving metadata", http.StatusInternalServerError)
		log.Println("Error in saving metadata:", err)
		return
	}

	err = s.contentService.Write(videoID, header.Filename, filedata)
	if err != nil {
		http.Error(w, "Error saving file to content service", http.StatusInternalServerError)
		log.Println("Content service write error:", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *server) handleVideo(w http.ResponseWriter, r *http.Request) {
	videoId := r.URL.Path[len("/videos/"):]
	log.Println("Video ID:", videoId)

	video, err := s.metadataService.Read(videoId)
	if err != nil || video == nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	type VideoData struct {
		Id         string
		UploadedAt string
	}

	data := VideoData{
		Id:         video.Id,
		UploadedAt: video.UploadedAt.Format("2006-01-02 15:04:05"),
	}

	tmpl, err := template.New("video").Parse(videoHTML)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		log.Println("Parse error:", err)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Println("Execute error:", err)
		return
	}
}

func (s *server) handleVideoContent(w http.ResponseWriter, r *http.Request) {
	// parse /content/<videoId>/<filename>
	videoId := r.URL.Path[len("/content/"):]
	parts := strings.Split(videoId, "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid content path", http.StatusBadRequest)
		return
	}
	videoId = parts[0]
	filename := parts[1]
	log.Println("Video ID:", videoId, "Filename:", filename)

	content, err := s.contentService.Read(videoId, filename)
	if err != nil {
		http.Error(w, "Error reading video content", http.StatusInternalServerError)
		log.Println("Content service read error:", err)
		return
	}

	// Optionally set Content-Type based on filename
	if strings.HasSuffix(filename, ".mpd") {
		w.Header().Set("Content-Type", "application/dash+xml")
	} else if strings.HasSuffix(filename, ".mp4") {
		w.Header().Set("Content-Type", "video/mp4")
	} else if strings.HasSuffix(filename, ".jpg") {
		w.Header().Set("Content-Type", "image/jpeg")
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)

}

func (s *server) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	videoId := r.URL.Path[len("/delete/"):]
	if videoId == "" {
		http.Error(w, "Video ID is required", http.StatusBadRequest)
		return
	}

	// First check if video exists
	video, err := s.metadataService.Read(videoId)
	if err != nil {
		http.Error(w, "Error checking video existence", http.StatusInternalServerError)
		log.Println("Error checking video:", err)
		return
	}
	if video == nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	err = s.metadataService.Delete(videoId)
	if err != nil {
		http.Error(w, "Error deleting video metadata", http.StatusInternalServerError)
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
		http.Error(w, "Error deleting some video files", http.StatusInternalServerError)
		log.Printf("Errors deleting files for video %s: %v", videoId, deleteErrors)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
