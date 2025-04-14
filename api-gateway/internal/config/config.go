package config

import (
    "os"
)

type Config struct {
    InventoryServiceAddr string
    OrderServiceAddr     string
    MongoURI             string
}

func Load() (*Config, error) {
    return &Config{
        InventoryServiceAddr: getEnv("INVENTORY_SERVICE_ADDR", "localhost:50051"),
        OrderServiceAddr:     getEnv("ORDER_SERVICE_ADDR", "localhost:50052"),
        MongoURI:             getEnv("MONGO_URI", "mongodb://localhost:27017"),
    }, nil
}

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}