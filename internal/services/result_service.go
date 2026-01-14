package services

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/janrusell-dev/distributed-file-processor/internal/cache"
	db "github.com/janrusell-dev/distributed-file-processor/internal/db/sqlc"
	"github.com/janrusell-dev/distributed-file-processor/proto/result"
	pb "github.com/janrusell-dev/distributed-file-processor/proto/result"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResultService struct {
	pb.UnimplementedResultServiceServer
	queries *db.Queries
	redis   *cache.RedisClient
}

func NewResultService(q *db.Queries, r *cache.RedisClient) *ResultService {
	return &ResultService{queries: q, redis: r}
}

func (s *ResultService) StoreResult(ctx context.Context,
	req *pb.StoreResultRequest) (*pb.StoreResultResponse, error) {

	if req.GetFileId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "file_id is required")
	}

	fileId, err := uuid.Parse(req.FileId)
	if err != nil {
		return nil, err
	}
	_, err = s.queries.CreateResult(ctx, db.CreateResultParams{
		FileID:     fileId,
		OutputData: req.ResultData,
	})

	if err != nil {
		log.Printf("Failed to save result: %v", err)
	}

	_ = s.redis.SetResult(ctx, req.FileId, req.ResultData)
	if err != nil {
		log.Printf("Redis cache failure: %v", err)
	}
	return &pb.StoreResultResponse{Success: true}, nil
}

func (s *ResultService) GetResult(ctx context.Context,
	req *result.GetResultRequest) (*result.GetResultResponse, error) {
	cachedData, err := s.redis.GetResult(ctx, req.FileId)
	if err == nil {
		log.Printf("Cache hit for file: %s", req.FileId)
		return &result.GetResultResponse{
			FileId:     req.FileId,
			ResultData: cachedData,
		}, nil
	}
	log.Printf("Cache miss for file: %s. Querying database...", req.FileId)
	fileId, err := uuid.Parse(req.FileId)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID format")
	}

	res, err := s.queries.GetResultByFileID(ctx, fileId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "result not found in database")
	}
	_ = s.redis.SetResult(ctx, req.FileId, res.OutputData)

	return &result.GetResultResponse{
		FileId:     req.FileId,
		ResultData: res.OutputData,
	}, nil
}
