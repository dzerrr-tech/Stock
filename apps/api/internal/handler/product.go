package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"Stock/apps/api/internal/model"
	"Stock/apps/api/internal/service"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

type createProductRequest struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	MinStock int    `json:"min_stock"`
}

type productResponse struct {
	ID        int64     `json:"id"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	Stock     int       `json:"stock"`
	MinStock  int       `json:"min_stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toProductResponse(p model.Product) productResponse {
	return productResponse{
		ID: p.ID, SKU: p.SKU, Name: p.Name, Stock: p.Stock,
		MinStock: p.MinStock, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func writeError(c *gin.Context, status int, code, message, field string) {
	body := gin.H{"code": code, "message": message}
	if field != "" {
		body["field"] = field
	}
	c.JSON(status, gin.H{"error": body})
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid_json", "request body must be valid JSON", "")
		return
	}

	p, err := h.service.Create(c.Request.Context(), service.CreateProductInput{
		SKU:      req.SKU,
		Name:     req.Name,
		MinStock: req.MinStock,
	})
	if err != nil {
		var ve service.ValidationError
		switch {
		case errors.As(err, &ve):
			writeError(c, http.StatusUnprocessableEntity, "validation_error", ve.Message, ve.Field)
		case errors.Is(err, service.ErrDuplicateSKU):
			writeError(c, http.StatusConflict, "duplicate_sku", "a product with this SKU already exists", "sku")
		default:
			log.Printf("create product: %v", err)
			writeError(c, http.StatusInternalServerError, "internal_error", "something went wrong", "")
		}
		return
	}

	c.JSON(http.StatusCreated, toProductResponse(p))
}
