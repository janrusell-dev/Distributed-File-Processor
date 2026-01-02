package services

import (
	"bytes"
	"io"

	pb "github.com/janrusell-dev/distributed-file-processor/proto/upload"
)

type UploadGRPCServer struct {
	pb.UnimplementedUploadServiceServer
	upload *UploadService
}

func (s *UploadGRPCServer) UploadFile(stream pb.UploadService_UploadFileServer) error {
	var buf bytes.Buffer
	var filename string
	var mimeType string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			id, err := s.upload.ProcessUpload(
				stream.Context(),
				filename,
				&buf,
				int64(buf.Len()),
				mimeType,
			)
			if err != nil {
				return err
			}

			return stream.SendAndClose(&pb.UploadFileResponse{Id: id})
		}
		if err != nil {
			return err
		}

		switch x := req.Data.(type) {
		case *pb.UploadFileRequest_Metadata:
			filename = x.Metadata.Filename
			mimeType = x.Metadata.MimeType
		case *pb.UploadFileRequest_Chunk:
			buf.Write(x.Chunk)
		}
	}
}
