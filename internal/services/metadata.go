package services

import (
	"context"

	"github.com/google/uuid"
	db "github.com/janrusell-dev/distributed-file-processor/internal/db/sqlc"
	pb "github.com/janrusell-dev/distributed-file-processor/proto/metadata"
)

type MetadataService struct {
	pb.UnimplementedMetadataServiceServer
	queries *db.Queries
}

func NewMetaDataService(q *db.Queries) *MetadataService {
	return &MetadataService{queries: q}
}

func (s *MetadataService) CreateMetadata(
	ctx context.Context, req *pb.CreateMetadataRequest,
) (*pb.CreateMetadataResponse, error) {
	id := uuid.New()

	err := s.queries.CreateFile(ctx, db.CreateFileParams{
		ID:       id,
		Filename: req.Filename,
		Size:     req.Size,
		MimeType: req.MimeType,
		Status:   "uploaded",
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateMetadataResponse{
		Id: id.String(),
	}, nil
}

func (s *MetadataService) GetMetadata(ctx context.Context,
	req *pb.GetMetadataRequest) (*pb.GetMetadataResponse, error) {
	file, err := s.queries.GetFile(ctx, uuid.MustParse(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.GetMetadataResponse{
		Id:       file.ID.String(),
		Filename: file.Filename,
		Size:     file.Size,
		MimeType: file.MimeType,
		Status:   file.Status,
	}, nil
}
