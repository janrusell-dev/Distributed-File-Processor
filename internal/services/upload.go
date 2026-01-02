package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/janrusell-dev/distributed-file-processor/internal/cache"
	pb "github.com/janrusell-dev/distributed-file-processor/proto/metadata"
)

type UploadService struct {
	metaClient pb.MetadataServiceClient
	redis      *cache.RedisClient
	uploadDir  string
}

func NewUploadService(mc pb.MetadataServiceClient, r *cache.RedisClient, dir string) *UploadService {
	return &UploadService{
		metaClient: mc,
		redis:      r,
		uploadDir:  dir,
	}
}

func (s *UploadService) ProcessUpload(ctx context.Context, filename string,
	content io.Reader, size int64, mime string) (string, error) {
	resp, err := s.metaClient.CreateMetadata(ctx, &pb.CreateMetadataRequest{
		Filename: filename,
		Size:     size,
		MimeType: mime,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create metadata: %w", err)
	}
	fileID := resp.Id

	dstPath := filepath.Join(s.uploadDir, fileID)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, content); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	if err := s.redis.PushTask(ctx, fileID); err != nil {
		return "", fmt.Errorf("failed to enqueue task: %w", err)
	}

	return fileID, nil
}
