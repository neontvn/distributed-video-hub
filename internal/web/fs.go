// Lab 7: Implement a local filesystem video content service

package web

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// FSVideoContentService implements VideoContentService using the local filesystem.
type FSVideoContentService struct {
	BaseDir string
}

// Read implements VideoContentService.
func (f *FSVideoContentService) Read(videoId string, filename string) ([]byte, error) {
	filePath := filepath.Join(f.BaseDir, videoId, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil
}

// Write implements VideoContentService.
func (f *FSVideoContentService) Write(videoId string, filename string, data []byte) error {

	videoDir := filepath.Join(f.BaseDir, videoId)
	err := os.MkdirAll(videoDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	tempInputFile := filepath.Join(videoDir, filename)
	err = os.WriteFile(tempInputFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write input file: %w", err)
	}
	defer os.Remove(tempInputFile)

	manifestPath := filepath.Join(videoDir, "manifest.mpd")

	cmd := exec.Command(
		"ffmpeg",
		"-i", tempInputFile, // Input file
		"-c:v", "libx264", // Video codec
		"-c:a", "aac", // Audio codec
		"-bf", "1", // Max 1 B-frame
		"-keyint_min", "120", // Minimum keyframe interval
		"-g", "120", // Keyframe every 120 frames
		"-sc_threshold", "0", // Scene change threshold
		"-b:v", "3000k", // Video bitrate
		"-b:a", "128k", // Audio bitrate
		"-f", "dash", // DASH format
		"-use_timeline", "1", // Use timeline
		"-use_template", "1", // Use template
		"-init_seg_name", "init-$RepresentationID$.m4s", // Init segment naming
		"-media_seg_name", "chunk-$RepresentationID$-$Number%05d$.m4s", // Media segment naming
		"-seg_duration", "4", // Segment duration in seconds
		manifestPath, // Output manifest file
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute ffmpeg command: %w", err)
	}

	return nil
}

// Uncomment the following line to ensure FSVideoContentService implements VideoContentService
var _ VideoContentService = (*FSVideoContentService)(nil)
