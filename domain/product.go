package domain

import "time"

// Product represents a product entity
type Product struct {
	ID        int64            `json:"id"`
	Name      string           `json:"name" validate:"required"`
	Merchant  Merchant         `json:"merchant"`
	Price     float64          `json:"price"` // Or use a specific type for money if needed
	Status    string           `json:"status" validate:"required"`
	CreatedAt time.Time        `json:"created_at"`
	Catalogs  []ProductCatalog `json:"catalogs"`
	Tags      []ProductTag     `json:"tags"`
}
