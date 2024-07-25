package domain

// ProductTag represents a product tag entity
type ProductTag struct {
	ID   int64  `json:"id"`
	Name string `json:"name" validate:"required"`
}
