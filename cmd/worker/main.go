package main

import (
	"context"
	"log"

	"github.com/janrusell-dev/distributed-file-processor/internal/cache"
	"github.com/janrusell-dev/distributed-file-processor/internal/config"
	"github.com/janrusell-dev/distributed-file-processor/internal/services"
	"github.com/janrusell-dev/distributed-file-processor/proto/metadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()

	conn, err := grpc.NewClient(cfg.MetadataAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	redisClient := cache.NewRedisClient(cfg.RedisAddr)
	metaClient := metadata.NewMetadataServiceClient(conn)

	worker := services.NewWorker(redisClient, metaClient)

	worker.Start(context.Background())
}
