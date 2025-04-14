package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/danikdaraboz/ecommerce/order-service/internal/config"
	deliverygrpc "github.com/danikdaraboz/ecommerce/order-service/internal/delivery/grpc"
	"github.com/danikdaraboz/ecommerce/order-service/internal/repository/mongodb"
	"github.com/danikdaraboz/ecommerce/order-service/internal/usecase"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize MongoDB repository
	repo, err := mongodb.NewOrderRepository(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := repo.Close(); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}
	}()

	// Initialize use case
	orderUC := usecase.NewOrderUseCase(repo)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	orderHandler := deliverygrpc.NewOrderHandler(orderUC)
	deliverygrpc.RegisterOrderServer(grpcServer, orderHandler)

	// Start server
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down server...")
		grpcServer.GracefulStop()
	}()

	log.Printf("Order service started on port %s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}