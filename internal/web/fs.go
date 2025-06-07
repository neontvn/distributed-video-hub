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

	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	tempDir := filepath.Join(projectRoot, "video-upload-"+videoId)
	err = os.Mkdir(tempDir, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	tempInputFile := filepath.Join(tempDir, filename)
	err = os.WriteFile(tempInputFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write input file: %w", err)
	}

	tempManifestPath := filepath.Join(tempDir, "manifest.mpd")

	cmd := exec.Command(
		"ffmpeg",
		"-i", tempInputFile,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-bf", "1",
		"-keyint_min", "120",
		"-g", "120",
		"-sc_threshold", "0",
		"-b:v", "3000k",
		"-b:a", "128k",
		"-f", "dash",
		"-use_timeline", "1",
		"-use_template", "1",
		"-init_seg_name", "init-$RepresentationID$.m4s",
		"-media_seg_name", "chunk-$RepresentationID$-$Number%05d$.m4s",
		"-seg_duration", "4",
		tempManifestPath,
	)
	cmd.Dir = tempDir

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute ffmpeg command: %w", err)
	}

	err = os.Remove(tempInputFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete temp input file: %w", err)
	}

	files, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read temp directory: %w", err)
	}
	for _, file := range files {
		src := filepath.Join(tempDir, file.Name())
		dst := filepath.Join(videoDir, file.Name())
		err := os.Rename(src, dst)
		if err != nil {
			return fmt.Errorf("failed to move file %s: %w", file.Name(), err)
		}
	}

	return nil
}

// Uncomment the following line to ensure FSVideoContentService implements VideoContentService
var _ VideoContentService = (*FSVideoContentService)(nil)
