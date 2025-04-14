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
	// Set default values
	cfg := &Config{
		GRPCPort:    getEnv("GRPC_PORT", "50051"),
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName: getEnv("MONGO_DB_NAME", "inventory"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}