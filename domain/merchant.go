package domain

import (
	"time"
)

// Merchant represents a merchant entity
type Merchant struct {
	ID           int64            `json:"id"`
	CountryCode  string           `json:"country_code" validate:"required"`
	MerchantName string           `json:"merchant_name" validate:"required"`
	CreatedAt    time.Time        `json:"created_at"`
	Admin        User             `json:"admin"`
	Periods      []MerchantPeriod `json:"periods"`
}
