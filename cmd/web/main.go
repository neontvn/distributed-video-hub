package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"strings"
	"tritontube/internal/proto"
	"tritontube/internal/web"

	"google.golang.org/grpc"
)

// printUsage prints the usage information for the application
func printUsage() {
	fmt.Println("Usage: ./program [OPTIONS] METADATA_TYPE METADATA_OPTIONS CONTENT_TYPE CONTENT_OPTIONS")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  METADATA_TYPE         Metadata service type (sqlite, etcd)")
	fmt.Println("  METADATA_OPTIONS      Options for metadata service (e.g., db path)")
	fmt.Println("  CONTENT_TYPE          Content service type (fs, nw)")
	fmt.Println("  CONTENT_OPTIONS       Options for content service (e.g., base dir, network addresses)")
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("Example: ./program sqlite db.db fs /path/to/videos")
}

func startAdminServer(networkService *web.NetworkVideoContentService, grpcServerAddr string) error {
	lis, err := net.Listen("tcp", grpcServerAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", grpcServerAddr, err)
	}

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc.MaxSendMsgSize(math.MaxInt32),
	)

	// Register the network service as the admin service implementation
	proto.RegisterVideoContentAdminServiceServer(grpcServer, networkService)

	// Start the gRPC server in a goroutine
	go func() {
		log.Printf("Admin gRPC server listening on %s", grpcServerAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("Admin gRPC server error: %v", err)
		}
	}()

	return nil
}

func main() {
	// Define flags
	port := flag.Int("port", 8080, "Port number for the web server")
	host := flag.String("host", "localhost", "Host address for the web server")

	// Set custom usage message
	flag.Usage = printUsage

	// Parse flags
	flag.Parse()

	// Check if the correct number of positional arguments is provided
	if len(flag.Args()) != 4 {
		fmt.Println("Error: Incorrect number of arguments")
		printUsage()
		return
	}

	// Parse positional arguments
	metadataServiceType := flag.Arg(0)
	metadataServiceOptions := flag.Arg(1)
	contentServiceType := flag.Arg(2)
	contentServiceOptions := flag.Arg(3)

	// Validate port number (already an int from flag, check if positive)
	if *port <= 0 {
		fmt.Println("Error: Invalid port number:", *port)
		printUsage()
		return
	}

	// Construct metadata service
	var metadataService web.VideoMetadataService
	fmt.Println("Creating metadata service of type", metadataServiceType, "with options", metadataServiceOptions)
	// TODO: Implement metadata service creation logic
	switch metadataServiceType {
	case "sqlite":
		dbInstance, err := sql.Open("sqlite3", metadataServiceOptions)
		if err != nil {
			fmt.Println("Error: Initializing sqlite3", metadataServiceType)
			return
		}
		defer dbInstance.Close()
		metadataService = &web.SQLiteVideoMetadataService{Instance: dbInstance}
	default:
		fmt.Println("Error: Unsupported metadata service type:", metadataServiceType)
		printUsage()
		return
	}

	// Construct content service
	var contentService web.VideoContentService
	fmt.Println("Creating content service of type", contentServiceType, "with options", contentServiceOptions)
	switch contentServiceType {
	case "fs":
		contentService = &web.FSVideoContentService{BaseDir: contentServiceOptions}
	case "nw":
		addresses := strings.Split(contentServiceOptions, ",")
		if len(addresses) < 2 {
			fmt.Println("Error: Network content service requires at least two addresses (gRPC server and one storage server)")
			printUsage()
			return
		}

		grpcServerAddr := addresses[0]
		storageServers := addresses[1:]

		networkService := web.NewNetworkVideoContentService(
			storageServers,
		)

		fmt.Printf("Starting admin gRPC server on %s...\n", grpcServerAddr)
		if err := startAdminServer(networkService, grpcServerAddr); err != nil {
			fmt.Printf("Error: Failed to start admin gRPC server: %v\n", err)
			return
		}
		contentService = networkService
		fmt.Printf("Network content service initialized with gRPC server at %s and %d storage servers\n",
			grpcServerAddr, len(storageServers))

	default:
		fmt.Println("Error: Unsupported content service type:", contentServiceType)
		printUsage()
		return
	}

	// Start the server
	server := web.NewServer(metadataService, contentService)
	listenAddr := fmt.Sprintf("%s:%d", *host, *port)
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("Error starting listener:", err)
		return
	}
	defer lis.Close()

	fmt.Println("Starting web server on", listenAddr)
	err = server.Start(lis)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
