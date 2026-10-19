package main

import (
	"context"
	"log"

	"Stock/apps/api/internal/config"
	"Stock/apps/api/internal/database"
	"Stock/apps/api/internal/handler"
	"Stock/apps/api/internal/repository"
	"Stock/apps/api/internal/router"
	"Stock/apps/api/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := database.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	log.Println("Database connection established")

	productRepo := repository.NewProductRepository(pool)
	productService := service.NewProductService(productRepo)

	r := router.New(
		handler.NewHealthHandler(pool),
		handler.NewProductHandler(productService),
	)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
