package middleware

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Simplified auth (e.g., check header or token)
        // For demo, allow all requests
        c.Next()
    }
}