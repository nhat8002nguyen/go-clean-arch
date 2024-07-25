package domain

import (
	"time"
)

// User represents a user entity
type User struct {
	ID          int64     `json:"id"`
	FullName    string    `json:"full_name" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	CountryCode string    `json:"country_code" validate:"required"`
}
