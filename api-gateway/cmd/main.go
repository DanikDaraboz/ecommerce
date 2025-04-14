package main

import (
    "api-gateway/internal/config"
    "api-gateway/internal/handler"
    "api-gateway/internal/middleware"
    "github.com/gin-gonic/gin"
    "log"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    router := gin.Default()
    router.Use(middleware.Logger())

    // Initialize gRPC clients
    inventoryClient, err := handler.NewInventoryClient(cfg.InventoryServiceAddr)
    if err != nil {
        log.Fatalf("Failed to connect to inventory service: %v", err)
    }
    orderClient, err := handler.NewOrderClient(cfg.OrderServiceAddr)
    if err != nil {
        log.Fatalf("Failed to connect to order service: %v", err)
    }

    // Setup handlers
    handler.SetupProductRoutes(router, inventoryClient)
    handler.SetupOrderRoutes(router, orderClient)

    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}