package model

import (
	"github.com/ashah360/foodquest-api/internal/api/db"
	"time"
)

type Restaurant struct {
	ID             string        `json:"id" db:"id"`
	OwnerID        string        `json:"ownerId" db:"owner_id"`
	RestaurantName string        `json:"restaurantName" db:"restaurant_name"`
	Category       string        `json:"category" db:"category"`
	OpenTime       *time.Time    `json:"openTime" db:"open_time"`
	CloseTime      *time.Time    `json:"closeTime" db:"close_time"`
	AddressLine1   string        `json:"addressLine1" db:"address_line_1"`
	AddressLine2   db.NullString `json:"addressLine2,omitempty" db:"address_line_2"`
	AddressLine3   db.NullString `json:"addressLine3,omitempty" db:"address_line_3"`
	City           string        `json:"city" db:"city"`
	State          string        `json:"state" db:"state"`
	PostalCode     string        `json:"postalCode" db:"postal_code"`
	Country        string        `json:"country" db:"country"`
	ImageURL       string        `json:"imageUrl" db:"image_url"`
}

type RestaurantPageData struct {
	Restaurant
	Stars      *float64 `json:"stars" db:"stars"`
	NumRatings *int     `json:"numRatings" db:"num_ratings"`
}

type BestSellingItem struct {
	ID        string `json:"id" db:"id"`
	ItemName  string `json:"itemName" db:"item_name"`
	Quantity  int    `json:"quantity" db:"quantity"`
	TotalCost int    `json:"totalCost" db:"total_cost"`
}
