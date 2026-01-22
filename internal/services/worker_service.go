package services

import (
	"context"
	"crypto/sha256"
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/janrusell-dev/distributed-file-processor/internal/cache"
	"github.com/janrusell-dev/distributed-file-processor/proto/metadata"
	"github.com/janrusell-dev/distributed-file-processor/proto/result"
	pb "github.com/janrusell-dev/distributed-file-processor/proto/result"
)

type Worker struct {
	redis        *cache.RedisClient
	metaClient   metadata.MetadataServiceClient
	resultClient result.ResultServiceClient
	uploadDir    string
}

func NewWorker(r *cache.RedisClient, mc metadata.MetadataServiceClient,
	rc result.ResultServiceClient, ud string) *Worker {
	return &Worker{redis: r, metaClient: mc, resultClient: rc, uploadDir: ud}
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
			go func(fileId string) {
				defer func() {
					<-sem
				}()
				w.processFile(ctx, fileId)
			}(fileID)

		}
	}
}

func (w *Worker) processFile(ctx context.Context, fileId string) {
	startTime := time.Now()
	defer func() {
		log.Printf("Finished processing %s in %v", fileId, time.Since(startTime))
	}()
	_, err := w.metaClient.UpdateStatus(ctx, &metadata.UpdateStatusRequest{
		Id:     fileId,
		Status: "processing",
	})

	path := fmt.Sprintf("%s/%s", w.uploadDir, fileId)

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read file %s: %v", fileId, err)
		return
	}

	mimeType := http.DetectContentType(data)
	var detail string

	if strings.HasPrefix(mimeType, "text/") {
		words := len(strings.Fields(string(data)))
		detail = fmt.Sprintf("%d words", words)
	} else if strings.HasPrefix(mimeType, "image/") {
		cfg, _, err := image.DecodeConfig(strings.NewReader(string(data)))
		if err != nil {
			detail = "Image detected (decoded failed)"
		} else {
			detail = fmt.Sprintf("%dx%d", cfg.Width, cfg.Height)
		}
	} else {
		detail = "Binary data"
	}

	hash := sha256.Sum256(data)

	processData := fmt.Sprintf("Type: %s | Analysis: %s | SHA256: %x", mimeType, detail, hash[:8])

	_, err = w.resultClient.StoreResult(ctx, &pb.StoreResultRequest{
		FileId:     fileId,
		ResultData: processData,
	})
	if err != nil {
		log.Printf("StoreResult failed %s: %v", fileId, err)
	}
	log.Printf("Processing file: %s", fileId)

	_, err = w.metaClient.UpdateStatus(ctx, &metadata.UpdateStatusRequest{
		Id:     fileId,
		Status: "completed",
	})

	if err != nil {
		log.Printf("Failed to update status to completed for %s: %v", fileId, err)
	}

}
