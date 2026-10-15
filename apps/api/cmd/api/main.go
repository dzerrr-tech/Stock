package main

import (
	"Stock/apps/api/internal/config"
	"log"

	"Stock/apps/api/internal/database"
	"context"

	"github.com/gin-gonic/gin"

	"Stock/apps/api/internal/handler"
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

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
