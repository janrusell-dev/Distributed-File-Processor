package services

import (
	"context"
	"log"
	"time"

	"github.com/janrusell-dev/distributed-file-processor/internal/cache"
	"github.com/janrusell-dev/distributed-file-processor/proto/metadata"
)

type Worker struct {
	redis      *cache.RedisClient
	metaClient metadata.MetadataServiceClient
}

func NewWorker(r *cache.RedisClient, mc metadata.MetadataServiceClient) *Worker {
	return &Worker{redis: r, metaClient: mc}
}

func (w *Worker) Start(ctx context.Context) {
	log.Println("Worker started. Watching Redis for tasks...")

	for {
		fileID, err := w.redis.PopTask(ctx)
		if err != nil {
			log.Printf("Error pulling task: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if fileID != "" {
			w.processFile(fileID)
		}
	}
}

func (w *Worker) processFile(id string) {
	log.Printf("Processing file: %s", id)

	time.Sleep(5 * time.Second)

	log.Printf("Finished processing: %s", id)
}
