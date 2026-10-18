package tests

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"Stock/apps/api/internal/database"
	"Stock/apps/api/internal/repository"
	"Stock/apps/api/internal/service"
)

func newTestDeps(t *testing.T) (*service.ProductService, *pgxpool.Pool) {
	t.Helper()

	_ = godotenv.Load("../.env")
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping integration test")
	}

	pool, err := database.NewPool(context.Background(), url)
	if err != nil {
		t.Fatalf("connect test database: %v", err)
	}
	t.Cleanup(pool.Close)

	return service.NewProductService(repository.NewProductRepository(pool)), pool
}

func uniqueSKU(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func cleanupSKU(t *testing.T, pool *pgxpool.Pool, sku string) {
	t.Helper()
	t.Cleanup(func() {
		_, err := pool.Exec(context.Background(), "DELETE FROM products WHERE sku = $1", sku)
		if err != nil {
			t.Errorf("cleanup %s: %v", sku, err)
		}
	})
}

func TestCreateProduct_NormalizesAndRejectsDuplicate(t *testing.T) {
	s, pool := newTestDeps(t)
	ctx := context.Background()

	input := uniqueSKU("it") // huruf kecil, untuk membuktikan normalisasi
	want := strings.ToUpper(input)
	cleanupSKU(t, pool, want)

	p, err := s.Create(ctx, service.CreateProductInput{SKU: input, Name: "  Baut Uji  ", MinStock: 3})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if p.ID == 0 {
		t.Error("id should be set by database")
	}
	if p.SKU != want {
		t.Errorf("sku = %q, want %q", p.SKU, want)
	}
	if p.Name != "Baut Uji" {
		t.Errorf("name = %q, want trimmed", p.Name)
	}
	if p.Stock != 0 {
		t.Errorf("stock = %d, want 0", p.Stock)
	}
	if p.MinStock != 3 {
		t.Errorf("min_stock = %d, want 3", p.MinStock)
	}

	// SKU yang sama, versi huruf besar: harus bentrok.
	_, err = s.Create(ctx, service.CreateProductInput{SKU: want, Name: "Duplikat"})
	if !errors.Is(err, service.ErrDuplicateSKU) {
		t.Fatalf("want ErrDuplicateSKU, got %v", err)
	}
}

func TestProducts_DatabaseRejectsNegativeStock(t *testing.T) {
	_, pool := newTestDeps(t)

	sku := uniqueSKU("NEG")
	cleanupSKU(t, pool, sku)

	// Sengaja langsung ke database, melewati service.
	_, err := pool.Exec(context.Background(),
		"INSERT INTO products (sku, name, stock) VALUES ($1, $2, -1)", sku, "Negatif")

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.ConstraintName != "products_stock_non_negative" {
		t.Fatalf("want products_stock_non_negative violation, got %v", err)
	}
}
