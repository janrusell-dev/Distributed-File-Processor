package services

import (
	"context"
	"log"
	"time"

	"github.com/janrusell-dev/distributed-file-processor/internal/cache"
	"github.com/janrusell-dev/distributed-file-processor/proto/metadata"
	"github.com/janrusell-dev/distributed-file-processor/proto/result"
)

type Worker struct {
	redis        *cache.RedisClient
	metaClient   metadata.MetadataServiceClient
	resultClient result.ResultServiceClient
}

func NewWorker(r *cache.RedisClient, mc metadata.MetadataServiceClient,
	rc result.ResultServiceClient) *Worker {
	return &Worker{redis: r, metaClient: mc, resultClient: rc}
}

func (w *Worker) Start(ctx context.Context) {

	sem := make(chan struct{}, 5)

	log.Println("Worker started. Watching Redis for tasks...")

	for {
		fileID, err := w.redis.PopTask(ctx)
		if err != nil {
			log.Printf("Error pulling task: %v", err)
			time.Sleep(time.Second)
			continue
		}

		if fileID != "" {
			sem <- struct{}{}
			go func(id string) {
				defer func() {
					<-sem
				}()
				w.processFile(ctx, id)
			}(fileID)

		}
	}
}

func (w *Worker) processFile(ctx context.Context, id string) {
	_, err := w.metaClient.UpdateStatus(ctx, &metadata.UpdateStatusRequest{
		Id:     id,
		Status: "processing",
	})
	if err != nil {
		log.Printf("Failed to update status to processing for %s: %v", id, err)
	}
	log.Printf("Processing file: %s", id)

	_, err = w.metaClient.UpdateStatus(ctx, &metadata.UpdateStatusRequest{
		Id:     id,
		Status: "completed",
	})

	if err != nil {
		log.Printf("Failed to update status to completed for %s: %v", id, err)
	}

	log.Printf("Finished processing: %s", id)
}
