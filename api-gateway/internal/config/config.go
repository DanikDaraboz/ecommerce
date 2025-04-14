package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	InventoryServiceAddress string
	OrderServiceAddress     string
	GatewayPort            string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		InventoryServiceAddress: getEnv("INVENTORY_SERVICE_ADDRESS", "localhost:50051"),
		OrderServiceAddress:     getEnv("ORDER_SERVICE_ADDRESS", "localhost:50052"),
		GatewayPort:             getEnv("GATEWAY_PORT", "8080"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
