package main

import (
    "order-service/internal/config"
    "order-service/internal/repository"
    "order-service/internal/server"
    "order-service/internal/service"
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

    db := client.Database("orders")
    repo := repository.NewOrderRepository(db)
    svc := service.NewOrderService(repo)
    srv := server.NewOrderServer(svc)

    if err := srv.Run(cfg.GRPCAddr); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}