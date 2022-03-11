package model

import (
	"github.com/ashah360/foodquest-api/internal/api/db"
	"time"
)

type Rating struct {
	UserID        string        `json:"userId" db:"user_id"`
	RestaurantID  string        `json:"restaurantId" db:"restaurant_id"`
	Stars         int           `json:"stars" db:"stars"`
	RatingComment db.NullString `json:"ratingComment,omitempty" db:"rating_comment"`
	CreatedAt     *time.Time    `json:"createdAt" db:"created_at"`
}
