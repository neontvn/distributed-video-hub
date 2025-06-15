// Lab 8: Implement a network video content service (client using consistent hashing)

package web

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"tritontube/internal/proto"
)

type NetworkVideoContentService struct {
	proto.UnimplementedVideoContentAdminServiceServer
	StorageServers []string
	hashRing       []uint64
	serverMap      map[uint64]string
	mu             sync.RWMutex // Protects StorageServers, hashRing, and serverMap
}

func NewNetworkVideoContentService(servers []string) *NetworkVideoContentService {
	service := &NetworkVideoContentService{
		StorageServers: servers,
		serverMap:      make(map[uint64]string),
	}
	service.initHashRing()

	return service
}

func (n *NetworkVideoContentService) initHashRing() {
	n.hashRing = make([]uint64, 0, len(n.StorageServers))

	for _, server := range n.StorageServers {
		hash := hashStringToUint64(server)
		n.hashRing = append(n.hashRing, hash)
		n.serverMap[hash] = server
	}

	sort.Slice(n.hashRing, func(i, j int) bool {
		return n.hashRing[i] < n.hashRing[j]
	})
}

func hashStringToUint64(s string) uint64 {
	sum := sha256.Sum256([]byte(s))
	return binary.BigEndian.Uint64(sum[:8])
}

func (n *NetworkVideoContentService) getServerForKey(videoId string, filename string) string {
	n.mu.RLock()
	defer n.mu.RUnlock()

	key := fmt.Sprintf("%s/%s", videoId, filename)
	hash := hashStringToUint64(key)
	idx := sort.Search(len(n.hashRing), func(i int) bool {
		return n.hashRing[i] >= hash
	})

	if idx == len(n.hashRing) {
		idx = 0
	}

	return n.serverMap[n.hashRing[idx]]
}

// Admin Service Implementation
// I have implemented the three methods that were expected:
// ListNodes, AddNode, RemoveNode
func (n *NetworkVideoContentService) ListNodes(ctx context.Context, req *proto.ListNodesRequest) (*proto.ListNodesResponse, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	listOfServers := make([]string, len(n.hashRing))
	for i, hash := range n.hashRing {
		listOfServers[i] = n.serverMap[hash]
	}

	return &proto.ListNodesResponse{Nodes: listOfServers}, nil
}

func (n *NetworkVideoContentService) AddNode(ctx context.Context, req *proto.AddNodeRequest) (*proto.AddNodeResponse, error) {
	nodeAddr := req.NodeAddress
	log.Printf("Adding node: %s", nodeAddr)

	n.mu.Lock()
	defer n.mu.Unlock()

	for _, server := range n.StorageServers {
		if server == nodeAddr {
			return nil, fmt.Errorf("node %s already exists", nodeAddr)
		}
	}

	allFiles, err := n.getAllFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get all files: %w", err)
	}

	oldMapping := make(map[string]string)
	for server, files := range allFiles {
		for _, file := range files {
			key := fmt.Sprintf("%s/%s", file.VideoId, file.Filename)
			oldMapping[key] = server
		}
	}

	n.StorageServers = append(n.StorageServers, nodeAddr)
	n.initHashRing()
	migratedCount, err := n.migrateFiles(allFiles, oldMapping)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate files: %w", err)
	}

	log.Printf("Successfully added node %s, migrated %d files", nodeAddr, migratedCount)
	return &proto.AddNodeResponse{MigratedFileCount: int32(migratedCount)}, nil
}

func (n *NetworkVideoContentService) RemoveNode(ctx context.Context, req *proto.RemoveNodeRequest) (*proto.RemoveNodeResponse, error) {
	nodeAddr := req.NodeAddress
	fmt.Printf("Removing node: %s", nodeAddr)

	n.mu.Lock()
	defer n.mu.Unlock()

	nodeIndex := -1
	for i, server := range n.StorageServers {
		if server == nodeAddr {
			nodeIndex = i
			break
		}
	}
	if nodeIndex == -1 {
		return nil, fmt.Errorf("Node %s not found", nodeAddr)
	}

	allFiles, err := n.getAllFiles()
	if err != nil {
		return nil, fmt.Errorf("Failed to get all files: %w", err)
	}
	oldMapping := make(map[string]string)
	for server, files := range allFiles {
		for _, file := range files {
			key := fmt.Sprintf("%s/%s", file.VideoId, file.Filename)
			oldMapping[key] = server
		}
	}
	n.StorageServers = append(n.StorageServers[:nodeIndex], n.StorageServers[nodeIndex+1:]...)
	n.initHashRing()

	migratedCount, err := n.migrateFiles(allFiles, oldMapping)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate files: %w", err)
	}

	log.Printf("Successfully removed node %s, migrated %d files", nodeAddr, migratedCount)
	return &proto.RemoveNodeResponse{MigratedFileCount: int32(migratedCount)}, nil
}

func (n *NetworkVideoContentService) getAllFiles() (map[string][]*proto.FileInfo, error) {
	allFiles := make(map[string][]*proto.FileInfo)

	for _, server := range n.StorageServers {
		files, err := n.listFilesOnServer(server)
		if err != nil {
			return nil, fmt.Errorf("failed to list files on server %s: %w", server, err)
		}
		allFiles[server] = files
	}

	return allFiles, nil
}

func (n *NetworkVideoContentService) listFilesOnServer(server string) ([]*proto.FileInfo, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32),
			grpc.MaxCallSendMsgSize(math.MaxInt32),
		),
	}

	conn, err := grpc.NewClient(server, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to storage server %s: %w", server, err)
	}
	defer conn.Close()

	client := proto.NewVideoContentStorageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, err := client.ListFiles(ctx, &proto.ListFilesRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return response.Files, nil
}

func (n *NetworkVideoContentService) migrateFiles(allFiles map[string][]*proto.FileInfo, oldMapping map[string]string) (int, error) {
	migratedCount := 0

	// TODO: iterate through all the files and for every file check current server and the new server and move the file
	for _, files := range allFiles {
		for _, file := range files {
			key := fmt.Sprintf("%s/%s", file.VideoId, file.Filename)

			oldServer := oldMapping[key]
			newServer := n.identifyServerForGivenKey(file.VideoId, file.Filename)

			if oldServer != newServer {

				log.Printf("Migrating file %s from %s to %s", key, oldServer, newServer)
				if err := n.moveFile(file, oldServer, newServer); err != nil {
					return migratedCount, fmt.Errorf("Error: Failed to move file %s from %s to %s: %w",
						key, oldServer, newServer, err)
				}
				migratedCount++
			}
		}
	}

	return migratedCount, nil
}

func (n *NetworkVideoContentService) identifyServerForGivenKey(videoId string, filename string) string {
	if len(n.hashRing) == 0 {
		return ""
	}

	key := fmt.Sprintf("%s/%s", videoId, filename)
	hash := hashStringToUint64(key)
	idx := sort.Search(len(n.hashRing), func(i int) bool {
		return n.hashRing[i] >= hash
	})

	if idx == len(n.hashRing) {
		idx = 0
	}

	return n.serverMap[n.hashRing[idx]]
}

func (n *NetworkVideoContentService) moveFile(file *proto.FileInfo, fromServer, toServer string) error {

	// TODO: read, write and then delete
	data, err := n.readFileFromServer(file.VideoId, file.Filename, fromServer)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := n.writeFileToServer(file.VideoId, file.Filename, data, toServer); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	if err := n.deleteFileFromServer(file.VideoId, file.Filename, fromServer); err != nil {
		return fmt.Errorf("failed to delete file from source: %w", err)
	}

	return nil
}

func (n *NetworkVideoContentService) readFileFromServer(videoId, filename, server string) ([]byte, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32),
			grpc.MaxCallSendMsgSize(math.MaxInt32),
		),
	}

	conn, err := grpc.NewClient(server, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to storage server %s: %w", server, err)
	}
	defer conn.Close()

	client := proto.NewVideoContentStorageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	response, err := client.Read(ctx, &proto.ReadRequest{
		VideoId:  videoId,
		Filename: filename,
	})
	if err != nil {
		return nil, fmt.Errorf("storage server read failed: %w", err)
	}

	return response.Data, nil
}

func (n *NetworkVideoContentService) writeFileToServer(videoId, filename string, data []byte, server string) error {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32),
			grpc.MaxCallSendMsgSize(math.MaxInt32),
		),
	}

	conn, err := grpc.NewClient(server, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to storage server %s: %w", server, err)
	}
	defer conn.Close()

	client := proto.NewVideoContentStorageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = client.Write(ctx, &proto.WriteRequest{
		VideoId:  videoId,
		Filename: filename,
		Data:     data,
	})
	if err != nil {
		return fmt.Errorf("storage server write failed: %w", err)
	}

	return nil
}

func (n *NetworkVideoContentService) deleteFileFromServer(videoId, filename, server string) error {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32),
			grpc.MaxCallSendMsgSize(math.MaxInt32),
		),
	}

	conn, err := grpc.NewClient(server, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to storage server %s: %w", server, err)
	}
	defer conn.Close()

	client := proto.NewVideoContentStorageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = client.DeleteFile(ctx, &proto.DeleteFileRequest{
		VideoId:  videoId,
		Filename: filename,
	})
	if err != nil {
		return fmt.Errorf("storage server delete failed: %w", err)
	}

	return nil
}

// Existing VideoContentService methods
// Read and Write methods
// writeToStorageServer is a helper func which I have used in Write method
// to write to the storage server consistent hash
func (n *NetworkVideoContentService) Read(videoId string, filename string) ([]byte, error) {
	server := n.getServerForKey(videoId, filename)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32),
			grpc.MaxCallSendMsgSize(math.MaxInt32),
		),
	}

	conn, err := grpc.NewClient(server, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to storage server %s: %w", server, err)
	}
	defer conn.Close()

	client := proto.NewVideoContentStorageServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := client.Read(ctx, &proto.ReadRequest{
		VideoId:  videoId,
		Filename: filename,
	})
	if err != nil {
		return nil, fmt.Errorf("storage server read failed for %s: %w", filename, err)
	}

	return response.Data, nil
}

func (n *NetworkVideoContentService) Write(videoId string, filename string, data []byte) error {
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

	// Generate thumbnail
	tempThumbnailPath := filepath.Join(tempDir, "thumbnail.jpg")
	thumbnailCmd := exec.Command(
		"ffmpeg",
		"-i", tempInputFile,
		"-ss", "00:00:01", // Take frame at 1 second
		"-vframes", "1", // Take only one frame
		"-q:v", "2", // High quality
		"-y", // Overwrite output file if it exists
		tempThumbnailPath,
	)
	if err := thumbnailCmd.Run(); err != nil {
		return fmt.Errorf("failed to generate thumbnail: %w", err)
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
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(tempDir, file.Name())
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read generated file %s: %w", file.Name(), err)
		}

		err = n.writeToStorageServer(videoId, file.Name(), fileData)
		if err != nil {
			return fmt.Errorf("failed to write file %s to storage: %w", file.Name(), err)
		}
	}

	return nil
}

func (n *NetworkVideoContentService) writeToStorageServer(videoId string, filename string, data []byte) error {
	server := n.getServerForKey(videoId, filename)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32),
			grpc.MaxCallSendMsgSize(math.MaxInt32),
		),
	}

	conn, err := grpc.NewClient(server, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to storage server %s: %w", server, err)
	}
	defer conn.Close()

	client := proto.NewVideoContentStorageServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = client.Write(ctx, &proto.WriteRequest{
		VideoId:  videoId,
		Filename: filename,
		Data:     data,
	})
	if err != nil {
		return fmt.Errorf("storage server write failed for %s: %w", filename, err)
	}

	return nil
}

// Delete implements VideoContentService.
func (n *NetworkVideoContentService) Delete(videoId string, filename string) error {
	server := n.getServerForKey(videoId, filename)
	return n.deleteFileFromServer(videoId, filename, server)
}

// ListFiles implements VideoContentService.
func (n *NetworkVideoContentService) ListFiles(videoId string) ([]string, error) {
	allFiles, err := n.getAllFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get all files: %w", err)
	}

	var files []string
	for _, serverFiles := range allFiles {
		for _, file := range serverFiles {
			if file.VideoId == videoId {
				files = append(files, file.Filename)
			}
		}
	}
	return files, nil
}

var _ VideoContentService = (*NetworkVideoContentService)(nil)
var _ proto.VideoContentAdminServiceServer = (*NetworkVideoContentService)(nil)
