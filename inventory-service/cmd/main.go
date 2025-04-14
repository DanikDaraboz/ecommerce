package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/danikdaraboz/ecommerce/inventory-service/internal/config"
	deliverygrpc "github.com/danikdaraboz/ecommerce/inventory-service/internal/delivery/grpc"
	"github.com/danikdaraboz/ecommerce/inventory-service/internal/domain" // Add this import
	"github.com/danikdaraboz/ecommerce/inventory-service/internal/repository/mongodb"
	"github.com/danikdaraboz/ecommerce/inventory-service/internal/usecase"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize MongoDB repository as interface type
	var repo domain.ProductRepository
	repo, err = mongodb.NewProductRepository(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Type assertion for Close() with safety check
	defer func() {
		if mongoRepo, ok := repo.(interface{ Close() error }); ok {
			if err := mongoRepo.Close(); err != nil {
				log.Printf("Error closing MongoDB connection: %v", err)
			}
		}
	}()

	// Initialize use case
	productUC := usecase.NewProductUseCase(repo)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	inventoryHandler := deliverygrpc.NewInventoryHandler(productUC)
	deliverygrpc.RegisterInventoryServer(grpcServer, inventoryHandler)

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

	log.Printf("Server started on port %s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}