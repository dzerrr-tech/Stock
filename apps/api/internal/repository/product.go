package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

type ListProductsParams struct {
	Search string
	Limit  int
	Offset int
}

func (r *ProductRepository) List(ctx context.Context, p ListProductsParams) ([]model.Product, int, error) {
	pattern := "%" + escapeLike(p.Search) + "%"

	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT count(*) FROM products WHERE sku ILIKE $1 OR name ILIKE $1`, pattern).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, sku, name, stock, min_stock, created_at, updated_at
		FROM products
		WHERE sku ILIKE $1 OR name ILIKE $1
		ORDER BY id DESC
		LIMIT $2 OFFSET $3`, pattern, p.Limit, p.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	products := make([]model.Product, 0, p.Limit)
	for rows.Next() {
		var m model.Product
		if err := rows.Scan(&m.ID, &m.SKU, &m.Name, &m.Stock, &m.MinStock, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, m)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate products: %w", err)
	}
	return products, total, nil
}

// escapeLike menetralkan karakter khusus LIKE agar input user dicari apa adanya.
func escapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}
