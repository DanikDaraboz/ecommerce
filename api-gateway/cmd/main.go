package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/danikdaraboz/ecommerce/api-gateway/internal/config"
	"github.com/danikdaraboz/ecommerce/api-gateway/internal/service"
	"github.com/danikdaraboz/ecommerce/api-gateway/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создание gRPC клиентов
	inventoryService := service.NewInventoryClient(cfg.InventoryServiceAddress)
	orderService := service.NewOrderClient(cfg.OrderServiceAddress)

	// Настройка маршрутов Gin
	r := router.SetupRouter(inventoryService, orderService)

	// Запуск сервера
	go func() {
		if err := r.Run(":" + cfg.GatewayPort); err != nil {
			log.Fatalf("Failed to start gateway server: %v", err)
		}
	}()

	// Обработка завершения работы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}
