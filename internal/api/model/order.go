package model

import (
	"github.com/ashah360/foodquest-api/internal/api/db"
	"time"
)

type OrderDetails struct {
	ID            string     `json:"id" db:"id"`
	UserID        string     `json:"userId" db:"user_id"`
	RestaurantId  string     `json:"restaurantId" db:"restaurant_id"`
	CreatedAt     *time.Time `json:"createdAt" db:"created_at"`
	Total         int        `json:"total" db:"total"`
	Confirmed     bool       `json:"confirmed" db:"confirmed"`
	Status        string     `json:"status" db:"status"`
	EstimatedMins *int       `json:"estimatedMins" db:"estimated_mins"`
}

type OrderInvoice struct {
	OrderDetails
	RestaurantName string           `json:"restaurantName" db:"restaurant_name"`
	LineItems      []*OrderLineItem `json:"lineItems"`
}

type OrderLineItem struct {
	OrderID         string        `json:"orderId" db:"order_id"`
	MenuItemID      string        `json:"menuItemId" db:"menu_item_id"`
	Title           string        `json:"title" db:"title"`
	ItemDescription db.NullString `json:"itemDescription" db:"item_description"`
	Quantity        int           `json:"quantity" db:"quantity"`
	Price           int           `json:"price" db:"price"`
	ImageURL        db.NullString `json:"imageUrl" db:"image_url"`
}
