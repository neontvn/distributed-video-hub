// Lab 8: Implement a network video content service (server)

package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"tritontube/internal/proto"
)

type StorageServer struct {
	proto.UnimplementedVideoContentStorageServiceServer
	BaseDir string
}

func NewStorageServer(baseDir string, port int) *StorageServer {
	storageDir := filepath.Join(baseDir)
	err := os.MkdirAll(storageDir, 0755)
	if err != nil {
		log.Printf("Warning: Failed to create storage directory: %v", err)
	}
	return &StorageServer{
		BaseDir: baseDir,
	}
}

func (s *StorageServer) Write(ctx context.Context, req *proto.WriteRequest) (*proto.WriteResponse, error) {

	videoDir := filepath.Join(s.BaseDir, req.VideoId)
	err := os.MkdirAll(videoDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(videoDir, req.Filename)
	err = os.WriteFile(filePath, req.Data, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &proto.WriteResponse{Success: true}, nil
}

func (s *StorageServer) Read(ctx context.Context, req *proto.ReadRequest) (*proto.ReadResponse, error) {
	filePath := filepath.Join(s.BaseDir, req.VideoId, req.Filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return &proto.ReadResponse{Data: data}, nil
}

func (s *StorageServer) ListFiles(ctx context.Context, req *proto.ListFilesRequest) (*proto.ListFilesResponse, error) {
	var files []*proto.FileInfo

	err := filepath.Walk(s.BaseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(s.BaseDir, path)
			if err != nil {
				return err
			}
			videoId := filepath.Dir(relPath)
			filename := filepath.Base(relPath)
			files = append(files, &proto.FileInfo{
				VideoId:  videoId,
				Filename: filename,
			})
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return &proto.ListFilesResponse{Files: files}, nil
}

func (s *StorageServer) DeleteFile(ctx context.Context, req *proto.DeleteFileRequest) (*proto.DeleteFileResponse, error) {
	filePath := filepath.Join(s.BaseDir, req.VideoId, req.Filename)
	err := os.Remove(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to delete file: %w", err)
	}

	videoDir := filepath.Join(s.BaseDir, req.VideoId)
	os.Remove(videoDir)
	return &proto.DeleteFileResponse{Success: true}, nil
}
