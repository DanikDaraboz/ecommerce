package handler

import (
    "api-gateway/internal/client"
    pb "ecommerce/proto/inventory"
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
)

func SetupProductRoutes(router *gin.Engine, client *client.InventoryClient) {
    products := router.Group("/products")
    {
        products.POST("", createProduct(client))
        products.GET(":id", getProductByID(client))
        products.PUT(":id", updateProduct(client))
        products.DELETE(":id", deleteProduct(client))
        products.GET("", listProducts(client))
    }
}

type CreateProductInput struct {
    Name        string  `json:"name" binding:"required"`
    Description string  `json:"description"`
    Price       float64 `json:"price" binding:"required"`
    Stock       int32   `json:"stock" binding:"required"`
    Category    string  `json:"category" binding:"required"`
}

func createProduct(client *client.InventoryClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input CreateProductInput
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        req := &pb.CreateProductRequest{
            Name:        input.Name,
            Description: input.Description,
            Price:       input.Price,
            Stock:       input.Stock,
            Category:    input.Category,
        }

        resp, err := client.CreateProduct(c, req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, resp)
    }
}

func getProductByID(client *client.InventoryClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        req := &pb.GetProductRequest{Id: id}

        resp, err := client.GetProductByID(c, req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, resp)
    }
}

func updateProduct(client *client.InventoryClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        var input CreateProductInput
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        req := &pb.UpdateProductRequest{
            Id:          id,
            Name:        input.Name,
            Description: input.Description,
            Price:       input.Price,
            Stock:       input.Stock,
            Category:    input.Category,
        }

        resp, err := client.UpdateProduct(c, req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, resp)
    }
}

func deleteProduct(client *client.InventoryClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        req := &pb.DeleteProductRequest{Id: id}

        _, err := client.DeleteProduct(c, req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
    }
}

func listProducts(client *client.InventoryClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
        pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

        req := &pb.ListProductsRequest{
            Page:     int32(page),
            PageSize: int32(pageSize),
        }

        resp, err := client.ListProducts(c, req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, resp)
    }
}