package domain

// OrderItem represents an order item entity
type OrderItem struct {
	ID       int64   `json:"id"`
	Order    Order   `json:"order"`
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}
