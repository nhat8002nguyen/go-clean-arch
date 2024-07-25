package domain

import "time"

// MerchantPeriod represents a merchant period value object
type MerchantPeriod struct {
	Merchant    Merchant  `json:"merchant"`
	CountryCode string    `json:"country_code"` // Or use a specific type for country codes
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}
