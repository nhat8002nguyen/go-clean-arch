package domain

import (
	"time"
)

// Order represents an order entity
type Order struct {
	ID          int64       `json:"id"`
	User        User        `json:"user"`
	Status      string      `json:"status" validate:"required"`
	CreatedAt   time.Time   `json:"created_at"`
	TotalAmount float64     `json:"total_amount"`
	Items       []OrderItem `json:"items"`
}
