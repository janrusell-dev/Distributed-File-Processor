package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	pb "github.com/janrusell-dev/distributed-file-processor/proto/upload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/client/main.go testdata/sample.txt")
	}

	filepath := os.Args[1]

	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewUploadServiceClient(conn)
	uploadFile(client, filepath)
}

func uploadFile(client pb.UploadServiceClient, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	stream, err := client.UploadFile(ctx)
	if err != nil {
		log.Fatal(err)
	}

	stream.Send(&pb.UploadFileRequest{
		Data: &pb.UploadFileRequest_Metadata{
			Metadata: &pb.Metadata{
				Filename: filepath.Base(filePath),
				MimeType: "application/octet-stream",
			},
		},
	})

	buf := make([]byte, 32*1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		stream.Send(&pb.UploadFileRequest{
			Data: &pb.UploadFileRequest_Chunk{
				Chunk: buf[:n],
			},
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Success file ID: %s\n", res)
}
