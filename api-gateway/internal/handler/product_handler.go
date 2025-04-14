package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/danikdaraboz/ecommerce/api-gateway/internal/service"
)

type ProductHandler struct {
	inventoryService *service.InventoryClient
}

func NewProductHandler(inventoryService *service.InventoryClient) *ProductHandler {
	return &ProductHandler{inventoryService: inventoryService}
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	productID := c.Param("id")
	product, err := h.inventoryService.GetProductByID(c, productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}
