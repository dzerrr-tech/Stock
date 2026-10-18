package service

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestProductService_Create_Validation(t *testing.T) {
	s := NewProductService(nil)

	tests := []struct {
		name      string
		in        CreateProductInput
		wantField string
	}{
		{"empty sku", CreateProductInput{SKU: "", Name: "Baut"}, "sku"},
		{"blank sku", CreateProductInput{SKU: "   ", Name: "Baut"}, "sku"},
		{"sku with space inside", CreateProductInput{SKU: "AB 12", Name: "Baut"}, "sku"},
		{"sku with symbol", CreateProductInput{SKU: "AB#12", Name: "Baut"}, "sku"},
		{"sku too long", CreateProductInput{SKU: strings.Repeat("A", 51), Name: "Baut"}, "sku"},
		{"blank name", CreateProductInput{SKU: "AB-1", Name: "  "}, "name"},
		{"name too long", CreateProductInput{SKU: "AB-1", Name: strings.Repeat("x", 201)}, "name"},
		{"negative min stock", CreateProductInput{SKU: "AB-1", Name: "Baut", MinStock: -1}, "min_stock"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.Create(context.Background(), tc.in)

			var ve *ValidationError
			if !errors.As(err, &ve) {
				t.Fatalf("want *ValidationError, got %v", err)
			}
			if ve.Field != tc.wantField {
				t.Errorf("field = %q, want %q", ve.Field, tc.wantField)
			}
		})
	}
}
