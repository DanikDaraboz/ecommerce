package router

import (
	"github.com/gin-gonic/gin"
	"github.com/danikdaraboz/ecommerce/api-gateway/internal/handler"
	"github.com/danikdaraboz/ecommerce/api-gateway/internal/service"
)

func SetupRouter(inventoryService *service.InventoryClient, orderService *service.OrderClient) *gin.Engine {
	r := gin.Default()

	productHandler := handler.NewProductHandler(inventoryService)
	r.GET("/products/:id", productHandler.GetProductByID)

	orderHandler := handler.NewOrderHandler(orderService)
	r.POST("/orders", orderHandler.CreateOrder)

	return r
}
