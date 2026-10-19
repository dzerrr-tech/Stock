package router

import (
	"github.com/gin-gonic/gin"

	"Stock/apps/api/internal/handler"
)

func New(health *handler.HealthHandler, product *handler.ProductHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/api/v1/health", health.HealthCheck)
	r.GET("/api/v1/health/db", health.DatabaseCheck)
	r.POST("/api/v1/products", product.Create)
	r.GET("/api/v1/products", product.List)

	return r
}
