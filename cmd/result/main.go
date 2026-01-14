package main

import (
	"log"
	"net"

	"github.com/janrusell-dev/distributed-file-processor/internal/cache"
	"github.com/janrusell-dev/distributed-file-processor/internal/config"
	"github.com/janrusell-dev/distributed-file-processor/internal/db"
	"github.com/janrusell-dev/distributed-file-processor/internal/db/sqlc"
	"github.com/janrusell-dev/distributed-file-processor/internal/services"
	"github.com/janrusell-dev/distributed-file-processor/proto/result"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	conn, err := db.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	queries := sqlc.New(conn)
	redisClient := cache.NewRedisClient(cfg.RedisAddr)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()

	healthServer := health.NewServer()
	healthServer.SetServingStatus(
		"result.ResultService",
		grpc_health_v1.HealthCheckResponse_SERVING,
	)

	grpc_health_v1.RegisterHealthServer(server, healthServer)

	resultService := services.NewResultService(queries, redisClient)

	result.RegisterResultServiceServer(server, resultService)

	reflection.Register(server)

	log.Printf("Result gRPC Server starting on port %s..", cfg.GRPCPort)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)

	}
}
