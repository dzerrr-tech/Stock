package model

import "time"

type Product struct {
	ID        int64
	SKU       string
	Name      string
	Stock     int
	MinStock  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
