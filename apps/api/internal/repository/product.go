package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"Stock/apps/api/internal/model"
)

var ErrDuplicateSKU = errors.New("sku already exists")

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) Create(ctx context.Context, sku, name string, minStock int) (model.Product, error) {
	const query = `
		INSERT INTO products (sku, name, min_stock)
		VALUES ($1, $2, $3)
		RETURNING id, sku, name, stock, min_stock, created_at, updated_at`

	var p model.Product
	err := r.pool.QueryRow(ctx, query, sku, name, minStock).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Stock, &p.MinStock, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "products_sku_unique" {
			return model.Product{}, ErrDuplicateSKU
		}
		return model.Product{}, fmt.Errorf("insert product: %w", err)
	}
	return p, nil
}
