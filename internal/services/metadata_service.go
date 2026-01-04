package services

import (
	"context"
	"log"

	"github.com/google/uuid"
	db "github.com/janrusell-dev/distributed-file-processor/internal/db/sqlc"
	pb "github.com/janrusell-dev/distributed-file-processor/proto/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	log.Printf("Received Metadata creation for file: %s (%d bytes)", req.Filename, req.Size)

	err := s.queries.CreateFile(ctx, db.CreateFileParams{
		ID:       id,
		Filename: req.Filename,
		Size:     req.Size,
		MimeType: req.MimeType,
		Status:   "uploaded",
	})
	if err != nil {
		log.Printf("Error in creating Metadata: %v", err)
		return nil, err
	}

	return &pb.CreateMetadataResponse{
		Id: id.String(),
	}, nil
}

func (s *MetadataService) GetMetadata(ctx context.Context,
	req *pb.GetMetadataRequest) (*pb.GetMetadataResponse, error) {
	parsedID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID format: %v", err)
	}
	file, err := s.queries.GetFile(ctx, parsedID)
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

func (s *MetadataService) UpdateStatus(
	ctx context.Context, req *pb.UpdateStatusRequest,
) (*pb.UpdateStatusResponse, error) {
	parsedID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID: %v", err)
	}

	log.Printf("Updating file %s status to: %s", req.Id, req.Status)

	err = s.queries.UpdateFileStatus(ctx, db.UpdateFileStatusParams{
		ID:     parsedID,
		Status: req.Status,
	})

	if err != nil {
		log.Printf("Failed to update status for %s: %v", req.Id, err)
		return &pb.UpdateStatusResponse{Success: false}, err
	}
	log.Printf("Status succesfully updated for %s", req.Id)
	return &pb.UpdateStatusResponse{Success: true}, nil
}
