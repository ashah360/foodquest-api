package model

import "github.com/ashah360/foodquest-api/internal/api/db"

type Menu struct {
	ID           string        `json:"id" db:"id"`
	RestaurantID string        `json:"restaurantId" db:"restaurant_id"`
	MenuName     string        `json:"menuName" db:"menu_name"`
	Description  db.NullString `json:"description" db:"description"`
}

type MenuItem struct {
	ID              string        `json:"id" db:"id"`
	MenuID          string        `json:"menuId" db:"menu_id"`
	Section         db.NullString `json:"section" db:"section"`
	Title           string        `json:"title" db:"title"`
	ItemDescription db.NullString `json:"itemDescription" db:"item_description"`
	Price           int           `json:"price" db:"price"`
	ImageURL        db.NullString `json:"imageUrl" db:"image_url"`
}

type RestaurantMenuItem struct {
	RestaurantID    string        `json:"restaurantId" db:"restaurant_id"`
	MenuID          string        `json:"menuId" db:"menu_id"`
	MenuName        string        `json:"menuName" db:"menu_name"`
	MenuDescription db.NullString `json:"menuDescription" db:"menu_description"`
	MenuItemID      string        `json:"menuItemId" db:"menu_item_id"`
	Title           string        `json:"title" db:"title"`
	Section         db.NullString `json:"section" db:"section"`
	ItemDescription db.NullString `json:"itemDescription" db:"item_description"`
	Price           int           `json:"price" db:"price"`
	ImageURL        db.NullString `json:"imageUrl" db:"image_url"`
}
