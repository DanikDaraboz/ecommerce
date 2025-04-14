package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/danikdaraboz/ecommerce/api-gateway/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderClient
}

func NewOrderHandler(orderService *service.OrderClient) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var request struct {
		UserID string                `json:"user_id"`
		Items  []*proto.OrderItem    `json:"items"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	order, err := h.orderService.CreateOrder(c, request.UserID, request.Items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusOK, order)
}
