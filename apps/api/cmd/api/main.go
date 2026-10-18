package main

import (
	"Stock/apps/api/internal/config"
	"log"

	"Stock/apps/api/internal/database"
	"context"

	"github.com/gin-gonic/gin"

	"Stock/apps/api/internal/handler"
	"Stock/apps/api/internal/repository"
	"Stock/apps/api/internal/service"
)

func main() {
	cfg, err := config.Load()
	ctx := context.Background()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	log.Println("Database connection established")
	r := gin.Default()

	health := handler.NewHealthHandler(pool)
	r.GET("/api/v1/health", health.HealthCheck)
	r.GET("/api/v1/health/db", health.DatabaseCheck)

	productRepo := repository.NewProductRepository(pool)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	r.POST("/api/v1/products", productHandler.Create)
	r.GET("/api/v1/products", productHandler.List)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}

}
