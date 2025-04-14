package config

import (
	"os"
)

type Config struct {
	GRPCPort    string
	MongoURI    string
	MongoDBName string
}

func Load() (*Config, error) {
	return &Config{
		GRPCPort:    getEnv("GRPC_PORT", "50052"),
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName: getEnv("MONGO_DB_NAME", "order-service"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}