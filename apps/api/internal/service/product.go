package service

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"

	"Stock/apps/api/internal/model"
	"Stock/apps/api/internal/repository"
)

var ErrDuplicateSKU = errors.New("sku already exists")

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

var skuPattern = regexp.MustCompile(`^[A-Z0-9][A-Z0-9_-]*$`)

type CreateProductInput struct {
	SKU      string
	Name     string
	MinStock int
}

type ProductService struct {
	products *repository.ProductRepository
}

func NewProductService(products *repository.ProductRepository) *ProductService {
	return &ProductService{products: products}
}

func (s *ProductService) Create(ctx context.Context, in CreateProductInput) (model.Product, error) {
	sku := strings.ToUpper(strings.TrimSpace(in.SKU))
	name := strings.TrimSpace(in.Name)

	if sku == "" {
		return model.Product{}, &ValidationError{Field: "sku", Message: "is required"}
	}
	if len(sku) > 50 || !skuPattern.MatchString(sku) {
		return model.Product{}, &ValidationError{Field: "sku", Message: "must be 1-50 chars: letters, digits, '-' or '_'"}
	}
	if name == "" {
		return model.Product{}, &ValidationError{Field: "name", Message: "is required"}
	}
	if utf8.RuneCountInString(name) > 200 {
		return model.Product{}, &ValidationError{Field: "name", Message: "must be at most 200 characters"}
	}
	if in.MinStock < 0 {
		return model.Product{}, &ValidationError{Field: "min_stock", Message: "must not be negative"}
	}

	p, err := s.products.Create(ctx, sku, name, in.MinStock)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateSKU) {
			return model.Product{}, ErrDuplicateSKU
		}
		return model.Product{}, err
	}
	return p, nil
}

const (
	defaultPageSize = 20
	maxPageSize     = 100
	maxPage         = 1_000_000
	maxSearchLength = 100
)

type ListProductsInput struct {
	Search   string
	Page     int
	PageSize int
}

type ListProductsResult struct {
	Products []model.Product
	Total    int
	Page     int
	PageSize int
}

func (s *ProductService) List(ctx context.Context, in ListProductsInput) (ListProductsResult, error) {
	search := strings.TrimSpace(in.Search)
	if utf8.RuneCountInString(search) > maxSearchLength {
		return ListProductsResult{}, &ValidationError{Field: "search", Message: "must be at most 100 characters"}
	}

	page := in.Page
	if page == 0 {
		page = 1
	}
	if page < 1 || page > maxPage {
		return ListProductsResult{}, &ValidationError{Field: "page", Message: "must be between 1 and 1000000"}
	}

	pageSize := in.PageSize
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	if pageSize < 1 || pageSize > maxPageSize {
		return ListProductsResult{}, &ValidationError{Field: "page_size", Message: "must be between 1 and 100"}
	}

	products, total, err := s.products.List(ctx, repository.ListProductsParams{
		Search: search,
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return ListProductsResult{}, err
	}

	return ListProductsResult{Products: products, Total: total, Page: page, PageSize: pageSize}, nil
}
