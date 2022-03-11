package model

import "github.com/ashah360/foodquest-api/internal/api/db"

type UserAddress struct {
	UserID       string        `json:"userId" db:"user_id"`
	AddressLine1 string        `json:"addressLine1" db:"address_line_1"`
	AddressLine2 db.NullString `json:"addressLine2,omitempty" db:"address_line_2"`
	AddressLine3 db.NullString `json:"addressLine3,omitempty" db:"address_line_3"`
	City         string        `json:"city" db:"city"`
	State        string        `json:"state" db:"state"`
	PostalCode   string        `json:"postalCode" db:"postal_code"`
	Country      string        `json:"country" db:"country"`
}
