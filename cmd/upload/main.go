package main

import (
	"log"
	"net"

	"github.com/janrusell-dev/distributed-file-processor/internal/cache"
	"github.com/janrusell-dev/distributed-file-processor/internal/config"
	"github.com/janrusell-dev/distributed-file-processor/internal/services"
	"github.com/janrusell-dev/distributed-file-processor/proto/metadata"
	"github.com/janrusell-dev/distributed-file-processor/proto/upload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()

	conn, err := grpc.NewClient(cfg.MetadataAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to metadata service: %v", err)
	}
	defer conn.Close()

	metaClient := metadata.NewMetadataServiceClient(conn)
	redisClient := cache.NewRedisClient(cfg.RedisAddr)
	uploadLogic := services.NewUploadService(metaClient, redisClient, cfg.UploadDir)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal("Error listening", err)
	}
	server := grpc.NewServer()
	upload.RegisterUploadServiceServer(server, &services.UploadGRPCServer{
		Upload: uploadLogic,
	})
	log.Println("Succesfully connected to Metadata Service")
	server.Serve(lis)
}
