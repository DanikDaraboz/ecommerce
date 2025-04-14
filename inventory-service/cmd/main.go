package main

import (
    "inventory-service/internal/config"
    "inventory-service/internal/repository"
    "inventory-service/internal/server"
    "inventory-service/internal/service"
    "log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "context"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    defer client.Disconnect(context.Background())

    db := client.Database("inventory")
    repo := repository.NewProductRepository(db)
    svc := service.NewProductService(repo)
    srv := server.NewInventoryServer(svc)

    if err := srv.Run(cfg.GRPCAddr); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}