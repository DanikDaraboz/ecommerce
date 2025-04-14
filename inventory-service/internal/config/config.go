package config

import (
    "os"
)

type Config struct {
    GRPCAddr string
    MongoURI string
}

func Load() (*Config, error) {
    return &Config{
        GRPCAddr: getEnv("GRPC_ADDR", ":50051"),
        MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
    }, nil
}

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}