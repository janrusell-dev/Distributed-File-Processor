package main

import (
	"context"
	"log"

	"github.com/janrusell-dev/distributed-file-processor/internal/cache"
	"github.com/janrusell-dev/distributed-file-processor/internal/config"
	"github.com/janrusell-dev/distributed-file-processor/internal/services"
	"github.com/janrusell-dev/distributed-file-processor/proto/metadata"
	"github.com/janrusell-dev/distributed-file-processor/proto/result"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()

	conn, err := grpc.NewClient(cfg.MetadataAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	resultConn, err := grpc.NewClient(cfg.ResultAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Result conn fail: %v", err)
	}

	redisClient := cache.NewRedisClient(cfg.RedisAddr)
	metaClient := metadata.NewMetadataServiceClient(conn)
	resultClient := result.NewResultServiceClient(resultConn)

	worker := services.NewWorker(redisClient, metaClient, resultClient)

	worker.Start(context.Background())
}
