package domain

// ProductCatalog represents a product catalog entity
type ProductCatalog struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}
