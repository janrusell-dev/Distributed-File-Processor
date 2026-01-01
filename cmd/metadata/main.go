package main

import (
	"log"
	"net"

	"github.com/janrusell-dev/distributed-file-processor/internal/config"
	"github.com/janrusell-dev/distributed-file-processor/internal/db"
	"github.com/janrusell-dev/distributed-file-processor/internal/db/sqlc"
	"github.com/janrusell-dev/distributed-file-processor/internal/services"
	pb "github.com/janrusell-dev/distributed-file-processor/proto/metadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	cfg := config.Load()

	conn, err := db.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	queries := sqlc.New(conn)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	healthServer := health.NewServer()
	healthServer.SetServingStatus(
		"metadata.MetadataService",
		grpc_health_v1.HealthCheckResponse_SERVING,
	)

	grpc_health_v1.RegisterHealthServer(server, healthServer)

	pb.RegisterMetadataServiceServer(
		server, services.NewMetaDataService(queries),
	)

	log.Println("Metadata service running on :" + cfg.GRPCPort)
	server.Serve(lis)

}
