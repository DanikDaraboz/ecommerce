package handler

import (
    "api-gateway/internal/client"
    pb "ecommerce/proto/order"
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
)

func SetupOrderRoutes(router *gin.Engine, client *client.OrderClient) {
    orders := router.Group("/orders")
    {
        orders.POST("", createOrder(client))
        orders.GET(":id", getOrderByID(client))
        orders.PUT(":id/status", updateOrderStatus(client))
        orders.GET("", listUserOrders(client))
    }
}

type OrderItemInput struct {
    ProductID string `json:"product_id" binding:"required"`
    Quantity  int32  `json:"quantity" binding:"required"`
}

type CreateOrderInput struct {
    UserID string           `json:"user_id" binding:"required"`
    Items  []OrderItemInput `json:"items" binding:"required"`
}

type UpdateOrderStatusInput struct {
    Status string `json:"status" binding:"required"`
}

func createOrder(client *client.OrderClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input CreateOrderInput
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        req := &pb.CreateOrderRequest{
            UserId: input.UserID,
        }
        for _, item := range input.Items {
            req.Items = append(req.Items, &pb.OrderItem{
                ProductId: item.ProductID,
                Quantity:  item.Quantity,
            })
        }

        resp, err := client.CreateOrder(c, req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, resp)
    }
}

func getOrderByID(client *client.OrderClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        req := &pb.GetOrderRequest{Id: id}

        resp, err := client.GetOrderByID(c, req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, resp)
    }
}

func updateOrderStatus(client *client.OrderClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        var input UpdateOrderStatusInput
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        req := &pb.UpdateOrderStatusRequest{
            Id:     id,
            Status: input.Status,
        }

        resp, err := client.UpdateOrderStatus(c, req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, resp)
    }
}

func listUserOrders(client *client.OrderClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.Query("user_id")
        page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
        pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

        req := &pb.ListOrdersRequest{
            UserId:   userID,
            Page:     int32(page),
            PageSize: int32(pageSize),
        }

        resp, err := client.ListUserOrders(c, req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, resp)
    }
}