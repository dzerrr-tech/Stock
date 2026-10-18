package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
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
		var ve *service.ValidationError
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

type paginationMeta struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type listProductsResponse struct {
	Data []productResponse `json:"data"`
	Meta paginationMeta    `json:"meta"`
}

// queryInt membaca query param bertipe integer. Param yang tidak ada dihitung 0.
func queryInt(c *gin.Context, name string) (int, error) {
	raw := c.Query(name)
	if raw == "" {
		return 0, nil
	}
	return strconv.Atoi(raw)
}

func (h *ProductHandler) List(c *gin.Context) {
	page, err := queryInt(c, "page")
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid_query", "page must be an integer", "page")
		return
	}
	pageSize, err := queryInt(c, "page_size")
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid_query", "page_size must be an integer", "page_size")
		return
	}

	result, err := h.service.List(c.Request.Context(), service.ListProductsInput{
		Search:   c.Query("search"),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		var ve *service.ValidationError
		if errors.As(err, &ve) {
			writeError(c, http.StatusUnprocessableEntity, "validation_error", ve.Message, ve.Field)
			return
		}
		log.Printf("list products: %v", err)
		writeError(c, http.StatusInternalServerError, "internal_error", "something went wrong", "")
		return
	}

	data := make([]productResponse, 0, len(result.Products))
	for _, p := range result.Products {
		data = append(data, toProductResponse(p))
	}

	c.JSON(http.StatusOK, listProductsResponse{
		Data: data,
		Meta: paginationMeta{Page: result.Page, PageSize: result.PageSize, Total: result.Total},
	})
}
