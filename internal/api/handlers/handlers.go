package handlers

import (
	"github.com/ashah360/foodquest-api/internal/api/service"
	"github.com/jmoiron/sqlx"
)

type HandlerGroup struct {
	db                *sqlx.DB
	restaurantService service.RestaurantService
	userService       service.UserService
}

func NewHandlerGroup(db *sqlx.DB, rs service.RestaurantService, us service.UserService) *HandlerGroup {
	return &HandlerGroup{
		db,
		rs,
		us,
	}
}
